package ui

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var root = chi.NewRouter()
var site = NewSite("static/data/site.json")

var (
	listenAddr string
	apiServer  string
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Root...\n"))
}

func Server() {
	http.ListenAndServe(listenAddr, root)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	_, tmpl := path.Split(r.URL.Path)
	RenderTemplate(w, tmpl, site)
}

func RenderTemplate(w http.ResponseWriter, tmpl string, s *Site) {
	if tmpl == "" {
		tmpl = "root"
	}
	t, err := template.ParseFiles(
		"templates/head.html",
		"templates/"+tmpl+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Msgf("renderTemplate: %s\n", err)
		return
	}
	err = t.ExecuteTemplate(w, tmpl+".html", s)
	if err != nil {
		log.Error().Msgf("renderTemplate: Execute: %s\n", err)
		return
	}
}

func init() {
	initConfig()
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
	root.Get("/", Handler)
	root.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}
