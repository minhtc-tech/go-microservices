package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// * run locally
const webPort = "80"

var brokerURL = "http://localhost:8080"

// * deploy docker swarm - with caddy
// const webPort = "8081"

// var brokerURL = os.Getenv("BROKER_URL")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Printf("Starting Front-end service on port %s\n", webPort)
	webAddr := fmt.Sprintf(":%s", webPort)
	err := http.ListenAndServe(webAddr, nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, t string) {
	partials := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data struct {
		BrokerURL string
	}
	data.BrokerURL = brokerURL

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
