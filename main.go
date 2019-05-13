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

var l *logrus.Logger

func init() {
	l = logrus.New()

	// Log as JSON instead of the default ASCII formatter
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	l.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	l.SetLevel(logrus.TraceLevel)

	// Display caller in log trace
	l.SetReportCaller(true)

}

func main() {
	xpb.MustExecute(l)
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
