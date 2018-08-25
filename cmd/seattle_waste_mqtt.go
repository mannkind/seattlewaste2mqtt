package cmd

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/mannkind/seattle_waste_mqtt/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var reload = make(chan bool)

// SeattleWasteMQTTCmd - The root Mysb commands
var SeattleWasteMQTTCmd = &cobra.Command{
	Use:   "seattle_waste_mqtt",
	Short: "Publish Seattle Waste pickup via MQTT",
	Long:  "Publish Seattle Waste pickup via MQTT",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			log.Printf("Creating the MQTT transport handler")
			controller := handlers.SeattleWasteMQTT{}
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

// Execute - Adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := SeattleWasteMQTTCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(func() {
		viper.SetConfigFile(cfgFile)
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("Configuration Changed: %s", e.Name)
			reload <- true
		})

		log.Printf("Loading Configuration %s", cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error Loading Configuration: %s ", err)
		}
		log.Printf("Loaded Configuration %s", cfgFile)
	})

	SeattleWasteMQTTCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", ".mysb.yaml", "The path to the configuration file")
}
