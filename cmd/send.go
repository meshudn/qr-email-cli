/*
Meshu Deb Nath
*/
package cmd

import (
	"fmt"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

var (
	to      string
	subject string
	body    string
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send emails",
	Long:  `For example: qremail send --to="bob@example.com" --subject="Hello" --body="Hello, Bob!"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Sending email to: %s\n", to)

		err := qrcode.WriteFile(body, qrcode.Medium, 256, "qr.png")
		if err != nil {
			SendEmail(to, subject)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Here you can define flags and configuration settings.
	// Cobra supports persistent flags and local flags.
	sendCmd.Flags().StringVarP(&to, "to", "t", "", "Recipient's email address (required)")
	sendCmd.Flags().StringVarP(&subject, "subject", "s", "", "Email subject (required)")
	sendCmd.Flags().StringVarP(&body, "body", "b", "", "Email body (required)")

	// Mark required flags
	sendCmd.MarkFlagRequired("to")
	sendCmd.MarkFlagRequired("subject")
	sendCmd.MarkFlagRequired("body")
}
