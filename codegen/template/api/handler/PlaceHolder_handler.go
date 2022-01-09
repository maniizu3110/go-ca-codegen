package handler

import (
	"errors"
	"go-ca-codegen/codegen/template/api/repository"
	"go-ca-codegen/codegen/template/api/services"
	"go-ca-codegen/codegen/template/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func AssignPlaceHolderHandlers(g *echo.Group) {
	g = g.Group("", func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			db := c.Get("Tx").(*gorm.DB)
			r := repository.NewPlaceHolderRepository(db)
			s := services.NewPlaceHolderService(r)
			c.Set("Service", s)
			return handler(c)
		}
	})
	g.POST("", CreatePlaceHolderHandler)
	g.PUT("/:id", UpdatePlaceHolderHandler)
	g.DELETE("/:id", DeletePlaceHolderHandler)
	g.PUT("/:id/restore", RestorePlaceHolderHandler)
	g.GET("/:id", GetPlaceHolderByIDHandler)
	g.GET("", GetPlaceHolderListHandler)
}

type CreatePlaceHolderHandlerBaseCallbackFunc func(services.PlaceHolderService, *models.PlaceHolder) (*models.PlaceHolder, error)

func CreatePlaceHolderHandlerBase(c echo.Context, params interface{}, callback CreatePlaceHolderHandlerBaseCallbackFunc) (err error) {
	service := c.Get("Service").(services.PlaceHolderService)

	data := &models.PlaceHolder{}
	if err != nil {
		return err
	}
	if err = c.Bind(data); err != nil {
		return err
	}
	if params != nil {
		if err = c.Bind(params); err != nil {
			return err
		}
	}
	if err = c.Validate(data); err != nil {
		return err
	}
	m, err := callback(service, data)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, m)
}

func CreatePlaceHolderHandler(c echo.Context) (err error) {
	return CreatePlaceHolderHandlerBase(c, nil, func(service services.PlaceHolderService, data *models.PlaceHolder) (*models.PlaceHolder, error) {
		return service.Create(data)
	})
}

type UpdatePlaceHolderHandlerBaseCallbackFunc func(services.PlaceHolderService, uint, *models.PlaceHolder) (*models.PlaceHolder, error)

func UpdatePlaceHolderHandlerBase(c echo.Context, params interface{}, callback UpdatePlaceHolderHandlerBaseCallbackFunc) (err error) {
	service := c.Get("Service").(services.PlaceHolderService)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return errors.New(err.Error())
	}
	data, err := service.GetByID(uint(id))
	if err != nil {
		return err
	}
	if err = c.Bind(data); err != nil {
		return errors.New(err.Error())
	}
	if params != nil {
		if err = c.Bind(params); err != nil {
			return errors.New(err.Error())
		}
	}
	if err = c.Validate(data); err != nil {
		return err
	}
	m, err := callback(service, uint(id), data)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, m)
}

func UpdatePlaceHolderHandler(c echo.Context) (err error) {
	return UpdatePlaceHolderHandlerBase(c, nil, func(service services.PlaceHolderService, id uint, data *models.PlaceHolder) (*models.PlaceHolder, error) {
		return service.Update(uint(id), data)
	})
}

type DeletePlaceHolderHandlerBaseCallbackFunc func(services.PlaceHolderService, uint) (*models.PlaceHolder, error)

func DeletePlaceHolderHandlerBase(c echo.Context, params interface{}, callback DeletePlaceHolderHandlerBaseCallbackFunc) (err error) {
	service := c.Get("Service").(services.PlaceHolderService)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return errors.New(err.Error())
	}
	if params != nil {
		if err = c.Bind(params); err != nil {
			return errors.New(err.Error())
		}
	}
	data, err := callback(service, uint(id))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func DeletePlaceHolderHandler(c echo.Context) (err error) {
	var param struct {
		HardDelete bool
	}
	return DeletePlaceHolderHandlerBase(c, &param, func(service services.PlaceHolderService, id uint) (*models.PlaceHolder, error) {
		if param.HardDelete {
			return service.HardDelete(id)
		}
		return service.SoftDelete(id)
	})
}

type RestorePlaceHolderHandlerBaseCallbackFunc func(services.PlaceHolderService, uint) (*models.PlaceHolder, error)

func RestorePlaceHolderHandlerBase(c echo.Context, params interface{}, callback RestorePlaceHolderHandlerBaseCallbackFunc) (err error) {
	service := c.Get("Service").(services.PlaceHolderService)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return errors.New(err.Error())
	}
	if params != nil {
		if err = c.Bind(params); err != nil {
			return errors.New(err.Error())
		}
	}
	m, err := callback(service, uint(id))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, m)
}

func RestorePlaceHolderHandler(c echo.Context) (err error) {
	return RestorePlaceHolderHandlerBase(c, nil, func(service services.PlaceHolderService, id uint) (*models.PlaceHolder, error) {
		return service.Restore(id)
	})
}

type GetPlaceHolderByIDHandlerBaseCallbackFunc func(services.PlaceHolderService, uint) (*models.PlaceHolder, error)

func GetPlaceHolderByIDHandlerBase(c echo.Context, params interface{}, callback GetPlaceHolderByIDHandlerBaseCallbackFunc) (err error) {
	service := c.Get("Service").(services.PlaceHolderService)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return errors.New(err.Error())
	}
	if params != nil {
		if err = c.Bind(params); err != nil {
			return errors.New(err.Error())
		}
	}
	m, err := callback(service, uint(id))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, m)
}

func GetPlaceHolderByIDHandler(c echo.Context) (err error) {
	var param struct {
		Expand []string
	}
	return GetPlaceHolderByIDHandlerBase(c, &param, func(service services.PlaceHolderService, id uint) (*models.PlaceHolder, error) {
		return service.GetByID(id, param.Expand...)
	})
}

type GetPlaceHolderListResponse struct {
	AllCount uint `json:",omitempty"`
	Offset   uint `json:",omitempty"`
	Limit    uint `json:",omitempty"`
	Data     []*models.PlaceHolder
}

type GetPlaceHolderListHandlerBaseCallbackFunc func(services.PlaceHolderService) (*GetPlaceHolderListResponse, error)

func GetPlaceHolderListHandlerBase(c echo.Context, params interface{}, callback GetPlaceHolderListHandlerBaseCallbackFunc) (err error) {
	service := c.Get("Service").(services.PlaceHolderService)
	if params != nil {
		if err = c.Bind(params); err != nil {
			return errors.New(err.Error())
		}
	}
	response, err := callback(service)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response)
}

func GetPlaceHolderListHandler(c echo.Context) (err error) {
	var param struct {
		services.GetAllConfig
	}

	return GetPlaceHolderListHandlerBase(c, &param, func(service services.PlaceHolderService) (*GetPlaceHolderListResponse, error) {
		m, count, err := service.GetAll(param.GetAllConfig)
		if err != nil {
			return nil, err
		}
		return &GetPlaceHolderListResponse{
			AllCount: count,
			Limit:    param.Limit,
			Offset:   param.Offset,
			Data:     m,
		}, nil
	})
}
