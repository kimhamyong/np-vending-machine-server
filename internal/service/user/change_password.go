package user

import (
	"errors"
	"fmt"
	"strings"
	"vending-system/internal/repository/user"

	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(repo user.Repository, userID, oldPassword, newPassword string) error {
	hashedPassword, err := repo.FindByUserID(userID)
	if err != nil {
		return errors.New("사용자를 찾을 수 없습니다")
	}

	oldPassword = strings.TrimSpace(oldPassword)

	fmt.Printf("oldPassword(raw): [%x]\n", oldPassword)
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(oldPassword)); err != nil {
		return errors.New("기존 비밀번호가 일치하지 않습니다")
	}

	// 새 비밀번호 해싱
	newHashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("비밀번호 암호화 실패")
	}

	return repo.UpdatePassword(userID, string(newHashed))
}
