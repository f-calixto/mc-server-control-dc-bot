package bot

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	instance "github.com/coding-kiko/mc-server-control-dc-bot/gcp-compute-instance"
	"github.com/coding-kiko/mc-server-control-dc-bot/playerCount"
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
			b.logger.Fatal(err)
		}

		b.logger.Println("Starting server")
		s.ChannelMessageSend(m.ChannelID, "Server starting... This could take a few seconds")
		b.WaitForInactivity(s, m.ChannelID)
	}
}

// waits 30 minutes and checks playerCount every 2 minutes.
// if the
func (b *Bot) WaitForInactivity(s *discordgo.Session, channelId string) {
	for {
		time.Sleep(2 * time.Minute)
		n, err := b.playerCountClient.Get()
		if err != nil {
			b.logger.Fatalln("error getting server player count")
		}

		if n > 0 {
			continue
		}

		var a int
		b.logger.Println("started 30 min counter")
		for i := 0; i < 14; i++ {
			time.Sleep(2 * time.Minute)
			a, err = b.playerCountClient.Get()
			if err != nil {
				b.logger.Fatalln("error getting server player count")
			}
			if a > 0 {
				break
			}
		}
		if a > 0 { // means that interval has been interrupted by activity
			b.logger.Println("30 minute interval interrupted")
			continue
		}

		time.Sleep(2 * time.Minute)
		c, err := b.playerCountClient.Get()
		if err != nil {
			b.logger.Fatalln("error getting server player count")
		}

		if c == 0 {
			b.instanceController.Stop()
			b.logger.Println("Stopping server")
			s.ChannelMessageSend(channelId, "30 minutes of inactivity - Stopping server...")
			s.ChannelMessageSend(channelId, "Use `start-server` to spin up the server again")
			return
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
		go b.WaitForInactivity(dgSession, dcChanId)
	}

	b.logger.Println("Bot is now running.")
	return dgSession
}
