package monitor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func Ping(ipaddr string) bool {
	var (
		icmp  ICMP
		laddr = net.IPAddr{IP: net.ParseIP("0.0.0.0")}
		raddr = net.IPAddr{IP: net.ParseIP(ipaddr)}
	)

	conn, err := net.DialIP("ip4:icmp", &laddr, &raddr)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer conn.Close()

	icmp.Type = 8
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = CheckSum(buffer.Bytes())
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)

	recv := make([]byte, 1024)

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		fmt.Println(err.Error())
		return false
	}

	conn.SetReadDeadline((time.Now().Add(time.Second * 2)))
	_, err = conn.Read(recv)

	if err != nil {
		fmt.Println("请求超时")
		return false
	}

	return true
}
