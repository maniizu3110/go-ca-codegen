package models

import (
	"errors"
	"time"
)

type PlaceHolder struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (m *PlaceHolder) GetID() uint {
	return m.ID
}

func (m *PlaceHolder) SetID(id uint) error {
	if id == 0 {
		return errors.New("idが０です")
	}
	m.ID = id
	return nil
}

func (m *PlaceHolder) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *PlaceHolder) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m *PlaceHolder) GetDeletedAt() *time.Time {
	return m.DeletedAt
}

func (m *PlaceHolder) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *PlaceHolder) SetDeletedAt(t time.Time) {
	t = t.UTC()
	m.DeletedAt = &t
}

func (m *PlaceHolder) UnsetDeleted() {
	m.DeletedAt = nil
}

func (m *PlaceHolder) SetCreatedAt(t time.Time) {
	m.CreatedAt = t.UTC()
}

func (m *PlaceHolder) SetUpdatedAt(t time.Time) {
	m.UpdatedAt = t.UTC()
}
