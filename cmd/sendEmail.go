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
	from := "meshu.uiu@gmail.com"
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
	headers.Set("Subject", "Email with PNG attachment in message body")

	// Create a multipart/related message
	multipartWriter := multipart.NewWriter(msg)
	defer multipartWriter.Close()

	// Set the Content-Type header for the multipart message
	headers.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=%s", multipartWriter.Boundary()))

	// Write the headers to the message
	for key, values := range headers {
		for _, value := range values {
			fmt.Fprintf(msg, "%s: %s\r\n", key, value)
		}
	}
	fmt.Fprint(msg, "\r\n")

	// Add the HTML body with the inline image
	htmlBody := `<html><body><p>This is the body of the email with an inline PNG image:<br><img src="qr.png"></p></body></html>`
	fmt.Fprintf(msg, "--%s\r\n", multipartWriter.Boundary())
	fmt.Fprintf(msg, "Content-Type: text/html; charset=utf-8\r\n")
	fmt.Fprintf(msg, "\r\n%s\r\n", htmlBody)

	// Add the PNG attachment
	encodedAttachment := base64.StdEncoding.EncodeToString(pngContent)
	fmt.Fprintf(msg, "\r\n--%s\r\n", multipartWriter.Boundary())
	fmt.Fprintf(msg, "Content-Type: image/png\r\n")
	fmt.Fprintf(msg, "Content-Transfer-Encoding: base64\r\n")
	fmt.Fprintf(msg, "Content-Disposition: inline; filename=\"logo.png\"\r\n")
	fmt.Fprintf(msg, "Content-ID: <logo.png>\r\n")
	fmt.Fprintf(msg, "\r\n%s\r\n", encodedAttachment)

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", from, password, smtpServer)
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{to}, msg.Bytes())
	if err != nil {
		log.Fatal("Error sending email:", err)
	} else {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")
}
