package client

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	manager "github.com/punkystone/twitch-token-manager"
)

const defaultReconnectInterval = time.Second * 5

type Client struct {
	IRCClient         *twitch.Client
	tokenManager      *manager.TokenManager
	accessToken       string
	refreshToken      string
	reconnectInterval time.Duration
	ErrorChannel      chan error
}

func NewClient(clientID string, clientSecret string, username string, accessToken string, refreshToken string, reconnectInterval *time.Duration) *Client {
	tokeManager := manager.NewTokenManager(clientID, clientSecret)
	client := twitch.NewClient(username, fmt.Sprintf("oauth:%s", accessToken))
	client.TLS = true
	interval := defaultReconnectInterval
	if reconnectInterval != nil {
		interval = *reconnectInterval
	}
	return &Client{
		IRCClient:         client,
		tokenManager:      tokeManager,
		accessToken:       accessToken,
		refreshToken:      refreshToken,
		reconnectInterval: interval,
		ErrorChannel:      make(chan error),
	}
}

func (client *Client) Connect() {
	for {
		err := client.IRCClient.Connect()
		if errors.Is(err, twitch.ErrLoginAuthenticationFailed) {
			success, refreshedAccessToken, refreshedRefreshToken, err := client.tokenManager.ValidateAndRefreshToken(client.accessToken, client.refreshToken)
			if err != nil {
				client.sendError(fmt.Errorf("failed to refresh token: %w", err))
				time.Sleep(client.reconnectInterval)
				continue
			}
			if !success {
				client.sendError(fmt.Errorf("failed to refresh token: %w", err))
				time.Sleep(client.reconnectInterval)
				continue
			}
			client.accessToken = refreshedAccessToken
			client.refreshToken = refreshedRefreshToken
			err = client.IRCClient.Disconnect()
			if err != nil {
				client.sendError(fmt.Errorf("failed to disconnect: %w", err))
			}
			client.IRCClient.SetIRCToken(fmt.Sprintf("oauth:%s", refreshedAccessToken))
			continue
		} else if err != nil {
			client.sendError(fmt.Errorf("failed to connect: %w", err))
			time.Sleep(client.reconnectInterval)
			continue
		}
	}
}

func (client *Client) sendError(err error) {
	select {
	case client.ErrorChannel <- err:
	default:
		log.Println("go-twitch-irc error channel is full, dropping error")
	}
}
