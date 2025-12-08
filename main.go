package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

func main() {
	// Charger .env si présent
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("⚠️ Aucun fichier .env trouvé, utilisation variables système : %v", err)
	}

	// Vérifier si lancement en service Windows
	isService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("Erreur IsWindowsService: %v", err)
	}

	if isService {
		runWindowsService()
		return
	}

	// Mode console (développement ou test)
	log.Println("Lancement SyncApp en mode console…")
	StartSyncLoop()
}
