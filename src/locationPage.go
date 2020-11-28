package main

import (
	"net/http"
	"os"
	"text/template"

	"github.com/beevik/etree"
)

var locationTmpl = template.Must(template.ParseFiles(assetTemplates + "location.html"))

type locationPageStruct struct {
	GameDir string
}

func locationPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		path := r.FormValue("game_dir")

		if e := locDoc.FindElement("locations/loc[@id='game_loc']"); e == nil {
			// if game_loc does not exist
			// parent
			p := locDoc.FindElement("locations")
			// create new loc[@id='game_loc']
			e = etree.NewElement("loc")
			e.CreateAttr("id", "game_loc")
			e.SetText(path)
			// add new element to parent
			p.AddChild(e)
		} else {
			// if game_loc exists
			// update content
			e.SetText(path)
		}

		// save locdoc
		locDoc.WriteToFile(assetLocsFile)

		// Check that given path exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic(err)
		}

		// redirect to unitselect
		http.Redirect(w, r, "/unitselect", http.StatusFound)

	} else {

		// default
		d := locDoc.FindElement("locations/def/loc[@id='game_loc']")
		// element
		e := locDoc.FindElement("locations/loc[@id='game_loc']")
		// struct
		strct := new(locationPageStruct)

		if e == nil {
			strct.GameDir = d.Text()
		} else {
			strct.GameDir = e.Text()
		}

		// Execute template
		if err := locationTmpl.Execute(w, strct); err != nil {
			panic(err)
		}

	}
}
