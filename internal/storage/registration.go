package storage

import (
	"github.com/google/uuid"
	"study-REST-API-PSUTI/internal/model"
)

func (s *Storage) Registration(u *model.UserInfo) (string, error) {
	u.Id = uuid.New().String()
	_, err := s.db.Exec(`
	INSERT INTO users (id, first_name, last_name, user_type, login, password) 
	VALUES ($1, $2, $3, $4, $5, $6)`,
		u.Id, u.FirstName, u.LastName, u.UserType, u.Login, u.Password)
	if err != nil {
		return "", err
	}

	return u.Id, nil
}
