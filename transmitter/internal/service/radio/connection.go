package radio

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	ConnectionBufferLength int = 1000
	// In Microseconds
	ZeroSleepInterval                int = 0
	LowBufferUsageSleepInterval      int = 250
	MediumBufferUsageSleepInterval   int = 750
	HighBufferUsageSleepInterval     int = 1000
	CriticalBufferUsageSleepInterval int = 2500
)

type Connection struct {
	ID               uuid.UUID
	DataChan         chan []byte
	DoneChan         chan struct{}
	BroadcastChannel chan []byte
}

func NewConnection() *Connection {
	return &Connection{
		ID:               uuid.New(),
		DataChan:         make(chan []byte, ConnectionBufferLength),
		DoneChan:         make(chan struct{}),
		BroadcastChannel: make(chan []byte, ConnectionBufferLength),
	}
}

func (c *Connection) Debug() {
	fmt.Printf("DEBUG BC: %.2f%%\n", float64(len(c.BroadcastChannel))/float64(cap(c.BroadcastChannel)))
	fmt.Printf("DEBUG DATA: %.2f%%\n", float64(len(c.DataChan))/float64(cap(c.DataChan)))
}

func (c *Connection) ConnectionLoop() error {
	for {
		select {
		case data := <-c.BroadcastChannel:
			dataCopy := make([]byte, len(data))
			copy(dataCopy, data)

			sleepInterval := c.calculateSleepInteval()
			time.Sleep(sleepInterval)

			select {
			case c.DataChan <- dataCopy:
				// Successfully sent to client
			default:
				// Client buffer is full, skip this chunk to avoid blocking
			}
		case <-c.DoneChan:
			return nil
		}
	}
}

func (c *Connection) calculateSleepInteval() time.Duration {
	bufferUsage := float64(len(c.DataChan)) / float64(cap(c.DataChan))

	switch {
	case bufferUsage < 0.15:
		return time.Duration(ZeroSleepInterval) * time.Microsecond
	case bufferUsage < 0.3:
		return time.Duration(LowBufferUsageSleepInterval) * time.Microsecond
	case bufferUsage < 0.6:
		return time.Duration(MediumBufferUsageSleepInterval) * time.Microsecond
	case bufferUsage < 0.9:
		return time.Duration(HighBufferUsageSleepInterval) * time.Microsecond
	default:
		return time.Duration(CriticalBufferUsageSleepInterval) * time.Microsecond
	}
}
