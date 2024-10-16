package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func LoginZammad(page playwright.Page) {

	// Go to Zammad login page
	if _, err := page.Goto("https://zammad.login.no/#login"); err != nil {
		log.Fatalf("could not go to the login page: %v", err)
	}
	button := page.Locator("button[type='submit']") // Adjust the selector if necessary
	err := button.WaitFor()
	if err != nil {
		log.Fatalf("could not find submit button: %v", err)
	}

	// Click the submit button
	err = button.Click()
	if err != nil {
		log.Fatalf("could not click the submit button: %v", err)
	}

	ClearScreen()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Login in to Zammad")
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	usernameField := page.Locator("input[name='uidField']")
	err = usernameField.WaitFor()
	if err != nil {
		log.Fatalf("could not find username field: %v", err)
	}

	passwordField := page.Locator("input[name='password']")
	err = passwordField.WaitFor()
	if err != nil {
		log.Fatalf("could not find password field: %v", err)
	}

	err = usernameField.Fill(username)
	if err != nil {
		log.Fatalf("could not fill username field: %v", err)
	}

	err = passwordField.Fill(password)
	if err != nil {
		log.Fatalf("could not fill password field: %v", err)
	}

	button = page.Locator("button[type='submit']")
	err = button.Click()
	if err != nil {
		log.Fatalf("could not click the submit button: %v", err)
	}

	ClearScreen()

	codeField := page.Locator("input[name='code']")
	successSelector := "div.navigation"

	// Wait for either the 2FA code field or the navigation div to appear
	for {
		twoFactorVisible, err := codeField.IsVisible()
		if err != nil {
			log.Fatalf("could not check 2FA field visibility: %v", err)
		}

		navigationVisible, err := page.Locator(successSelector).IsVisible()
		if err != nil {
			log.Fatalf("could not check navigation div visibility: %v", err)
		}

		if twoFactorVisible {
			// If 2FA field is visible, prompt for the code
			fmt.Print("2FA code: ")
			twoFactorCode, _ := reader.ReadString('\n')
			twoFactorCode = strings.TrimSpace(twoFactorCode)

			// Fill in the two-factor authentication code
			err = codeField.Fill(twoFactorCode)
			if err != nil {
				log.Fatalf("could not fill 2FA field: %v", err)
			}

			// Click the submit button for 2FA
			button = page.Locator("button[type='submit']")
			err = button.Click()
			if err != nil {
				log.Fatalf("could not click the submit button for 2FA: %v", err)
			}
			break
		} else if navigationVisible {
			fmt.Println("Logged in successfully without 2FA.")
			break
		}

		// Sleep for a short time to avoid busy-waiting
		time.Sleep(500 * time.Millisecond) // Adjust the interval as necessary
	}

	ClearScreen()

	fmt.Print("This can take up to 10 seconds...")
	ticketButton := page.Locator("a.list-button.fit.horizontal.centered[href='#ticket/create']")
	err = ticketButton.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(10000), // 10 seconds timeout
	})
	if err != nil {
		fmt.Print("\nUsername, password or 2fa code was wrong... try agien")
		time.Sleep(2 * time.Second)
		LoginZammad(page)
	}
}
