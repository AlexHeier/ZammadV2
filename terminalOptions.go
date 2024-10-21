package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type Company struct {
	Emails string
	CC     []string
}

func TerminalOptions(page playwright.Page) (string, string, []string, []string, []Company) {
	mailGroup := ""
	var mailOwner []string
	mailTitle := ""
	var mailText []string
	mailCustomerPath := ""
	reader := bufio.NewReader(os.Stdin)
	var companies []Company

	for {
		// Assuming emailOptions is defined elsewhere and works as expected
		emailOptions(mailGroup, mailTitle, mailCustomerPath, mailOwner, len(mailText) > 0)

		answer, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		answer = strings.ToLower(strings.TrimSpace(answer))
		switch answer {
		case "1":
			mailOwner = setMailOwner(mailOwner)
		case "2":
			mailGroup = setEmailGroup(mailGroup)
		case "3":
			mailTitle = setEmailTitle(mailTitle)
		case "4":
			mailText = setEmailContent(mailText)
		case "5":
			mailCustomerPath, companies = setCsvPath(mailCustomerPath, companies)
		case "s":
			confirm := sendMailConfirm(mailGroup, mailTitle, mailOwner, mailText, companies, page)
			if confirm {
				return mailGroup, mailTitle, mailOwner, mailText, companies
			}
		case "q":
			ClearScreen()
			fmt.Print("\nTakk for at du brukte Zammad V2 laget av Alexander Engebrigtsen Heier :)\n\n")
			os.Exit(0)
		default:

		}
	}
}

func emailOptions(mailGroup, mailTitle, mailCustomerPath string, mailOwner []string, hasText bool) {
	ClearScreen()
	fmt.Printf("\n1) Set email owner(s). Current:                      %s", mailOwner)
	fmt.Printf("\n2) Set email group. Current:                         %s", mailGroup)
	fmt.Printf("\n3) Set email title. Current:                         %s", mailTitle)
	fmt.Printf("\n4) Set email content or path. Have (bool):           %v", hasText)
	fmt.Printf("\n5) Set email customer and CC file path. Current:     %s", mailCustomerPath)
	fmt.Printf("\nS) to the send mail(s)")
	fmt.Printf("\nQ) to quit the program")
	fmt.Print("\n\nChoice: ")
}

func setMailOwner(currentOwners []string) []string {
	reader := bufio.NewReader(os.Stdin)
	for {
		ClearScreen() // Assuming ClearScreen is defined elsewhere
		fmt.Printf("\nCurrent: %v", currentOwners)
		fmt.Print("\nYou can have more than one owner; this will assign owners using Round Robin.")
		fmt.Print("\n\n1) Add a name (has to be full name correctly written as in Zammad\n2) Remove a name\nQ) exit\n\nChoice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			for {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("New owner (or Q to quit): ")
				newOwner, _ := reader.ReadString('\n')
				newOwner = strings.TrimSpace(newOwner)

				// Check for exit condition
				if strings.ToUpper(newOwner) == "Q" {
					break
				}

				// Capitalize the first letter of each word
				newOwner = capitalize(newOwner)

				// Append the new owner to the slice
				currentOwners = append(currentOwners, newOwner)
			}
		case "2":
			for {
				reader := bufio.NewReader(os.Stdin)
				ClearScreen() // Assuming ClearScreen is defined elsewhere
				fmt.Print("Which do you want to remove: ")

				// Display the current items with their indices
				for i, name := range currentOwners {
					fmt.Printf("\n%d) %s", i+1, name) // +1 for user-friendly indexing
				}
				fmt.Print("\n0) exit\n\nChoice: ")

				selection, _ := reader.ReadString('\n')
				selection = strings.TrimSpace(selection)
				selectionInt, err := strconv.Atoi(selection)

				// Validate input
				if err != nil || selectionInt < 0 || selectionInt > len(currentOwners) {
					continue
				}

				// Exit option
				if selectionInt == 0 {
					break
				}

				// Remove the selected item
				indexToRemove := selectionInt - 1                                                         // Convert to 0-based index
				currentOwners = append(currentOwners[:indexToRemove], currentOwners[indexToRemove+1:]...) // Remove the item
			}
		case "q":
			return currentOwners
		default:
			continue
		}
	}
}

