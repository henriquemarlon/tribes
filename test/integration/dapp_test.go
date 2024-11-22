package main

import (
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
	"github.com/tribeshq/tribes/cmd/tribes-rollup/root"
	"github.com/tribeshq/tribes/configs"
)

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(DAppSuite))
}

type DAppSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *DAppSuite) SetupTest() {
	db, err := configs.SetupSQlite(":memory:")
	if err != nil {
		slog.Error("Failed to setup in-memory SQLite database", "error", err)
		os.Exit(1)
	}
	ah, err := root.NewAdvanceHandlers(db)
	if err != nil {
		slog.Error("Failed to setup advance handlers", "error", err)
		os.Exit(1)
	}
	ih, err := root.NewInspectHandlers(db)
	if err != nil {
		slog.Error("Failed to setup inspect handlers", "error", err)
		os.Exit(1)
	}
	ms, err := root.NewMiddlewares(db)
	if err != nil {
		slog.Error("Failed to setup middlewares", "error", err)
		os.Exit(1)
	}
	app := root.NewDApp(ah, ih, ms)
	s.tester = rollmelette.NewTester(app)
}

func (s *DAppSuite) TestItCreatedCrowdfundingAndFinishdCrowdfundingWithoutPartialSellingAndPayingAllOrderder() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
	investor01 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	investor02 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	investor03 := common.HexToAddress("0x0000000000000000000000000000000000000003")
	investor04 := common.HexToAddress("0x0000000000000000000000000000000000000004")
	investor05 := common.HexToAddress("0x0000000000000000000000000000000000000005")

	appAddressResult := s.tester.RelayAppAddress(common.HexToAddress("0xdadadadadadadadadadadadadadadadadadadada"))
	s.Nil(appAddressResult.Err)

	createUserInput := []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"creator"}}`, creator))
	expectedOutput := fmt.Sprintf(`user created - {"id":2,"role":"creator","address":"0x0000000000000000000000000000000000000007","investment_limit":"0","debt_issuance_limit":"15000000","created_at":%d}`, time.Now().Unix())
	result := s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"qualified_investor"}}`, investor01))
	expectedOutput = fmt.Sprintf(`user created - {"id":3,"role":"qualified_investor","address":"0x0000000000000000000000000000000000000001","investment_limit":"115792089237316195423570985008687907853269984665640564039457584007913129639935","debt_issuance_limit":"0","created_at":%d}`, time.Now().Unix())
	result = s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"qualified_investor"}}`, investor02))
	expectedOutput = fmt.Sprintf(`user created - {"id":4,"role":"qualified_investor","address":"0x0000000000000000000000000000000000000002","investment_limit":"115792089237316195423570985008687907853269984665640564039457584007913129639935","debt_issuance_limit":"0","created_at":%d}`, time.Now().Unix())
	result = s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"non_qualified_investor"}}`, investor03))
	expectedOutput = fmt.Sprintf(`user created - {"id":5,"role":"non_qualified_investor","address":"0x0000000000000000000000000000000000000003","investment_limit":"20000","debt_issuance_limit":"0","created_at":%d}`, time.Now().Unix())
	result = s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"non_qualified_investor"}}`, investor04))
	expectedOutput = fmt.Sprintf(`user created - {"id":6,"role":"non_qualified_investor","address":"0x0000000000000000000000000000000000000004","investment_limit":"20000","debt_issuance_limit":"0","created_at":%d}`, time.Now().Unix())
	result = s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	createUserInput = []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"non_qualified_investor"}}`, investor05))
	expectedOutput = fmt.Sprintf(`user created - {"id":7,"role":"non_qualified_investor","address":"0x0000000000000000000000000000000000000005","investment_limit":"20000","debt_issuance_limit":"0","created_at":%d}`, time.Now().Unix())
	result = s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000009"}}`)
	// expectedOutput = fmt.Sprintf(`contract created - {"id":1,"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000009","created_at":%d}`, time.Now().Unix())
	// result = s.tester.Advance(admin, createContractInput)
	// s.Len(result.Notices, 1)
	// s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// createContractInput = []byte(`{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000009"}}`)
	// expectedOutput = fmt.Sprintf(`contract created - {"id":1,"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000009","created_at":%d}`, time.Now().Unix())
	// result = s.tester.Advance(admin, createContractInput)
	// s.Len(result.Notices, 1)
	// s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// createCrowdfundingInput := []byte(fmt.Sprintf(`{"path":"createdCrowdfunding","payload":{"max_interest_rate":"10","expires_at":%d,"debt_issued":%d}}`, time.Now().Add(5*time.Second).Unix(), 2020))
	// result = s.tester.Advance(creator, createCrowdfundingInput)
	// expectedOutput = fmt.Sprintf(`auction created - {"id":1,"creator":"vitalik","debt_issued":"2020","max_interest_rate":"10","state":"ongoing","expires_at":%d,"created_at":%d}`, time.Now().Add(5*time.Second).Unix(), time.Now().Unix())
	// s.Len(result.Notices, 1)
	// s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// createOrderInput := []byte(`{"path": "createOrder", "payload": {"creator": "vitalik","interest_rate":"9"}}`)
	// result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), investor01, big.NewInt(600), createOrderInput)
	// expectedOutput = fmt.Sprintf(`bid created - {"id":1,"auction_id":1,"investor":"0x0000000000000000000000000000000000000001","amount":"600","interest_rate":"9","state":"pending","created_at":%d}`, time.Now().Unix())
	// s.Len(result.Notices, 1)
	// s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// createOrderInput = []byte(`{"path": "createOrder", "payload": {"creator": "vitalik","interest_rate":"8"}}`)
	// result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), investor02, big.NewInt(520), createOrderInput)
	// expectedOutput = fmt.Sprintf(`bid created - {"id":2,"auction_id":1,"investor":"0x0000000000000000000000000000000000000002","amount":"520","interest_rate":"8","state":"pending","created_at":%d}`, time.Now().Unix())
	// s.Len(result.Notices, 1)
	// s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// createOrderInput = []byte(`{"path": "createOrder", "payload": {"creator": "vitalik","interest_rate":"4"}}`)
	// result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), investor03, big.NewInt(200), createOrderInput)
	// expectedOutput = fmt.Sprintf(`bid created - {"id":3,"auction_id":1,"investor":"0x0000000000000000000000000000000000000003","amount":"200","interest_rate":"4","state":"pending","created_at":%d}`, time.Now().Unix())
	// s.Len(result.Notices, 1)
	// s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// createOrderInput = []byte(`{"path": "createOrder", "payload": {"creator": "vitalik","interest_rate":"6"}}`)
	// result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), investor04, big.NewInt(300), createOrderInput)
	// expectedOutput = fmt.Sprintf(`bid created - {"id":4,"auction_id":1,"investor":"0x0000000000000000000000000000000000000004","amount":"300","interest_rate":"6","state":"pending","created_at":%d}`, time.Now().Unix())
	// s.Len(result.Notices, 1)
	// s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// createOrderInput = []byte(`{"path": "createOrder", "payload": {"creator": "vitalik","interest_rate":"4"}}`)
	// result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), investor05, big.NewInt(400), createOrderInput)
	// expectedOutput = fmt.Sprintf(`bid created - {"id":5,"auction_id":1,"investor":"0x0000000000000000000000000000000000000005","amount":"400","interest_rate":"4","state":"pending","created_at":%d}`, time.Now().Unix())
	// s.Len(result.Notices, 1)
	// s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// closeCrowdfundingInput := []byte(`{"path": "createOrder", "payload": {"creator": "vitalik","interest_rate":"4"}}`)
	// result = s.tester.Advance(admin, closeCrowdfundingInput)
	// expectedOutput = fmt.Sprintf(`bid created - {"id":5,"auction_id":1,"investor":"0x0000000000000000000000000000000000000005","amount":"400","interest_rate":"4","state":"pending","created_at":%d}`, time.Now().Unix())
	// s.Len(result.Notices, 1)
	// s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// settleCrowdfundingInput := []byte(`{"path":"finishdCrowdfunding", "payload":{"creator":"vitalik"}}`)
	// result = s.tester.DepositERC20(admin, common.HexToAddress("0x0000000000000000000000000000000000000009"), big.NewInt(0), settleCrowdfundingInput)
	// expectedOutput = fmt.Sprintf(`auction finished - {"id":1,"creator":"vitalik","debt_issued":"2020","max_interest_rate":"10","state":"finished","bids":[{"id":3,"auction_id":1,"investor":"0x0000000000000000000000000000000000000003","amount":"200","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},{"id":5,"auction_id":1,"investor":"0x0000000000000000000000000000000000000005","amount":"400","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},{"id":4,"auction_id":1,"investor":"0x0000000000000000000000000000000000000004","amount":"300","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},{"id":2,"auction_id":1,"investor":"0x0000000000000000000000000000000000000002","amount":"520","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},{"id":1,"auction_id":1,"investor":"0x0000000000000000000000000000000000000001","amount":"600","interest_rate":"9","state":"accepted","created_at":%d,"updated_at":%d}],"expires_at":%d,"created_at":%d,"updated_at":%d}`, time.Now().Add(-5*time.Second).Unix(), time.Now().Unix(), time.Now().Add(-5*time.Second).Unix(), time.Now().Unix(), time.Now().Add(-5*time.Second).Unix(), time.Now().Unix(), time.Now().Add(-5*time.Second).Unix(), time.Now().Unix(), time.Now().Add(-5*time.Second).Unix(), time.Now().Unix(), time.Now().Unix(), time.Now().Add(-5*time.Second).Unix(), time.Now().Unix())
	// s.Len(result.Notices, 1)
	// s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// creatorWithdrawInput := []byte(`{"path":"withdraw"}`)
	// expectedWithdrawVoucherPayload := make([]byte, 0, 4+32+32)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, creator[:]...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(1919).FillBytes(make([]byte, 32))...)
	// withdrawResult := s.tester.Advance(creator, creatorWithdrawInput)
	// s.Len(withdrawResult.Vouchers, 1)
	// s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)
	// s.Equal(common.HexToAddress("0x0000000000000000000000000000000000000009"), withdrawResult.Vouchers[0].Destination)

	// investor01WithdrawInput := []byte(`{"path":"withdraw"}`)
	// expectedWithdrawVoucherPayload = make([]byte, 0, 4+32+32)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, investor01[:]...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(1919).FillBytes(make([]byte, 32))...)
	// withdrawResult = s.tester.Advance(investor01, investor01WithdrawInput)
	// s.Len(withdrawResult.Vouchers, 1)
	// s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)
	// s.Equal(common.HexToAddress("0x0000000000000000000000000000000000000009"), withdrawResult.Vouchers[0].Destination)

	// investor02WithdrawInput := []byte(`{"path":"withdraw"}`)
	// expectedWithdrawVoucherPayload = make([]byte, 0, 4+32+32)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, investor02[:]...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(1919).FillBytes(make([]byte, 32))...)
	// withdrawResult = s.tester.Advance(investor02, investor02WithdrawInput)
	// s.Len(withdrawResult.Vouchers, 1)
	// s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)
	// s.Equal(common.HexToAddress("0x0000000000000000000000000000000000000009"), withdrawResult.Vouchers[0].Destination)

	// investor03WithdrawInput := []byte(`{"path":"withdraw"}`)
	// expectedWithdrawVoucherPayload = make([]byte, 0, 4+32+32)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, investor03[:]...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(1919).FillBytes(make([]byte, 32))...)
	// withdrawResult = s.tester.Advance(investor03, investor03WithdrawInput)
	// s.Len(withdrawResult.Notices, 1)
	// s.Len(withdrawResult.Vouchers, 1)
	// s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)

	// investor04WithdrawInput := []byte(`{"path":"withdraw"}`)
	// expectedWithdrawVoucherPayload = make([]byte, 0, 4+32+32)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, investor04[:]...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(1919).FillBytes(make([]byte, 32))...)
	// withdrawResult = s.tester.Advance(investor04, investor04WithdrawInput)
	// s.Len(withdrawResult.Vouchers, 1)
	// s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)
	// s.Equal(common.HexToAddress("0x0000000000000000000000000000000000000009"), withdrawResult.Vouchers[0].Destination)

	// investor05WithdrawInput := []byte(`{"path":"withdraw"}`)
	// expectedWithdrawVoucherPayload = make([]byte, 0, 4+32+32)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, investor05[:]...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(1919).FillBytes(make([]byte, 32))...)
	// withdrawResult = s.tester.Advance(investor05, investor05WithdrawInput)
	// s.Len(withdrawResult.Vouchers, 1)
	// s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)
	// s.Equal(common.HexToAddress("0x0000000000000000000000000000000000000009"), withdrawResult.Vouchers[0].Destination)

	// adminWithdrawInput := []byte(`{"path":"withdrawApp"}`)
	// expectedWithdrawVoucherPayload = make([]byte, 0, 4+32+32)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, admin[:]...)
	// expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(101).FillBytes(make([]byte, 32))...)
	// withdrawResult = s.tester.Advance(admin, adminWithdrawInput)
	// s.Len(withdrawResult.Vouchers, 1)
	// s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)
	// s.Equal(common.HexToAddress("0x0000000000000000000000000000000000000009"), withdrawResult.Vouchers[0].Destination)
}

func (s *DAppSuite) checkBalance(expected int64, payload []byte) {
	balance := new(big.Int).SetBytes(payload)
	s.Zerof(balance.Cmp(big.NewInt(expected)), "expected %v; got %v", expected, balance)
}
