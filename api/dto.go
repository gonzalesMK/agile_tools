package main

type Players struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status int8   `json:"status"`
}

type State struct {
	Players []Players `json:"players"`
}
