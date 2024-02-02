package repo

type OrgUser struct {
	ID    int64  `gorm:"primaryKey"`
	Uid   uint64 `gorm:"index:uid"`
	OrgId int    `gorm:"index:uid"`
	Role  int    `gorm:"index:uid"`
}

func (receiver OrgUser) TableName() string {
	return GetTablePrefix() + "org_user"
}
