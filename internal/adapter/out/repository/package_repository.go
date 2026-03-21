package repository

import (
	"context"
	"phone-tokens/internal/model"

	"gorm.io/gorm"
)

type PackageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return PackageRepository{db: db}
}

func (r *PackageRepository) AddPackage(ctx context.Context, pkg *model.Package) (err error) {
	return r.db.WithContext(ctx).Save(&pkg).Error
}

func (r *PackageRepository) DeletePackage(ctx context.Context, pkg *model.Package) (err error) {
	return r.db.WithContext(ctx).Delete(&pkg).Error
}

func (r *PackageRepository) GetPackageByID(ctx context.Context, id string) (*model.Package, error) {
	var pkg model.Package
	err := r.db.WithContext(ctx).Find(&pkg, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *PackageRepository) GetPackageByName(ctx context.Context, name string) (*model.Package, error) {
	var pkg model.Package
	err := r.db.WithContext(ctx).Find(&pkg, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *PackageRepository) GetPackages(ctx context.Context) ([]model.Package, error) {
	packages := []model.Package{}
	err := r.db.WithContext(ctx).Find(&packages).Error
	return packages, err
}
