# ZammadV2

ZammadV2 is a command-line application for sending emails through the Zammad ticketing system. This program automates the process of logging into Zammad, filling out ticket details, and sending emails, making it easier for users to bulk send emails.

## Features

- **Automated Login**: Automatically logs into the Zammad system using your credentials.
- **Two-Factor Authentication**: Supports Zammad's 2FA.
- **Send Emails**: Interface to bulk send emails.
- **User-Friendly Interface**: Terminal prompts guide users through the process.

## Requirements

- Go (1.23 or later)
- Playwright Go library

## Installation

1. Clone the repository:
   ```bash
   git clone https://gitlab.com/bot2310121/zammadv2.git
   cd zammadv2
   ```
2. Install dependencies
   ```bash
   go get github.com/playwright-community/playwright-go
   ```
3. Build or run zammadv2
   ```bash
   go build .
   or
   go run .
   ```


### Author
Alexander Engebrigtsen Heier
