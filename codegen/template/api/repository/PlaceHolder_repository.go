package repository

import (
	"errors"
	"go-ca-codegen/codegen/template/api/repository/querybuilder"
	"go-ca-codegen/codegen/template/api/services"
	"go-ca-codegen/codegen/template/models"
	"time"

	"gorm.io/gorm"
)

type placeHolderRepositoryImpl struct {
	db *gorm.DB
	services.PlaceHolderRepository
}

func NewPlaceHolderRepository(db *gorm.DB) services.PlaceHolderRepository {
	res := &placeHolderRepositoryImpl{}
	res.PlaceHolderRepository = NewPlaceHolderRepository(db)
	return res
}

type PlaceHolderRepositoryImpl struct {
	db        *gorm.DB
	companyID uint
	cache     map[uint]*models.PlaceHolder
	now       func() time.Time
}

func (m *PlaceHolderRepositoryImpl) GetByID(id uint, expand ...string) (*models.PlaceHolder, error) {
	if cache, ok := m.cache[id]; ok && cache != nil && len(expand) == 0 {
		return cache, nil
	}
	data := &models.PlaceHolder{}
	db := m.db.Unscoped()
	db, err := querybuilder.BuildExpandQuery(&models.PlaceHolder{}, expand, db, func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	})
	if err != nil {
		return nil, err
	}
	if err := db.Unscoped().Where("id = ?", id).First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

type GetAllPlaceHolderBaseQueryBuildFunc func(db *gorm.DB) (*gorm.DB, error)

func GetAllPlaceHolderBase(config services.GetAllConfig, db *gorm.DB, companyId uint, queryBuildFunc GetAllPlaceHolderBaseQueryBuildFunc) ([]*models.PlaceHolder, uint, error) {
	var limit int = GetAllMaxLimit
	var offset int = 0
	var allCount int64
	var (
		err   error
		model []*models.PlaceHolder = []*models.PlaceHolder{}
		q     *gorm.DB              = db.Model(&models.PlaceHolder{})
	)
	if config.Limit > 0 {
		limit = int(config.Limit)
	}
	if config.Offset > 0 {
		offset = int(config.Offset)
	}
	if config.IncludeDeleted {
		q = q.Unscoped()
	}
	if config.OnlyDeleted {
		q = q.Unscoped().Where("deleted_at is not null")
	}
	q, err = querybuilder.BuildQueryQuery(&models.PlaceHolder{}, config.Query, q)
	if err != nil {
		return nil, 0, err
	}
	q, err = querybuilder.BuildOrderQuery(&models.PlaceHolder{}, config.Order, q)
	if err != nil {
		return nil, 0, err
	}
	q, err = querybuilder.BuildExpandQuery(&models.PlaceHolder{}, config.Expand, q, func(db *gorm.DB) *gorm.DB {
		return db.Where("company_id = ?", companyId).Unscoped()
	})
	if err != nil {
		return nil, 0, err
	}
	if queryBuildFunc != nil {
		q, err = queryBuildFunc(q)
		if err != nil {
			return nil, 0, err
		}
	}
	// 最大10000件ずつでちょっとずつ読み込む
	load := func() (bool, error) {
		var sub []models.PlaceHolder
		subLimit := GetAllSubLimit
		if limit <= subLimit {
			subLimit = limit + 1
		}
		if err := q.Offset(offset).Limit(subLimit).Find(&sub).Error; err != nil {
			return false, err
		}
		var size int
		offset += size
		limit -= size
		return size < subLimit || limit < 0, nil
	}
	for {
		shouldEnd, err := load()
		if err != nil {
			return nil, 0, err
		}
		if shouldEnd {
			break
		}
	}

	if (config.Limit > 0 && uint(len(model)) > config.Limit) || config.Offset > 0 {
		if err := q.Model(&models.PlaceHolder{}).Count(&allCount).Error; err != nil {
			return nil, 0, err
		}
	} else {
		allCount = int64(len(model))
	}
	if config.Limit > 0 && uint(len(model)) > config.Limit {
		model = model[:config.Limit]
	}
	if len(model) > GetAllMaxLimit {
		return nil, 0, errors.New("データ数が多すぎるため取得できません")
	}
	return model, uint(allCount), nil
}

func (m *PlaceHolderRepositoryImpl) GetAll(config services.GetAllConfig) ([]*models.PlaceHolder, uint, error) {
	return GetAllPlaceHolderBase(config, m.db, m.companyID, nil)
}

func (m *PlaceHolderRepositoryImpl) Create(data *models.PlaceHolder) (*models.PlaceHolder, error) {
	data = ShallowCopy(data).(*models.PlaceHolder)
	now := m.now()
	data.SetUpdatedAt(now)
	data.SetCreatedAt(now)
	if err := m.db.
		Set("gorm:save_associations", false).
		Set("gorm:association_save_reference", false).
		Create(data).Error; err != nil {
		return nil, err
	}
	data, err := m.GetByID(data.GetID())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *PlaceHolderRepositoryImpl) Update(id uint, data *models.PlaceHolder) (*models.PlaceHolder, error) {
	orgData, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}
	if data.GetID() != orgData.GetID() {
		return nil, errors.New("IDは変更できません")
	}
	if data.GetCreatedAt().UTC().Unix() != orgData.GetCreatedAt().UTC().Unix() {
		return nil, errors.New("作成日時は変更できません")
	}
	if data.GetUpdatedAt().UTC().Unix() != orgData.GetUpdatedAt().UTC().Unix() {
		return nil, errors.New("更新日時は変更できません")
	}
	if data.GetDeletedAt() != orgData.GetDeletedAt() {
		if data.GetDeletedAt() == nil && orgData.GetDeletedAt() != nil {
		} else if data.GetDeletedAt() == nil || orgData.GetDeletedAt() == nil {
			return nil, errors.New("削除日時は変更できません")
		} else if data.GetDeletedAt().UTC().Unix() != orgData.GetDeletedAt().UTC().Unix() {
			return nil, errors.New("削除日時は変更できません")
		}
	}
	data.SetUpdatedAt(m.now())
	if err := m.db.
		Set("gorm:save_associations", false).
		Set("gorm:association_save_reference", false).
		Set("gorm:update_column", false).
		Unscoped().Save(data).Error; err != nil {
		return nil, err
	}
	data, err = m.GetByID(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *PlaceHolderRepositoryImpl) SoftDelete(id uint) (*models.PlaceHolder, error) {
	data, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}
	data.SetDeletedAt(m.now())
	if err := m.db.
		Set("gorm:save_associations", false).
		Set("gorm:association_save_reference", false).
		Unscoped().Save(data).Error; err != nil {
		return nil, err
	}
	data, err = m.GetByID(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *PlaceHolderRepositoryImpl) HardDelete(id uint) (*models.PlaceHolder, error) {
	data, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}
	if !data.IsDeleted() {
		return nil, errors.New("指定のデータは削除されていないため，完全に削除できません")
	}
	if err := m.db.Unscoped().Delete(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (m *PlaceHolderRepositoryImpl) Restore(id uint) (*models.PlaceHolder, error) {
	data, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}
	if err := m.db.Unscoped().Save(data).Error; err != nil {
		return nil, err
	}
	data, err = m.GetByID(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}
