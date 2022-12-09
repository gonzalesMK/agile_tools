package main

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeToRoom(t *testing.T) {
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

	playerSubscribe := PlayerSubscribeMock{}.AllFields()

	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(123))).
		Return(users, nil)

	subsFunc, err := service.Subscribe(playerSubscribe)

	assert.Nil(t, err)
	assert.NotNil(t, subsFunc)

	// Use closure to get channel result
	b := new(bytes.Buffer)
	w := bufio.NewWriter(b)

	go subsFunc(w)

	service.broadcaster.SendMessage(123, []byte(""))

	// Assert channel has message
	assert.Equal(t, "data: {\"players\":[{\"id\":1,\"name\":\"Santhia Witchy\",\"status\":-3}]}\n\ndata: \n\n", string(b.Bytes()))
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
		DeleteById(gomock.AssignableToTypeOf(&Users{}), uint(1)).
		Return(nil)

	users := []Users{
		{ID: 1, Name: "Santhia Witchy", Status: -3},
	}
	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(123))).
		Return(users, nil)

	service.DeletePlayer(1, 123)

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

func TestUpsertPlayer(t *testing.T) {

	user := UserMocks{}.AllFields()
	request := PlayerRequestMocks{}.AllFields()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepoInterface(ctrl)
	channel := make(chan []byte, 1)
	service := ServiceMock.NewServiceWithChannel(repo, &channel)

	repo.
		EXPECT().
		Save(gomock.Eq(user)).
		Return(nil)

	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(123))).
		Return(nil, nil)

	player, err := service.UpsertPlayer(request)

	assert.Nil(t, err)
	assert.Equal(t, uint(123), player.ID)

	resp := <-(*&channel)
	assert.Equal(t, "{\"players\":[]}", string(resp))
}
