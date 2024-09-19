package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[2J\033[H")
	}
}

func GetTerminalSize() (width, height int) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "$Host.UI.RawUI.WindowSize.Width")
		cmd.Stdin = os.Stdin
		out, err := cmd.Output()
		if err != nil {
			return 80, 24
		}
		fmt.Sscanf(string(out), "%d", &width)
		width = width / 2

		cmd = exec.Command("powershell", "$Host.UI.RawUI.WindowSize.Height")
		cmd.Stdin = os.Stdin
		out, err = cmd.Output()
		if err != nil {
			return 80, 24
		}
		fmt.Sscanf(string(out), "%d", &height)
		height = height - 5
	} else {
		cmd := exec.Command("tput", "cols")
		out, err := cmd.Output()
		if err != nil {
			width = 80
		} else {
			fmt.Sscanf(string(out), "%d", &width)
		}

		cmd = exec.Command("tput", "lines")
		out, err = cmd.Output()
		if err != nil {
			height = 24
		} else {
			fmt.Sscanf(string(out), "%d", &height)
		}
	}
	return
}
