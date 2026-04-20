package repository

import (
	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/model"
)

type FileObjectRepository struct {
	db *gorm.DB
}

func NewFileObjectRepository(db *gorm.DB) *FileObjectRepository {
	return &FileObjectRepository{db: db}
}

func (r *FileObjectRepository) GetByID(id uint) (*model.FileObject, error) {
	var item model.FileObject
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *FileObjectRepository) FindByObjectKey(objectKey string) (*model.FileObject, error) {
	var item model.FileObject
	err := r.db.Where("object_key = ?", objectKey).First(&item).Error
	return &item, err
}

func (r *FileObjectRepository) Create(item *model.FileObject) error {
	return r.db.Create(item).Error
}

func (r *FileObjectRepository) Update(item *model.FileObject) error {
	return r.db.Save(item).Error
}