func capitalize(s string) string {
	words := strings.Fields(s) // Split into words
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ") // Join words back into a single string
}

func setEmailGroup(oldMailGroup string) (newMailGroup string) {

	reader := bufio.NewReader(os.Stdin)

	ClearScreen()
	fmt.Print("Set email group")
	fmt.Printf("\nCurrent: %v", oldMailGroup)
	fmt.Println("\n1) Arrangementer")
	fmt.Println("2) BedKom")
	fmt.Println("3) Bedriftpresentasjon")
	fmt.Println("4) Cyberdagene")
	fmt.Println("5) Felles")
	fmt.Println("6) Karrieredag")
	fmt.Println("7) Stillingsutlysing")

	for {

		fmt.Print("\n\nChoise: ")
		choise, _ := reader.ReadString('\n')
		choise = strings.TrimSpace(choise)

		switch choise {
		case "1":
			newMailGroup = "Arrangementer"
			return
		case "2":
			newMailGroup = "BedKom"
			return
		case "3":
			newMailGroup = "Bedriftpresentasjon"
			return
		case "4":
			newMailGroup = "Cyberdagene"
			return
		case "5":
			newMailGroup = "Felles"
			return
		case "6":
			newMailGroup = "Karrieredag"
			return
		case "7":
			newMailGroup = "Stillingsutlysing"
			return
		default:
			continue
		}
	}
}

func setEmailTitle(oldTitle string) (newTitle string) {
	ClearScreen()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("The title of the email\n")
	fmt.Printf("\nCurrent: %v", oldTitle)
	fmt.Print("\nNew (Q to exit): ")
	newTitle, _ = reader.ReadString('\n')
	newTitle = strings.TrimSpace(newTitle)
	if strings.ToLower(newTitle) == "q" {
		return oldTitle
	}
	return newTitle
}

func setCsvPath(oldCsvPath string, oldCompanies []Company) (csvPath string, companiesObject []Company) {
	ClearScreen()

	fmt.Print("\nCSV format: Emails, CC")
	fmt.Print("\nExample line: x@stud.ntnu.no, y@stud.ntnu.no z@stud.ntnu.no ")
	fmt.Print("\nThere has be 1 customer email and 0...n CC\n")
	var companies []Company

	read := bufio.NewReader(os.Stdin)
	fmt.Printf("\nCurrent absolute path: %v", oldCsvPath)
	fmt.Print("\nNew absolute path (Q to exit): ")
	newCsvPath, _ := read.ReadString('\n')
	newCsvPath = strings.TrimSpace(newCsvPath)

	if strings.ToLower(newCsvPath) == "q" {
		return oldCsvPath, oldCompanies
	}

	file, err := os.Open(newCsvPath)
	if err != nil {
		return "Didn't find CSV path", nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return "Error reading CSV", nil
	}

	for i, record := range records {
		// Skip the header row
		if i == 0 {
			continue
		}

		CCList := strings.Split(record[1], " ")

		// Create a Company struct and populate it
		company := Company{
			Emails: record[0],
			CC:     CCList,
		}

		// Add the company to the list
		companies = append(companies, company)
	}

	return newCsvPath, companies
}

