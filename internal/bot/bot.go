package bot

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	instance "github.com/coding-kiko/mc-server-control-dc-bot/internal/gcp-compute-instance"
	"github.com/coding-kiko/mc-server-control-dc-bot/internal/playerCount"
)

const (
	statusStaging      = "STAGING"
	statusRunning      = "RUNNING"
	statusStopping     = "STOPPING"
	startServerMessage = "start-server"
)

type Bot struct {
	mu                 sync.Mutex
	logger             log.Logger
	instanceController instance.InstanceController
	playerCountClient  playerCount.Client
}

func (b *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == startServerMessage {
		b.mu.Lock()
		defer b.mu.Unlock()

		switch b.instanceController.GetStatus() {
		case statusStaging:
			b.logger.Println("Start attempt while staging")
			s.ChannelMessageSend(m.ChannelID, "Wait... Server is starting")
			return
		case statusRunning:
			b.logger.Println("Start attempt while running")
			s.ChannelMessageSend(m.ChannelID, "Server is already running")
			return
		case statusStopping:
			b.logger.Println("Start attempt while stopping")
			s.ChannelMessageSend(m.ChannelID, "Server is stopping, wait a moment before starting again")
			return
		}

		if err := b.instanceController.Start(); err != nil {
			b.logger.Fatalln(err)
		}

		b.logger.Println("Starting server")
		s.ChannelMessageSend(m.ChannelID, "Server starting... This could take a few seconds")
		go b.waitForInactivity(s, m.ChannelID)
	}
}

// Waits for 30 minutes of server inactivity before stopping the server instance.
// It uses a ticker to check the player count every 2 minutes, and if no one is playing,
// it increments a minutesInactive by 2. If the counter reaches 30, the server is stopped.
// If anyone starts playing before the 30-minute threshold, the counter is reset.
func (b *Bot) waitForInactivity(s *discordgo.Session, channelId string) {
	var (
		err                    error
		minutesInactive, count int
	)

	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	// every 2 minutes:
	for range ticker.C {
		if minutesInactive == 30 {
			b.instanceController.Stop()
			b.logger.Println("Stopping server")
			s.ChannelMessageSend(channelId, "30 minutes of inactivity - Stopping server...")
			s.ChannelMessageSend(channelId, "Use `start-server` to spin up the server again")
			return
		}

		if count, err = b.playerCountClient.Get(); err != nil {
			b.logger.Fatalln("error getting server player count")
		}

		switch count {
		case 0:
			minutesInactive += 2
			b.logger.Println("No one is playing - minutes inactive:", minutesInactive)
		default:
			minutesInactive = 0
		}
	}
}

func New(logger log.Logger, it instance.InstanceController, pcc playerCount.Client) *Bot {
	return &Bot{
		logger:             logger,
		instanceController: it,
		playerCountClient:  pcc,
	}
}

func (b *Bot) Init(dcBotTkn, dcChanId string) *discordgo.Session {
	dgSession, err := discordgo.New("Bot " + dcBotTkn)
	if err != nil {
		b.logger.Fatalln("error creating Discordgo session,", err)
	}

	dgSession.AddHandler(b.onMessage)

	// we only care about receiving message events.
	dgSession.Identify.Intents = discordgo.IntentsGuildMessages

	if err = dgSession.Open(); err != nil {
		b.logger.Fatalln("error opening connection,", err)
	}

	// case: server started before bot started listening
	if b.instanceController.GetStatus() == statusRunning {
		go b.waitForInactivity(dgSession, dcChanId)
	}

	b.logger.Println("Bot is now running.")
	return dgSession
}
