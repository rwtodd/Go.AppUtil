// +build windows

// Package password provides functionality to read a password
// from the current tty/console. Importantly, it does so
// even if STDIN and STDOUT have been redirected. On
// windows, it uses syscalls to open "CONIN$" and "CONOUT$".
// On unix, it opens "/dev/tty".
package password

import (
	"bytes"
	"fmt"
	"io"
	"syscall"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

const (
	enableLineInput       = 2
	enableEchoInput       = 4
	enableProcessedInput  = 1
	enableProcessedOutput = 1
)

// GetConsoleHandles opens CONIN$ and CONOUT$ on windows, which
// can be used to read and write from the console even if
// stdin and stdout have been redirected.  These handles need to
// be closed with syscall.CloseHandle().
func GetConsoleHandles() (input syscall.Handle, output syscall.Handle, err error) {
	infl, _ := syscall.UTF16PtrFromString("CONIN$")
	input, err = syscall.CreateFile(infl,
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ,
		nil,
		syscall.OPEN_EXISTING,
		0,
		0)
	if err != nil {
		return
	}

	outfl, _ := syscall.UTF16PtrFromString("CONOUT$")
	output, err = syscall.CreateFile(outfl,
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		0,
		0)
	if err != nil {
		syscall.CloseHandle(input)
		input = 0
	}

	return
}

// Read issues the given prompt, and then reads
// a line of input with no echo to the screen.  If times is
// greater than 1, it will ask the user to retype the line.
// Realistically, users will never set times greater than 2,
// but they can if they want.
func Read(prompt string, times int) (string, error) {
	retype, pbyte, newline := []byte("(retype) "), []byte(prompt), []byte("\n")
	in, out, err := GetConsoleHandles()
	if err != nil {
		return "", err
	}
	defer syscall.CloseHandle(in)
	defer syscall.CloseHandle(out)

	syscall.Write(out, pbyte)
	ans, err := readNoEcho(in)
	syscall.Write(out, newline)
	if err != nil {
		return "", err
	}

	for times > 1 {
		times--

		syscall.Write(out, retype)
		syscall.Write(out, pbyte)
		ans2, err2 := readNoEcho(in)
		syscall.Write(out, newline)
		if err2 != nil {
			return "", err2
		}

		if !bytes.Equal(ans, ans2) {
			return "", fmt.Errorf("Answers did not match!")
		}
	}

	return string(ans), nil
}

func setstate(fd syscall.Handle, mode uint32) (err error) {
	_, _, e := syscall.Syscall(procSetConsoleMode.Addr(), 2, uintptr(fd), uintptr(mode), 0)
        if e != 0 {
		err = error(e)	
        }
	return  
}

// readNoEcho reads a line of input from a terminal without local echo.  This
// is commonly used for inputting passwords and other sensitive data. The slice
// returned does not include the \n.
func readNoEcho(fd syscall.Handle) ([]byte, error) {
	var st uint32
	var err error
	if err = syscall.GetConsoleMode(fd, &st); err != nil {
		return nil, err
	}
	old := st

	st &^= (enableEchoInput)
	st |= (enableProcessedInput | enableLineInput | enableProcessedOutput)
        if err = setstate(fd, st); err != nil {
		return nil, err
	}
	defer setstate(fd, old)

	var buf [16]byte
	var ret []byte
	for {
		n, err := syscall.Read(fd, buf[:])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			if len(ret) == 0 {
				return nil, io.EOF
			}
			break
		}
		if buf[n-1] == '\n' {
			n--
		}
		if n > 0 && buf[n-1] == '\r' {
			n--
		}
		ret = append(ret, buf[:n]...)
		if n < len(buf) {
			break
		}
	}

	return ret, nil
}
