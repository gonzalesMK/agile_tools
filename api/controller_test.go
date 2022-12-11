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

	request := PlayerSubscribeMock{}.AllFields()

	service.
		EXPECT().
		Subscribe(gomock.Eq(request)).
		Return(func(b *bufio.Writer) { b.Write([]byte("AB")); b.Flush() }, nil)

	app := fiber.New()
	app.Get("/", controller.UpdateState)

	req := httptest.NewRequest("GET", "/?room=12&name=Santhia%20Witchy", nil)

	resp, err := app.Test(req, 1000)

	assert.Nil(t, err)
	bytes, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "AB", string(bytes))
}

func TestUpdateRoomController(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockServiceInterface(ctrl)
	controller := Controller{
		s: service,
	}

	request := RoomRequestMock{}.AllFields()
	response := RoomResponseMocks{}.AllFields()

	service.
		EXPECT().
		UpdateRoom(gomock.Eq(request)).
		Return(response, nil)

	app := fiber.New()
	app.Post("/", controller.UpdateRoom)

	content, err := json.Marshal(request)
	assert.Nil(t, err)
	fmt.Println(string(content))
	req := httptest.NewRequest("POST", "/", bytes.NewReader(content))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, 1000)

	assert.Nil(t, err)

	bytes, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "{\"id\":12}", string(bytes))
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
	app.Post("/", controller.UpdatePlayer)

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
