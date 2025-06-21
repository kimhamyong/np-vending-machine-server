package user

import (
    "errors"
    "vending-system/internal/repository/user"
	"fmt"

    "golang.org/x/crypto/bcrypt"
)

func Login(repo user.Repository, userid string, password string) error {
    hashedPassword, err := repo.FindByUserID(userid)
    if err != nil {
        return errors.New("사용자를 찾을 수 없습니다")
    }

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        return errors.New("비밀번호가 일치하지 않습니다")
    }

    return nil
}