package ibc_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"ibc-fee/ibc"
)

// TestIBCRelayingSuccess tests a successful relaying between the two chains
func TestIBCRelayingSuccess(t *testing.T) {
	// Initialize the apps
	chainA := ibc.NewApp()
	chainB := ibc.NewApp()

	// Send the packet
	tx1 := []byte("tx1")
	timeout := 5 * time.Second
	_ = chainA.SendPacket(chainB, tx1, 10, timeout)

	// Receive the packet
	err := chainB.RecvPacket(chainA)
	require.NoError(t, err)

	// Now we check the states
	require.NotEmpty(t, chainB.Data)
}

// TestIBCDifferentTime tests the relaying of packets with different block times
// First we test the normal flow
// In a second moment we test a expired packet
func TestIBCDifferentTime(t *testing.T) {
	// Initialize the apps
	chainA := ibc.NewApp()
	chainB := ibc.NewApp()

	// TEST #1 - Normal flow withing time

	// Send the packet
	tx1 := []byte("tx1")
	timeout := 50 * time.Millisecond
	_ = chainA.SendPacket(chainB, tx1, 10, timeout)

	// Chain B advance in time and check for timeout
	time.Sleep(time.Millisecond * 20)
	chainB.AdvanceHeight()
	chainB.CheckTimeout()

	// Chain B get the packet
	err := chainB.RecvPacket(chainA)
	require.NoError(t, err)

	// The packet should be received
	require.NotEmpty(t, chainB.Data)

	// TEST #2 - Package reach the timeout
	tx2 := []byte("tx2")
	_ = chainA.SendPacket(chainB, tx2, 10, timeout) // Reuse the same timeout

	// Chain B advance in time and check for timeout
	// This simulates a delay
	time.Sleep(time.Millisecond * 60)
	chainB.AdvanceHeight()
	chainB.CheckTimeout()

	// Now try to receive the packet
	err = chainB.RecvPacket(chainA)
	require.ErrorContains(t, err, "no packet to receive") // It should return the no packet to receive error

	// The data should only contain the first packet data
	require.Equal(t, 1, len(chainB.Data))
	// And it must be that the tx 1
	tx1Hash := ibc.HashTx([]byte("tx1"))
	require.Equal(t, []byte("tx1"), chainB.Data[tx1Hash])
}

// TestPacketOutOfOrder tests the ordering of packets
func TestPacketOutOfOrder(t *testing.T) {
	// Initialize the apps
	chainA := ibc.NewApp()
	chainB := ibc.NewApp()

	// Build two txs
	tx1 := []byte("tx1")
	tx2 := []byte("tx2")
	timeout := 5 * time.Second

	// Broadcast the two txs
	_ = chainA.SendPacket(chainB, tx1, 10, timeout)
	packet2 := chainA.SendPacket(chainB, tx2, 10, timeout)

	// Force the packet2 to get in the queue first
	chainB.PacketQueue = append([]*ibc.Packet{packet2}, chainB.PacketQueue...)

	// Receive the packet
	err := chainB.RecvPacket(chainA)
	require.ErrorContains(t, err, "packet sequence mismatch")

	// List of txs should be empty
	require.Empty(t, chainB.Data)
}
