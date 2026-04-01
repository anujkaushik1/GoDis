package main

import (
	"fmt"
	"syscall"
	"time"
)

func singleThreadedNonBlocking() {
	// Open file
	fd, err := syscall.Open("test.txt", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("fddd = ", fd)
	kq, err := syscall.Kqueue()
	if err != nil {
		panic(err)
	}

	event := syscall.Kevent_t{
		Ident:  uint64(fd),           // which file
		Filter: syscall.EVFILT_VNODE, // file events
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE | syscall.EV_CLEAR,
		Fflags: syscall.NOTE_WRITE, // notify on write
	}

	_, err = syscall.Kevent(kq, []syscall.Kevent_t{event}, nil, nil)
	if err != nil {
		panic(err)
	}

	count := 1

	for count < 10000 {
		fmt.Println("Watching file...")

		events := make([]syscall.Kevent_t, 1)

		timeout := syscall.Timespec{Sec: 1, Nsec: 0}

		n, err := syscall.Kevent(kq, nil, events, &timeout)
		if err != nil {
			panic(err)
		}

		if n > 0 {
			fmt.Println("🚨 File modified!")
		}

		fmt.Println("hel =", count)
		count++
		fmt.Println("Doing heavy stuff....")
		time.Sleep(10 * time.Second)
	}
}

func main() {
	singleThreadedNonBlocking()
	// // Open file
	// fd, err := syscall.Open("test.txt", syscall.O_RDONLY, 0)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("fddd = ", fd)
	// kq, err := syscall.Kqueue()
	// if err != nil {
	// 	panic(err)
	// }

	// event := syscall.Kevent_t{
	// 	Ident:  uint64(fd),           // which file
	// 	Filter: syscall.EVFILT_VNODE, // file events
	// 	Flags:  syscall.EV_ADD | syscall.EV_ENABLE | syscall.EV_CLEAR,
	// 	Fflags: syscall.NOTE_WRITE, // notify on write
	// }

	// _, err = syscall.Kevent(kq, []syscall.Kevent_t{event}, nil, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Watching file...")

	// for {
	// 	events := make([]syscall.Kevent_t, 1)

	// 	_, err := syscall.Kevent(kq, nil, events, nil)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println("File modified!")
	// }

	// count := 1

	// for count < 10000 {
	// 	fmt.Println("hel = ", count)
	// 	count++

	// 	time.Sleep(1 * time.Second)

	// }

}
