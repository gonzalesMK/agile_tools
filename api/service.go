package main

import (
	"encoding/json"
	"sync"

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
	DeleteById(*Users) error
	GetPlayersFromRoom(uint) ([]Users, error)
}

func (s *Service) CreatePlayer(name string, status int8, roomID uint) (*Players, *chan []byte, error) {

	user := Users{
		Name:   name,
		Status: status,
		Room:   roomID,
	}

	err := s.repo.Save(&user)
	channel := s.broadcaster.AddSubscriber(roomID, user.ID)

	var player Players
	deepcopier.Copy(&user).To(&player)

	s.BroadcastRoomStatus(roomID)
	return &player, channel, err
}

func (s *Service) DeletePlayer(player *Players, roomID uint) error {
	var user Users
	deepcopier.Copy(player).To(&user)
	if err := s.repo.DeleteById(&user); err != nil {
		return err
	}
	s.broadcaster.DeleteSubscriber(roomID, player.ID)
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
