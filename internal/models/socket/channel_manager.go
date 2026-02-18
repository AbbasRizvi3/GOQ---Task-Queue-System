package socket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/olahol/melody"
)

type ChannelManager struct {
	Channels map[string]*Channel
	mu       sync.RWMutex
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		Channels: make(map[string]*Channel),
	}
}

func (cm *ChannelManager) AddClientToChannel(userID string, s *melody.Session) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.Channels[userID] == nil {
		cm.Channels[userID] = NewChannel()
	}
	cm.Channels[userID].AddClient(s)
	totalConns := 0
	for _, ch := range cm.Channels {
		totalConns += len(ch.Clients)
	}
	log.Printf("[WS] User %s connected. Active Users: %d | Total Tabs: %d\n", userID, len(cm.Channels), totalConns)
}

func (cm *ChannelManager) RemoveClientFromChannel(userID string, s *melody.Session) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if channel, exists := cm.Channels[userID]; exists {
		channel.RemoveClient(s)
		if len(channel.Clients) == 0 {
			delete(cm.Channels, userID)
		}
	}
}
func (cm *ChannelManager) BroadcastJSON(userID string, v interface{}) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	channel, exists := cm.Channels[userID]
	if !exists {
		return
	}

	msg, err := json.Marshal(v)
	if err != nil {
		log.Printf("Failed to encode websocket message for user %s: %v", userID, err)
		return
	}

	channel.Broadcast(msg)
}
