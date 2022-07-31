package nudge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	mod                     = windows.NewLazyDLL("user32.dll")
	procGetWindowText       = mod.NewProc("GetWindowTextW")
	procGetWindowTextLength = mod.NewProc("GetWindowTextLengthW")
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

func getWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

func getWindowText(hwnd HWND) string {
	textLen := getWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func getWindow(funcName string) uintptr {
	proc := mod.NewProc(funcName)
	hwnd, _, _ := proc.Call()
	return hwnd
}

func FindActiveWindowTitle() string {
	if hwnd := getWindow("GetForegroundWindow"); hwnd != 0 {
		text := getWindowText(HWND(hwnd))
		return text
	}

	return ""
}
