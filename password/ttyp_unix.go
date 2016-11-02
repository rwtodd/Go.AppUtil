
// +build !windows

// Package password provides functionality to read a password
// from the current tty/console. Importantly, it does so
// even if STDIN and STDOUT have been redirected. On
// windows, it uses syscalls to open "CONIN$" and "CONOUT$".
// On unix, it opens "/dev/tty".
package password

import (
	"bytes"
        "os"
	"fmt"
	"io"
	"syscall"
        "unsafe"
)

// Read issues the given prompt, and then reads
// a line of input with no echo to the screen.  If times is
// greater than 1, it will ask the user to retype the line.
// Realistically, users will never set times greater than 2,
// but they can if they want.
func Read(prompt string, times int) (string, error) {
        con, err := os.OpenFile("/dev/tty", os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	defer con.Close()

        con.WriteString(prompt)
	ans, err := readNoEcho(con.Fd())
	con.WriteString("\n")
	if err != nil {
		return "", err
	}

	for times > 1 {
		times--

                fmt.Fprintf(con, "(retype) %s",prompt)
		ans2, err2 := readNoEcho(con.Fd())
	        con.WriteString("\n")
		if err2 != nil {
			return "", err2
		}

		if !bytes.Equal(ans, ans2) {
			return "", fmt.Errorf("Answers did not match!")
		}
	}

	return string(ans), nil
}

func setstate(fd uintptr, state *syscall.Termios) error {
    _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, ioctlWriteTermios, uintptr(unsafe.Pointer(state)), 0, 0, 0)
    if err != 0 {
       return error(err)
    }
    return nil
}

func getstate(fd uintptr, state *syscall.Termios) error {
    if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, ioctlReadTermios, uintptr(unsafe.Pointer(state)), 0, 0, 0); err != 0 {
            return err
    }

    return nil
}

// readNoEcho reads a line of input from a terminal without local echo.  This
// is commonly used for inputting passwords and other sensitive data. The slice
// returned does not include the \n.
func readNoEcho(fd uintptr) ([]byte, error) {
	var st syscall.Termios 
	var err error
	err = getstate(fd, &st)
	if err != nil {
		return nil, err
	}
	old := st

        st.Lflag &^= syscall.ECHO
        st.Lflag |= syscall.ICANON | syscall.ISIG
        st.Iflag |= syscall.ICRNL
        if err = setstate(fd, &st); err != nil {
                fmt.Fprintf(os.Stderr,"Setting terminal\n")
		return nil, err
	}
	defer setstate(fd, &old)

	var buf [16]byte
	var ret []byte
	for {
		n, err := syscall.Read(int(fd), buf[:])
		if err != nil {
                        fmt.Fprintf(os.Stderr,"Reading terminal\n")
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

		ret = append(ret, buf[:n]...)
		if n < len(buf) {
			break
		}
	}

	return ret, nil
}
