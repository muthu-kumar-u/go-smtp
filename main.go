package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

// SendEmail function sends an email using the specified SMTP server.
func SendEmail(smtpServer string, smtpPort string, smtpUser string, smtpPass string, fromEmail string, toEmail string, subject string, body string) error {
	msg := []byte("To: " + toEmail + "\r\n" +
		"From: " + fromEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := smtpServer + ":" + smtpPort
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpServer)

	err := smtp.SendMail(addr, auth, fromEmail, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// LambdaHandler handles email sending in Lambda.
func LambdaHandler(ctx context.Context, event map[string]interface{}) (Response, error) {
	return triggerEmail(event)
}

// triggerEmail is a reusable function to send email, used in both HTTP and Lambda handler.
func triggerEmail(event map[string]interface{}) (Response, error) {
	fromEmail := os.Getenv("MAIL_FROM")
	toEmail, ok := event["to"].(string)
	if !ok {
		return Response{StatusCode: 400, Body: "Missing 'to' field"}, nil
	}

	subject, ok := event["subject"].(string)
	if !ok {
		return Response{StatusCode: 400, Body: "Missing 'subject' field"}, nil
	}

	body, ok := event["body"].(string)
	if !ok {
		return Response{StatusCode: 400, Body: "Missing 'body' field"}, nil
	}

	smtpServer := os.Getenv("MAIL_HOST")
	if smtpServer == "" {
		fmt.Println("SMTP server not configured")
		return Response{StatusCode: 500, Body: "Internal server error"}, nil
	}

	smtpPort := os.Getenv("MAIL_PORT")
	if smtpPort == "" {
		fmt.Println("SMTP port not configured")
		return Response{StatusCode: 500, Body: "Internal server error"}, nil
	}

	smtpUser := os.Getenv("MAIL_USERNAME")
	smtpPass := os.Getenv("MAIL_PASSWORD")

	if smtpUser == "" || smtpPass == "" || fromEmail == "" {
		fmt.Println("SMTP username and password not configured")
		return Response{StatusCode: 500, Body: "Internal server error"}, nil
	}

	err := SendEmail(smtpServer, smtpPort, smtpUser, smtpPass, fromEmail, toEmail, subject, body)
	if err != nil {
		fmt.Println(err.Error())
		return Response{StatusCode: 500, Body: "Error sending email"}, nil
	}

	return Response{
		StatusCode: 200,
		Body:       "Email sent successfully!",
	}, nil
}

// httpHandler handles the HTTP requests in non-production mode.
func httpHandler(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request or missing fields"})
		return
	}

	response, err := triggerEmail(request)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(response.StatusCode, gin.H{"message": response.Body})
		return
	}

	c.JSON(response.StatusCode, gin.H{"message": response.Body})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	// Check if the app is running in production
	if os.Getenv("PRODUCTION") == "true" {
		fmt.Println("Running in production mode with AWS Lambda handler")
		lambda.Start(LambdaHandler)
	} else {
		fmt.Println("Running in development mode with Gin HTTP server")

		router := gin.New()
		router.Use(gin.Logger())
		router.Use(gin.Recovery())
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

		version := os.Getenv("VERSION")
		if version == "" {
			version = "v1"
		}

		// Create a route group with the version as a prefix
		apiGroup := router.Group("api/" + version)
		apiGroup.POST("/send-email", httpHandler)

		router.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "path not found",
			})
		})

		router.Run(":8080")
	}
}
