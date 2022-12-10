package main

type Users struct {
	ID     uint  `gorm:"primarykey"`
	RoomID uint  `gorm:"index"`
	Room   Rooms `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name   string
	Status int8
}

type Rooms struct {
	ID   uint `gorm:"primarykey"`
	Show bool
}
