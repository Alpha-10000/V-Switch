package main

import "encoding/binary"
import "errors"
import "fmt"

func handleFrame(data []byte, port Port, fib map[string]Port) (*Port, error) {
	idx := 0
	mac_src := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", data[idx], data[idx+1], data[idx+2], data[idx+3], data[idx+4], data[idx+5])
	idx += 6
	mac_dst := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", data[idx], data[idx+1], data[idx+2], data[idx+3], data[idx+4], data[idx+5])
	idx += 6
	eth_type := binary.BigEndian.Uint16(data[idx:idx+2])
	fmt.Printf("%s %s 0x%x\n", mac_src, mac_dst, eth_type)

	// broadcast and multicast
	if (mac_src != "ff:ff:ff:ff:ff:ff") && ((data[0] & 1) != 1) {
		fib[mac_src] = port
	}

	if port, ok := fib[mac_dst]; ok {
		return &port, nil 
	} else {
		return nil, errors.New("No FIB entry found.")
	}

	
}

