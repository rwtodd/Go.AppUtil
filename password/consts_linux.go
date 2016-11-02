// I lifted this kind-of from util_linux from /x/crypto/ssh/terminal

// +build linux

package password

import "syscall"

const ioctlReadTermios = syscall.TCGETS
const ioctlWriteTermios = syscall.TCSETS
