package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/MedBou33/gold-tracker/internal/alert"
	"github.com/MedBou33/gold-tracker/internal/notify"
	"github.com/MedBou33/gold-tracker/internal/provider"
	"github.com/MedBou33/gold-tracker/internal/store"
)

func check(ctx context.Context, client *provider.Client, notifier *notify.Ntfy, state *store.FileStore, target float64) error {
	quote, err := client.GetQuote(ctx, "XAU", "EUR")
	if err != nil {
		return fmt.Errorf("get price: %w", err)
	}
	log.Printf("current price: %.2f %s", quote.Price, quote.Currency)

	last, found, err := state.Load()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	if found && alert.Crossed(last.LastPrice, quote.Price, target) {
		msg := fmt.Sprintf("Gold is now %.2f %s (target %.2f)", quote.Price, quote.Currency, target)
		if err := notifier.Send(ctx, "Gold alert", msg); err != nil {
			return errors.New("send notification failed")
		}
		log.Println("alert sent")
	}

	if err := state.Save(store.State{LastPrice: quote.Price}); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	return nil
}

func main() {
	topic := os.Getenv("NTFY_TOPIC")
	if topic == "" {
		log.Fatal("NTFY_TOPIC is not set")
	}

	target, err := strconv.ParseFloat(os.Getenv("TARGET_PRICE"), 64)
	if err != nil {
		log.Fatal("TARGET_PRICE must be a number, for example 3500")
	}

	statePath := os.Getenv("STATE_FILE")
	if statePath == "" {
		statePath = "state.json"
	}

	interval := 5 * time.Minute
	if v := os.Getenv("CHECK_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d < time.Minute {
			log.Fatal("CHECK_INTERVAL must be at least 1m, for example 5m")
		}
		interval = d
	}

	client := provider.NewClient("https://api.gold-api.com")
	notifier := notify.NewNtfy("https://ntfy.sh", topic)
	state := store.NewFileStore(statePath)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("tracking gold: target %.2f EUR, checking every %s", target, interval)

	if err := check(ctx, client, notifier, state, target); err != nil {
		log.Printf("check failed: %v", err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("stopping")
			return
		case <-ticker.C:
			if err := check(ctx, client, notifier, state, target); err != nil {
				log.Printf("check failed: %v", err)
			}
		}
	}
}
