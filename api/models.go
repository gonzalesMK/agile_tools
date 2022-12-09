package main

type Users struct {
	ID     uint `gorm:"primarykey"`
	RoomID uint `gorm:"index"`
	Name   string
	Status int8
}
