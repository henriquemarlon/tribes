package repository

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
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
	// Create crowdfunding entity
	err := r.Db.Create(&input).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create crowdfunding: %w", err)
	}
	return r.FindCrowdfundingById(input.Id)
}

func (r *CrowdfundingRepositorySqlite) FindAllCrowdfundings() ([]*entity.Crowdfunding, error) {
	// Manually map fields for each crowdfunding
	var results []map[string]interface{}
	err := r.Db.Model(&entity.Crowdfunding{}).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load all crowdfundings: %w", err)
	}

	var crowdfundings []*entity.Crowdfunding
	for _, result := range results {
		crowdfunding := &entity.Crowdfunding{
			Id:              uint(result["id"].(int64)),
			Creator:         common.HexToAddress(result["creator"].(string)),
			DebtIssued:      uint256.MustFromHex(result["debt_issued"].(string)),
			MaxInterestRate: uint256.MustFromHex(result["max_interest_rate"].(string)),
			State:           entity.CrowdfundingState(result["state"].(string)),
			ExpiresAt:       result["expires_at"].(int64),
			MaturityAt:      result["maturity_at"].(int64),
			CreatedAt:       result["created_at"].(int64),
			UpdatedAt:       result["updated_at"].(int64),
		}

		// Manually map orders
		var orderResults []map[string]interface{}
		err = r.Db.Model(&entity.Order{}).Where("crowdfunding_id = ?", crowdfunding.Id).Find(&orderResults).Error
		if err != nil {
			return nil, fmt.Errorf("failed to load orders for crowdfunding: %w", err)
		}
		for _, orderResult := range orderResults {
			order := &entity.Order{
				Id:             uint(orderResult["id"].(int64)),
				CrowdfundingId: uint(orderResult["crowdfunding_id"].(int64)),
				Investor:       common.HexToAddress(orderResult["investor"].(string)),
				Amount:         uint256.MustFromHex(orderResult["amount"].(string)),
				InterestRate:   uint256.MustFromHex(orderResult["interest_rate"].(string)),
				State:          entity.OrderState(orderResult["state"].(string)),
				CreatedAt:      orderResult["created_at"].(int64),
				UpdatedAt:      orderResult["updated_at"].(int64),
			}
			crowdfunding.Orders = append(crowdfunding.Orders, order)
		}

		crowdfundings = append(crowdfundings, crowdfunding)
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingById(id uint) (*entity.Crowdfunding, error) {
	var result map[string]interface{}
	err := r.Db.Model(&entity.Crowdfunding{}).Where("id = ?", id).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, fmt.Errorf("failed to find crowdfunding by id: %w", err)
	}

	crowdfunding := &entity.Crowdfunding{
		Id:              uint(result["id"].(int64)),
		Creator:         common.HexToAddress(result["creator"].(string)),
		DebtIssued:      uint256.MustFromHex(result["debt_issued"].(string)),
		MaxInterestRate: uint256.MustFromHex(result["max_interest_rate"].(string)),
		State:           entity.CrowdfundingState(result["state"].(string)),
		ExpiresAt:       result["expires_at"].(int64),
		MaturityAt:      result["maturity_at"].(int64),
		CreatedAt:       result["created_at"].(int64),
		UpdatedAt:       result["updated_at"].(int64),
	}

	var orderResults []map[string]interface{}
	err = r.Db.Model(&entity.Order{}).Where("crowdfunding_id = ?", crowdfunding.Id).Find(&orderResults).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load orders for crowdfunding: %w", err)
	}
	for _, orderResult := range orderResults {
		order := &entity.Order{
			Id:             uint(orderResult["id"].(int64)),
			CrowdfundingId: uint(orderResult["crowdfunding_id"].(int64)),
			Investor:       common.HexToAddress(orderResult["investor"].(string)),
			Amount:         uint256.MustFromHex(orderResult["amount"].(string)),
			InterestRate:   uint256.MustFromHex(orderResult["interest_rate"].(string)),
			State:          entity.OrderState(orderResult["state"].(string)),
			CreatedAt:      orderResult["created_at"].(int64),
			UpdatedAt:      orderResult["updated_at"].(int64),
		}
		crowdfunding.Orders = append(crowdfunding.Orders, order)
	}

	return crowdfunding, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingsByCreator(creator common.Address) ([]*entity.Crowdfunding, error) {
	var results []map[string]interface{}
	err := r.Db.Model(&entity.Crowdfunding{}).Where("creator = ?", creator.String()).Find(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, fmt.Errorf("failed to find crowdfundings by creator: %w", err)
	}

	var crowdfundings []*entity.Crowdfunding
	for _, result := range results {
		crowdfunding := &entity.Crowdfunding{
			Id:              uint(result["id"].(int64)),
			Creator:         common.HexToAddress(result["creator"].(string)),
			DebtIssued:      uint256.MustFromHex(result["debt_issued"].(string)),
			MaxInterestRate: uint256.MustFromHex(result["max_interest_rate"].(string)),
			State:           entity.CrowdfundingState(result["state"].(string)),
			ExpiresAt:       result["expires_at"].(int64),
			MaturityAt:      result["maturity_at"].(int64),
			CreatedAt:       result["created_at"].(int64),
			UpdatedAt:       result["updated_at"].(int64),
		}

		var orderResults []map[string]interface{}
		err = r.Db.Model(&entity.Order{}).Where("crowdfunding_id = ?", crowdfunding.Id).Find(&orderResults).Error
		if err != nil {
			return nil, fmt.Errorf("failed to load orders for crowdfunding: %w", err)
		}
		for _, orderResult := range orderResults {
			order := &entity.Order{
				Id:             uint(orderResult["id"].(int64)),
				CrowdfundingId: uint(orderResult["crowdfunding_id"].(int64)),
				Investor:       common.HexToAddress(orderResult["investor"].(string)),
				Amount:         uint256.MustFromHex(orderResult["amount"].(string)),
				InterestRate:   uint256.MustFromHex(orderResult["interest_rate"].(string)),
				State:          entity.OrderState(orderResult["state"].(string)),
				CreatedAt:      orderResult["created_at"].(int64),
				UpdatedAt:      orderResult["updated_at"].(int64),
			}
			crowdfunding.Orders = append(crowdfunding.Orders, order)
		}

		crowdfundings = append(crowdfundings, crowdfunding)
	}
	return crowdfundings, nil
}

