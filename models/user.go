package models

import (
	"github.com/lempiy/echo_api/types"
	"github.com/lempiy/echo_api/utils"
)

type user struct{}
var User *user

func (u *user) Create(user *types.User) error {
	encryptPass := utils.EncryptPassword(user.Password)
	sqlQuery := `INSERT INTO person(username, password, login, age, telephone, created_date)
		VALUES($1,$2,$3,$4,$5,now());`
	err := Database.SingleQuery(sqlQuery, user.Username, encryptPass, user.Login,
		user.Age, user.Telephone)
	return err
}

func (u *user) Read(id int) (*types.User, error) {
	var user types.User
	sqlQuery := `SELECT * FROM person WHERE id=$1;`
	rows := Database.Query(sqlQuery, id)
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&user.ID,&user.Username,&user.Password,&user.Login,&user.Age,&user.Telephone,&user.CreatedDate)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (u *user) ReadByLogin(login string) (*types.User, error) {
	var user types.User
	sqlQuery := `SELECT p.id, p.username, p.password, p.login,
		p.age, p.telephone, p.created_date FROM person p WHERE login=$1;`
	rows := Database.Query(sqlQuery, login)
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&user.ID,&user.Username,&user.Password,&user.Login,&user.Age,&user.Telephone,&user.CreatedDate)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

