package main

import (
	"TgBot/config"
	"context"
	"fmt"
	"os"
	"os/signal"

	"TgBot/content"

	provider2 "github.com/Hogants/tg-bot-proto/gen/go/content"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var providers map[string]*content.Provider

func main() {
	cfg := config.LoadConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	providers = make(map[string]*content.Provider)
	for _, p := range cfg.ContentProviders {
		provider, err := content.New(p)
		if err != nil {
			panic("Failed to create provider: " + err.Error())
		}
		providers[p.Name] = provider
	}
	fmt.Println(providers)
	b, err := bot.New(cfg.TelegramToken, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)

}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	for _, provider := range providers {
		features, err := provider.Client.Features(ctx, &provider2.FeaturesRequest{
			UserId: update.Message.Chat.ID,
		})
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, msg := range features.Features {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   msg.Name + msg.Text,
			})
		}
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