func (r *CrowdfundingRepositorySqlite) FindCrowdfundingsByInvestor(investor common.Address) ([]*entity.Crowdfunding, error) {
	var results []map[string]interface{}
	err := r.Db.Raw(`
		SELECT DISTINCT c.* 
		FROM crowdfundings c
		JOIN orders o ON c.id = o.crowdfunding_id
		WHERE o.investor = ?
	`, investor.String()).Scan(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, fmt.Errorf("failed to find crowdfundings by investor: %w", err)
	}

	var crowdfundings []*entity.Crowdfunding
	for _, result := range results {
		crowdfunding := &entity.Crowdfunding{
			Id:              uint(result["id"].(int64)),
			Creator:         common.HexToAddress(result["creator"].(string)),
			DebtIssued:      uint256.MustFromHex(result["debt_issued"].(string)),
			MaxInterestRate: uint256.MustFromHex(result["max_interest_rate"].(string)),
			State:           entity.CrowdfundingState(result["state"].(string)),
			ExpiresAt:       result["expires_at"].(int64),
			MaturityAt:      result["maturity_at"].(int64),
			CreatedAt:       result["created_at"].(int64),
			UpdatedAt:       result["updated_at"].(int64),
		}

		var orderResults []map[string]interface{}
		err = r.Db.Model(&entity.Order{}).Where("crowdfunding_id = ?", crowdfunding.Id).Find(&orderResults).Error
		if err != nil {
			return nil, fmt.Errorf("failed to load orders for crowdfunding: %w", err)
		}
		for _, orderResult := range orderResults {
			order := &entity.Order{
				Id:             uint(orderResult["id"].(int64)),
				CrowdfundingId: uint(orderResult["crowdfunding_id"].(int64)),
				Investor:       common.HexToAddress(orderResult["investor"].(string)),
				Amount:         uint256.MustFromHex(orderResult["amount"].(string)),
				InterestRate:   uint256.MustFromHex(orderResult["interest_rate"].(string)),
				State:          entity.OrderState(orderResult["state"].(string)),
				CreatedAt:      orderResult["created_at"].(int64),
				UpdatedAt:      orderResult["updated_at"].(int64),
			}
			crowdfunding.Orders = append(crowdfunding.Orders, order)
		}

		crowdfundings = append(crowdfundings, crowdfunding)
	}
	return crowdfundings, nil
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

	if input.DebtIssued != nil {
		crowdfundingJSON["debt_issued"] = input.DebtIssued.Hex()
	}
	if input.MaxInterestRate != nil {
		crowdfundingJSON["max_interest_rate"] = input.MaxInterestRate.Hex()
	}
	if input.State != "" {
		crowdfundingJSON["state"] = input.State
	}
	if input.ExpiresAt != 0 {
		crowdfundingJSON["expires_at"] = input.ExpiresAt
	}
	if input.MaturityAt != 0 {
		crowdfundingJSON["maturity_at"] = input.MaturityAt
	}
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

	res := r.Db.Save(&crowdfunding)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update crowdfunding: %w", res.Error)
	}
	return &crowdfunding, nil
}

func (r *CrowdfundingRepositorySqlite) DeleteCrowdfunding(id uint) error {
	res := r.Db.Delete(&entity.Crowdfunding{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete crowdfunding: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrCrowdfundingNotFound
	}
	return nil
}

func (r *CrowdfundingRepositorySqlite) CloseCrowdfunding(crowdfundingId uint) ([]*entity.Crowdfunding, error) {
	return nil, nil
}

func (r *CrowdfundingRepositorySqlite) SettleCrowdfunding(crowdfundingId uint) ([]*entity.Crowdfunding, error) {
	return nil, nil
}
