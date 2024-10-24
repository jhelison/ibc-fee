// Simplified implementation of packets and the IBC modules
// This implementation doesn't consider relayers usage and has a simple flow
package ibc

import (
	"crypto/sha256"
	"fmt"
	"time"
)

const (
	// Errors for the app
	ErrReceivePacketWithEmptyQueue = "no packet to receive"
	ErrRcvPacketWithWrongSequence  = "packet sequence mismatch. expected: %d, got: %d"
)

var _ ICSInterface = (*App)(nil)

// App is representation of a App with the package handling
type App struct {
	// The height of the app
	Height uint64
	// The list of packets
	PacketQueue []*Packet
	// Data is the received Tx
	Data map[string][]byte
	// The number of the next packet to send
	NextSeqSend uint64
	// The number of the next packet to receive
	NextSeqRecv uint64
}

// NewApp returns a new app for the simulation
func NewApp() *App {
	return &App{
		Height:      0,
		PacketQueue: []*Packet{},
		NextSeqSend: 1,
		NextSeqRecv: 1,
		Data:        make(map[string][]byte),
	}
}

// SendPacket simulates a send packet action from the IBC module
// Original interface can be found at ibc-go/modules/core/05-port/types/module.go
func (sourceApp *App) SendPacket(destApp *App, data []byte, timeoutHeight uint64, timeout time.Duration) *Packet {
	// Create a new packet with the data
	newPacket := NewPacket(sourceApp.NextSeqSend, data, timeoutHeight, timeout)

	// Update the internal sequence counter for send
	sourceApp.NextSeqSend++

	// We can image the TX will be added on the queue to of packets to receive on the dst chain
	destApp.PacketQueue = append(destApp.PacketQueue, newPacket)

	return newPacket
}

// RecvPacket receiver a packet in the destination chain
func (destApp *App) RecvPacket(sourceApp *App) error {
	// Check if the chain is expecting a packet
	if len(destApp.PacketQueue) == 0 {
		return fmt.Errorf(ErrReceivePacketWithEmptyQueue)
	}

	packet := destApp.PacketQueue[0]
	// Check if the package number is correct
	if packet.Sequence != destApp.NextSeqRecv {
		return fmt.Errorf(ErrRcvPacketWithWrongSequence, destApp.NextSeqRecv, packet.Sequence)
	}

	// Remove the packet from the queue
	destApp.PacketQueue = destApp.PacketQueue[1:]
	destApp.NextSeqRecv++

	// Register the packet data as a tx
	txHash := HashTx(packet.Data)
	destApp.Data[txHash] = packet.Data

	// Already send the ack
	return sourceApp.AcknowledgePacket(packet)
}

// AcknowledgePacket tag the packet as acknowledged
func (sourceApp *App) AcknowledgePacket(packet *Packet) error {
	// Flip the flag
	packet.Acknowledged = true

	return nil
}

// CheckTimeout checks if any packet has timed out
func (app *App) CheckTimeout() {
	for i, packet := range app.PacketQueue {
		if time.Now().After(packet.TimeoutTimestamp) { // Disconsider timezones for this simulation
			// Remove the timed out packet
			app.PacketQueue = append(app.PacketQueue[:i], app.PacketQueue[i+1:]...)
		}
	}
}

// AdvanceHeight advance the chain heigh for simulations
func (app *App) AdvanceHeight() {
	app.Height++
}

// hashTx hashes a TX with Sh256 to generate a key for the state
// Instead of using this we could use a strategy design pattern to make the hashing more flexible
func HashTx(tx []byte) string {
	// Create a new sha256 and write the TX
	h := sha256.New()
	h.Write(tx)

	// Return the stringified version of the tx
	return string(h.Sum(nil))
}
