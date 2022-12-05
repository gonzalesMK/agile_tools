package main

import (
	"sync"
)

type Broadcaster struct {
	rooms map[uint]*Room
	mu    sync.RWMutex
}

type Room struct {
	channels map[uint]*chan []byte
	mu       sync.RWMutex
}

func NewBroadcaster() *Broadcaster {

	return &Broadcaster{
		rooms: map[uint]*Room{},
		mu:    sync.RWMutex{},
	}
}

func (s *Broadcaster) AddSubscriber(roomId uint, userId uint) *chan []byte {
	// handle subscriber already exists

	s.mu.RLock()
	room, exist := s.rooms[roomId]
	s.mu.RUnlock()

	if !exist {
		room = &Room{
			channels: map[uint]*chan []byte{},
			mu:       sync.RWMutex{},
		}
		s.mu.Lock()
		s.rooms[roomId] = room
		s.mu.Unlock()
	}

	channel := make(chan []byte, 1)
	room.mu.Lock()
	room.channels[userId] = &channel
	room.mu.Unlock()

	return &channel
}

func (b *Broadcaster) SendMessage(roomId uint, message []byte) {

	room := b.rooms[roomId]

	room.mu.RLock()
	defer room.mu.RUnlock()
	for key := range room.channels {
		(*room.channels[key]) <- message
	}

}

func (b *Broadcaster) DeleteSubscriber(roomId uint, userId uint) {
	// handle subscriber does not exist

	room := b.rooms[roomId]

	room.mu.Lock()

	channel := room.channels[userId]
	close(*channel)
	delete(room.channels, userId)
	room.mu.Unlock()
}
