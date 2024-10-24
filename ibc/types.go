package ibc

import "time"

// Packet is a simplified implementation of a pack based on:
// ibc-go/modules/core/04-channel/types/channel.pb.go
type Packet struct {
	// This is the current packet sequence
	Sequence uint64
	// The timeout height of the packet
	TimeoutHeight uint64
	// The timeout timestamp of the packet
	TimeoutTimestamp time.Time
	// The data the package has
	Data []byte
	// If the packet has been received on the destination chain
	Acknowledged bool
}

// NewPacket returns a new packet
// The timeout is calculated based on the now with no timezone
func NewPacket(sequence uint64, data []byte, timeoutHeight uint64, timeout time.Duration) *Packet {
	return &Packet{
		Sequence:         sequence,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: time.Now().Add(timeout), // Simplification with not timezone consideration
		Data:             data,
	}
}
