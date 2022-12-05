package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlayer(t *testing.T) {
	//
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockRepoInterface(ctrl)

	service := Service{
		repo:        repo,
		broadcaster: NewBroadcaster(),
	}

	repo.
		EXPECT().
		Save(gomock.AssignableToTypeOf(&Users{})).
		DoAndReturn(func(u *Users) error {
			u.ID = 1
			return nil
		})

	users := []Users{
		{ID: 1, Name: "Santhia Witchy", Status: -3},
	}
	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(123))).
		Return(users, nil)

	player, channel, err := service.CreatePlayer("Santhia Witchy", -3, 123)

	assert.Nil(t, err)
	assert.NotNil(t, channel)

	// Assert player created correctly
	assert.Exactly(t, &Players{ID: 1, Name: "Santhia Witchy", Status: -3}, player)
	assert.Equal(t, channel, service.broadcaster.rooms[123].channels[1])

	// Assert channel has message
	result := <-*channel
	assert.Equal(t, "{\"players\":[{\"id\":1,\"name\":\"Santhia Witchy\",\"status\":-3}]}", string(result))
}

func TestDeletePlayer(t *testing.T) {
	//
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepoInterface(ctrl)
	channel := make(chan []byte, 1)
	service := ServiceMock.NewServiceWithChannel(repo, &channel)

	repo.
		EXPECT().
		DeleteById(gomock.AssignableToTypeOf(&Users{})).
		DoAndReturn(func(u *Users) error {
			assert.Equal(t, &Users{ID: 1, Name: "Santhia Witchy", Status: -3}, u)
			return nil
		})
	users := []Users{
		{ID: 1, Name: "Santhia Witchy", Status: -3},
	}
	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(123))).
		Return(users, nil)

	player := &Players{ID: 1, Name: "Santhia Witchy", Status: -3}

	service.DeletePlayer(player, 123)

	assert.Equal(t, len(channel), 0)
	assert.Nil(t, service.broadcaster.rooms[123].channels[1])
}

func TestBroadcastRoomStatus(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepoInterface(ctrl)
	channel := make(chan []byte, 1)
	service := ServiceMock.NewServiceWithChannel(repo, &channel)

	users := []Users{
		Users{ID: 1, Name: "Santhia Witchy", Status: -3},
	}
	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(123))).
		Return(users, nil)

	err := service.BroadcastRoomStatus(123)

	result := <-channel

	assert.Nil(t, err)
	assert.Equal(t, "{\"players\":[{\"id\":1,\"name\":\"Santhia Witchy\",\"status\":-3}]}", string(result))
}

func TestBroadcastRoomStatusEmpty(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepoInterface(ctrl)
	channel := make(chan []byte, 1)
	service := ServiceMock.NewServiceWithChannel(repo, &channel)

	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(123))).
		Return(nil, nil)

	err := service.BroadcastRoomStatus(123)

	result := <-channel

	assert.Nil(t, err)
	assert.Equal(t, "{\"players\":[]}", string(result))
}
