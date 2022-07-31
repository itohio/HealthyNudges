package idle

import (
	"syscall"
	"time"
	"unsafe"
)

var user32 = syscall.MustLoadDLL("user32.dll")
var getLastInputInfo = user32.MustFindProc("GetLastInputInfo")

func lastInputTicks() (uint32, error) {
	// set up struct to contain returned information
	var lastInputInfo struct {
		cbSize uint32
		dwTime uint32
	}
	lastInputInfo.cbSize = uint32(unsafe.Sizeof(lastInputInfo))

	// call the WinAPI function with a pointer to the struct
	errBool, _, err := getLastInputInfo.Call(
		uintptr(unsafe.Pointer(&lastInputInfo)),
	)

	if errBool == 0 {
		return 0, err
	}

	return lastInputInfo.dwTime, nil
}

var kernel32 = syscall.MustLoadDLL("kernel32.dll")
var getTickCount = kernel32.MustFindProc("GetTickCount")

func currentTicks() uint32 {
	ticks, _, _ := getTickCount.Call()
	return uint32(ticks)
}

// Duration returns the time since the last user input
func Get() (time.Duration, error) {
	lastInput, err := lastInputTicks()
	if err != nil {
		return 0, err
	}
	idleTicks := currentTicks() - lastInput
	idleDuration := time.Duration(idleTicks) * time.Millisecond
	return idleDuration, nil
}
