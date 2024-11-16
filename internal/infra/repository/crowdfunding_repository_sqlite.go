package repository

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"gorm.io/gorm"
)

type CrowdfundingRepositorySqlite struct {
	Db *gorm.DB
}

func NewCrowdfundingRepositorySqlite(db *gorm.DB) *CrowdfundingRepositorySqlite {
	return &CrowdfundingRepositorySqlite{
		Db: db,
	}
}

func (r *CrowdfundingRepositorySqlite) CreateCrowdfunding(input *entity.Crowdfunding) (*entity.Crowdfunding, error) {
	err := r.Db.Create(&input).Error
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingsByCreator(creator common.Address) ([]*entity.Crowdfunding, error) {
	var crowdfundings []*entity.Crowdfunding
	err := r.Db.Preload("Orders").Where("creator = ?", creator).Find(&crowdfundings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingsByInvestor(investor common.Address) ([]*entity.Crowdfunding, error) {
	return nil, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingById(id uint) (*entity.Crowdfunding, error) {
	var crowdfunding entity.Crowdfunding
	err := r.Db.Preload("Orders").First(&crowdfunding, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
	}
	return &crowdfunding, nil
}

func (r *CrowdfundingRepositorySqlite) FindAllCrowdfundings() ([]*entity.Crowdfunding, error) {
	var crowdfundings []*entity.Crowdfunding
	err := r.Db.Preload("Orders").Find(&crowdfundings).Error
	if err != nil {
		return nil, err
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) CloseCrowdfunding(crowdfundingId uint) ([]*entity.Order, error) {
	return nil, nil
}

func (r *CrowdfundingRepositorySqlite) SettleCrowdfunding(crowdfundingId uint) ([]*entity.Order, error) {
	return nil, nil
}

func (r *CrowdfundingRepositorySqlite) UpdateCrowdfunding(input *entity.Crowdfunding) (*entity.Crowdfunding, error) {
	var crowdfundingJSON map[string]interface{}
	err := r.Db.Where("id = ?", input.Id).First(&crowdfundingJSON).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
	}

	crowdfundingJSON["max_interest_rate"] = input.MaxInterestRate
	crowdfundingJSON["state"] = input.State
	crowdfundingJSON["expires_at"] = input.ExpiresAt
	crowdfundingJSON["updated_at"] = input.UpdatedAt

	crowdfundingBytes, err := json.Marshal(crowdfundingJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal crowdfunding: %w", err)
	}

	var crowdfunding entity.Crowdfunding
	if err = json.Unmarshal(crowdfundingBytes, &crowdfunding); err != nil {
		return nil, err
	}
	if err = crowdfunding.Validate(); err != nil {
		return nil, err
	}

	res := r.Db.Save(crowdfundingJSON)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update crowdfunding: %w", res.Error)
	}
	return &crowdfunding, nil
}

func (r *CrowdfundingRepositorySqlite) DeleteCrowdfunding(id uint) error {
	err := r.Db.Delete(&entity.Crowdfunding{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}
