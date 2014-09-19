package testlib

/*
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

int is_sock(int mode)
{
	return S_ISSOCK(mode);
}
*/
import "C"
import (
	// "fmt"
	// "os"
	"syscall"
)

// from syscall
func fcntl(fd int, cmd int, arg int) (val int, err error) {
	r0, _, e1 := syscall.Syscall(syscall.SYS_FCNTL, uintptr(fd), uintptr(cmd), uintptr(arg))
	val = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

func IsValidFd(fd int) bool {
	ret, err := fcntl(fd, C.F_GETFD, 0)
	if ret == -1 && err == syscall.EBADF {
		return false
	}
	return true
}

func IsSock(mode uint32) bool {
	if ret, _ := C.is_sock(C.int(mode)); ret == 0 {
		return false
	}
	return true
}
