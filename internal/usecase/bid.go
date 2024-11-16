package usecase

import (
	"errors"
	"log/slog"
	"tender_management/internal/entity"
)

type BidService struct {
	bidRepo BidRepo
	log     *slog.Logger
}

func NewBidUseCase(bidRepo BidRepo, log *slog.Logger) *BidService {
	return &BidService{bidRepo: bidRepo, log: log}
}

func (b *BidService) SubmitBid(bid entity.Bid) (entity.Bid, error) {

	b.log.Info("started submitting bid", "tender_id", bid.TenderID)
	defer b.log.Info("ended submitting bid", "tender_id", bid.TenderID)

	// Validate bid data
	if bid.Price <= 0 {
		b.log.Error("error submitting bid", "error", errors.New("price must be greater than 0"))
		return entity.Bid{}, errors.New("price must be greater than 0")
	}
	if bid.DeliveryTime <= 0 {
		b.log.Error("error submitting bid", "error", errors.New("delivery time must be greater than 0"))
		return entity.Bid{}, errors.New("delivery time must be greater than 0")
	}

	// Submit the bid to the repo
	newBid, err := b.bidRepo.SubmitBid(bid)
	if err != nil {
		b.log.Error("error submitting bid", "error", err)
		return entity.Bid{}, err
	}

	return newBid, nil
}

func (b *BidService) GetBids(in entity.ListBidReq) ([]entity.Bid, error) {
	b.log.Info("started getting bids", "tender_id", in.TenderID)
	defer b.log.Info("ended getting bids", "tender_id", in.TenderID)

	// Fetch all bids for the given tender
	bids, err := b.bidRepo.GetBids(in)
	if err != nil {
		b.log.Error("error getting bids", "error", err)
		return nil, err
	}
	return bids, nil
}
