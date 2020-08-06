package params

type CaptchaParam struct {
	CaptchaID    string `json:"captcha_id" validate:"required"`
	CaptchaValue string `json:"captcha_value" validate:"required"`
}

type LoginParams struct {
	CaptchaParam
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password"  validate:"required"`
}

type ForgetPasswordParams struct {
	CaptchaParam
	Email string `json:"email"  validate:"required,email"`
}

type ResetPasswordParams struct {
	Password string `json:"password"  validate:"required,gte=6,lt=12"`
}
