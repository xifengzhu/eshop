package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AdminUser struct {
	BaseModel

	Email    string `gorm:"type: varchar(100); not null" json:"email"`
	Status   string `gorm:"type: varchar(10);" json:"status"`
	Password string `gorm:"type: varchar(120); not null" json:"-"`

	Roles       []string `gorm:"-" json:"roles"`
	Permissions []string `gorm:"-" json:"permissions"`
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

func (admin *AdminUser) BeforeSave(scope *gorm.Scope) (err error) {
	if pw, err := HashPassword(admin.Password); err == nil {
		scope.SetColumn("Password", pw)
	}
	return
}

func GetAdminUserById(ID int) (admin AdminUser, err error) {
	err = db.First(&admin, ID).Error
	return
}

func (admin *AdminUser) AuthKey() string {
	return admin.Email
}

func (admin *AdminUser) GetRoles() (roles []string) {
	return GetRolesByKey(admin.AuthKey())
}

func (admin *AdminUser) GetPermissions() (permits []string) {
	return GetPermissionByKey(admin.AuthKey())
}
