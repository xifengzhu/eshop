package models

import (
	"golang.org/x/crypto/bcrypt"
)

type AdminUser struct {
	BaseModel

	WxappId  string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Email    string `gorm:"type: varchar(100); not null" json:"email"`
	Status   int    `gorm:"type: tinyint;" json:"status"`
	Password string `gorm:"type: varchar(120); not null" json:"password"`
	Role     string `gorm:"type: varchar(50);" json:"role"`
}

func (AdminUser) TableName() string {
	return "admin_user"
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (admin AdminUser) Authenticate(password string) bool {
	return CompareHashAndPassword(password, admin.Password)
}

func CompareHashAndPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (admin *AdminUser) GetAdminUserByEmail(email string) (err error) {
	err = db.Where("email = ?", email).First(&admin).Error
	return
}

func (admin *AdminUser) Create() (err error) {
	var password string
	password, err = HashPassword(admin.Password)
	if err != nil {
		return
	}
	admin.Password = password
	err = db.Create(&admin).Error
	return
}

func GetAdminUserById(ID int) (admin AdminUser, err error) {
	err = db.First(&admin, ID).Error
	return
}
