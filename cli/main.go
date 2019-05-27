package main

import (
	// "fmt"
	// "strings"

	// "github.com/spf13/cobra"
	// "github.com/spf13/viper"

	xpb "github.com/swoldemi/xpb"
)

const (
	version = "1.0.0-rc.1"
)

func main() {
	config := &xpb.Config{
		HostKeyFilePath:  "./xpb-host.json",
		GuestKeyFilePath: "./xpb-guest.json",
	}
	xpb.MustExecute(config)

	// var xpb = &cobra.Command{
	// 	Use:   "xpb -f [config_yaml]",
	// 	Short: "",
	// 	Long: `echo things multiple times back to the user by providing
	// a count and a string.`,
	// 	Args: cobra.MinimumNArgs(1),
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		for i := 0; i < 3; i++ {
	// 			fmt.Println("Echo: " + strings.Join(args, " "))
	// 		}
	// 	},
	// }
}
