package storage

import (
	"database/sql"
	"study-REST-API-PSUTI/internal/model"
)

func (s *Storage) GetUsers() ([]model.UserInfo, error) {
	var users []model.UserInfo
	err := s.db.Select(&users, `
		select * from users
		where `)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}
