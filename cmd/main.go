package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"discord-gopilot/pkg/discord"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("DISCORD_TOKEN is not set.")
	}

	// Initialize and start the Discord bot
	bot, err := discord.NewBot(token)
	if err != nil {
		log.Fatalf("Failed to create Discord bot: %v", err)
	}

	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	// Wait for a termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stop

	// Clean up resources
	bot.Stop()
}
