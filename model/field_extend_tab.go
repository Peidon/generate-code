package model
type UserFieldExtendTabModel struct {
Id int64 `json:"id"`
SaasId int64 `json:"saas_id"`
Region string `json:"region"`
CsUserId int64 `json:"cs_user_id"`
FieldId int64 `json:"field_id"`
FieldValue string `json:"field_value"`
StatusFlag int8 `json:"status_flag"`
Mtime int32 `json:"mtime"`
Ctime int32 `json:"ctime"`
}
