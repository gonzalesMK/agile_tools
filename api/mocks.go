package main

type ServiceMocks struct {
}

var ServiceMock = ServiceMocks{}

func (s *ServiceMocks) NewServiceWithChannel(repo RepoInterface, channel *chan []byte) Service {

	return Service{
		repo: repo,
		broadcaster: &Broadcaster{
			rooms: map[uint]*Room{
				123: {
					channels: map[uint]*chan []byte{
						1: channel,
					},
				},
			},
		},
	}
}
