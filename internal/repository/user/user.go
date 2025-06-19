package user

import (
    "database/sql"
    "vending-system/internal/model/user"
)

type Repository interface {
    Create(user model.User) error
    ExistsByUserID(userid string) (bool, error)
}

type repoImpl struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
    return &repoImpl{db: db}
}

func (r *repoImpl) Create(user model.User) error {
    _, err := r.db.Exec(`INSERT INTO users (user_id, password) VALUES (?, ?)`, user.UserID, user.Password)
    return err
}

func (r *repoImpl) ExistsByUserID(userid string) (bool, error) {
    var count int
    err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE user_id = ?`, userid).Scan(&count)
    return count > 0, err
}
