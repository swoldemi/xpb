package main

import (
	// "fmt"
	// "strings"

	// "github.com/spf13/cobra"
	// "github.com/spf13/viper"
	"os"

	"github.com/sirupsen/logrus"
	xpb "github.com/swoldemi/xpb/pkg"
)

const (
	// Version is redefined a build time
	Version = "1.0.0-rc.1"
)

var (
	l *logrus.Logger
	e *logrus.Entry
)

func init() {
	l = logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	l.SetOutput(os.Stdout)

	// log the trace severity
	// TODO: Make this a config flag
	l.SetLevel(logrus.TraceLevel)

	// Display caller in log trace
	l.SetReportCaller(false)

	// Set shared fields
	e = l.WithField("v", Version)
}

func main() {
	config := &xpb.Config{
		AddressedProjectID: "nickel-api",
		HostKeyFilePath:    "./xpb-host.json",
		GuestKeyFilePath:   "./xpb-guest.json",
	}
	xpb.MustExecute(e, config)

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
