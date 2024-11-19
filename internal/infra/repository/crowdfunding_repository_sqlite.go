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
	err := r.Db.Model(&entity.Crowdfunding{}).Create(map[string]interface{}{
		"creator":           input.Creator.String(),
		"debt_issued":       input.DebtIssued.Hex(),
		"max_interest_rate": input.MaxInterestRate.Hex(),
		"state":             input.State,
		"expires_at":        input.ExpiresAt,
		"created_at":        input.CreatedAt,
		"updated_at":        input.UpdatedAt,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create crowdfunding: %w", err)
	}
	return r.FindCrowdfundingById(input.Id)
}

func (r *CrowdfundingRepositorySqlite) FindAllCrowdfundings() ([]*entity.Crowdfunding, error) {
	var results []map[string]interface{}
	err := r.Db.Raw("SELECT id, creator, debt_issued, max_interest_rate, state, expires_at, created_at, updated_at FROM crowdfundings").Scan(&results).Error
	if err != nil {
		return nil, err
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
			CreatedAt:       result["created_at"].(int64),
			UpdatedAt:       result["updated_at"].(int64),
		}

		var orderResults []map[string]interface{}
		err = r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE crowdfunding_id = ?", crowdfunding.Id).Scan(&orderResults).Error
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
	err := r.Db.Raw("SELECT id, creator, debt_issued, max_interest_rate, state, expires_at, created_at, updated_at FROM crowdfundings WHERE id = ? LIMIT 1", id).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
	}

	crowdfunding := &entity.Crowdfunding{
		Id:              uint(result["id"].(int64)),
		Creator:         common.HexToAddress(result["creator"].(string)),
		DebtIssued:      uint256.MustFromHex(result["debt_issued"].(string)),
		MaxInterestRate: uint256.MustFromHex(result["max_interest_rate"].(string)),
		State:           entity.CrowdfundingState(result["state"].(string)),
		ExpiresAt:       result["expires_at"].(int64),
		CreatedAt:       result["created_at"].(int64),
		UpdatedAt:       result["updated_at"].(int64),
	}

	var orderResults []map[string]interface{}
	err = r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE crowdfunding_id = ?", crowdfunding.Id).Scan(&orderResults).Error
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
	err := r.Db.Raw("SELECT id, creator, debt_issued, max_interest_rate, state, expires_at, created_at, updated_at FROM crowdfundings WHERE creator = ?", creator.String()).Scan(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
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
			CreatedAt:       result["created_at"].(int64),
			UpdatedAt:       result["updated_at"].(int64),
		}

		var orderResults []map[string]interface{}
		err = r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE crowdfunding_id = ?", crowdfunding.Id).Scan(&orderResults).Error
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
		SELECT c.id, c.creator, c.debt_issued, c.max_interest_rate, c.state, c.expires_at, c.created_at, c.updated_at 
		FROM crowdfundings c
		JOIN orders o ON c.id = o.crowdfunding_id
		WHERE o.investor = ?
	`, investor.String()).Scan(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCrowdfundingNotFound
		}
		return nil, err
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
			CreatedAt:       result["created_at"].(int64),
			UpdatedAt:       result["updated_at"].(int64),
		}

		var orderResults []map[string]interface{}
		err = r.Db.Raw("SELECT id, crowdfunding_id, investor, amount, interest_rate, state, created_at, updated_at FROM orders WHERE crowdfunding_id = ?", crowdfunding.Id).Scan(&orderResults).Error
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
