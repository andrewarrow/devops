package main

import (
	"embed"
	_ "embed"
	"math/rand"
	"os"
	"time"
	"web/app"

	"github.com/andrewarrow/feedback/network"
	"github.com/andrewarrow/feedback/prefix"
	"github.com/andrewarrow/feedback/router"
)

//go:embed app/feedback.json
var embeddedFile []byte

//go:embed views/*.html
var embeddedTemplates embed.FS

//go:embed assets/**/*.*
var embeddedAssets embed.FS

var buildTag string

func main() {

	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 1 {
		return
	}
	arg := os.Args[1]

	prefix.FeedbackName = ""
	network.BaseUrl = os.Getenv("BACKEND")
	network.BaseHeaderKey = "AA-Service"
	network.BaseHeaderValue = "backend"
	router.BuildTag = buildTag
	router.EmbeddedTemplates = embeddedTemplates
	router.EmbeddedAssets = embeddedAssets
	r := router.NewRouter("DATABASE_URL", embeddedFile)

	if arg == "init" {
		//router.InitNewApp()
	} else if arg == "run" {
		r.Paths["/"] = app.HandleWelcome
		r.Paths["table"] = app.HandleTable
		r.ListenAndServe(":" + os.Args[2])
	} else if arg == "tr" {
	}
}
