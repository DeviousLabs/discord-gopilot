package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func StartBot(token string) {
	// Create and initialize the bot
	bot, err := NewBot(token)
	if err != nil {
		fmt.Printf("Error creating bot: %s\n", err)
		return
	}

	// Start the bot
	if err := bot.Start(); err != nil {
		fmt.Printf("Error starting bot: %s\n", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	handleSigInt(bot)
}

func handleSigInt(bot *Bot) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("Shutting down the bot...")
	bot.Stop()
}
