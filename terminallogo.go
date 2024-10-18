package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode/utf8"
)

func ClearScreen() {
	var clearCmd *exec.Cmd

	// Check OS type and use the appropriate clear command
	if runtime.GOOS == "windows" {
		clearCmd = exec.Command("cmd", "/c", "cls")
	} else {
		clearCmd = exec.Command("clear")
	}
	clearCmd.Stdout = os.Stdout
	clearCmd.Run()

	terminalHeader()
}

// GetTerminalWidth gets the current terminal width in characters
func GetTerminalWidth() (int, error) {
	var cmd *exec.Cmd

	// Use OS-specific commands to get terminal width
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", "($Host.UI.RawUI.WindowSize.Width)")
	} else {
		cmd = exec.Command("tput", "cols")
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var width int
	fmt.Sscanf(string(output), "%d", &width)
	return width, nil
}

// terminalHeader prints the header centered horizontally
func terminalHeader() {
	zammadHeader := `
__________                                  .___       ________  
\____    /____    _____   _____ _____     __| _/ ___  _\_____  \ 
  /     /\__  \  /     \ /     \\__  \   / __ |  \  \/ //  ____/ 
 /     /_ / __ \|  Y Y  \  Y Y  \/ __ \_/ /_/ |   \   //       \ 
/_______ (____  /__|_|  /__|_|  (____  /\____ |    \_/ \_______ \
        \/    \/      \/      \/     \/      \/                \/
+---------------------------------------------------------------------------------------+
	`
	// Split header into lines
	lines := strings.Split(zammadHeader, "\n")

	// Get terminal width
	width, err := GetTerminalWidth()
	if err != nil {
		fmt.Println("Could not get terminal width:", err)
		width = 80 // default to 80 if error occurs
	}

	for _, line := range lines {
		lineLength := utf8.RuneCountInString(line)
		if lineLength < width {
			padding := (width - lineLength) / 2
			fmt.Printf("%s%s\n", strings.Repeat(" ", padding), line)
		} else {
			fmt.Println(line) // Print as-is if it doesn't fit
		}
	}
}