func setEmailContent(oldEmailText []string) []string {

	ClearScreen()
	print("1) See current email content\n2) Set email content with file\n3) Type new email content\nQ) exit\n")
	var emailText []string
	for {
		fmt.Printf("\nChoise: ")
		reader := bufio.NewReader(os.Stdin)
		choise, _ := reader.ReadString('\n')
		choise = strings.ToLower(strings.TrimSpace(choise))
		switch choise {
		case "1":
			if len(oldEmailText) <= 0 && len(emailText) <= 0 {
				fmt.Print("There is no email content yet")
				continue
			}

			if len(emailText) > 0 {
				for _, lines := range emailText {
					fmt.Println(lines)
				}
			} else {
				for _, lines := range oldEmailText {
					fmt.Println(lines)
				}
			}
			continue

		case "2":
			fmt.Print("\nNew absolute path: ")
			csvPath, _ := reader.ReadString('\n')
			csvPath = strings.TrimSpace(csvPath)

			file, err := os.Open(csvPath)
			if err != nil {
				fmt.Printf("Error opening file: %v\n", err)
				return emailText
			}
			defer file.Close()

			// Create a scanner to read each line
			scanner := bufio.NewScanner(file)

			// Read each line and append it to the slice
			for scanner.Scan() {
				emailText = append(emailText, scanner.Text())
			}

			// Check for errors during scanning
			if err := scanner.Err(); err != nil {
				fmt.Printf("Error reading file: %v\n", err)
			}

			fmt.Println("\nHere is the new email content:\n\n")
			for _, line := range emailText {
				fmt.Println(line)
			}
			continue

		case "3":

			reader := bufio.NewReader(os.Stdin)
			var inputLines []string

			fmt.Println("Enter your email content text here. Type 'DONE' on a new line when finished. \nType 'DONE' with no new lines to exit without saving:\n")

			for {
				fmt.Print("> ")
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(line)

				if strings.ToUpper(line) == "DONE" {
					break
				}

				inputLines = append(inputLines, line)
			}
			if len(inputLines) == 0 {
				continue
			} else {
				emailText = inputLines
			}

		case "q":
			if len(emailText) <= 0 {
				return oldEmailText
			}
			return emailText

		default:
			continue
		}

	}
}

