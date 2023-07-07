/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/lab42/gha-keda-webhook/counter"
	"github.com/lab42/gha-keda-webhook/handler"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gha-webhook-server",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		counter, err := counter.NewRedisCounter(viper.GetString("REDIS_ADDRESS"), viper.GetString("REDIS_PASSWORD"), viper.GetInt("REDIS_DATABASE"))
		cobra.CheckErr(err)

		handler := handler.Handler{Counter: counter}

		// Create an Echo instance
		e := echo.New()
		e.HideBanner = true

		// GZIP compression
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}))

		// Add middleware to gather metrics
		e.Use(echoprometheus.NewMiddleware("keda_actions_webhook"))
		// Add route to serve gathered metrics
		e.GET("/metrics", echoprometheus.NewHandler())
		log.Info("Registered '/metrics' endpoint")

		// Request ID middleware
		e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() string {
				return uuid.Must(uuid.NewRandom()).String()
			},
		}))

		// Kubernetes probe endpoint. Can be used for all probes
		e.GET("/healthz", handler.Probes)
		log.Info("Registered '/healtz' endpoint")

		// Define the webhook endpoint handler
		e.POST("/webhook", handler.Webhook)
		log.Info("Registered '/webhook' endpoint")

		// Start the web server
		log.Infof("Starting webhook server on %s", viper.GetString("SERVER_ADDRESS"))

		// Start server
		go func() {
			if err := e.Start(viper.GetString("SERVER_ADDRESS")); err != nil && err != http.ErrServerClosed {
				e.Logger.Fatal("shutting down the server")
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
		// Use a buffered channel to avoid missing signals as recommended for signal.Notify
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.env")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigType("env")
		viper.SetConfigName(".env")
	}

	viper.SetDefault("SERVER_ADDRESS", ":1234")
	viper.SetDefault("SECRET_TOKEN", "CHANGE_ME!")
	viper.SetDefault("REDIS_ADDRESS", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file:", viper.ConfigFileUsed())
	}
}
