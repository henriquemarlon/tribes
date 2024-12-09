package repository

import (
	"context"
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
	"gorm.io/gorm"
)

type CrowdfundingRepositorySqlite struct {
	Db *gorm.DB
}

func NewCrowdfundingRepositorySqlite(db *gorm.DB) *CrowdfundingRepositorySqlite {
	return &CrowdfundingRepositorySqlite{Db: db}
}

func (r *CrowdfundingRepositorySqlite) CreateCrowdfunding(ctx context.Context, input *entity.Crowdfunding) (*entity.Crowdfunding, error) {
	if err := r.Db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create crowdfunding: %w", err)
	}
	return input, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingById(ctx context.Context, id uint) (*entity.Crowdfunding, error) {
	var crowdfunding entity.Crowdfunding
	if err := r.Db.WithContext(ctx).
		Preload("Orders").
		First(&crowdfunding, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, fmt.Errorf("failed to find crowdfunding by id: %w", err)
	}
	return &crowdfunding, nil
}

func (r *CrowdfundingRepositorySqlite) FindAllCrowdfundings(ctx context.Context) ([]*entity.Crowdfunding, error) {
	var crowdfundings []*entity.Crowdfunding
	if err := r.Db.WithContext(ctx).
		Preload("Orders").
		Find(&crowdfundings).Error; err != nil {
		return nil, fmt.Errorf("failed to find all crowdfundings: %w", err)
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingsByInvestor(ctx context.Context, investor custom_type.Address) ([]*entity.Crowdfunding, error) {
	var crowdfundings []*entity.Crowdfunding
	if err := r.Db.WithContext(ctx).
		Joins("JOIN orders ON orders.crowdfunding_id = crowdfundings.id").
		Where("orders.investor = ?", investor).
		Preload("Orders").
		Find(&crowdfundings).Error; err != nil {
		return nil, fmt.Errorf("failed to find crowdfundings by investor: %w", err)
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingsByCreator(ctx context.Context, creator custom_type.Address) ([]*entity.Crowdfunding, error) {
	var crowdfundings []*entity.Crowdfunding
	if err := r.Db.WithContext(ctx).
		Where("creator = ?", creator).
		Preload("Orders").
		Find(&crowdfundings).Error; err != nil {
		return nil, fmt.Errorf("failed to find crowdfundings by creator: %w", err)
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) UpdateCrowdfunding(ctx context.Context, input *entity.Crowdfunding) (*entity.Crowdfunding, error) {
	if err := r.Db.WithContext(ctx).Updates(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to update crowdfunding: %w", err)
	}
	crowdfunding, err := r.FindCrowdfundingById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	return crowdfunding, nil
}

func (r *CrowdfundingRepositorySqlite) DeleteCrowdfunding(ctx context.Context, id uint) error {
	res := r.Db.WithContext(ctx).Delete(&entity.Crowdfunding{}, id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete crowdfunding: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrCrowdfundingNotFound
	}
	return nil
}
