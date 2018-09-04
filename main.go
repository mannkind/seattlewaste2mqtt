package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var reload = make(chan bool)

// Version - Set during compilation when using included Makefile
var Version = "X.X.X"

// SeattleWasteMQTTCmd - The root Mysb commands
var SeattleWasteMQTTCmd = &cobra.Command{
	Use:   "seattlewaste2mqtt",
	Short: "Publish Seattle Waste pickup via MQTT",
	Long:  "Publish Seattle Waste pickup via MQTT",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			log.Printf("Creating the MQTT transport handler")
			controller := CollectionLookup{}
			if err := viper.Unmarshal(&controller); err != nil {
				log.Panicf("Error unmarshaling configuration: %s", err)
			}

			if err := controller.Start(); err != nil {
				log.Panicf("Error starting MQTT transport handler: %s", err)
			}

			<-reload
			log.Printf("Received Reload Signal")
			controller.Stop()
		}
	},
}

func init() {
	cobra.OnInitialize(func() {
		viper.SetConfigFile(cfgFile)
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("Configuration Changed: %s", e.Name)
			reload <- true
		})
		viper.SetDefault("control.alertwithin", "24h")
		viper.SetDefault("control.lookupinterval", "8h")

		log.Printf("Loading Configuration %s", cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error Loading Configuration: %s ", err)
		}
		log.Printf("Loaded Configuration %s", cfgFile)
	})

	SeattleWasteMQTTCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", ".seattlewaste2mqtt.yaml", "The path to the configuration file")
}

func main() {
	log.Printf("Seattle Waste Version: %s", Version)
	Execute()
}

// Execute - Adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := SeattleWasteMQTTCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
