package main

import (
	"encoding/json"
	"log"
	"os"
)

const (
	contentDir = "content/"
	tabDir     = "tabs/"
)

// tabs are the intermediary between the main model and upgrades
type tab struct {
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	selection int
	Upgrades  []upgrade `json:"upgrades"`
}

// upgrades hold the core math behind this game !
type upgrade struct {
	Description string  `json:"name"`
	Cost        float64 `json:"cost"`
	GrowthRate  float64 `json:"growthRate"`
	Production  float64 `json:"production"`
	owned       uint64
}

// base settings (fps/grid dims)
type config struct {
	Fps  int `json:"fps"`
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

// loads tabs from ./resources/tabs/
func loadBaseModel() model {
	tabNames := [5]string{"main", "sales", "kitchen", "ingred", "settings"}
	tabs := make([]tab, len(tabNames))

	for i, name := range tabNames {
		content, err := os.ReadFile(contentDir + tabDir + name + "Tab.json")
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(content, &tabs[i])
		if err != nil {
			log.Fatal(err)
		}
	}

	return model{tabs, 0}
}

func loadConfig() config {
	var out config

	content, err := os.ReadFile(contentDir + "config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, &out)
	if err != nil {
		log.Fatal(err)
	}

	return out
}
