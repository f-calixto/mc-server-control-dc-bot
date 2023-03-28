# A discord bot to control a minecraft server

The server is running on a GCP compute engine instance, in order to save money, I created this bot that shuts the machine down whenever the server remains inactive for 30 minutes. If someone wants to play, the person simply needs to write "start-server" on the discord channel where the bot is listening. 

The minecraft server runs in the virtual machine as a service, so whenever the machine stops, it stops. When the machine starts the server starts.

## Example

#### Starting the server
![starting the server from discord](https://github.com/coding-kiko/mc-server-control-dc-bot/blob/main/docs/start_server.png?raw=true)


#### After 30 minutes of inactivity
![30 minutes of inactivity](https://github.com/coding-kiko/mc-server-control-dc-bot/blob/main/docs/shut_down_server.png?raw=true)

## The main logic of the bot

![enter image description here](https://github.com/coding-kiko/mc-server-control-dc-bot/blob/main/docs/server_bot_discord_logic.svg?raw=true)

The machine can be turned on and off easily using the [Compute Engine API](https://pkg.go.dev/google.golang.org/api/compute/v1) and the server player count is retrieved from https://mcapi.us/. 
The reason for waiting in intervals of 2 minutes before checking the amount of players connected is because the API maintainers specifically tell us to make at most one request per minute as the data is stored server side, and as I didn't really need data every second I decided to call every 2 minutes and prevent any kind of issue.
