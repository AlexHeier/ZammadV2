package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

func terminalHeader() {
	zammmadHeader := `
	__________                                  .___       ________  
	\____    /____    _____   _____ _____     __| _/ ___  _\_____  \ 
	  /     /\__  \  /     \ /     \\__  \   / __ |  \  \/ //  ____/ 
	 /     /_ / __ \|  Y Y  \  Y Y  \/ __ \_/ /_/ |   \   //       \ 
	/_______ (____  /__|_|  /__|_|  (____  /\____ |    \_/ \_______ \
    		\/    \/      \/      \/     \/      \/                \/
+-------------------------------------------------------------------------------+
	`
	fmt.Println(zammmadHeader)
}
