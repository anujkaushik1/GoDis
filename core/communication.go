package core

import "syscall"

type FileDescriptor struct {
	FD int
}

func (fd FileDescriptor) Read(buffer []byte) (int, error) {
	return syscall.Read(fd.FD, buffer)
}

func (fd FileDescriptor) Write(buffer []byte) (int, error) {
	return syscall.Write(fd.FD, buffer)
}

func (fd FileDescriptor) Close() error {
	return syscall.Close(fd.FD)
}
