// I lifted this directly from util_bsd from /x/crypto/ssh/terminal

// +build darwin dragonfly freebsd netbsd openbsd

package password

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA
const ioctlWriteTermios = syscall.TIOCSETA
