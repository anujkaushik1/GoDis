package main

import (
	"fmt"
	"syscall"
	"time"
)

func singleThreadedNonBlockingMultipleFiles() {
	fd1, err := syscall.Open("test1.txt", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	fd2, err := syscall.Open("test2.txt", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("fd1 =", fd1, "fd2 =", fd2)

	kq, err := syscall.Kqueue()
	if err != nil {
		panic(err)
	}

	event1 := syscall.Kevent_t{
		Ident:  uint64(fd1),
		Filter: syscall.EVFILT_VNODE,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE | syscall.EV_CLEAR,
		Fflags: syscall.NOTE_WRITE,
	}

	event2 := syscall.Kevent_t{
		Ident:  uint64(fd2),
		Filter: syscall.EVFILT_VNODE,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE | syscall.EV_CLEAR,
		Fflags: syscall.NOTE_WRITE,
	}

	_, err = syscall.Kevent(kq, []syscall.Kevent_t{event1, event2}, nil, nil)
	if err != nil {
		panic(err)
	}

	count := 1

	for count < 100 {

		fmt.Println("Watching files...")

		events := make([]syscall.Kevent_t, 2)

		timeout := syscall.Timespec{Sec: 20, Nsec: 0}

		n, err := syscall.Kevent(kq, nil, events, &timeout)
		if err != nil {
			panic(err)
		}

		// 🔥 Check which file triggered
		for i := 0; i < n; i++ {
			if events[i].Ident == uint64(fd1) {
				fmt.Println("🚨 file1.txt modified")
			}
			if events[i].Ident == uint64(fd2) {
				fmt.Println("🚨 file2.txt modified")
			}
		}

		// your normal work
		fmt.Println("count =", count)
		count++

		time.Sleep(2 * time.Second)
	}
}

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
	singleThreadedNonBlockingMultipleFiles()
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
