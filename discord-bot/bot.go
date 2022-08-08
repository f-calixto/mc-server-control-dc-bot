package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	instance "github.com/coding-kiko/mc-server-control-dc-bot/gcp-compute-instance"
	ssc "github.com/coding-kiko/mc-server-control-dc-bot/sv-status-client"
)

const defaultChannelId = "1005923849684123678" // start-server channel id

type Bot struct {
	Token    string
	Instance *instance.Instance
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "start-server" {
		if b.Instance.GetStatus() == "RUNNING" {
			s.ChannelMessageSend(m.ChannelID, "Server is already running    chupapija")
			return
		}
		if b.Instance.GetStatus() == "STOPPING" {
			s.ChannelMessageSend(m.ChannelID, "Server is stopping, wait a moment before starting again")
			return
		}
		err := b.Instance.Start()
		if err != nil {
			log.Fatal(err)
		}
		s.ChannelMessageSend(m.ChannelID, "Server starting... This could take a few seconds")
		b.WaitForInactivity(s, m.ChannelID)
	}
}

func (b *Bot) WaitForInactivity(s *discordgo.Session, channelId string) {
	for {
		time.Sleep(2 * time.Minute)
		n, err := ssc.GetPlayerCount()
		if err != nil {
			log.Fatalln("error getting server player count")
		}

		if n == 0 {
			log.Println("started 30 min counter")
			time.Sleep(30 * time.Minute)
			n, err = ssc.GetPlayerCount()
			if err != nil {
				log.Fatalln("error getting server player count")
			}
			if n == 0 {
				b.Instance.Stop()
				s.ChannelMessageSend(channelId, "30 minutes of inactivity - Stopping server...")
				s.ChannelMessageSend(channelId, "Use `start-server` to spin up the server again")
				return
			}
		}
	}
}

func (b *Bot) Init() *discordgo.Session {
	dg, err := discordgo.New("Bot " + b.Token)
	if err != nil {
		log.Fatalln("error creating Discord session,", err)
	}

	dg.AddHandler(b.messageCreate)

	// we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
	}

	if b.Instance.GetStatus() == "RUNNING" { // case: server started before bot started listening
		go b.WaitForInactivity(dg, defaultChannelId)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	return dg
}

func New(it *instance.Instance, tk string) *Bot {
	return &Bot{
		Token:    tk,
		Instance: it,
	}
}
