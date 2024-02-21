package todo

// добавляем к полям структуры теги чтобы корректно принимать и выводить данные на http-запросах
type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
