package models

import (
	config "github.com/xifengzhu/eshop/initializers"
)

func GetPermissionByKey(key string) (permits []string) {
	permissions, _ := Enforcer.GetImplicitPermissionsForUser(key)
	if len(permissions) > 0 {
		for _, permissions := range permissions {
			permits = append(permits, permissions[1])
		}
	}
	return
}

func AllPermissions() (permissions []string) {
	return config.Permissions
}
