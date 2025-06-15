# go-twitch-irc

Based on [go-twitch-irc](https://github.com/gempir/go-twitch-irc) by gempir but with a focus on refreshing tokens and reconnecting to Twitch IRC.

## Methods

```go
func NewClient(clientID string, clientSecret string, username string, accessToken string, refreshToken string, reconnectInterval *time.Duration) *Client
func (client *Client) Connect()
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v4"
	client "github.com/punkystone/go-twitch-irc"
)

func main() {
	client := client.NewClient("clientID", "clientSecret", "username", "accessToken", "refreshToken", nil)

	go func() {
        for err := range client.ErrorChannel {
			fmt.Println(err)
		}
	}()

	client.IRCClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Message)
	})

	client.IRCClient.OnConnect(func() {
		client.IRCClient.Say("punkystone", "TriHard")
	})

    client.IRCClient.Join("punkystone")

	client.Connect()
}
```
