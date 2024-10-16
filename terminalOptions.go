package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type Company struct {
	Emails string
	CC     []string
}

func TerminalOptions(page playwright.Page) {
	mailGroup := ""
	var mailOwner []string
	mailTitle := ""
	var mailText []string
	mailCustomerPath := ""
	reader := bufio.NewReader(os.Stdin)
	var companies []Company

	for {
		emailOptions(mailGroup, mailTitle, mailCustomerPath, mailOwner, (len(mailText) > 0))
		answer, _ := reader.ReadString('\n')
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
			fmt.Print(companies)
			fmt.Print(mailOwner)
			time.Sleep(10 * time.Second)
		case "q":
			os.Exit(0)

		default:
			continue
		}
	}
}

func emailOptions(mailGroup, mailTitle, mailCustomerPath string, mailOwner []string, hasText bool) {
	ClearScreen()
	fmt.Printf("\n1) Set email owner(s). Current:                  %s", mailOwner)
	fmt.Printf("\n2) Set email group. Current:                     %s", mailGroup)
	fmt.Printf("\n3) Set email title. Current:                     %s", mailTitle)
	fmt.Printf("\n4) Set email content or path. Have (bool):       %v", hasText)
	fmt.Printf("\n5) Set email customer and CC file path. Current: %s", mailCustomerPath)
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
	print("1) See current email content\n2) Set email content with file\n3) Type new email content\nQ) exit")
	var emailText []string
	for {
		fmt.Print("\n\nChoise: ")
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
			if len(inputLines) == 1 {
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
