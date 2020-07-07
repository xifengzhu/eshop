package models

type Role struct {
	BaseModel

	Name    string `gorm:"type: varchar(50); not null; unique_index" json:"name"`
	Creator string `json:"creator"`
	Remark  string `json:"remark"`

	Permissions []string `gorm:"-" json:"permissions"`
}

func (role *Role) AuthKey() string {
	return role.Name
}

func GetRolesByKey(key string) (roles []string) {
	roles, _ = Enforcer.GetImplicitRolesForUser(key)
	return
}

func (r *Role) GetPermissions() (permits []string) {
	return GetPermissionByKey(r.AuthKey())
}
