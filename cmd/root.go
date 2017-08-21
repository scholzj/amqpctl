// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/scholzj/amqpctl/mgmtlink"
)

//var cfgFile string
var amqpCfg mgmtlink.AmqpConfiguration

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "amqpctl",
	Short: "amqpctl is command line based client for AMQP Management protocol",
	Long: `amqpctl is a command line based client for the AMQP Management protocol. It
is written in Go language and tested against Apache Qpid Dispatch AMQP router.
It should be compatible with all implementations of AMQP Management Working
Draft 9 and 11. amqpctl is using Apache Qpid Proton as its underlying AMQP client.
Apache Qpid Proton has to be installed to use amqpctl. See amqpctl --help for
more details about the usage and supported operations.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() { 
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.amqpctl.yaml and ./.amqpctl.yaml)")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	RootCmd.PersistentFlags().StringVarP(&amqpCfg.AmqpHostname, "hostname", "b","localhost", "AMQP hostname (default localhost)")
	viper.BindPFlag("hostname", RootCmd.PersistentFlags().Lookup("hostname"))
	RootCmd.PersistentFlags().IntVarP(&amqpCfg.AmqpPort, "port","p", 5672, "AMQP port (default 5672)")
	viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))
	RootCmd.PersistentFlags().StringVar(&amqpCfg.AmqpUsername, "username","", "AMQP username")
	viper.BindPFlag("username", RootCmd.PersistentFlags().Lookup("username"))
	RootCmd.PersistentFlags().StringVar(&amqpCfg.AmqpPassword, "password","", "AMQP password")
	viper.BindPFlag("password", RootCmd.PersistentFlags().Lookup("password"))
	RootCmd.PersistentFlags().StringVar(&amqpCfg.SaslMechanism, "sasl-mechanism","", "AMQP SASL mechanism")
	viper.BindPFlag("sasl-mechanism", RootCmd.PersistentFlags().Lookup("sasl-mechanism"))
	RootCmd.PersistentFlags().StringVar(&amqpCfg.SslCaFile, "ssl-ca","", "SSL certification authority certificate(s)")
	viper.BindPFlag("ssl-ca", RootCmd.PersistentFlags().Lookup("ssl-ca"))
	RootCmd.PersistentFlags().StringVar(&amqpCfg.SslCertFile, "ssl-cert","", "SSL certificate for client authentication")
	viper.BindPFlag("ssl-cert", RootCmd.PersistentFlags().Lookup("ssl-cert"))
	RootCmd.PersistentFlags().StringVar(&amqpCfg.SslKeyFile, "ssl-key","", "SSL private key for client authentication")
	viper.BindPFlag("ssl-key", RootCmd.PersistentFlags().Lookup("ssl-key"))
	RootCmd.PersistentFlags().BoolVar(&amqpCfg.SslSkipHostnameVerification, "ssl-skip-verify",false, "Skip hostname verification")
	viper.BindPFlag("ssl-skip-verify", RootCmd.PersistentFlags().Lookup("ssl-skip-verify"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	if viper.GetString("config") != "" {
		// Use config file from the flag.
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".amqpctl" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(".amqpctl")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())

		amqpCfg.AmqpHostname = viper.GetString("hostname")
		amqpCfg.AmqpPort = viper.GetInt("port")
		amqpCfg.AmqpUsername = viper.GetString("username")
		amqpCfg.AmqpPassword = viper.GetString("password")
		amqpCfg.SaslMechanism = viper.GetString("sasl-mechanism")
		amqpCfg.SslCaFile = viper.GetString("ssl-ca")
		amqpCfg.SslCertFile = viper.GetString("ssl-cert")
		amqpCfg.SslKeyFile = viper.GetString("ssl-key")
		amqpCfg.SslSkipHostnameVerification = viper.GetBool("ssl-skip-verify")

	}
}
