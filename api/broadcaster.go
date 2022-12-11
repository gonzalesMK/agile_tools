package main

import (
	"errors"
	"fmt"
	"sync"
)

type Broadcaster struct {
	rooms map[uint]*RoomMap
	mu    sync.RWMutex
}

type RoomMap struct {
	channels map[uint]*chan []byte
	mu       sync.RWMutex
}

func NewBroadcaster() *Broadcaster {

	return &Broadcaster{
		rooms: map[uint]*RoomMap{},
		mu:    sync.RWMutex{},
	}
}

func (s *Broadcaster) AddSubscriber(roomId uint, userId uint) *chan []byte {
	// handle subscriber already exists

	s.mu.RLock()
	room, exist := s.rooms[roomId]
	s.mu.RUnlock()

	if !exist {
		room = &RoomMap{
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

func (b *Broadcaster) SendMessage(roomId uint, message []byte) error {

	room, exists := b.rooms[roomId]

	if !exists {
		return errors.New(fmt.Sprint(roomId) + " not found")
	}
	room.mu.RLock()
	defer room.mu.RUnlock()
	for key := range room.channels {
		(*room.channels[key]) <- message
	}

	return nil
}

func (b *Broadcaster) DeleteSubscriber(roomId uint, userId uint) error {
	// handle subscriber does not exist

	room, exists := b.rooms[roomId]
	if !exists {
		return errors.New(fmt.Sprint(roomId) + " not found")
	}

	room.mu.Lock()

	channel := room.channels[userId]
	close(*channel)
	delete(room.channels, userId)
	room.mu.Unlock()

	return nil
}
