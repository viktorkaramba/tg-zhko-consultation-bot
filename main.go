package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Get the bot type from environment variable
	botType := os.Getenv("BOT_TYPE")

	if botType == "" {
		// Default to running both bots
		log.Println("No specific bot type specified, starting both bots")

		// For Docker, this would be handled by supervisord
		log.Println("If running in Docker, please ensure supervisord is properly configured")
		return
	}

	var cmd *exec.Cmd

	switch botType {
	case "consultation":
		log.Println("Starting consultation bot...")
		cmd = exec.Command("./consultation_bot/consultation_bot")
	case "applications":
		log.Println("Starting applications bot...")
		cmd = exec.Command("./applications_bot/applications_bot")
	default:
		log.Printf("Unknown bot type: %s", botType)
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	fmt.Println("Bot stopped")
}
