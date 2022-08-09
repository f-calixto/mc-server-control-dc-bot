package bot

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	instance "github.com/coding-kiko/mc-server-control-dc-bot/gcp-compute-instance"
	ssc "github.com/coding-kiko/mc-server-control-dc-bot/sv-status-client"
)

const defaultChannelId = "1005923849684123678" // start-server channel id

type Bot struct {
	logger   *log.Logger
	Token    string
	Instance *instance.Instance
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "start-server" {
		status := b.Instance.GetStatus()
		if status == "STAGING" {
			b.logger.Println("Start attempt while staging")
			s.ChannelMessageSend(m.ChannelID, "Wait... Server is starting")
			return
		}
		if status == "RUNNING" {
			b.logger.Println("Start attempt while running")
			s.ChannelMessageSend(m.ChannelID, "Server is already running")
			return
		}
		if status == "STOPPING" {
			b.logger.Println("Start attempt while stopping")
			s.ChannelMessageSend(m.ChannelID, "Server is stopping, wait a moment before starting again")
			return
		}
		err := b.Instance.Start()
		if err != nil {
			b.logger.Fatal(err)
		}
		b.logger.Println("Starting server")
		s.ChannelMessageSend(m.ChannelID, "Server starting... This could take a few seconds")
		b.WaitForInactivity(s, m.ChannelID)
	}
}

func (b *Bot) WaitForInactivity(s *discordgo.Session, channelId string) {
	for {
		time.Sleep(2 * time.Minute)
		n, err := ssc.GetPlayerCount()
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
			a, err = ssc.GetPlayerCount()
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
		c, err := ssc.GetPlayerCount()
		if err != nil {
			b.logger.Fatalln("error getting server player count")
		}

		if c == 0 {
			b.Instance.Stop()
			b.logger.Println("Stopping server")
			s.ChannelMessageSend(channelId, "30 minutes of inactivity - Stopping server...")
			s.ChannelMessageSend(channelId, "Use `start-server` to spin up the server again")
			return
		}
	}
}

func (b *Bot) Init() *discordgo.Session {
	dg, err := discordgo.New("Bot " + b.Token)
	if err != nil {
		b.logger.Fatalln("error creating Discord session,", err)
	}

	dg.AddHandler(b.messageCreate)

	// we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		b.logger.Fatalln("error opening connection,", err)
	}

	if b.Instance.GetStatus() == "RUNNING" { // case: server started before bot started listening
		go b.WaitForInactivity(dg, defaultChannelId)
	}

	b.logger.Println("Bot is now running.")
	return dg
}

func New(logger *log.Logger, it *instance.Instance, tk string) *Bot {
	return &Bot{
		logger:   logger,
		Token:    tk,
		Instance: it,
	}
}
