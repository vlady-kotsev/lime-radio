package radio

import (
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

func NewConnection(broadcastChan chan []byte) *Connection {
	return &Connection{
		ID:               uuid.New(),
		DataChan:         make(chan []byte, ConnectionBufferLength),
		DoneChan:         make(chan struct{}),
		BroadcastChannel: broadcastChan,
	}
}

func (c *Connection) ConnectionLoop() error {
	for {
		select {
		case data := <-c.BroadcastChannel:
			sleepInterval := c.calculateSleepInteval()
			time.Sleep(sleepInterval)
			c.DataChan <- data
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
