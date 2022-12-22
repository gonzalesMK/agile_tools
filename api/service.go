package main

import (
	"bufio"
	"encoding/json"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/ulule/deepcopier"
)

var (
	DEFAULT_STATUS = int8(-2)
)

type Service struct {
	repo        RepoInterface
	mu          sync.RWMutex
	broadcaster *Broadcaster
}

//go:generate mockgen -source=$GOFILE -destination=mock_repository_test.go -package=main
type RepoInterface interface {
	GetOneById(model interface{}, id uint) error
	Save(interface{}) error
	UpdateFieldById(id uint, content interface{}) error
	DeleteById(model interface{}, id uint) error
	GetPlayersFromRoom(uint) ([]Users, error)
	ClearPlayerStatusInRoom(roomID uint, statusID int8) error
}

func (s *Service) ClearRoom(request *RoomClearRequest) error {

	if err := s.repo.ClearPlayerStatusInRoom(request.RoomID, DEFAULT_STATUS); err != nil {
		return err
	}

	if err := s.repo.UpdateFieldById(request.RoomID, &Room{
		ID:   request.RoomID,
		Show: false,
	}); err != nil {
		return err
	}
	s.BroadcastRoomStatus(request.RoomID)

	return nil
}

func (s *Service) UpdateRoom(roomRequest *RoomRequest) (*RoomResponse, error) {

	var room Room
	deepcopier.Copy(roomRequest).To(&room)

	err := s.repo.UpdateFieldById(room.ID, &room)

	var response RoomResponse
	deepcopier.Copy(&room).To(&response)

	s.BroadcastRoomStatus(room.ID)

	return &response, err
}

func (s *Service) UpsertPlayer(playerRequest *PlayerRequest) (*PlayerResponse, error) {

	var user Users
	deepcopier.Copy(playerRequest).To(&user)

	err := s.repo.UpdateFieldById(user.ID, &user)

	var player PlayerResponse
	deepcopier.Copy(&user).To(&player)

	s.BroadcastRoomStatus(user.RoomID)

	return &player, err
}

func (s *Service) Subscribe(playerSubscribe *PlayerSubscribe) (func(w *bufio.Writer), error) {

	var user Users
	deepcopier.Copy(playerSubscribe).To(&user)
	user.Status = DEFAULT_STATUS
	user.Room = Room{ID: user.RoomID}
	err := s.repo.Save(&user)

	if err != nil {
		log.Errorf("Not able to save new user %s", err)
		return nil, err
	}

	channel := s.broadcaster.AddSubscriber(playerSubscribe.RoomID, user.ID)

	closure := func(w *bufio.Writer) {
		timeout := make(chan bool, 1)
		go func() {
			for {
				timeout <- true
				time.Sleep(time.Millisecond * 100)
			}
		}()

		// Send own player info
		var player PlayerResponse
		deepcopier.Copy(&user).To(&player)
		content, _ := json.Marshal(player)

		w.Write([]byte("data: "))
		w.Write(content)
		w.Write([]byte("\n\n"))
		if err := w.Flush(); err != nil {
			if err := s.DeletePlayer(user.ID, playerSubscribe.RoomID); err != nil {
				log.Error(err)
			}
		}

	Loop:
		for {
			select {

			case content, opened := <-(*channel):
				w.Write([]byte("data: "))
				w.Write(content)
				w.Write([]byte("\n\n"))
				err := w.Flush()

				if (err != nil) || (!opened) {
					break Loop
				}

			case <-timeout:
				w.Write([]byte("\n\n"))
				err := w.Flush()
				if err != nil {
					break Loop
				}
			}
		}

		if err := s.DeletePlayer(user.ID, playerSubscribe.RoomID); err != nil {
			log.Error(err)
		}
	}

	s.BroadcastRoomStatus(user.RoomID)
	return closure, nil
}

func (s *Service) DeletePlayer(playerID, roomID uint) error {

	if err := s.UnsubscribeChannel(playerID, roomID); err != nil {
		return err
	}

	return s.repo.DeleteById(&Users{}, playerID)
}

func (s *Service) UnsubscribeChannel(playerID, roomID uint) error {
	if err := s.broadcaster.DeleteSubscriber(roomID, playerID); err != nil {
		return err
	}
	return s.BroadcastRoomStatus(roomID)
}

func (s *Service) BroadcastRoomStatus(roomId uint) error {

	message, sErr := s.GetRoomStatus(roomId)
	if sErr != nil {
		return sErr
	}

	content, err := json.Marshal(message)

	if err != nil {
		return err
	}

	err = s.broadcaster.SendMessage(roomId, content)
	return err
}

func (s *Service) GetRoomStatus(roomId uint) (*State, error) {

	room := new(Room)
	if err := s.repo.GetOneById(room, roomId); err != nil {
		log.Errorf("Not able to found room with Id: %d. Error %s ", roomId, err)
		return nil, err
	}

	users, err := s.repo.GetPlayersFromRoom(roomId)
	if err != nil {
		return nil, err
	}

	// In case the room is not in show status, all status are -1
	if !room.Show {
		for i := 0; i < len(users); i++ {
			if users[i].Status >= 0 {
				users[i].Status = -1
			}
		}
	}
	state := s.convertUsersToState(users)

	return state, nil
}

func (s *Service) convertUsersToState(users []Users) *State {

	players := []Players{}
	for _, user := range users {
		player := Players{
			ID:     user.ID,
			Name:   user.Name,
			Status: user.Status,
		}
		players = append(players, player)
	}

	state := State{
		Players: players,
	}

	return &state
}
