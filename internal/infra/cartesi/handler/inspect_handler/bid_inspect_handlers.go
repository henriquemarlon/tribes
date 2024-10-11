package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Mugen-Builders/devolt/internal/domain/entity"
	"github.com/Mugen-Builders/devolt/internal/usecase/bid_usecase"
	"github.com/Mugen-Builders/devolt/pkg/router"
	"github.com/rollmelette/rollmelette"
)

type BidInspectHandlers struct {
	BidRepository entity.BidRepository
}

func NewBidInspectHandlers(bidRepository entity.BidRepository) *BidInspectHandlers {
	return &BidInspectHandlers{
		BidRepository: bidRepository,
	}
}

func (h *BidInspectHandlers) FindBidByIdHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	id, err := strconv.Atoi(router.PathValue(ctx, "id"))
	if err != nil {
		return fmt.Errorf("failed to parse id into int: %v", router.PathValue(ctx, "id"))
	}
	findBidById := bid_usecase.NewFindBidByIdUseCase(h.BidRepository)
	res, err := findBidById.Execute(&bid_usecase.FindBidByIdInputDTO{
		Id: uint(id),
	})
	if err != nil {
		return fmt.Errorf("failed to find bid: %w", err)
	}
	bid, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal bid: %w", err)
	}
	env.Report(bid)
	return nil
}

func (h *BidInspectHandlers) FindBisdByAuctionIdHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	id, err := strconv.Atoi(router.PathValue(ctx, "id"))
	if err != nil {
		return fmt.Errorf("failed to parse id into int: %v", router.PathValue(ctx, "id"))
	}
	findBidsByAuctionId := bid_usecase.NewFindBidsByAuctionIdUseCase(h.BidRepository)
	res, err := findBidsByAuctionId.Execute(&bid_usecase.FindBidsByAuctionIdInputDTO{
		AuctionId: uint(id),
	})
	if err != nil {
		return fmt.Errorf("failed to find bids by auction id: %v", err)
	}
	bids, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal bids: %w", err)
	}
	env.Report(bids)
	return nil
}

func (h *BidInspectHandlers) FindAllBidsHandler(env rollmelette.EnvInspector, ctx context.Context) error {
	findAllBids := bid_usecase.NewFindAllBidsUseCase(h.BidRepository)
	res, err := findAllBids.Execute()
	if err != nil {
		return fmt.Errorf("failed to find all bids: %w", err)
	}
	allBids, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all bids: %w", err)
	}
	env.Report(allBids)
	return nil
}
