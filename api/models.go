package main

type Users struct {
	ID     uint `gorm:"primarykey"`
	RoomID uint `gorm:"index"`
	Room   Room `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name   string
	Status int8
}

type Room struct {
	ID   uint `gorm:"primarykey"`
	Show bool
}
