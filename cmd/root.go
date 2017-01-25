package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var logFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "status",
	Short: "Checks your site's status",
	Run: func(cmd *cobra.Command, args []string) {
		initializeFlags(cmd)
		configureLogging()
		execWithArgs(args)
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

func execWithArgs(args []string) {
	u, err := url.Parse(args[0])
	if err != nil {
		log.Fatal(err)
	}
	site := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	t0 := time.Now()
	res, err := http.Get(site)
	t1 := time.Now()

	if err != nil {
		log.Fatal(err)
	}
	log.Println(u.Host + " " + res.Status + " (" + t1.Sub(t0).String() + ")")
}

func init() {
	RootCmd.PersistentFlags().StringVar(&logFile, "logFile", "", "Log File path")
}

func initializeFlags(cmd *cobra.Command) {
	persFlagKeys := []string{"logFile"}
	for _, key := range persFlagKeys {
		setValueFromFlag(cmd.PersistentFlags(), key)
	}
}

func setValueFromFlag(flags *pflag.FlagSet, key string) {
	if flagChanged(flags, key) {
		f := flags.Lookup(key)
		viper.Set(key, f.Value.String())
	}
}

func flagChanged(flags *pflag.FlagSet, key string) bool {
	flag := flags.Lookup(key)
	if flag == nil {
		return false
	}
	return flag.Changed
}

func configureLogging() {
	if viper.IsSet("logFile") && viper.GetString("logFile") != "" {
		path := viper.GetString("logFile")
		logFile, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
		}
		log.SetOutput(logFile)
	}
}
