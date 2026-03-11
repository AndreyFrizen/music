package modeluser

// User represents a user in the system
type User struct {
	ID                int64  `db:"id" redis:"id"`
	Username          string `db:"username" redis:"username"`
	Password          string
	EncryptedPassword string `db:"password"`
	Email             string `db:"email" redis:"email"`
}
