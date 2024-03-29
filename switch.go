package main

// #include <arpa/inet.h>
// #include <net/if.h>
// #include <net/ethernet.h>
// #include <linux/if_packet.h>
import "C"
import "fmt"
import "os"
import S "syscall"

type Port struct {
	intfIdx int
	num int
	sockFd int
	vlanId uint8
	mode Mode
}

const bufSize = 1514

func listen(intfIndex int) (int) {

	sockFd, err := S.Socket(S.AF_PACKET, S.SOCK_RAW | S.SOCK_NONBLOCK, int(C.htons(S.ETH_P_ALL)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = S.SetsockoptInt(sockFd, S.SOL_PACKET, int(C.PACKET_AUXDATA), 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sockAddr := S.SockaddrLinklayer {
		Protocol: uint16(C.htons(S.ETH_P_ALL)),
		Ifindex: intfIndex,
		Halen: C.ETH_ALEN,
	}

	err = S.Bind(sockFd, &sockAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return sockFd
}

func transfer(data []byte, port *Port) {
	if port.mode == Access {
		if getVlanId(data) != port.vlanId {
			return
		}
		data = untagFrame(data)
	}

	outSockAddr := S.SockaddrLinklayer {
		Ifindex: port.intfIdx,
		Halen: C.ETH_ALEN,
	}

	fmt.Println("SENDTO", port, len(data))
	err := S.Sendto(port.sockFd, data, 0, &outSockAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func dispatch(opts *Opts) {
	sockToPorts := make(map[int]Port)
	FIB := make(map[string]Port)
	
	epollFd, err := S.EpollCreate1(0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for idxNum, intfName := range(opts.intfs) {

		intfIdx := int((C.if_nametoindex(C.CString(intfName))))
		sockFd := listen(intfIdx)
		vlanId, mode := configPort(intfName, opts)
		sockToPorts[sockFd] = Port{intfIdx, idxNum, sockFd, vlanId, mode}
		fmt.Println(sockToPorts[sockFd], intfName)

		event := S.EpollEvent {
			Events: (S.EPOLLIN | (S.EPOLLET & 0xffffffff)),
			Fd: int32(sockFd),
		}

		err := S.EpollCtl(epollFd, S.EPOLL_CTL_ADD, sockFd, &event)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	events := make([]S.EpollEvent, len(opts.intfs))
	for true {
		nfds, err := S.EpollWait(epollFd, events, -1)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for i := 0; i < nfds; i++ {
			inFd := int(events[i].Fd)
			buf := make([]byte, bufSize)
			ctrl := make([]byte,1024)

			n, oob, _, _, err := S.Recvmsg(inFd, buf, ctrl, 0)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("RECVFROM", sockToPorts[inFd], n, oob)

			buf, port, err := handleFrame(buf, ctrl[:oob], sockToPorts[inFd], FIB)
			if err == nil {
				transfer(buf, port)
				continue
			}
			for outFd, outPort := range(sockToPorts) {
				if outFd == inFd {
					continue
				}
				transfer(buf, &outPort)
			}
		}
	}
}

