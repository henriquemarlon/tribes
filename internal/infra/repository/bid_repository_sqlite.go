package db

import (
	"fmt"

	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"gorm.io/gorm"
)

type BidRepositorySqlite struct {
	Db *gorm.DB
}

func NewBidRepositorySqlite(db *gorm.DB) *BidRepositorySqlite {
	return &BidRepositorySqlite{
		Db: db,
	}
}

func (r *BidRepositorySqlite) CreateBid(input *entity.Bid) (*entity.Bid, error) {
	err := r.Db.Create(input).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}
	return input, nil
}

func (r *BidRepositorySqlite) FindBidsByState(auctionId uint, state string) ([]*entity.Bid, error) {
	var bids []*entity.Bid
	err := r.Db.Where("auction_id = ? AND state = ?", auctionId, state).Find(&bids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find bids by state: %w", err)
	}
	return bids, nil
}

func (r *BidRepositorySqlite) FindBidById(id uint) (*entity.Bid, error) {
	var bid entity.Bid
	err := r.Db.First(&bid, "bid_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find bid by ID: %w", err)
	}
	return &bid, nil
}

func (r *BidRepositorySqlite) FindBidsByAuctionId(id uint) ([]*entity.Bid, error) {
	var bids []*entity.Bid
	err := r.Db.Where("auction_id = ?", id).Find(&bids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find bids by auction ID: %w", err)
	}
	return bids, nil
}

func (r *BidRepositorySqlite) FindAllBids() ([]*entity.Bid, error) {
	var bids []*entity.Bid
	err := r.Db.Find(&bids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all bids: %w", err)
	}
	return bids, nil
}

func (r *BidRepositorySqlite) UpdateBid(input *entity.Bid) (*entity.Bid, error) {
	var bid entity.Bid
	err := r.Db.First(&bid, "id = ?", input.Id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrBidNotFound
		}
		return nil, err
	}

	bid.AuctionId = input.AuctionId
	bid.Bidder = input.Bidder
	bid.Amount = input.Amount
	bid.InterestRate = input.InterestRate
	bid.State = input.State
	bid.UpdatedAt = input.UpdatedAt

	res := r.Db.Save(&bid)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to update bid: %w", res.Error)
	}
	return input, nil
}

func (r *BidRepositorySqlite) DeleteBid(id uint) error {
	res := r.Db.Delete(&entity.Bid{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return entity.ErrBidNotFound
	}
	return nil
}
