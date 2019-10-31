package main

import "encoding/binary"
import "errors"
import "fmt"
import "syscall"

func getMacDestination(frame []byte) (string) {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		frame[0], frame[1], frame[2], frame[3], frame[4], frame[5])
}

func getMacSource(frame []byte) (string) {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		frame[6], frame[7], frame[8], frame[9], frame[10], frame[11])
}

func getEthernetType(frame []byte) (uint16) {
	return binary.BigEndian.Uint16(frame[12:14])
}

func isMulticast(frame []byte) (bool) {
	return (frame[0] & 1) == 1
}

func isBroadcast(frame []byte) (bool) {
	return getMacSource(frame) == "ff:ff:ff:ff:ff:ff"
}

func tagFrame(frame []byte, vlanId uint8) ([]byte){
	var tpidh byte
	tpidh = (syscall.ETH_P_8021Q >> 8) & 0xff
	var tpidl byte
	tpidl = syscall.ETH_P_8021Q & 0xff

	tag := []byte { tpidh, tpidl, 0x00, vlanId }
	return append(frame[:12], append(tag, frame[12:]...)...)
}

func getVlanId(frame []byte) (uint8) {
	return frame[15]
}

func untagFrame(frame []byte) ([]byte) {
	return append(frame[:12], frame[16:]...)
}

func kernelUntagged(auxdata []byte) (bool, uint8) {
	vlanId := auxdata[32]
	return vlanId != 0, vlanId
}

func handleFrame(frame []byte, auxdata []byte, port Port, fib map[string]Port) ([]byte, *Port, error) {
	mac_dst := getMacDestination(frame)
	mac_src := getMacSource(frame)

	if !isMulticast(frame) && !isBroadcast(frame) {
		fib[mac_src] = port
	}

	//HACK
	untagged, vlanId := kernelUntagged(auxdata)
	if untagged {
		frame = tagFrame(frame, vlanId)
	}

	if (port.mode == Access) {
		frame = tagFrame(frame, port.vlanId)
	}

	if port, ok := fib[mac_dst]; ok {
		return frame, &port, nil
	} else {
		return frame, nil, errors.New("No FIB entry found.")
	}
}

