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

// ////////////// User ///////////////////

func (s *AppSuite) TestItCreateUser() {
	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
	input := []byte(`{"path":"createUser","payload":{"address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","role":"admin","username":"vitalik"}}`)
	expectedOutput := fmt.Sprintf(`user created - {"id":3,"role":"admin","username":"vitalik","address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","created_at":%d}`, time.Now().Unix())
	result := s.tester.Advance(admin, input)
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))
}

////////////// Auction //////////////////

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
	expectedOutput = fmt.Sprintf(`created auction - {"id":1,"creator":"0x0000000000000000000000000000000000000007","debt_issued":"2020","max_interest_rate":"10","state":"ongoing","expires_at":%d,"created_at":%d}`, time.Now().Add(5*time.Second).Unix(), time.Now().Unix())
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
	expectedOutput = `finished auction with - id: 1, required amount: 2020 and max interest rate: 10`
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

}
