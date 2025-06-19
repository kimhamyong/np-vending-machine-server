package model

// User 모델 구조체
type User struct {
    ID       int    `json:"id"`
    UserID   string `json:"userid"`
    Password string `json:"password"`
}
