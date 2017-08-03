package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var filePushCmd = &cobra.Command{
	Use:   "push SOURCE DESTINATION [lane]",
	Short: "Pushes SOURCE to DESTINATION on all [lane] instances",

	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Pushing files...")
	},
}
