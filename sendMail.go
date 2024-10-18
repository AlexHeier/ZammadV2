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

func SendMails(mailGroup, mailTitle string, mailOwner, mailText []string, companies []Company, page playwright.Page) {

	total := len(companies)
	page.SetDefaultTimeout(30000)

	for i, comp := range companies {
		if i == 0 {
			continue
		}
		owner := mailOwner[(i-1)%len(mailOwner)]
		ClearScreen()

		dots := (i % 3) + 1
		fmt.Print("\nSending " + strings.Repeat(".", dots) + " ")

		// Calculate progress
		progress := ((i + 1) * 71 / total)
		bar := fmt.Sprintf("|%s%s| (%d / %d)",
			strings.Repeat("#", progress),
			strings.Repeat("-", 71-progress),
			i+1,
			total)

		// Print the progress bar
		fmt.Print("\n" + bar + "\n\n")

		button := page.Locator("a.list-button.fit.horizontal.centered[href='#ticket/create']")
		err := button.WaitFor()
		if err != nil {
			log.Fatalf("could not find the button: %v", err)
		}

		err = button.Click()
		if err != nil {
			log.Fatalf("could not click the button: %v", err)
		}

		_, err = page.Reload(playwright.PageReloadOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle, // Wait until no network connections for at least 500 ms
		})

		if err != nil {
			log.Fatalf("could not reload page without cache: %v", err)
		}

		inputTitle := page.Locator(`.ticket-form-top .input.form-group.is-required input[name="title"]`).First()

		err = inputTitle.WaitFor()
		if err != nil {
			log.Fatalf("could not find the inputTitle field within timeout: %v", err)
		}

		err = inputTitle.Fill(mailTitle)
		if err != nil {
			log.Fatalf("could not fill the input field: %v", err)
		}

		inputGroup := page.Locator(`input[name="group_id"]`)

		// Wait for the input field
		err = inputGroup.WaitFor()
		if err != nil {
			return // Handle error appropriately if needed
		}

		// Fill the input with the mail group and press Enter
		if err = inputGroup.Fill(mailGroup); err == nil {
			err = inputGroup.Press("Enter")
		}

		// Wait indefinitely for the owner options to load
		_, err = page.WaitForFunction(`() => document.querySelector('select[name="owner_id"]').options.length > 1`, playwright.PageWaitForFunctionOptions{
			Timeout: playwright.Float(0), // 0 means wait indefinitely
		})
		if err != nil {
			fmt.Printf("Error finding owner '%s': %v\n", owner, err)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
			return
		}

		// Proceed with the rest of your code when successful

		// Define the owner select element
		ownerSelect := page.Locator("select[name='owner_id']")

		// Select the found mailOwner in the dropdown menu

		// Attempt to select the option by its label
		if _, err := ownerSelect.SelectOption(playwright.SelectOptionValues{
			Labels: &[]string{owner}, // Use the label field here
		}); err != nil {
			// Print error message and wait for user input before returning
			fmt.Printf("Error selecting owner '%s': %v\n", owner, err)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
			return                                     // Return to exit the function
		}

		// Enter the email into the input field
		emailInput := page.Locator(`*[name="customer_id_completion"]`)
		err = emailInput.WaitFor()
		if err != nil {
			fmt.Println("Error waiting for email input:", err)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
			return
		}
		err = emailInput.Fill(comp.Emails)
		if err != nil {
			fmt.Println("Error filling email input:", err)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
			return
		}

		// If there are CCs, concatenate them into a string separated by commas
		if len(comp.CC) > 0 {
			ccString := strings.Join(comp.CC, " ")

			// Enter CCs into the input field
			ccInput := page.Locator(`div[data-attribute-name="cc"] .token-input.ui-autocomplete-input`) // Target the correct input
			err = ccInput.WaitFor()
			if err != nil {
				fmt.Println("Error waiting for CC input:", err)
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
				return
			}
			err = ccInput.Fill(ccString)
			if err != nil {
				fmt.Println("Error filling CC input:", err)
				fmt.Println("Press Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
				return
			}
		}

		// Locate the content editable div by its data-name attribute
		textField := page.Locator(`div[data-name="body"]`)

		// Wait for the text field to be visible
		err = textField.WaitFor()
		if err != nil {
			fmt.Println("Error waiting for text field:", err)
			time.Sleep(5 * time.Second)
			return
		}

		// Clear the existing value in the text field
		err = textField.Fill("") // This removes all existing content
		if err != nil {
			fmt.Println("Error clearing text field:", err)
			time.Sleep(5 * time.Second)
			return
		}

		// Create the new content string with selectedOwner
		fromGreating := fmt.Sprintf(`

Med vennlig hilsen,
%s
%s

Login - Linjeforeningen for IT
Teknologivegen 22, 2815 Gjøvik
Email: kontakt@login.no
Web: https://login.no/`, owner, mailGroup)

		newContent := strings.Join(mailText, "\n")
		combinedContent := newContent + "\n" + fromGreating

		// Set the new content back to the text field
		err = textField.Fill(combinedContent)
		if err != nil {
			fmt.Println("Error setting new input value:", err)
			time.Sleep(5 * time.Second)
			return
		}

		time.Sleep(8 * time.Second)

		// Locate the button by its class or type and click it
		createButton := page.Locator(`button[type="submit"].btn.btn--success.js-submit`)
		err = createButton.WaitFor()
		if err != nil {
			log.Printf("could not find the Create button: %v", err)
			time.Sleep(5 * time.Second)
		}
		err = createButton.Click()
		if err != nil {
			log.Printf("could not click the Create button: %v", err)
			time.Sleep(5 * time.Second)
		}
		time.Sleep(1 * time.Second)
	}
}
