package main

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) Save(model interface{}) error {

	return r.db.Create(model).Error

}
func (r *Repository) DeleteById(model *Users) error {

	return r.db.Delete(&Users{}, model.ID).Error

}

func (r *Repository) GetPlayersFromRoom(roomId uint) ([]Users, error) {

	var users []Users

	err := r.db.Where(&Users{Room: roomId}).Find(&users).Error

	return users, err
}
