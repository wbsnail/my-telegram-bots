package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/magicconch"
	"github.com/wbsnail/my-telegram-bots/pkg/log"
	"io"
	"math/rand"
	"os"
	"time"
)

// Version will be replaced when building via "-ldflags"
var Version = "unknown"

const (
	EnvDebug = "debug"

	DebugConfigFile    = "config/config.local.yaml"
	ConfigFileTemplate = "config/config.tmpl.yaml"
)

var cfgFile string
var env string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bots",
	Short: "This is my-telegram-bots command line tool",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	cobra.OnInitialize(initConfig)

	f := rootCmd.PersistentFlags()
	f.StringVarP(&cfgFile, "config", "c", "", "config file")
	f.StringVarP(&env, "env", "e", EnvDebug, "environment")
	magicconch.Must(viper.BindPFlags(f))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.Infof("environment: %s", env)

	if env == EnvDebug {
		if _, exist := os.LookupEnv("LOG_LEVEL"); !exist {
			log.Logger.SetLevel(logrus.DebugLevel)
		}

		if cfgFile == "" {
			// get debug config file in place automatically
			fs := afero.NewOsFs()
			exist, err := afero.Exists(fs, DebugConfigFile)
			magicconch.Must(err)
			if !exist {
				tmpl, err := fs.Open(ConfigFileTemplate)
				magicconch.Must(err)
				defer func() {
					magicconch.Must(tmpl.Close())
				}()
				debugConfig, err := fs.Create(DebugConfigFile)
				magicconch.Must(err)
				defer func() {
					magicconch.Must(debugConfig.Close())
				}()
				_, err = io.Copy(debugConfig, tmpl)
				magicconch.Must(err)
			}
			log.Infof("config file not specified, using default for debugging")
			cfgFile = DebugConfigFile
		}
	}

	if cfgFile != "" {
		viper.AutomaticEnv()
		viper.SetConfigFile(cfgFile)

		err := viper.ReadInConfig()
		if err != nil {
			log.Warn(errors.Wrapf(err, "read in config error, file: %s", viper.ConfigFileUsed()))
		} else {
			log.Infof("using config file: %s", viper.ConfigFileUsed())
		}
	}

	// default GIN_MODE is release
	if _, exist := os.LookupEnv("GIN_MODE"); !exist {
		gin.SetMode(gin.ReleaseMode)
	}
}
