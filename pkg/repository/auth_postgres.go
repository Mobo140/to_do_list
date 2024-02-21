package repository

import (
	"fmt"

	todo "github.com/Mobo140/projects/to_do_list"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable) //знак $ - плейс холдеры, в которые будут подставлены значения которые мы передадим в качестве аргумента функции для выполнения запроса к базе данных. В конце функции добавлен returning id - возвращает id записи, который будет возвращать id новой записи после операции insert.

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password) //row хранит в себе информацию о возвращаемой строке из базы. В нашем случае запрос из базы возвращает одну строку со значением поля Id.
	if err := row.Scan(&id); err != nil {                                //используя Scan() записываем это значение в переменную id обязательно передавая её по ссылке.
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err

}
