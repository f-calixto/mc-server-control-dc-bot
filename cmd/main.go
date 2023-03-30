package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	bot "github.com/coding-kiko/mc-server-control-dc-bot/internal/discord-bot"
	instance "github.com/coding-kiko/mc-server-control-dc-bot/internal/gcp-compute-instance"
	"github.com/coding-kiko/mc-server-control-dc-bot/internal/playerCount"
)

var (
	projectId         = os.Getenv("PROJECT_ID")
	instanceZone      = os.Getenv("INSTANCE_ZONE")
	instanceName      = os.Getenv("INSTANCE_NAME")
	credFileBase64    = os.Getenv("GCP_CREDS_JSON_BASE64")
	discordBotToken   = os.Getenv("DC_BOT_TOKEN")
	discordChannelId  = os.Getenv("DC_CHANNEL_ID")
	minecraftServerIp = os.Getenv("MC_SERVER_IP")
)

func main() {
	logger := log.Default()
	logger.SetFlags(log.LstdFlags)

	playerCountClient := playerCount.NewClient(minecraftServerIp)

	instanceController := instance.New(projectId, instanceZone, instanceName, credFileBase64)

	b := bot.New(*logger, instanceController, playerCountClient)

	dgSession := b.Init(discordBotToken, discordChannelId)
	defer dgSession.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
