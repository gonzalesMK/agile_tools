package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockServiceInterface(ctrl)
	controller := Controller{
		s: service,
	}

	request := &PlayerSubscribe{
		Name:   "developer",
		RoomID: uint(123),
	}

	service.
		EXPECT().
		Subscribe(gomock.Eq(request)).
		Return(func(b *bufio.Writer) { b.Write([]byte("AB")); b.Flush() }, nil)

	app := fiber.New()
	app.Get("/", controller.UpdateState)

	req := httptest.NewRequest("GET", "/?room=123&name=developer", nil)

	resp, err := app.Test(req, 1000)

	assert.Nil(t, err)
	bytes, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "AB", string(bytes))
}

func TestUpsertPlayerController(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockServiceInterface(ctrl)
	controller := Controller{
		s: service,
	}

	request := PlayerRequestMocks{}.AllFields()
	response := PlayerResponseMocks{}.AllFields()

	service.
		EXPECT().
		UpsertPlayer(gomock.Eq(request)).
		Return(response, nil)

	app := fiber.New()
	app.Post("/", controller.UpsertPlayer)

	content, err := json.Marshal(request)
	assert.Nil(t, err)
	fmt.Println(string(content))
	req := httptest.NewRequest("POST", "/", bytes.NewReader(content))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, 1000)

	assert.Nil(t, err)

	bytes, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "{\"id\":123}", string(bytes))
}
