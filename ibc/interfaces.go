package ibc

import "time"

// ICSInterface is the interface for the packet broadcasting
// It simplifies the original interface from ibc-go/modules/core/05-port/types/module.go
type ICSInterface interface {
	SendPacket(destApp *App, data []byte, timeoutHeight uint64, timeout time.Duration) *Packet
	RecvPacket(sourceApp *App) error
	AcknowledgePacket(packet *Packet) error
}
