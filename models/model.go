package models

import (
	"errors"
	"time"
)

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (m *Model) GetID() uint {
	return m.ID
}

func (m *Model) SetID(id uint) error {
	if id == 0 {
		return errors.New("ID は正の整数である必要があります")
	}
	m.ID = id
	return nil
}

func (m *Model) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *Model) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m *Model) GetDeletedAt() *time.Time {
	return m.DeletedAt
}

func (m *Model) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *Model) SetDeletedAt(t time.Time) {
	t = t.UTC()
	m.DeletedAt = &t
}

func (m *Model) UnsetDeleted() {
	m.DeletedAt = nil
}

func (m *Model) SetCreatedAt(t time.Time) {
	m.CreatedAt = t.UTC()
}

func (m *Model) SetUpdatedAt(t time.Time) {
	m.UpdatedAt = t.UTC()
}
