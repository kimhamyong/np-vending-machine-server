package user

import (
    "database/sql"
    "vending-system/internal/model/user"
    "errors"
)

type Repository interface {
    Create(user model.User) error
    ExistsByUserID(userid string) (bool, error)
    FindByUserID(userid string) (string, error)
    UpdatePassword(userid string, newPassword string) error
    DeleteByUserID(userid string) error
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

func (r *repoImpl) FindByUserID(userid string) (string, error) {
    var password string
    err := r.db.QueryRow("SELECT password FROM users WHERE user_id = ?", userid).Scan(&password)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", errors.New("존재하지 않는 사용자입니다")
        }
        return "", err
    }
    return password, nil
}

func (r *repoImpl) UpdatePassword(userid string, hashedPassword string) error {
	_, err := r.db.Exec(`UPDATE users SET password = ? WHERE user_id = ?`, hashedPassword, userid)
	return err
}

func (r *repoImpl) DeleteByUserID(userid string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE user_id = ?`, userid)
	return err
}