package main

import (
	"net/http"

	"github.com/beevik/etree"
)

const (
	// network
	port = "3000"

	// local
	assetLocsFile  = "assets/xml/locs.xml"
	assetTemplates = "assets/templates/"
	assetStatic    = "assets/static"
)

var (
	game   *Game
	locDoc *etree.Document
)

func homePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/location", http.StatusFound)
}

func main() {

	// locations
	locDoc = etree.NewDocument()
	if err := locDoc.ReadFromFile(assetLocsFile); err != nil {
		panic(err)
	}

	// -- web stuff --

	mux := http.NewServeMux()

	// Load static
	fs := http.FileServer(http.Dir(assetStatic))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// All available pages
	mux.HandleFunc("/location", locationPage)
	mux.HandleFunc("/unitselect", unitselectPage)
	mux.HandleFunc("/unit", unitPage)

	mux.HandleFunc("/", homePage)

	http.ListenAndServe(":"+port, mux)
}
