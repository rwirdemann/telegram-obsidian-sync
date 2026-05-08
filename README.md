# telegram-obsidian-sync

Polls a Telegram bot and saves each text message as a Markdown note
in your Obsidian inbox.

## Prerequisites

- Go 1.22+
- A Telegram account
- A Dropbox (or any local) folder used as your Obsidian inbox

## Step 1 — Create a Telegram bot

1. Open Telegram and search for **@BotFather**.
2. Send `/newbot` and follow the prompts (choose a name and a username).
3. BotFather replies with a **token** that looks like:
   ```
   123456789:AAFxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   ```
   Keep this token — you need it in the config file.

## Step 2 — Find your chat ID

The service filters messages by chat ID so it only processes messages
you send to your own bot.

1. Open Telegram and search for **@userinfobot**.
2. Send `/start`. It replies with your numeric user ID, e.g. `987654321`.
3. Note that number — it is your chat ID for direct messages.

> **Tip:** To use the bot like a personal inbox, just open a chat with
> your own bot and send messages there. Your user ID equals the chat ID
> of that conversation.

## Step 3 — Configure the service

```bash
mkdir -p ~/.config/telegram-obsidian-sync
cp config.toml.example ~/.config/telegram-obsidian-sync/config.toml
```

Edit the file:

```toml
[telegram]
token    = "123456789:AAFxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
chat_ids = [987654321]

[sync]
inbox_dir = "~/Dropbox/Zettelkasten/Inbox"
```

| Key               | Description                                            |
|-------------------|--------------------------------------------------------|
| `token`           | The token from BotFather                               |
| `chat_ids`        | List of chat IDs to process. Empty list = all chats.  |
| `inbox_dir`       | Directory where Markdown files are written             |

## Step 4 — Build and run

```bash
go build -o telegram-obsidian-sync .
./telegram-obsidian-sync
```

Or run directly without building:

```bash
go run .
```

The service logs each saved file to stdout:

```
2026/05/08 14:32:01 authorized as @myobsidianbot
2026/05/08 14:32:45 saved 2026-05-08-143245.md
```

## Output format

Each message becomes a separate Markdown file in your inbox:

**Filename:** `2026-05-08-143245.md`

```markdown
# 2026-05-08 14:32

Your message text here.
```

## State persistence

The last processed update offset is saved to
`~/.config/telegram-obsidian-sync/state.json`. This ensures that
restarting the service never creates duplicate notes.

## Run as a background service (macOS)

Create `~/Library/LaunchAgents/com.rwirdemann.telegram-obsidian-sync.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.rwirdemann.telegram-obsidian-sync</string>
  <key>ProgramArguments</key>
  <array>
    <string>/usr/local/bin/telegram-obsidian-sync</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
  <key>StandardOutPath</key>
  <string>/tmp/telegram-obsidian-sync.log</string>
  <key>StandardErrorPath</key>
  <string>/tmp/telegram-obsidian-sync.log</string>
</dict>
</plist>
```

Then install it:

```bash
# copy the binary to a stable location first
cp telegram-obsidian-sync /usr/local/bin/

launchctl load ~/Library/LaunchAgents/com.rwirdemann.telegram-obsidian-sync.plist
```

Check logs:

```bash
tail -f /tmp/telegram-obsidian-sync.log
```

To stop the service:

```bash
launchctl unload ~/Library/LaunchAgents/com.rwirdemann.telegram-obsidian-sync.plist
```
