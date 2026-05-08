package markdown

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Note is a markdown note ready to be written to disk.
type Note struct {
	Filename string
	Content  string
}

// Convert transforms a Telegram message into a Note.
func Convert(msg *tgbotapi.Message) Note {
	t := time.Unix(int64(msg.Date), 0).Local()
	filename := t.Format("2006-01-02-150405") + ".md"
	title := t.Format("2006-01-02 15:04")
	content := fmt.Sprintf("# %s\n\n%s\n", title, msg.Text)
	return Note{Filename: filename, Content: content}
}
