package main

import (
	"bufio"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
	"github.com/valyala/fasthttp"
)

type Controller struct {
	s ServiceInterface
}

type PlayerResponse struct {
	ID uint `json:"id"`
}

type PlayerRequest struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	RoomID uint   `json:"room"`
	Status int8   `json:"status"`
}

type PlayerSubscribe struct {
	Name   string `query:"name"`
	RoomID uint   `query:"room"`
}

//go:generate mockgen -source=$GOFILE -destination=mock_service_test.go -package=main
type ServiceInterface interface {
	UpsertPlayer(playerRequest *PlayerRequest) (*PlayerResponse, error)
	Subscribe(playerSubscribe *PlayerSubscribe) (func(w *bufio.Writer), error)
}

func (s *Controller) UpsertPlayer(c *fiber.Ctx) error {
	request := new(PlayerRequest)

	if err := c.BodyParser(request); err != nil {
		log.Error(err)
		return err
	}

	playerResponse, err := s.s.UpsertPlayer(request)
	if err != nil {
		return err
	}

	return c.JSON(playerResponse)
}

func (s *Controller) UpdateState(c *fiber.Ctx) error {
	// Config for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	conf := new(PlayerSubscribe)
	if err := c.QueryParser(conf); err != nil {
		log.Error(err)
		return err
	}

	updater, dbErr := s.s.Subscribe(conf)
	if dbErr != nil {
		log.Error(dbErr)
		return dbErr
	}

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(updater))

	return nil
}
