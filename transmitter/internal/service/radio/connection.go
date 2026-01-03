package radio

import (
	"fmt"

	"github.com/google/uuid"
)

const ConnectionBufferLength int = 10_000

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
	fmt.Printf("DEBUG BC: %.2f%% %d\\%d\n", float64(len(c.BroadcastChannel))/float64(cap(c.BroadcastChannel)), len(c.BroadcastChannel), cap(c.BroadcastChannel))
	fmt.Printf("DEBUG DATA: %.2f%% %d\\%d\n", float64(len(c.DataChan))/float64(cap(c.DataChan)), len(c.DataChan), cap(c.DataChan))
}

func (c *Connection) ConnectionLoop() error {
	for {
		select {
		case data := <-c.BroadcastChannel:
			dataCopy := make([]byte, len(data))
			copy(dataCopy, data)

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
