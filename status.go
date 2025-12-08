package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var statusFile = "data/last_import.json"

func LoadLastImport() (time.Time, error) {
	data, err := os.ReadFile(statusFile)
	if err != nil {
		return time.Time{}, err
	}

	var t time.Time
	err = json.Unmarshal(data, &t)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func SaveLastImport() {
	t := time.Now()
	data, _ := json.Marshal(t)
	os.MkdirAll("data", 0777)
	err := os.WriteFile(statusFile, data, 0666)
	if err != nil {
		log.Printf("⚠️ Impossible de sauvegarder last_import.json : %v", err)
	} else {
		log.Printf("📁 Dernier import sauvegardé : %s", t.Format(time.RFC3339))
	}
}
