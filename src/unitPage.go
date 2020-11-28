package main

import (
	"fmt"
	"net/http"
	"text/template"
)

var unitTmpl = template.Must(template.ParseFiles(assetTemplates + "unit.html"))

type skillStruct struct {
	Name  string
	Value string
}

type equipmentStruct struct {
	Slot  string
	Value string
}

type equipmentSetStruct struct {
	Name       string
	Equipments []*equipmentStruct
}

type unitStruct struct {
	Name          string
	Skills        []*skillStruct
	Equipments    []*equipmentStruct
	EquipmentSets []*equipmentSetStruct
}

func loadUnit(unitName string) *unitStruct {
	xpath := fmt.Sprintf("//NPCCharacter[@id='%s']", unitName)
	character := game.SpNpc.FindElement(xpath)

	if character == nil {
		panic("Character not found!")
	}

	strct := new(unitStruct)
	strct.Name = unitName

	// get skills
	for _, s := range character.FindElements("skills/skill") {
		skill := new(skillStruct)
		skill.Name = s.SelectAttrValue("id", "ERROR")
		skill.Value = s.SelectAttrValue("value", "ERROR")
		strct.Skills = append(strct.Skills, skill)
	}

	// get other equipments
	for _, e := range character.FindElements("equipment") {
		equipment := new(equipmentStruct)
		equipment.Slot = e.SelectAttrValue("slot", "ERROR")
		equipment.Value = e.SelectAttrValue("id", "ERROR")
		strct.Equipments = append(strct.Equipments, equipment)
	}

	// get equipment sets
	for i, es := range character.FindElements("equipmentSet") {
		equipmentSet := new(equipmentSetStruct)
		equipmentSet.Name = "Equipment_set_" + fmt.Sprint(i)

		for _, e := range es.FindElements("equipment") {
			equipment := new(equipmentStruct)
			equipment.Slot = e.SelectAttrValue("slot", "ERROR")
			equipment.Value = e.SelectAttrValue("id", "ERROR")
			equipmentSet.Equipments = append(equipmentSet.Equipments, equipment)
		}

		strct.EquipmentSets = append(strct.EquipmentSets, equipmentSet)
	}

	return strct
}

func updateStruct(strct *unitStruct, r *http.Request) *unitStruct {

	xpath := fmt.Sprintf("//NPCCharacter[@id='%s']", strct.Name)
	character := game.SpNpc.FindElement(xpath)

	if character == nil {
		panic("Character not found!")
	}

	// update skills
	for _, s := range character.FindElements("skills/skill") {
		name := s.SelectAttrValue("id", "ERROR")

		var skill *skillStruct
		for _, skill = range strct.Skills {
			if skill.Name == name {
				break
			}
		}
		if skill == nil || skill.Name != name {
			panic("Couldn't find element! Trying to find " + name)
		}

		skill.Value = r.FormValue("skill-" + name)

	}

	// update other equipments
	for _, e := range character.FindElements("equipment") {
		slot := e.SelectAttrValue("slot", "ERROR")

		var equipment *equipmentStruct
		for _, equipment = range strct.Equipments {
			if equipment.Slot == slot {
				break
			}
		}
		if equipment == nil || equipment.Slot != slot {
			panic("Couldn't find element! Trying to find " + slot)
		}

		equipment.Value = r.FormValue("equipment-other-" + slot)
	}

	// update equipment sets
	for i, es := range character.FindElements("equipmentSet") {
		equipmentSetName := "Equipment_set_" + fmt.Sprint(i)

		var equipmentSet *equipmentSetStruct
		for _, equipmentSet = range strct.EquipmentSets {
			if equipmentSet.Name == equipmentSetName {
				break
			}
		}
		if equipmentSet.Name != equipmentSetName {
			panic("Couldn't find element! Trying to find " + equipmentSetName)
		}

		for _, e := range es.FindElements("equipment") {
			slot := e.SelectAttrValue("slot", "ERROR")

			var equipment *equipmentStruct
			for _, equipment = range equipmentSet.Equipments {
				if equipment.Slot == slot {
					break
				}
			}
			if equipment == nil || equipment.Slot != slot {
				panic("Couldn't find element! Trying to find " + equipmentSetName + "-" + slot)
			}

			equipment.Value = r.FormValue("equipment-other-" + slot)
		}

	}

	return strct
}

func saveStruct(strct *unitStruct) {
	xpath := fmt.Sprintf("//NPCCharacter[@id='%s']", strct.Name)
	character := game.SpNpc.FindElement(xpath)

	if character == nil {
		panic("Character not found!")
	}

	// save skills
	for _, skill := range strct.Skills {
		xpath := fmt.Sprintf("skills/skill[@id='%s']", skill.Name)
		s := character.FindElement(xpath)
		s.SelectAttr("value").Value = skill.Value
	}

	// save other equipments
	for _, equipment := range strct.Equipments {
		xpath := fmt.Sprintf("equipment[@slot='%s']", equipment.Slot)
		e := character.FindElement(xpath)
		e.SelectAttr("id").Value = equipment.Value
	}

	// save equipment sets
	for i, es := range character.FindElements("equipmentSet") {
		equipmentSetName := "Equipment_set_" + fmt.Sprint(i)
		var equipmentSet *equipmentSetStruct
		for _, equipmentSet = range strct.EquipmentSets {
			if equipmentSet.Name == equipmentSetName {
				break
			}
		}
		if equipmentSet.Name != equipmentSetName {
			panic("Couldn't find element! Trying to find " + equipmentSetName)
		}

		for _, equipment := range equipmentSet.Equipments {
			xpath := fmt.Sprintf("equipment[@slot='%s']", equipment.Slot)
			e := es.FindElement(xpath)
			e.SelectAttr("id").Value = equipment.Value
		}

	}

}

func unitPage(w http.ResponseWriter, r *http.Request) {

	strct := loadUnit(r.URL.RawQuery)

	if r.Method == http.MethodPost {

		strct = updateStruct(strct, r)

		saveStruct(strct)

		http.Redirect(w, r, "/unitselect", http.StatusFound)

	} else {

		if err := unitTmpl.Execute(w, strct); err != nil {
			panic(err)
		}
	}

}
