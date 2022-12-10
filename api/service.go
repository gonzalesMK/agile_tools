package main

import (
	"bufio"
	"encoding/json"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/ulule/deepcopier"
)

type Service struct {
	repo        RepoInterface
	mu          sync.RWMutex
	broadcaster *Broadcaster
}

//go:generate mockgen -source=$GOFILE -destination=mock_repository_test.go -package=main
type RepoInterface interface {
	Save(interface{}) error
	UpdateFieldById(id uint, content interface{}) error
	DeleteById(model interface{}, id uint) error
	GetPlayersFromRoom(uint) ([]Users, error)
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

	err := s.repo.Save(&user)
	if err != nil {
		return nil, err
	}

	channel := s.broadcaster.AddSubscriber(playerSubscribe.RoomID, user.ID)

	closure := func(w *bufio.Writer) {
		timeout := make(chan bool, 1)
		go func() {
			for {
				timeout <- true
				time.Sleep(time.Second)
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

			case content := <-(*channel):
				w.Write([]byte("data: "))
				w.Write(content)
				w.Write([]byte("\n\n"))
				err := w.Flush()

				if err != nil {
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

	if err := s.repo.DeleteById(&Users{}, playerID); err != nil {
		return err
	}
	s.broadcaster.DeleteSubscriber(roomID, playerID)
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

	s.broadcaster.SendMessage(roomId, content)
	return nil
}

func (s *Service) GetRoomStatus(roomId uint) (*State, error) {

	users, err := s.repo.GetPlayersFromRoom(roomId)

	if err != nil {
		return nil, err
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
