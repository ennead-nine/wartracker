package api

import (
	"fmt"
	"net/http"
	"os"
	"wartracker/pkg/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

var root = chi.NewRouter()

func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Root...\n"))
}

func Server() {
	http.ListenAndServe(":3000", root)
}

func init() {
	initConfig()
	initDb()
	initServer()
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.SetConfigName(".wartracker-api.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic(configError(err))
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		panic(configError(err))
	}
}

func initDb() {
	var err error
	db.Connection, err = db.Connect(viper.GetString("dbfile"))
	if err != nil {
		panic(err)
	}
}

func initServer() {
	root.Use(middleware.Logger)
	root.Get("/", GetRoot)
	root.Mount("/health", HealthRoutes())
	root.Mount("/alliance", AllianceRoutes())
}