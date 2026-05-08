package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rwirdemann/telegram-obsidian-sync/internal/config"
	"github.com/rwirdemann/telegram-obsidian-sync/internal/markdown"
	"github.com/rwirdemann/telegram-obsidian-sync/internal/state"
	"github.com/rwirdemann/telegram-obsidian-sync/internal/storage"
)

const longPollTimeout = 30

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatalf("telegram: %v", err)
	}
	log.Printf("authorized as @%s", bot.Self.UserName)

	home, _ := os.UserHomeDir()
	stateDir := filepath.Join(
		home, ".config", "telegram-obsidian-sync",
	)

	s, err := state.Load(stateDir)
	if err != nil {
		log.Fatalf("state: %v", err)
	}

	chatSet := make(map[int64]bool, len(cfg.Telegram.ChatIDs))
	for _, id := range cfg.Telegram.ChatIDs {
		chatSet[id] = true
	}

	writer := storage.NewWriter(cfg.Sync.InboxDir)
	offset := s.Offset

	for {
		u := tgbotapi.NewUpdate(offset)
		u.Timeout = longPollTimeout

		updates, err := bot.GetUpdates(u)
		if err != nil {
			log.Printf("poll error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for i := range updates {
			upd := &updates[i]
			if upd.UpdateID+1 > offset {
				offset = upd.UpdateID + 1
			}
			if upd.Message == nil || upd.Message.Text == "" {
				continue
			}
			if len(chatSet) > 0 && !chatSet[upd.Message.Chat.ID] {
				continue
			}
			note := markdown.Convert(upd.Message)
			if err := writer.Write(note); err != nil {
				log.Printf("write %s: %v", note.Filename, err)
				continue
			}
			log.Printf("saved %s", note.Filename)
		}

		if offset != s.Offset {
			s.Offset = offset
			if err := state.Save(stateDir, s); err != nil {
				log.Printf("state save: %v", err)
			}
		}
	}
}
