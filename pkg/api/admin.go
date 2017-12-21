package api

// Admin models the admin user of dock_server in the database
type Admin struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
