package main

import (
	"fmt"
	"net/http"
	"text/template"
)

var (
	unitselectTmpl = template.Must(template.ParseFiles(assetTemplates + "unitselect.html"))
	cultures       = []string{
		"empire",
		"aserai",
		"sturgia",
		"vlandia",
		"battania",
		"khuzait",
		"nord",
		"vakken",
		"darshi",
	}
	groups = []string{
		"Infantry",
		"infantry",
		"General",
		"Ranged",
		"Cavalry",
		"HorseArcher",
	}
)

type groupStruct struct {
	Name  string
	Units []string
}

type cultureStruct struct {
	Name   string
	Groups []*groupStruct
}

type unitselectStruct struct {
	Cultures []*cultureStruct
}

func unitselectPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		saveGame()
	}

	// load game
	loadgame()

	strct := new(unitselectStruct)

	for _, c := range cultures {
		// add a culture c
		a := new(cultureStruct)
		a.Name = c

		for _, g := range groups {
			// add a group g
			b := new(groupStruct)
			b.Name = g

			xpath := "//NPCCharacter[@culture='%s'][@default_group='%s']"
			xpath = fmt.Sprintf(xpath, "Culture."+c, g)
			for _, e := range game.SpNpc.FindElements(xpath) {
				// add a unit e
				b.Units = append(b.Units, e.SelectAttrValue("id", "ERROR"))
			}

			if b.Units != nil {
				// add group if there is units in it
				a.Groups = append(a.Groups, b)
			}
		}
		strct.Cultures = append(strct.Cultures, a)
	}

	if err := unitselectTmpl.Execute(w, strct); err != nil {
		panic(err)
	}
}
