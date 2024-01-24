package cmd

import (
	"os"

	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	setting data.SettingType
)

var rootCmd = &cobra.Command{
	Use:   "datasync",
	Short: "",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./datasync.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		dir, err := file.FindCurrentDir()
		if err != nil {
			return
		}
		viper.AddConfigPath(dir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("datasync")
	}

	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	err = viper.Unmarshal(&setting)
	cobra.CheckErr(err)
}
