package main

import (
	"fmt"
	"path/filepath"

	"github.com/beevik/etree"
)

// Game contains etree docs for all relevant xml trees.
type Game struct {
	SpNpc  *etree.Document
	SpCult *etree.Document

	Cultures []string
}

// loadGame get all etree docs in Game struct from game directory path.
func loadgame() {

	if game != nil {
		return
	}

	// handle path

	gameAbsPath := locDoc.FindElement("locations/loc[@id='game_loc']").Text()
	gameAbsPath = filepath.FromSlash(gameAbsPath)

	spnpcRelPath := locDoc.FindElement("locations/rel/loc[@id='spnpccharacters']").Text()
	spnpcRelPath = filepath.FromSlash(spnpcRelPath)

	spcultRelPath := locDoc.FindElement("locations/rel/loc[@id='spcultures']").Text()
	spcultRelPath = filepath.FromSlash(spcultRelPath)

	// create Game

	game = new(Game)

	game.SpNpc = getDoc(gameAbsPath, spnpcRelPath)
	game.SpCult = getDoc(gameAbsPath, spcultRelPath)

}

// saveGame saves the content of game at appropriate locations
func saveGame() {

	if game == nil {
		panic("Cannot save unitialised game.")
	}

	// handle path

	gameAbsPath := locDoc.FindElement("locations/loc[@id='game_loc']").Text()
	gameAbsPath = filepath.FromSlash(gameAbsPath)

	spnpcRelPath := locDoc.FindElement("locations/rel/loc[@id='spnpccharacters']").Text()
	spnpcRelPath = filepath.FromSlash(spnpcRelPath)

	path := filepath.Join(gameAbsPath, spnpcRelPath)

	// write to file

	if err := game.SpNpc.WriteToFile(path); err != nil {
		panic(err)
	}

}

func getDoc(absPath string, relPath string) *etree.Document {
	doc := etree.NewDocument()

	path := filepath.Join(absPath, relPath)

	if err := doc.ReadFromFile(path); err != nil {
		panic(err)
	}

	return doc
}

func getCultures(game *Game) []string {

	es := game.SpCult.FindElements("//Culture[@is_main_culture]")

	cs := make([]string, len(es))

	for i, e := range es {
		cs[i] = e.SelectAttr("id").Value
		fmt.Println(cs[i])
	}

	return cs
}
