package main

type ServiceMocks struct {
}

var ServiceMock = ServiceMocks{}

func (s *ServiceMocks) NewServiceWithChannel(repo RepoInterface, channel *chan []byte) Service {

	return Service{
		repo: repo,
		broadcaster: &Broadcaster{
			rooms: map[uint]*RoomMap{
				12: {
					channels: map[uint]*chan []byte{
						1: channel,
					},
				},
			},
		},
	}
}

type UserMocks struct {
}

func (UserMocks) AllFields() *Users {

	return &Users{
		Name:   "Santhia Witchy",
		ID:     123,
		Room:   Room{ID: 12},
		RoomID: 12,
		Status: -1,
	}
}

type PlayerRequestMocks struct {
}

func (PlayerRequestMocks) AllFields() *PlayerRequest {

	return &PlayerRequest{
		Name:   "Santhia Witchy",
		ID:     123,
		RoomID: 12,
		Status: -1,
	}
}

type PlayerSubscribeMock struct {
}

func (PlayerSubscribeMock) AllFields() *PlayerSubscribe {

	return &PlayerSubscribe{
		Name:   "Santhia Witchy",
		RoomID: 12,
	}
}

type PlayerResponseMocks struct {
}

func (PlayerResponseMocks) AllFields() *PlayerResponse {

	return &PlayerResponse{
		ID: 123,
	}
}

type RoomMock struct {
}

func (RoomMock) AllFields() *Room {

	return &Room{
		ID:   12,
		Show: true,
	}
}

type RoomRequestMock struct {
}

func (RoomRequestMock) AllFields() *RoomRequest {

	return &RoomRequest{
		ID:   12,
		Show: true,
	}
}

type RoomResponseMocks struct {
}

func (RoomResponseMocks) AllFields() *RoomResponse {

	return &RoomResponse{
		ID: 12,
	}
}
