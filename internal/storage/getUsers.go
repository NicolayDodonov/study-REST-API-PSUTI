package storage

import (
	"database/sql"
	"study-REST-API-PSUTI/internal/model"
)

func (s *Storage) GetUsers(p model.GetUserParams) ([]model.UserInfo, error) {
	var users []model.UserInfo
	err := s.db.Select(&users, `
		SELECT * FROM users 
		         WHERE 
		             height > $1 AND
		             weight > $2 AND
		             age > $3`,
		p.LargeHeight,
		p.LargeWeight,
		p.LargeAge)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}
