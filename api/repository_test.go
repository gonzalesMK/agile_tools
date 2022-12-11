package main

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type e2eTestSuite struct {
	suite.Suite
	db *gorm.DB
	t  *testing.T
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{t: t})
}

func (s *e2eTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)

	s.db = db

	db.AutoMigrate(&Users{}, &Room{})
	s.Require().NoError(err)

}

func (s *e2eTestSuite) TestSaveWorks() {

	repo := Repository{
		db: s.db,
	}

	user := UserMocks{}.AllFields()

	err := repo.Save(user)
	s.Assertions.Nil(err)

	savedUsers := new([]Users)
	result := s.db.Preload("Room").Find(savedUsers)
	s.Assertions.Equal(int64(1), result.RowsAffected)

	savedUser := (*savedUsers)[0]
	s.Assertions.Equal(savedUser.RoomID, savedUser.Room.ID)

}

func (s *e2eTestSuite) TestUpdateWorks() {

	repo := Repository{
		db: s.db,
	}

	user := UserMocks{}.AllFields()

	err := repo.Save(user)
	s.Assertions.Nil(err)

	user.Name = "New Name"
	err = repo.UpdateFieldById(user.ID, user)
	s.Assertions.Nil(err)

	savedUsers := new([]Users)
	result := s.db.Find(savedUsers)
	s.Assertions.Equal(int64(1), result.RowsAffected)

	s.Assertions.Equal("New Name", (*savedUsers)[0].Name)

}

func (s *e2eTestSuite) TestUpdateDoesNotCreateModel() {

	repo := Repository{
		db: s.db,
	}

	user := UserMocks{}.AllFields()

	err := repo.UpdateFieldById(1, user)
	s.Assertions.Nil(err)

	savedUsers := new([]Users)
	result := s.db.Find(savedUsers)
	s.Assertions.Equal(int64(0), result.RowsAffected)

}

func (s *e2eTestSuite) TestClearPlayerStatusInRoom() {

	repo := Repository{
		db: s.db,
	}

	user := UserMocks{}.AllFields()

	err := repo.Save(user)
	s.Assertions.Nil(err)

	err = repo.ClearPlayerStatusInRoom(user.RoomID, 10)
	s.Assertions.Nil(err)

	savedUsers := new([]Users)
	result := s.db.Find(savedUsers)
	s.Assertions.Equal(int64(1), result.RowsAffected)
	savedUser := (*savedUsers)[0]
	s.Assertions.Equal(int8(10), savedUser.Status)

}
