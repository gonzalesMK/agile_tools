package main

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeToRoomAndShow(t *testing.T) {
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
		GetPlayersFromRoom(gomock.Eq(uint(12))).
		Return(users, nil)
	repo.
		EXPECT().
		GetOneById(gomock.AssignableToTypeOf(new(Room)), gomock.Eq(uint(12))).
		DoAndReturn(func(model *Room, id uint) error {
			model.Show = true
			return nil
		})

	subsFunc, err := service.Subscribe(playerSubscribe)

	assert.Nil(t, err)
	assert.NotNil(t, subsFunc)

	// Use closure to get channel result
	b := new(bytes.Buffer)
	w := bufio.NewWriter(b)

	go subsFunc(w)

	service.broadcaster.SendMessage(12, []byte(""))

	// Assert channel has message
	assert.Equal(t, "data: {\"id\":1}\n\ndata: {\"players\":[{\"id\":1,\"name\":\"Santhia Witchy\",\"status\":-3}]}\n\ndata: \n\n", string(b.Bytes()))
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
		GetPlayersFromRoom(gomock.Eq(uint(12))).
		Return(users, nil)
	repo.
		EXPECT().
		GetOneById(gomock.AssignableToTypeOf(new(Room)), gomock.Eq(uint(12))).
		Return(nil)

	if err := service.DeletePlayer(1, 12); err != nil {
		assert.Nil(t, err)
		return
	}
	assert.Equal(t, len(channel), 0)
	room, exists := service.broadcaster.rooms[12]

	if !exists {
		assert.True(t, exists)
		return
	}
	assert.Nil(t, room.channels[1])
}

func TestBroadcastRoomStatus(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepoInterface(ctrl)
	channel := make(chan []byte, 1)
	service := ServiceMock.NewServiceWithChannel(repo, &channel)

	users := []Users{
		{ID: 1, Name: "Santhia Witchy", Status: -3},
		{ID: 2, Name: "Another One", Status: 0},
	}
	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(12))).
		Return(users, nil)

	repo.
		EXPECT().
		GetOneById(gomock.AssignableToTypeOf(new(Room)), gomock.Eq(uint(12))).
		Return(nil)

	if err := service.BroadcastRoomStatus(12); err != nil {
		assert.Nil(t, err)
		return
	}

	result := <-channel

	assert.Equal(t, "{\"players\":[{\"id\":1,\"name\":\"Santhia Witchy\",\"status\":-3},{\"id\":2,\"name\":\"Another One\",\"status\":-1}]}", string(result))
}

func TestBroadcastRoomStatusEmptyWorks(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepoInterface(ctrl)
	channel := make(chan []byte, 1)
	service := ServiceMock.NewServiceWithChannel(repo, &channel)

	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(12))).
		Return(nil, nil)
	repo.
		EXPECT().
		GetOneById(gomock.AssignableToTypeOf(new(Room)), gomock.Eq(uint(12))).
		Return(nil)

	err := service.BroadcastRoomStatus(12)
	assert.Nil(t, err)

	result := <-channel

	assert.Equal(t, "{\"players\":[]}", string(result))
}

func TestUpdatePlayerWorks(t *testing.T) {

	user := UserMocks{}.AllFields()
	user.Room = Room{}
	request := PlayerRequestMocks{}.AllFields()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepoInterface(ctrl)
	channel := make(chan []byte, 1)
	service := ServiceMock.NewServiceWithChannel(repo, &channel)

	repo.
		EXPECT().
		UpdateFieldById(gomock.Eq(uint(123)), gomock.Eq(user)).
		Return(nil)

	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(12))).
		Return(nil, nil)
	repo.
		EXPECT().
		GetOneById(gomock.AssignableToTypeOf(new(Room)), gomock.Eq(uint(12))).
		Return(nil)

	player, err := service.UpsertPlayer(request)

	assert.Nil(t, err)
	assert.Equal(t, uint(123), player.ID)

	resp := <-(*&channel)
	assert.Equal(t, "{\"players\":[]}", string(resp))
}

func TestUpdateRoomWorks(t *testing.T) {

	request := RoomRequestMock{}.AllFields()
	room := RoomMock{}.AllFields()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepoInterface(ctrl)
	channel := make(chan []byte, 1)
	service := ServiceMock.NewServiceWithChannel(repo, &channel)

	repo.
		EXPECT().
		UpdateFieldById(gomock.Eq(uint(12)), gomock.Eq(room)).
		Return(nil)

	repo.
		EXPECT().
		GetPlayersFromRoom(gomock.Eq(uint(12))).
		Return(nil, nil)
	repo.
		EXPECT().
		GetOneById(gomock.AssignableToTypeOf(new(Room)), gomock.Eq(uint(12))).
		Return(nil)

	player, err := service.UpdateRoom(request)

	assert.Nil(t, err)
	assert.Equal(t, uint(12), player.ID)

	resp := <-(*&channel)
	assert.Equal(t, "{\"players\":[]}", string(resp))
}
