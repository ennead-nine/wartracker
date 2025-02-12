package ui

import (
	"fmt"
	"net/http"
	"os"
	"wartracker/ui/wartracker-ui/handler"
	"wartracker/ui/wartracker-ui/site/rxk"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

var root = chi.NewRouter()

var (
	listenAddr string
	apiServer  string
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Root...\n"))
}

func Server() {
	fmt.Println(handler.SiteDataDir)
	http.ListenAndServe(listenAddr, root)
}

func init() {
	initConfig()
	initSiteDataDir()
	initApi()
	initServer()
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.SetConfigName(".wartracker-ui.yaml")
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

func initSiteDataDir() {
	handler.SiteDataDir = viper.GetString("siteDataDir")
	fmt.Println(handler.SiteDataDir)
}

func initApi() {
	apiServer = viper.GetString("apiURL")

	resp, err := http.Get(apiServer + "/health")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
}

func initServer() {
	listenAddr = viper.GetString("listenAddr")
	root.Use(middleware.Logger)

	h := handler.Handler{}
	root.Get("/", h.Default)
	root.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	root.Mount("/rxk", rxk.RxKRoutes())
}
