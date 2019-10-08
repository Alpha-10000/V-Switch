package main

// #include <arpa/inet.h>
// #include <net/if.h>
// #include <net/ethernet.h>
// #include <sys/epoll.h>
// #include <sys/socket.h>
// #include <linux/if_arp.h>
// #include <linux/if_packet.h>
import "C"
import "fmt"
import "os"
import S "syscall"

type Port struct {
	intfIdx int
	intfName string
	portNum int
}

const bufSize = 1514

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func listen(port Port) (int) {

	sockFd, err := S.Socket(S.AF_PACKET, S.SOCK_RAW | S.SOCK_NONBLOCK, int(htons(S.ETH_P_ALL)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sockAddr := S.SockaddrLinklayer {
		Protocol: htons(S.ETH_P_ALL),
		Ifindex: port.intfIdx,
		Halen: C.ETH_ALEN,
	}

	err = S.Bind(sockFd, &sockAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return sockFd
}

func handleFrame(inFd int, ports map[int]Port) {
	buf := make([]byte, bufSize)
	ctrl := make([]byte, 1024)
	n, _, _, _, _ := S.Recvmsg(inFd, buf, ctrl, 0)
	for outFd, p := range(ports) {
		if outFd == inFd {
			continue
		}
		
		outSockAddr := S.SockaddrLinklayer {
			Ifindex: p.intfIdx,
			Halen: C.ETH_ALEN,
		}

		err := S.Sendto(outFd, buf, 0, &outSockAddr)
		fmt.Println("SENDTO", ports[outFd], n)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func dispatch(opts *Opts) {
	ports := make(map[int]Port)
	
	epollFd, err := S.EpollCreate1(0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for idx, intfName := range(opts.intfs) {

		intfIdx := int((C.if_nametoindex(C.CString(intfName))))
		port := Port{intfIdx, intfName, idx}
		sockFd := listen(port)
		ports[sockFd] = port
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
			fd := int(events[i].Fd)
			fmt.Println("RECVFROM", ports[fd])
			handleFrame(fd, ports)
		}
	}
}

