package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	bot "github.com/coding-kiko/mc-server-control-dc-bot/discord-bot"
	instance "github.com/coding-kiko/mc-server-control-dc-bot/gcp-compute-instance"
)

var (
	projectId      = "minecraft-server-357716" //"minecraft-server-357716"
	instanceZone   = "southamerica-east1-a"    //"southamerica-east1-a"
	instanceName   = "minecraft-server"        //"minecraft-server"
	credFileBase64 = os.Getenv("GCP_CREDS_JSON_BASE64")
	botToken       = os.Getenv("BOT_TOKEN")
)

func main() {
	logger := log.Default()
	logger.SetFlags(log.LstdFlags)

	it := instance.New(projectId, instanceZone, instanceName, credFileBase64)
	bt := bot.New(logger, it, botToken)
	dg := bt.Init()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
