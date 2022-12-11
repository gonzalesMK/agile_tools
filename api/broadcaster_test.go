package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSubscriber(t *testing.T) {

	b := Broadcaster{
		rooms: map[uint]*RoomMap{},
	}

	channel1 := b.AddSubscriber(10, 15)
	b.AddSubscriber(10, 16)

	room, exists := b.rooms[10]
	assert.True(t, exists)
	assert.NotNil(t, room)
	assert.Equal(t, 1, len(b.rooms))

	channel, cExist := room.channels[15]
	assert.True(t, cExist)
	assert.NotNil(t, channel)
	assert.Equal(t, 2, len(room.channels))
	assert.Equal(t, channel1, channel)
}

func TestDeleteSubscriber(t *testing.T) {

	b := Broadcaster{
		rooms: map[uint]*RoomMap{},
	}

	b.AddSubscriber(10, 15)
	b.DeleteSubscriber(10, 15)

	room, exists := b.rooms[10]
	assert.True(t, exists)
	assert.NotNil(t, room)
	assert.Equal(t, 1, len(b.rooms))

	channel, cExist := room.channels[15]
	assert.False(t, cExist)
	assert.Nil(t, channel)
	assert.Equal(t, 0, len(room.channels))
}

func TestBroadcast(t *testing.T) {

	b := Broadcaster{
		rooms: map[uint]*RoomMap{},
	}
	channel1 := b.AddSubscriber(10, 15)
	channel2 := b.AddSubscriber(10, 16)
	channel3 := b.AddSubscriber(12, 16)

	bytes := []byte{'A', 'B', 'C'}
	b.SendMessage(10, bytes)

	by1 := <-(*channel1)
	by2 := <-(*channel2)

	assert.Equal(t, bytes, by1)
	assert.Equal(t, bytes, by2)
	assert.Len(t, *channel1, 0)
	assert.Len(t, *channel2, 0)
	assert.Len(t, *channel3, 0)
}
