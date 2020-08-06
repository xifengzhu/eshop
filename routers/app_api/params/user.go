package params

// Binding from JSON
type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthParams struct {
	Code          string `json:"code" validate:"required"`
	EncryptedData string `json:"encrypted_data" validate:"required"`
	IV            string `json:"iv" validate:"required"`
}

type UserInfo struct {
	Username string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
}
