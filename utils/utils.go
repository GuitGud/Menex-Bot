package utils

import (
	"encoding/json"
	"log"
	"os"

	handlers "Menex-bot/m/cmd" // Adjust if your project structure is different
)

// LoadEvents reads events.json and updates the events slice in the handlers package.
func LoadEvents() {
	file, err := os.Open("events.json")
	if err != nil {
		log.Println("No existing events file found, starting fresh.")
		return
	}
	defer file.Close()

	var loaded []handlers.Event
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&loaded)
	if err != nil {
		log.Printf("Error decoding events: %v", err)
		return
	}

	handlers.SetEvents(loaded)
}
