package sever

import (
	"fmt"
	"log"
	"net"
	"syscall"

	config "github.com/anujkaushik1/GoDis/Config"
)

func  RunAsyncTcpServer() error {
	fmt.Println("world")
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

	fmt.Println("kaka")

	return nil
}
