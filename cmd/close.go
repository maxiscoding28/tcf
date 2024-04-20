/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// closeCmd represents the close command
var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "Move ticket to the .closed/ directory",
	Run: func(cmd *cobra.Command, args []string) {
		envVars := getEnvVars()
		homeInfo, err := setHomeDirectory(envVars["TCK_HOME_DIR"], false)
		if err != nil {
			fatalError(err)
		}
		var ticket TicketStruct
		if err := ticket.setTicketId(args, envVars["TCK_ID"]); err != nil {
			fatalError(err)
		}

		ticketDirectoryPath := ticket.getPath(homeInfo.getTicketsPath())

		if err := fileOrDirectoryExists(ticketDirectoryPath); err != nil {
			fatalError(err)
		} else {
			err := os.Rename(ticketDirectoryPath, ticket.getPath(homeInfo.getClosedPath()))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			log(fmt.Sprintf("Ticket closed - %s", ticket.TicketId), "success")
		}
	},
}

func init() {
	rootCmd.AddCommand(closeCmd)
}
