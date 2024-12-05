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

func (s *DAppSuite) TestItCreatedCrowdfundingAndSettle() {
	// Set up addresses
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
	investor01 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	investor02 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	investor03 := common.HexToAddress("0x0000000000000000000000000000000000000003")
	investor04 := common.HexToAddress("0x0000000000000000000000000000000000000004")
	investor05 := common.HexToAddress("0x0000000000000000000000000000000000000005")

	appAddressResult := s.tester.RelayAppAddress(common.HexToAddress("0xdadadadadadadadadadadadadadadadadadadada"))
	s.Nil(appAddressResult.Err)

	baseTime := time.Now().Unix()

	// Create users
	createUserInput := []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"creator"}}`, creator))
	result := s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)

	createUserInput = []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"qualified_investor"}}`, investor01))
	result = s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)

	createUserInput = []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"qualified_investor"}}`, investor02))
	result = s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)

	createUserInput = []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"non_qualified_investor"}}`, investor03))
	result = s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)

	createUserInput = []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"non_qualified_investor"}}`, investor04))
	result = s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)

	createUserInput = []byte(fmt.Sprintf(`{"path":"createUser","payload":{"address":"%s","role":"non_qualified_investor"}}`, investor05))
	result = s.tester.Advance(admin, createUserInput)
	s.Len(result.Notices, 1)

	// Create contracts
	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000008"}}`)
	result = s.tester.Advance(admin, createContractInput)
	s.Len(result.Notices, 1)

	createContractInput = []byte(`{"path":"createContract","payload":{"symbol":"PINK","address":"0x0000000000000000000000000000000000000009"}}`)
	result = s.tester.Advance(admin, createContractInput)
	s.Len(result.Notices, 1)

	// Set closesAt and maturityAt to future timestamps
	closesAt := baseTime + 5
	maturityAt := baseTime + 10

	// Create crowdfunding
	createCrowdfundingInput := []byte(fmt.Sprintf(`{"path":"createCrowdfunding","payload":{"max_interest_rate":"10", "debt_issued":"100000", "fundraising_duration":10, "closes_at":%d,"maturity_at":%d}}`, closesAt, maturityAt))
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000009"), creator, big.NewInt(10000), createCrowdfundingInput)
	s.Len(result.Notices, 1)

	// Update crowdfunding state to ongoing
	updateCrowdfundingInput := []byte(`{"path":"updateCrowdfunding","payload":{"id":1,"state":"ongoing"}}`)
	result = s.tester.Advance(admin, updateCrowdfundingInput)
	s.Len(result.Notices, 1)

	orderCreatedAt := baseTime

	// Investors create orders
	createOrderInput := []byte(`{"path": "createOrder", "payload": {"crowdfunding_id":1,"interest_rate":"9"}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000008"), investor01, big.NewInt(60000), createOrderInput)
	s.Len(result.Notices, 1)

	createOrderInput = []byte(`{"path": "createOrder", "payload": {"crowdfunding_id":1,"interest_rate":"8"}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000008"), investor02, big.NewInt(52000), createOrderInput)
	s.Len(result.Notices, 1)

	createOrderInput = []byte(`{"path": "createOrder", "payload": {"crowdfunding_id":1,"interest_rate":"4"}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000008"), investor03, big.NewInt(2000), createOrderInput)
	s.Len(result.Notices, 1)

	createOrderInput = []byte(`{"path": "createOrder", "payload": {"crowdfunding_id":1,"interest_rate":"6"}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000008"), investor04, big.NewInt(3000), createOrderInput)
	s.Len(result.Notices, 1)

	createOrderInput = []byte(`{"path": "createOrder", "payload": {"crowdfunding_id":1,"interest_rate":"4"}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000008"), investor05, big.NewInt(400), createOrderInput)
	s.Len(result.Notices, 1)

	time.Sleep(5 * time.Second)

	closeCrowdfundingInput := []byte(fmt.Sprintf(`{"path": "closeCrowdfunding", "payload": {"creator": "%s"}}`, creator))
	result = s.tester.Advance(admin, closeCrowdfundingInput)
	s.Len(result.Notices, 1)

	updatedAt := baseTime + 5 // baseTime + sleep duration

	// Expected output for closing crowdfunding
	expectedOutput := fmt.Sprintf(`crowdfunding closed - {"id":1,"token":"0x0000000000000000000000000000000000000009","amount":"10000","creator":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108270","orders":[`+
		`{"id":1,"crowdfunding_id":1,"investor":"%s","amount":"42600","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"crowdfunding_id":1,"investor":"%s","amount":"52000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"crowdfunding_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"crowdfunding_id":1,"investor":"%s","amount":"3000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"crowdfunding_id":1,"investor":"%s","amount":"400","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"crowdfunding_id":1,"investor":"%s","amount":"17400","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"state":"closed","fundraising_duration":10,"closes_at":%d,"maturity_at":%d,"created_at":%d,"updated_at":%d}`,
		creator.Hex(),
		investor01.Hex(), orderCreatedAt, updatedAt, // Order 1
		investor02.Hex(), orderCreatedAt, updatedAt, // Order 2
		investor03.Hex(), orderCreatedAt, updatedAt, // Order 3
		investor04.Hex(), orderCreatedAt, updatedAt, // Order 4
		investor05.Hex(), orderCreatedAt, updatedAt, // Order 5
		investor01.Hex(), updatedAt, updatedAt, // Order 6 (rejected portion)
		closesAt, maturityAt, baseTime, updatedAt,
	)
	s.Equal(expectedOutput, string(result.Notices[0].Payload))

	// Settle crowdfunding before maturity date
	settleCrowdfundingInput := []byte(`{"path":"settleCrowdfunding", "payload":{"crowdfunding_id":1}}`)
	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000008"), creator, big.NewInt(108270), settleCrowdfundingInput)
	s.Len(result.Notices, 1)

	settledAt := updatedAt // baseTime

	// Expected output for settling crowdfunding
	expectedOutput = fmt.Sprintf(`crowdfunding settled - {"id":1,"token":"0x0000000000000000000000000000000000000009","amount":"10000","creator":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108270","orders":[`+
		`{"id":1,"crowdfunding_id":1,"investor":"%s","amount":"42600","interest_rate":"9","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":2,"crowdfunding_id":1,"investor":"%s","amount":"52000","interest_rate":"8","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":3,"crowdfunding_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":4,"crowdfunding_id":1,"investor":"%s","amount":"3000","interest_rate":"6","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":5,"crowdfunding_id":1,"investor":"%s","amount":"400","interest_rate":"4","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":6,"crowdfunding_id":1,"investor":"%s","amount":"17400","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"state":"settled","fundraising_duration":10,"closes_at":%d,"maturity_at":%d,"created_at":%d,"updated_at":%d}`,
		creator.Hex(),
		investor01.Hex(), orderCreatedAt, settledAt,
		investor02.Hex(), orderCreatedAt, settledAt,
		investor03.Hex(), orderCreatedAt, settledAt,
		investor04.Hex(), orderCreatedAt, settledAt,
		investor05.Hex(), orderCreatedAt, settledAt,
		investor01.Hex(), updatedAt, updatedAt, // Order 6 remains rejected with previous updatedAt
		closesAt, maturityAt, baseTime, settledAt,
	)

	s.Equal(expectedOutput, string(result.Notices[0].Payload))
}
