package sever

import (
	"fmt"
	"log"
	"net"
	"syscall"

	config "github.com/anujkaushik1/GoDis/Config"
)

func RunAsyncTcpServer() error {
	host := config.App.Host
	port := config.App.Port
	log.Println("Starting TCP server on: ", host, ":", port)

	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD)
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	ip4 := net.ParseIP(host).To4()
	err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	})
	if err != nil {
		return err
	}

	max_clients := 20000
	if err = syscall.Listen(serverFD, max_clients); err != nil {
		return err
	}

	kq, err := syscall.Kqueue()
	if err != nil {
	}

	event := syscall.Kevent_t{
		Ident:  uint64(serverFD),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
	}

	_, err = syscall.Kevent(kq, []syscall.Kevent_t{event}, nil, nil)
	if err != nil {
		return err

	}

	fmt.Println("Watching for client connections...")

	events := make([]syscall.Kevent_t, max_clients)

	for {
		_, err := syscall.Kevent(kq, nil, events, nil)
		if err != nil {
			fmt.Println("Error in event:: ", err.Error())
			continue
		}

		fmt.Println("multiple clients might be ready")

		for {
			clientFD, _, err := syscall.Accept(serverFD)

			if err != nil {
				if err == syscall.EAGAIN {
					break // queue empty
				}

				fmt.Println("accept error:", err.Error())
				break
			}

			fmt.Println(clientFD)
		}

	}

}
