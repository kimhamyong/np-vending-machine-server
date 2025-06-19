package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
)

// SyncCommand는 raw socket으로 주고받을 메시지 구조
type SyncCommand struct {
	Type  byte
	Action string
}

// RawSocketSender는 raw socket으로 메시지를 보내는 구조체
type RawSocketSender struct {
	DestAddr [4]byte
	FD       int
}

// InitRawSocketSender: 송신용 raw socket 초기화
func InitRawSocketSender(destIP string) (*RawSocketSender, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return nil, fmt.Errorf("raw socket 생성 실패: %v", err)
	}

	var addr [4]byte
	copy(addr[:], netToBytes(destIP))
	return &RawSocketSender{
		DestAddr: addr,
		FD:       fd,
	}, nil
}

// SendCommand: SyncCommand를 raw socket으로 전송
func (s *RawSocketSender) SendCommand(cmd SyncCommand) error {
	payload := encodeCommand(cmd)

	ipHeader := buildIPHeader(len(payload), s.DestAddr)
	packet := append(ipHeader, payload...)

	sa := &syscall.SockaddrInet4{Addr: s.DestAddr}
	return syscall.Sendto(s.FD, packet, 0, sa)
}

// StartRawListener: 패킷 수신 루프
func StartRawListener(handler func(cmd SyncCommand)) error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		return fmt.Errorf("raw 수신 socket 생성 실패: %v", err)
	}

	go func() {
		for {
			buf := make([]byte, 1500)
			n, _, err := syscall.Recvfrom(fd, buf, 0)
			if err != nil || n < 20+binary.Size(SyncCommand{}) {
				continue
			}

			// IP 헤더 스킵
			data := buf[20:n]

			cmd := decodeCommand(data)
			handler(cmd)
		}
	}()

	return nil
}

// ↓↓↓ Helper Functions ↓↓↓

func encodeCommand(cmd SyncCommand) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, cmd)
	return buf.Bytes()
}

func decodeCommand(data []byte) SyncCommand {
	var cmd SyncCommand
	binary.Read(bytes.NewReader(data), binary.BigEndian, &cmd)
	return cmd
}

func buildIPHeader(payloadLen int, dst [4]byte) []byte {
	totalLen := 20 + payloadLen

	ip := make([]byte, 20)
	ip[0] = 0x45                     // Version + Header Length
	ip[1] = 0x00                     // Type of Service
	ip[2] = byte(totalLen >> 8)     // Total Length
	ip[3] = byte(totalLen & 0xff)
	ip[4] = 0x00                    // ID
	ip[5] = 0x00
	ip[6] = 0x40                    // Flags + Fragment offset
	ip[7] = 0x00
	ip[8] = 64                      // TTL
	ip[9] = syscall.IPPROTO_TCP    // Protocol
	// [10,11] Checksum은 나중에
	copy(ip[12:16], []byte{127, 0, 0, 1}) // Source IP
	copy(ip[16:20], dst[:])              // Destination IP

	// Checksum 계산
	csum := calculateChecksum(ip)
	ip[10] = byte(csum >> 8)
	ip[11] = byte(csum & 0xff)

	return ip
}

func netToBytes(ip string) []byte {
	var out [4]byte
	fmt.Sscanf(ip, "%d.%d.%d.%d", &out[0], &out[1], &out[2], &out[3])
	return out[:]
}

func calculateChecksum(buf []byte) uint16 {
	var sum uint32
	for i := 0; i < len(buf)-1; i += 2 {
		sum += uint32(buf[i])<<8 | uint32(buf[i+1])
	}
	if len(buf)%2 == 1 {
		sum += uint32(buf[len(buf)-1]) << 8
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	return ^uint16(sum)
}
