package main

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}

type AppSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *AppSuite) SetupTest() {
	app := NewDAppMemory()
	s.tester = rollmelette.NewTester(app)
}

func (s *AppSuite) TestItCreateAuctionAndFinishAuctionWithoutPartialSellingAndPayingAllBidder() {
	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
	bidder01 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	bidder02 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	bidder03 := common.HexToAddress("0x0000000000000000000000000000000000000003")
	bidder04 := common.HexToAddress("0x0000000000000000000000000000000000000004")
	bidder05 := common.HexToAddress("0x0000000000000000000000000000000000000005")

	appAddressResult := s.tester.RelayAppAddress(common.HexToAddress("0xdadadadadadadadadadadadadadadadadadadada"))
	s.Nil(appAddressResult.Err)

	createUserInput := []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"creator","username":"vitalik"}}`, creator))
	expectedOutput := fmt.Sprintf(`user created - {"id":3,"role":"creator","username":"vitalik","address":"0x0000000000000000000000000000000000000007","created_at":%d}`, time.Now().Unix())
	result := s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000009"}}`)
	expectedOutput = fmt.Sprintf(`contract created - {"id":1,"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000009","created_at":%d}`, time.Now().Unix())
	result = s.tester.Advance(admin, createContractInput)
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createAuctionInput := []byte(fmt.Sprintf(`{"path":"createAuction","payload":{"max_interest_rate":"10","expires_at":%d,"debt_issued":%d}}`, time.Now().Add(5*time.Second).Unix(), 2020))
	result = s.tester.Advance(creator, createAuctionInput)
	expectedOutput = fmt.Sprintf(`auction created - {"id":1,"creator":"0x0000000000000000000000000000000000000007","debt_issued":"2020","max_interest_rate":"10","state":"ongoing","expires_at":%d,"created_at":%d}`, time.Now().Add(5*time.Second).Unix(), time.Now().Unix())
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createBidInput := []byte(`{"path": "createBid", "payload": {"interest_rate":"9"}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), bidder01, big.NewInt(600), createBidInput)
	expectedOutput = fmt.Sprintf(`bid created - {"id":1,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000001","amount":"600","interest_rate":"9","state":"pending","created_at":%d}`, time.Now().Unix())
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"8"}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), bidder02, big.NewInt(520), createBidInput)
	expectedOutput = fmt.Sprintf(`bid created - {"id":2,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000002","amount":"520","interest_rate":"8","state":"pending","created_at":%d}`, time.Now().Unix())
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"4"}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), bidder03, big.NewInt(200), createBidInput)
	expectedOutput = fmt.Sprintf(`bid created - {"id":3,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000003","amount":"200","interest_rate":"4","state":"pending","created_at":%d}`, time.Now().Unix())
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"6"}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), bidder04, big.NewInt(300), createBidInput)
	expectedOutput = fmt.Sprintf(`bid created - {"id":4,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000004","amount":"300","interest_rate":"6","state":"pending","created_at":%d}`, time.Now().Unix())
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"4"}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), bidder05, big.NewInt(400), createBidInput)
	expectedOutput = fmt.Sprintf(`bid created - {"id":5,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000005","amount":"400","interest_rate":"4","state":"pending","created_at":%d}`, time.Now().Unix())
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	time.Sleep(5 * time.Second)

	finishAuctionInput := []byte(`{"path":"finishAuction"}`)
	result = s.tester.Advance(admin, finishAuctionInput)
	expectedOutput = fmt.Sprintf(`auction finished - {"id":1,"creator":"0x0000000000000000000000000000000000000007","debt_issued":"2020","max_interest_rate":"10","state":"finished","bids":[{"id":3,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000003","amount":"200","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},{"id":5,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000005","amount":"400","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},{"id":4,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000004","amount":"300","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},{"id":2,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000002","amount":"520","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},{"id":1,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000001","amount":"600","interest_rate":"9","state":"accepted","created_at":%d,"updated_at":%d}],"expires_at":%d,"created_at":%d,"updated_at":%d}`, time.Now().Add(-5 * time.Second).Unix(), time.Now().Unix(), time.Now().Add(-5 * time.Second).Unix(), time.Now().Unix(), time.Now().Add(-5 * time.Second).Unix(), time.Now().Unix(),time.Now().Add(-5 * time.Second).Unix(), time.Now().Unix(),time.Now().Add(-5 * time.Second).Unix(), time.Now().Unix(), time.Now().Unix(), time.Now().Add(-5 * time.Second).Unix(), time.Now().Unix())
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	creatorWithdrawInput := []byte(`{"path":"withdraw"}`)
	expectedNoticePayload := `withdrawn STABLECOIN of 1919 from 0x0000000000000000000000000000000000000007 with voucher index: 1`
	expectedWithdrawVoucherPayload := make([]byte, 0, 4+32+32)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, creator[:]...)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(1919).FillBytes(make([]byte, 32))...)
	withdrawResult := s.tester.Advance(creator, creatorWithdrawInput)
	s.Len(withdrawResult.Notices, 1)
	s.Len(withdrawResult.Vouchers, 1)
	s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)
	s.Equal(common.HexToAddress("0x0000000000000000000000000000000000000009"), withdrawResult.Vouchers[0].Destination)
	s.Equal(expectedNoticePayload, string(withdrawResult.Notices[0].Payload))

	tribesProfitWithdrawInput := []byte(`{"path":"withdrawApp"}`)
	expectedNoticePayload = `withdrawn STABLECOIN of 101 from 0x0142f501EE21f4446009C3505c51d0043feC5c68 with voucher index: 1`
	expectedWithdrawVoucherPayload = make([]byte, 0, 4+32+32)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, admin[:]...)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(101).FillBytes(make([]byte, 32))...)
	withdrawResult = s.tester.Advance(admin, tribesProfitWithdrawInput)
	s.Len(withdrawResult.Notices, 1)
	s.Len(withdrawResult.Vouchers, 1)
	s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)
	s.Equal(common.HexToAddress("0x0000000000000000000000000000000000000009"), withdrawResult.Vouchers[0].Destination)
	s.Equal(expectedNoticePayload, string(withdrawResult.Notices[0].Payload))
}
