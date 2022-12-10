package main

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) Save(model interface{}) error {

	return r.db.Save(model).Error

}

func (r *Repository) UpdateFieldById(id uint, content interface{}) error {

	return r.db.Where("id", id).Updates(content).Error

}
func (r *Repository) DeleteById(model interface{}, id uint) error {

	return r.db.Delete(model, id).Error

}

func (r *Repository) GetPlayersFromRoom(roomId uint) ([]Users, error) {

	var users []Users

	err := r.db.Where(&Users{RoomID: roomId}).Find(&users).Error

	return users, err
}