func sendMailConfirm(mailGroup, mailTitle string, mailOwner, mailText []string, companies []Company, page playwright.Page) (accept bool) {

	page.SetDefaultTimeout(30000)
	ClearScreen()

	fmt.Print("\nTesting data...")

	accept = false
	reader := bufio.NewReader(os.Stdin)

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
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})

	fmt.Print("\nTitle")

	inputTitle := page.Locator(`.ticket-form-top .input.form-group.is-required input[name="title"]`).First()

	err = inputTitle.WaitFor()
	if err != nil {
		log.Fatalf("could not find the input field within timeout: %v", err)
	}

	err = inputTitle.Fill(mailTitle)
	if err != nil {
		log.Fatalf("could not fill the input field: %v", err)
	}

	inputGroup := page.Locator(`input[name="group_id"]`)

	fmt.Print("\nGroup")

	err = inputGroup.WaitFor()
	if err != nil {
		fmt.Printf("Error waiting for input field: %v\n", err)
		return
	}

	err = inputGroup.Fill(mailGroup)
	if err != nil {
		fmt.Printf("Error filling input field: %v\n", err)
		return
	}

	err = inputGroup.Press("Enter")
	if err != nil {
		fmt.Printf("Error pressing Enter: %v\n", err)
		return
	}

	fmt.Print("\nOwner")

	// Wait for the options to contain more than one entry
	if _, err := page.WaitForFunction(`() => document.querySelector('select[name="owner_id"]').options.length > 1`, nil); err != nil {
		fmt.Printf("Error waiting for options to have more than one entry: %v\n", err)
		fmt.Println("Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
		return
	}

	// Define the owner select element
	ownerSelect := page.Locator("select[name='owner_id']")

	// Get the options from the select element
	options := ownerSelect.Locator("option")
	optionCount, err := options.Count()
	if err != nil {
		fmt.Printf("Error retrieving options: %v\n", err)
		return
	}

	// Initialize an array to store the option values
	var optionValues []string

	// Store the existing options in the array
	for i := 0; i < optionCount; i++ {
		optionText, err := options.Nth(i).TextContent()
		if err != nil {
			fmt.Printf("Error getting option text: %v\n", err)
			continue
		}
		optionValues = append(optionValues, strings.TrimSpace(optionText))
	}

	var notFound []string
	var selectedOwner string // Variable to store the selected owner

	// Check if all mailOwners exist in existingOwners
	for _, owner := range mailOwner {
		found := false
		for _, existing := range optionValues {
			if owner == existing {
				found = true
				selectedOwner = strings.TrimSpace(owner) // Store the found owner
				break
			}
		}
		if !found {
			notFound = append(notFound, owner)
		}
	}

	// Print the not found owners
	if len(notFound) > 0 {
		fmt.Println("The following owners were not found:")
		for _, name := range notFound {
			fmt.Println(name)
		}
		fmt.Println("Delete them and try again. Going back to main menu")
		fmt.Println("Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
		return
	}

	// Select the found mailOwner in the dropdown menu

	if selectedOwner != "" {
		// Attempt to select the option by its label
		if _, err := ownerSelect.SelectOption(playwright.SelectOptionValues{
			Labels: &[]string{selectedOwner}, // Use the label field here
		}); err != nil {
			// Print error message and wait for user input before returning
			fmt.Printf("Error selecting owner '%s': %v\n", selectedOwner, err)
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
			return                                     // Return to exit the function
		}
	}

	fmt.Print("\nEmail content")

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
Web: https://login.no/`, selectedOwner, mailGroup)

	newContent := strings.Join(mailText, "\n")
	combinedContent := newContent + "\n" + fromGreating

	// Set the new content back to the text field
	err = textField.Fill(combinedContent)
	if err != nil {
		fmt.Println("Error setting new input value:", err)
		time.Sleep(5 * time.Second)
		return
	}

	fmt.Print("\nEmail and CC")
	// Get the first company
	firstCompany := companies[0]

	// Enter the email into the input field
	emailInput := page.Locator(`*[name="customer_id_completion"]`)
	err = emailInput.WaitFor()
	if err != nil {
		fmt.Println("Error waiting for email input:", err)
		fmt.Println("Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
		return
	}
	err = emailInput.Fill(firstCompany.Emails)
	if err != nil {
		fmt.Println("Error filling email input:", err)
		fmt.Println("Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n') // Wait for user input
		return
	}

	// If there are CCs, concatenate them into a string separated by commas
	if len(firstCompany.CC) > 0 {
		ccString := strings.Join(firstCompany.CC, " ")

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

	screenshotPath := filepath.Join(".", "new_ticket_screenshot.png")

	// Take a screenshot of the specific div
	if _, err = page.Locator(`div.newTicket`).Screenshot(playwright.LocatorScreenshotOptions{
		Path: playwright.String(screenshotPath),
	}); err != nil {
		fmt.Println("Error taking screenshot:", err)
	}

	var response string

	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", screenshotPath).Start()
	case "darwin": // macOS
		exec.Command("open", screenshotPath).Start()
	case "linux":
		exec.Command("xdg-open", screenshotPath).Start()
	default:
		// Unsupported OS
	}

	pwd, _ := os.Getwd()
	screenshotFullPath := filepath.Join(pwd, screenshotPath)

	for {
		fmt.Printf("\n\nScreenshot taken %s. Does everything look correct? (y/n): ", screenshotFullPath)
		response, _ = reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response)) // Convert to lowercase and trim whitespace

		// Check if the response is either 'y' or 'n'
		if response == "y" {
			break
		} else if response == "n" {
			break
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}

	// Delete the screenshot file
	_ = os.Remove(screenshotPath)

	if response != "y" {
		return false
	}

	createButton := page.Locator(`button[type="submit"].btn.btn--success.js-submit`)
	err = createButton.Click()
	if err != nil {
		log.Fatalf("could not click the Create button: %v", err)
		time.Sleep(5 * time.Second)
	}
	time.Sleep(1 * time.Second)
	return true
}
