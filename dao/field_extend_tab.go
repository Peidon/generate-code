package repo

import (
	"git.garena.com/shopee/seller-server/cs/user/internal/constant"
	"git.garena.com/shopee/seller-server/cs/user/internal/model/dbmodel"
	"git.garena.com/shopee/seller-server/cs/user/internal/repo/db"
	"github.com/jinzhu/gorm"
)

type UserFieldExtendTabRepository interface {
	NewTransaction(tx *gorm.DB, userId int64) UserFieldExtendTabRepository
	ChooseShardTab(userId int64, useSlave ...bool) UserFieldExtendTabRepository
	Add(*dbmodel.UserFieldExtendTabTab) error
	Update(t *dbmodel.UserFieldExtendTabTab) error
	UpdateByUid(uid int64, t *dbmodel.UserFieldExtendTabTab) error
	UpdateByUidWithMap(uid int64, data map[string]interface{}) error
	Save(t *dbmodel.UserFieldExtendTabTab) error
	Delete(IDs []int64) error
	GetByUserFieldExtendTabId(ID int64) (*dbmodel.UserFieldExtendTabTab, error)
	GetBaseInfoByID(ID int64) (*dbmodel.BaseUserFieldExtendTab, error)
	ExistByID(ID int64) (bool, error)
	GetAll(cursor int64, size int64) ([]*dbmodel.UserFieldExtendTabTab, error)
}

type UserFieldExtendTabRepo struct {
	selectDB *gorm.DB
	tabName  string
}

var UserFieldExtendTabRepo = new(UserFieldExtendTabRepo)

func (r *UserFieldExtendTabRepo) NewTransaction(tx *gorm.DB, userId int64) UserFieldExtendTabRepository {
	tab := &dbmodel.UserFieldExtendTabTab{}
	return &UserFieldExtendTabRepo{selectDB: tx, tabName: db.GenShardTabName(tab.TableName(), userId)}
}

func (r *UserFieldExtendTabRepo) ChooseShardTab(userId int64, useSlave ...bool) UserFieldExtendTabRepository {
	tab := &dbmodel.UserFieldExtendTabTab{}
	dbName := db.GenShardDBName(userId)
	if len(useSlave) > 0 && useSlave[0] {
		return &UserFieldExtendTabRepo{selectDB: db.ShardSlaveDbs[dbName], tabName: db.GenShardTabName(tab.TableName(), userId)}
	}
	return &UserFieldExtendTabRepo{selectDB: db.ShardDbs[dbName], tabName: db.GenShardTabName(tab.TableName(), userId)}
}

func (r *UserFieldExtendTabRepo) Add(t *dbmodel.UserFieldExtendTabTab) error {
	return r.selectDB.Table(r.tabName).Create(t).Error
}

func (r *UserFieldExtendTabRepo) Update(t *dbmodel.UserFieldExtendTabTab) error {
	return r.selectDB.Table(r.tabName).Model(t).Updates(t).Error
}

func (r *UserFieldExtendTabRepo) UpdateByUid(uid int64, t *dbmodel.UserFieldExtendTabTab) error {
	return r.selectDB.Table(r.tabName).
		Where("cs_user_id = ?", uid).
		Updates(t).Error
}
func (r *UserFieldExtendTabRepo) UpdateByUidWithMap(uid int64, data map[string]interface{}) error {
	return r.selectDB.Table(r.tabName).
		Where("cs_user_id = ?", uid).
		Updates(data).Error
}
func (r *UserFieldExtendTabRepo) Save(t *dbmodel.UserFieldExtendTabTab) error {
	return r.selectDB.Table(r.tabName).Model(t).Save(t).Error
}

func (r *UserFieldExtendTabRepo) Delete(IDs []int64) error {
	data := map[string]interface{}{
		"status_flag": constant.StatusFlagDelete,
	}
	return r.selectDB.Table(r.tabName).Model(&dbmodel.UserFieldExtendTabTab{}).Where("cs_user_id IN(?)", IDs).Scopes(ItemIsExist).Updates(data).Error
}
func (r *UserFieldExtendTabRepo) GetByUserFieldExtendTabId(ID int64) (*dbmodel.UserFieldExtendTabTab, error) {
	var res dbmodel.UserFieldExtendTabTab
	err := r.selectDB.Table(r.tabName).Where("cs_user_id = ?", ID).Scopes(ItemIsExist).Find(&res).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &res, nil
}
func (r *UserFieldExtendTabRepo) GetBaseInfoByID(ID int64) (*dbmodel.BaseUserFieldExtendTab, error) {
	var res dbmodel.BaseUserFieldExtendTab
	err := r.selectDB.Table(r.tabName).Where("cs_user_id = ? ", ID).Find(&res).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &res, nil
}

func (r *UserFieldExtendTabRepo) ExistByID(ID int64) (bool, error) {
	var res dbmodel.UserFieldExtendTabTab
	err := r.selectDB.Table(r.tabName).Select("cs_user_id").Where("cs_user_id = ?", ID).Scopes(ItemIsExist).First(&res).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}
	if res.UserFieldExtendTabId > 0 {
		return true, nil
	}

	return false, err
}

func (r *UserFieldExtendTabRepo) GetAll(cursor int64, size int64) ([]*dbmodel.UserFieldExtendTabTab, error) {
	var res []*dbmodel.UserFieldExtendTabTab
	err := r.selectDB.Table(r.tabName).Where("id > ?", cursor).Order("id asc").Limit(size).Scopes(ItemIsExist).Find(&res).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return res, nil
}
