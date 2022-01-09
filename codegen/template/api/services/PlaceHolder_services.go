package services

import (
	"go-ca-codegen/codegen/template/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=${GOPACKAGE}_mock/${GOFILE}.mock.go -package=${GOPACKAGE}_mock

type PlaceHolderRepository interface {
	GetByID(id uint, expand ...string) (*models.PlaceHolder, error)
	GetAll(config GetAllConfig) (data []*models.PlaceHolder, count uint, err error)
	Create(data *models.PlaceHolder) (*models.PlaceHolder, error)
	Update(id uint, data *models.PlaceHolder) (*models.PlaceHolder, error)
	SoftDelete(id uint) (*models.PlaceHolder, error)
	HardDelete(id uint) (*models.PlaceHolder, error)
	Restore(id uint) (*models.PlaceHolder, error)
}

type PlaceHolderService interface {
	GetByID(id uint, expand ...string) (*models.PlaceHolder, error)
	GetAll(config GetAllConfig) (data []*models.PlaceHolder, count uint, err error)
	Create(data *models.PlaceHolder) (*models.PlaceHolder, error)
	Update(id uint, data *models.PlaceHolder) (*models.PlaceHolder, error)
	SoftDelete(id uint) (*models.PlaceHolder, error)
	HardDelete(id uint) (*models.PlaceHolder, error)
	Restore(id uint) (*models.PlaceHolder, error)
}

type placeHolderServiceImpl struct {
	repo PlaceHolderRepository
	PlaceHolderService
}

func NewPlaceHolderService(repository PlaceHolderRepository) PlaceHolderService {
	res := &placeHolderServiceImpl{}
	res.repo = repository
	return res
}

func (c *placeHolderServiceImpl) GetByID(id uint, expand ...string) (*models.PlaceHolder, error) {
	return c.repo.GetByID(id, expand...)
}

func (c *placeHolderServiceImpl) GetAll(config GetAllConfig) ([]*models.PlaceHolder, uint, error) {
	return c.repo.GetAll(config)
}

func (c *placeHolderServiceImpl) Create(data *models.PlaceHolder) (*models.PlaceHolder, error) {
	return c.repo.Create(data)
}

func (c *placeHolderServiceImpl) Update(id uint, data *models.PlaceHolder) (*models.PlaceHolder, error) {
	return c.repo.Update(id, data)
}

func (c *placeHolderServiceImpl) SoftDelete(id uint) (*models.PlaceHolder, error) {
	return c.repo.SoftDelete(id)
}

func (c *placeHolderServiceImpl) HardDelete(id uint) (*models.PlaceHolder, error) {
	return c.repo.HardDelete(id)
}

func (c *placeHolderServiceImpl) Restore(id uint) (*models.PlaceHolder, error) {
	return c.repo.Restore(id)
}
