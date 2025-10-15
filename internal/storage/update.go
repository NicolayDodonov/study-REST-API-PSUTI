package storage

import (
	"database/sql"
	"study-REST-API-PSUTI/internal/model"
)

func (s *Storage) UpdateUserData(data *model.UserInfo) error {
	_, err := s.db.Exec(`
		UPDATE users 
		SET height=$1,
		    weight=$2,
		    age=$3,
		WHERE id=$4`,
		data.Height, data.Weight, data.Age,
		data.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CheckUserByID(uid string) (bool, error) {
	var User struct {
		uid string `db:"id"`
	}
	err := s.db.Get(&User, `SELECT * FROM users WHERE id=$1`, uid)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
