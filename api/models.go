package main

type Users struct {
	ID     uint `gorm:"primarykey"`
	Room   uint `gorm:"index"`
	Name   string
	Status int8
}
