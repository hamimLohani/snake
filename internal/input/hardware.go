package input

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"snake/internal/engine"
)

// HardwareProtocol defines supported transmission protocols.
type HardwareProtocol string

const (
	ProtocolUDP    HardwareProtocol = "udp"
	ProtocolSerial HardwareProtocol = "serial"
)

// HardwareInputListener listens for inputs from an ESP8266/ESP32 hardware controller.
type HardwareInputListener struct {
	Addr string
}

// NewHardwareInputListener creates a listener on given address (e.g., ":4242").
func NewHardwareInputListener(addr string) *HardwareInputListener {
	if addr == "" {
		addr = ":4242"
	}
	return &HardwareInputListener{Addr: addr}
}

// StartUDP starts listening on the UDP port and pushes decoded Directions to outChan.
func (h *HardwareInputListener) StartUDP(outChan chan<- engine.Direction, stopChan <-chan struct{}) error {
	pc, err := net.ListenPacket("udp", h.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind UDP port %s: %w", h.Addr, err)
	}

	go func() {
		<-stopChan
		pc.Close()
	}()

	go func() {
		buf := make([]byte, 128)
		for {
			n, _, err := pc.ReadFrom(buf)
			if err != nil {
				return // Closed
			}
			if n > 0 {
				dir := ParseControllerPacket(buf[:n])
				if dir != engine.DirNone {
					select {
					case outChan <- dir:
					default:
					}
				}
			}
		}
	}()

	return nil
}

// ParseControllerPacket converts binary or text payloads into engine Directions.
// Supported payloads:
// Binary: 0x01=UP, 0x02=DOWN, 0x03=LEFT, 0x04=RIGHT
// ASCII:  "UP", "DOWN", "LEFT", "RIGHT", "U", "D", "L", "R", "W", "S", "A", "D"
func ParseControllerPacket(data []byte) engine.Direction {
	if len(data) == 0 {
		return engine.DirNone
	}

	// 1. Binary single-byte protocol
	if len(data) == 1 {
		switch data[0] {
		case 0x01:
			return engine.DirUp
		case 0x02:
			return engine.DirDown
		case 0x03:
			return engine.DirLeft
		case 0x04:
			return engine.DirRight
		}
	}

	// 2. Text / ASCII protocol
	text := strings.TrimSpace(strings.ToUpper(string(data)))
	switch text {
	case "UP", "U", "W":
		return engine.DirUp
	case "DOWN", "D", "S":
		return engine.DirDown
	case "LEFT", "L", "A":
		return engine.DirLeft
	case "RIGHT", "R", "D_KEY":
		return engine.DirRight
	default:
		// Check first token in case of newline or space
		scanner := bufio.NewScanner(strings.NewReader(text))
		if scanner.Scan() {
			token := scanner.Text()
			switch token {
			case "UP", "U":
				return engine.DirUp
			case "DOWN", "D":
				return engine.DirDown
			case "LEFT", "L":
				return engine.DirLeft
			case "RIGHT", "R":
				return engine.DirRight
			}
		}
	}

	return engine.DirNone
}
