package entity
type UserFieldExtendTab00000000 struct {
Id int64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
SaasId int64 `gorm:"column:saas_id"`
Region string `gorm:"column:region"`
CsUserId int64 `gorm:"column:cs_user_id"`
FieldId int64 `gorm:"column:field_id"`
FieldValue string `gorm:"column:field_value"`
StatusFlag int8 `gorm:"column:status_flag"`
Mtime int32 `gorm:"column:mtime"`
Ctime int32 `gorm:"column:ctime"`
}
func (*UserFieldExtendTab00000000) TableName() string {
return "user_field_extend_tab_00000000"
}