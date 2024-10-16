package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
)

func main() {
	ClearScreen()
	hidden := true
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to see the program running or not? y or n: ")
	seeProgram, _ := reader.ReadString('\n')
	seeProgram = strings.ToLower(strings.TrimSpace(seeProgram))
	if seeProgram == "y" {
		hidden = false
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	// Launch a headless browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(hidden),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	// Create a new page
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	LoginZammad(page)

	TerminalOptions(page)
}
