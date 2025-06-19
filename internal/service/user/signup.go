package user

import (
    "errors"
    "vending-system/internal/model/user"
    "vending-system/internal/repository/user"
    "golang.org/x/crypto/bcrypt"
)

func Signup(repo user.Repository, userid string, password string) error {
    // 유효성 검사
    if userid == "" || password == "" {
        return errors.New("필수 필드가 누락되었습니다")
    }
    if !isPasswordValid(password) {
        return errors.New("비밀번호는 숫자와 특수문자를 포함해 8자 이상이어야 합니다")
    }

    // 아이디 중복 검사
    exists, err := repo.ExistsByUserID(userid)
    if err != nil {
        return err
    }
    if exists {
        return errors.New("이미 존재하는 사용자 ID입니다")
    }

    // 비밀번호 해싱
    hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // DB 저장
    user := model.User{
        UserID:   userid,
        Password: string(hashedPw),
    }
    return repo.Create(user)
}

// 비밀번호 유효성 검사 함수
func isPasswordValid(password string) bool {
    // 최소 8자, 숫자 포함 여부 체크
    if len(password) < 8 {
        return false
    }

    hasNumber := false
    hasSpecialChar := false

    for _, ch := range password {
        if '0' <= ch && ch <= '9' {
            hasNumber = true
        }
        if (ch >= 33 && ch <= 47) || (ch >= 58 && ch <= 64) || (ch >= 91 && ch <= 96) || (ch >= 123 && ch <= 126) {
            hasSpecialChar = true
        }
    }

    return hasNumber && hasSpecialChar
}
