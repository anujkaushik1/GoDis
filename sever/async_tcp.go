package sever

import (
	"fmt"
	"log"
	"net"
	"syscall"

	config "github.com/anujkaushik1/GoDis/config"
	"github.com/anujkaushik1/GoDis/core"
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
		return err
	}

	event := syscall.Kevent_t{ //notify when socket becomes (readable)
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

	noOfClients := 0

	for {
		n, err := syscall.Kevent(kq, nil, events, nil)
		if err != nil {
			continue
		}

		for i := 0; i < n; i++ {
			if events[i].Ident == uint64(serverFD) {
				for {
					clientFD, _, err := syscall.Accept(serverFD)
					if err != nil {
						if err == syscall.EAGAIN {
							break // queue empty
						}

						fmt.Println("accept error:", err.Error())
						break
					}

					if err = syscall.SetNonblock(clientFD, true); err != nil {
						return err
					}

					noOfClients++

					clientEvent := syscall.Kevent_t{
						Ident:  uint64(clientFD),
						Filter: syscall.EVFILT_READ,
						Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
					}

					_, err = syscall.Kevent(kq, []syscall.Kevent_t{clientEvent}, nil, nil)
					if err != nil {
						fmt.Println("Error adding client event:", err.Error())
						continue
					}
				}
			} else {
				clientFD := events[i].Ident
				clientFileDescriptorStruct := core.FileDescriptor{FD: int(clientFD)}
				redisCmd, err := ReadCommand(clientFileDescriptorStruct)

				if err != nil {
					clientFileDescriptorStruct.Close()
					noOfClients--
					log.Println("client disconnected")
					continue
				}

				err = Respond(clientFileDescriptorStruct, redisCmd)
				if err != nil {
					clientFileDescriptorStruct.Close()
					noOfClients--
					log.Println("client disconnected")
					continue
				}

			}

		}

	}

}
