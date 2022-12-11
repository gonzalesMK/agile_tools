package main

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetOneById(model interface{}, id uint) error {

	return r.db.First(model, id).Error

}

func (r *Repository) Save(model interface{}) error {

	return r.db.Save(model).Error
}

func (r *Repository) UpdateFieldById(id uint, content interface{}) error {

	return r.db.Where("id", id).Select("*").Updates(content).Error

}
func (r *Repository) DeleteById(model interface{}, id uint) error {

	return r.db.Delete(model, id).Error

}

func (r *Repository) GetPlayersFromRoom(roomId uint) ([]Users, error) {

	var users []Users

	err := r.db.Where(&Users{RoomID: roomId}).Find(&users).Error

	return users, err
}

func (r *Repository) ClearPlayerStatusInRoom(roomID uint, status int8) error {

	err := r.db.Model(&Users{}).Where("room_id = ?", roomID).Update("status", status).Error

	return err
}
