package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"regexp"
)

var (
	delimiter = ":-"
	regMac    = regexp.MustCompile(`^([0-9a-fA-F]{2}[` + delimiter + `]){5}([0-9a-fA-F]{2})$`)
)

// MACAddress represents a 6 byte network mac address.
type MACAddress [6]byte

// A MagicPacket is constituted of 6 bytes of 0xFF followed by 16-groups of the
// destination MAC address.
type MagicPacket struct {
	header  [6]byte
	payload [16]MACAddress
}

func Wol(addr string, mac string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	// Build the magic packet.
	mp, err := New(mac)
	if err != nil {
		return err
	}

	// Grab a stream of bytes to send.
	bs, err := mp.Marshal()
	if err != nil {
		return err
	}

	// Grab a UDP connection to send our packet of bytes.
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("Attempting to send a magic packet to MAC %s\n", mac)
	fmt.Printf("... Broadcasting to: %s\n", addr)
	n, err := conn.Write(bs)
	if err == nil && n != 102 {
		err = fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Magic packet sent successfully to %s\n", mac)
	return nil
}

// New returns a magic packet based on a mac address string.
func New(mac string) (*MagicPacket, error) {
	var packet MagicPacket
	var macAddr MACAddress

	hwAddr, err := net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}

	// We only support 6 byte MAC addresses since it is much harder to use the
	// binary.Write(...) interface when the size of the MagicPacket is dynamic.
	if !regMac.MatchString(mac) {
		return nil, fmt.Errorf("%s is not a IEEE 802 MAC-48 address", mac)
	}

	// Copy bytes from the returned HardwareAddr -> a fixed size MACAddress.
	for idx := range macAddr {
		macAddr[idx] = hwAddr[idx]
	}

	// Setup the header which is 6 repetitions of 0xFF.
	for idx := range packet.header {
		packet.header[idx] = 0xFF
	}

	// Setup the payload which is 16 repetitions of the MAC addr.
	for idx := range packet.payload {
		packet.payload[idx] = macAddr
	}

	return &packet, nil
}

// Marshal serializes the magic packet structure into a 102 byte slice.
func (mp *MagicPacket) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, mp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

const (
	PORT = ":9"                //wol端口，默认都是9
	MAC  = "70:85:C2:76:60:7F" //电脑mac地址
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("args too less, usage wol $ip $mac")
		return
	}
	fmt.Println(os.Args)
	err := Wol(os.Args[1]+PORT, os.Args[2])
	if err != nil {
		fmt.Println("wake on lan failed, error:" + err.Error())
	}
}
