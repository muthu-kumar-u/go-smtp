# go-smtp-mailer

### Go SMTP Mailer is a standalone Go application that sends emails using SMTP. This application is designed to be simple, lightweight, and efficient for sending emails programmatically.

## Features
  * Send emails using any SMTP server.
  * Configurable SMTP server, port, authentication, and sender/receiver details.
  * Easily integrable into any Go project or as an AWS Lambda function.
  * Lightweight and fast.

## Prerequisites
#### To use this Go SMTP Mailer, ensure you have the following installed:
  * Go 1.19 or higher
  * SMTP server credentials (e.g., Gmail, AWS SES, Mailgun, etc.)

## Installation
#### Clone the repository to your local machine:

    git clone <repository_url>

## Configuration
The mailer uses environment variables to configure the SMTP server and email details. You can set them in your environment or within a `.env` file in the root of your project as per the `.env.sample` file.

## Usage
### Running as a Standalone Program

* Update your .env file or environment variables with your SMTP credentials and desired sender information.
* Build and run the application:

      go run main.go

The application will start and wait for email sending events to be triggered via the provided input (such as a JSON object passed to the function).
