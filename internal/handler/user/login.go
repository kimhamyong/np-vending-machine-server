package user

import "database/sql"

func HandleLogin(req map[string]interface{}, db *sql.DB) []byte {
	// 실제로 DB를 사용하지 않는 경우, DB 인자는 그냥 추가만 해둡니다.
	return []byte(`{"success": true, "message": "Login successful (dummy)"}`)
}
