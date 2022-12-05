package main

import (
	"bufio"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
	"github.com/valyala/fasthttp"
)

type Controller struct {
	s ServiceInterface
}

type PlayerConfig struct {
	Name   string `query:"name"`
	RoomId uint   `query:"room"`
}

//go:generate mockgen -source=$GOFILE -destination=mock_service_test.go -package=main
type ServiceInterface interface {
	CreatePlayer(name string, status int8, roomId uint) (*Players, *chan []byte, error)
	DeletePlayer(player *Players, room uint) error
}

func (s *Controller) UpdateState(c *fiber.Ctx) error {
	// Config for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	conf := new(PlayerConfig)
	if err := c.QueryParser(conf); err != nil {
		log.Error(err)
		return err
	}

	player, channel, dbErr := s.s.CreatePlayer(conf.Name, -1, conf.RoomId)
	if dbErr != nil {
		log.Error(dbErr)
		return dbErr
	}

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		timeout := make(chan bool, 1)
		go func() {
			for {
				timeout <- true
				time.Sleep(time.Second)
			}
		}()
	Loop:
		for {
			select {

			case content := <-(*channel):
				fmt.Println("Sent data: ", player.ID)
				w.Write([]byte("data: "))
				w.Write(content)
				w.Write([]byte("\n\n"))
				err := w.Flush()
				fmt.Println("Flushed data", player.ID)

				if err != nil {
					break Loop
				}

			case <-timeout:
				w.Write([]byte("\n\n"))
				err := w.Flush()
				fmt.Println("HERE ", player.ID)
				if err != nil {
					fmt.Println("Error ", err)
					break Loop
				}
			}
		}

		log.Info("Closed a connection", player.ID)
		if err := s.s.DeletePlayer(player, conf.RoomId); err != nil {
			log.Error(err)
		}
	}))

	return nil
}
