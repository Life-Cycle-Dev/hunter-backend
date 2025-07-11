package entity

type RoleToPermission struct {
	RoleId       string `json:"role_id" gorm:"type:varchar(255);primary_key"`
	PermissionId string `json:"permission_id" gorm:"type:varchar(255);primary_key"`
}
