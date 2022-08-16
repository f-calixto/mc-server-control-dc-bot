# A discord bot to control a Minecraft server

The server is running on a GCP compute engine instance, in order to save money, I created this bot that shuts the machine down whenever the server remains inactive for 30 minutes. If someone wants to play, the person simply needs to write "start-server" on the discord channel where the bot is listening.

## Example

#### Starting the server
![starting the server from discord](https://admin.francisco-calixto.com/static/start_server.png)


#### After 30 minutes of inactivity
![30 minutes of inactivity](https://admin.francisco-calixto.com/static/shut_down.png)

## The main logic followed by the bot

![enter image description here](https://admin.francisco-calixto.com/static/server_bot_discord_logic.svg)

A more detailed blog about the project: _insert blog link_
