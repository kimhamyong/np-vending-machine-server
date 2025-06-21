package user

import (
	"errors"
	"vending-system/internal/repository/user"

	"golang.org/x/crypto/bcrypt"
)

func DeleteAccount(repo user.Repository, userid, password string) error {
	hashedPw, err := repo.FindByUserID(userid)
	if err != nil {
		return errors.New("사용자를 찾을 수 없습니다")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(password)); err != nil {
		return errors.New("비밀번호가 일치하지 않습니다")
	}

	return repo.DeleteByUserID(userid)
}
