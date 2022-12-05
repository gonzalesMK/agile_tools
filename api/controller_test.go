package main

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUpdateState(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMockServiceInterface(ctrl)
	controller := Controller{
		s: service,
	}
	channel := make(chan []byte, 1)
	player := &Players{
		ID:     1,
		Name:   "Santhia Witchy",
		Status: -3,
	}

	service.
		EXPECT().
		CreatePlayer(gomock.Eq("developer"), gomock.Eq(int8(-1)), gomock.Eq(uint(123))).
		Return(player, &channel, nil)
	service.
		EXPECT().
		DeletePlayer(gomock.Any(), gomock.Eq(uint(123))).
		Return(nil)

	app := fiber.New()
	app.Get("/", controller.UpdateState)

	req := httptest.NewRequest("GET", "/?room=123&name=developer", nil)

	// Test
	channel <- []byte{'A', 'B'}

	close(channel)

	resp, err := app.Test(req, 1000)

	assert.Nil(t, err)
	bytes, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "data: AB\n\n", string(bytes))
}
