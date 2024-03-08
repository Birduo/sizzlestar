package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"log"
	"os"
)

const (
	contentDir = "content/"
	tabDir     = "tabs/"
)

var state gameState

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

type upgradeState struct {
	Description string
	owned       uint64
}

type gameState struct {
	lastLogged float64

	yen float64 // sales -> selling fr for yen
	fr  float64 // fried rice -> selling veg for fr
	veg float64 // vegetables -> farming ingredients

	sales   []upgradeState
	kitchen []upgradeState
	ingred  []upgradeState
}

// loads tabs from ./content/tabs/
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

	var state gameState
	out := model{state, tabs, 0}
	out.loadGameState()

	return out
}

// loads associated json data from ./content/config.json
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

func (m *model) saveGameState() {
	content, err := m.state.MarshalGameState()
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(contentDir+"save.szs", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *model) loadGameState() {
	if _, err := os.Stat(contentDir + "save.zsz"); err == nil {
		content, readErr := os.ReadFile(contentDir + "save.szs")
		if readErr != nil {
			log.Fatal(err)
		}

		err = m.state.UnmarshalGameState(content)
		if err != nil {
			log.Fatal(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// file doens't exist, make a new gameState
		m.initGameState()
	} else {
		log.Fatal(err)
	}
}

func (m *model) initGameState() {
	m.state = gameState{
		lastLogged: 0,
		yen:        0,
		fr:         0,
		veg:        0,
		sales:      []upgradeState{},
		kitchen:    []upgradeState{},
		ingred:     []upgradeState{},
	}

	for _, t := range m.tabs {
		switch t.Name {
		case "Sales":
			for _, u := range t.Upgrades {
				m.state.sales = append(m.state.sales, upgradeState{u.Description, 0})
			}
		case "Kitchen":
			for _, u := range t.Upgrades {
				m.state.kitchen = append(m.state.kitchen, upgradeState{u.Description, 0})
			}
		case "Ingred":
			for _, u := range t.Upgrades {
				m.state.ingred = append(m.state.ingred, upgradeState{u.Description, 0})
			}
		}
	}
}

func (g gameState) MarshalGameState() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(g.lastLogged)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(g.yen)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(g.fr)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(g.veg)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(g.sales)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(g.kitchen)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(g.ingred)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *gameState) UnmarshalGameState(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&g.lastLogged)
	if err != nil {
		return err
	}
	err = dec.Decode(&g.yen)
	if err != nil {
		return err
	}
	err = dec.Decode(&g.fr)
	if err != nil {
		return err
	}
	err = dec.Decode(&g.veg)
	if err != nil {
		return err
	}
	err = dec.Decode(&g.sales)
	if err != nil {
		return err
	}
	err = dec.Decode(&g.kitchen)
	if err != nil {
		return err
	}
	err = dec.Decode(&g.ingred)
	if err != nil {
		return err
	}
	return nil
}
