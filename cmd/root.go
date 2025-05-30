package cmd

import (
	"os"
	"strings"

	"github.com/dragosboca/haws/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	prefix string
	region string
	dryRun bool
	logLevel string

	record string
	zoneId string
	path   string

	rootCmd = &cobra.Command{
		Use:   "haws",
		Short: "Hugo on AWS",
		Long:  "A cloudformation and template generator for running Hugo on AWS",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set log level based on flag
			switch strings.ToLower(logLevel) {
			case "debug":
				logger.SetLevel(logger.LevelDebug)
			case "info":
				logger.SetLevel(logger.LevelInfo)
			case "warn", "warning":
				logger.SetLevel(logger.LevelWarn)
			case "error":
				logger.SetLevel(logger.LevelError)
			default:
				// Keep default (info) level if invalid value provided
				logger.SetLevel(logger.LevelInfo)
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Command execution failed: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "Prefix for resources created. Can not be empty")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .haws.toml in current directory)")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "AWS region for the bucket and cloudfront distribution")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")

	rootCmd.PersistentFlags().StringVar(&record, "record", "", "Record name to be added to R53 zone")
	rootCmd.PersistentFlags().StringVar(&zoneId, "zone-id", "", "AWS Id of the zone used for SSL certificate validation and where the record should be added")
	rootCmd.PersistentFlags().StringVar(&path, "bucket-path", "", "Path prefix that will be appended by cloudfront to all requests (it should correspond to a sub-folder in the bucket)")

	if err := viper.BindPFlag("prefix", rootCmd.PersistentFlags().Lookup("prefix")); err != nil {
		logger.Fatal("Failed to bind prefix flag: %v", err)
	}

	if err := viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region")); err != nil {
		logger.Fatal("Failed to bind region flag: %v", err)
	}

	if err := viper.BindPFlag("record", rootCmd.PersistentFlags().Lookup("record")); err != nil {
		logger.Fatal("Failed to bind record flag: %v", err)
	}

	if err := viper.BindPFlag("zone_id", rootCmd.PersistentFlags().Lookup("zone-id")); err != nil {
		logger.Fatal("Failed to bind zone_id flag: %v", err)
	}

	if err := viper.BindPFlag("bucket_path", rootCmd.PersistentFlags().Lookup("bucket-path")); err != nil {
		logger.Fatal("Failed to bind bucket_path flag: %v", err)
	}

	if err := viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		logger.Fatal("Failed to bind log_level flag: %v", err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(cwd)
		viper.SetConfigType("toml")
		viper.SetConfigName(".haws")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Info("Using config file: %s", viper.ConfigFileUsed())
		} else {
			logger.Fatal("Fatal error config file: %v", err)
		}
	}
}
