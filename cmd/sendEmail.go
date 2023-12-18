package cmd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
)

func SendEmail(to string, subject string) {
	from := "test@example.com"
	password := "your-email-password"

	// SMTP server address and port
	smtpServer := "localhost"
	smtpPort := "1025"

	// Path to the PNG file
	pngFilePath := "qr.png"

	// Open the PNG file
	pngFile, err := os.Open(pngFilePath)
	if err != nil {
		log.Fatalf("Error opening PNG file: %v", err)
	}
	defer pngFile.Close()

	// Read the contents of the PNG file
	pngContent := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := pngFile.Read(buffer)
		if err != nil {
			break
		}
		pngContent = append(pngContent, buffer[:n]...)
	}

	msg := &bytes.Buffer{}

	// Write the headers
	headers := textproto.MIMEHeader{}
	headers.Set("From", from)
	headers.Set("To", to)
	headers.Set("Subject", subject)

	// Create a multipart/related message
	multipartWriter := multipart.NewWriter(msg)
	defer multipartWriter.Close()

	// Set the Content-Type header for the multipart message
	headers.Set("Content-Type", fmt.Sprintf("multipart/mixed; boundary=%s", multipartWriter.Boundary()))

	// Write the headers to the message
	for key, values := range headers {
		for _, value := range values {
			fmt.Fprintf(msg, "%s: %s\r\n", key, value)
		}
	}
	fmt.Fprint(msg, "\r\n")

	encodedImage := base64.StdEncoding.EncodeToString(pngContent)

	// Add the HTML body with the embedded image
	htmlBody := `<html><body><p>This is the body of the email with an embedded image:<br><img src="data:image/png;base64,` + encodedImage + `"></p></body></html>`
	fmt.Fprintf(msg, "--%s\r\n", multipartWriter.Boundary())
	fmt.Fprintf(msg, "Content-Type: text/html; charset=utf-8\r\n")
	fmt.Fprintf(msg, "\r\n%s\r\n", htmlBody)

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", from, password, smtpServer)
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{to}, msg.Bytes())
	if err != nil {
		log.Fatal("Error sending email:", err)
	}

	fmt.Println("Email sent successfully!")
}
