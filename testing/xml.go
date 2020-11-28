package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

func isin(a string, bs []string) bool {
	for _, b := range bs {
		if a == b {
			return true
		}
	}

	return false
}

func main() {

	locDoc := etree.NewDocument()
	if err := locDoc.ReadFromFile("assets/xml/locs.xml"); err != nil {
		panic(err)
	}

	ae := locDoc.FindElement("locations/def/loc[@id='game_loc']").Text()

	re := locDoc.FindElement("locations/rel/loc[@id='spnpccharacters']").Text()

	fp := filepath.Join(ae, re)
	fp = strings.Replace(fp, "C:", "/mnt/c", -1)

	// spculture.xml

	spcultDoc := etree.NewDocument()

	if err := spcultDoc.ReadFromFile(fp); err != nil {
		panic(err)
	}

	es := spcultDoc.FindElements("//NPCCharacter[@default_group]")

	ss := make([]string, 0)

	for _, e := range es {
		a := e.SelectAttrValue("default_group", "")

		if a != "" && !isin(a, ss) {
			ss = append(ss, a)
		}
	}

	for _, s := range ss {
		fmt.Println(s)
	}

}
