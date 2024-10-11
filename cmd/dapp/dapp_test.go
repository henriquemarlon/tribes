package main

// import (
// 	"fmt"
// 	"math/big"
// 	"testing"
// 	"time"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/rollmelette/rollmelette"
// 	"github.com/stretchr/testify/suite"
// )

// func TestAppSuite(t *testing.T) {
// 	suite.Run(t, new(AppSuite))
// }

// type AppSuite struct {
// 	suite.Suite
// 	tester *rollmelette.Tester
// }

// func (s *AppSuite) SetupTest() {
// 	app := NewDAppMemory()
// 	s.tester = rollmelette.NewTester(app)
// }

// ////////////// User ///////////////////

// func (s *AppSuite) TestItCreateUser() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"createUser","payload":{"address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","role":"admin"}}`)
// 	expectedOutput := fmt.Sprintf(`created user - {"id":3,"role":"admin","address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, input)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }

// func (s *AppSuite) TestItCreateUserWithoutPermissions() {
// 	admin := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
// 	input := []byte(`{"path":"createUser","payload":{"address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","role":"admin"}}`)
// 	expectedOutput := `failed to find user by address 0x1234567890AbcdEF1234567890aBcdef12345678: record not found`
// 	result := s.tester.Advance(admin, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItCreateUserWithInvalidData() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"createUser","payload":{"address":"","role":""}}`)
// 	expectedOutput := `invalid user`
// 	result := s.tester.Advance(admin, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItDeleteUser() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")

// 	createUserInput := []byte(`{"path":"createUser","payload":{"address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","role":"admin"}}`)
// 	expectedOutput := fmt.Sprintf(`created user - {"id":3,"role":"admin","address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createUserInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	deleteUserInput := []byte(`{"path":"deleteUser","payload":{"address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8"}}`)
// 	expectedOutput = `deleted user with - {"address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8"}`
// 	result = s.tester.Advance(admin, deleteUserInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }

// func (s *AppSuite) TestItDeleteUserWithoutPermissions() {
// 	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
// 	input := []byte(`{"path":"deleteUser","payload":{"address":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8"}}`)
// 	expectedOutput := `failed to find user by address 0x1234567890AbcdEF1234567890aBcdef12345678: record not found`
// 	result := s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItDeleteNonExistentUser() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"deleteUser","payload":{"address":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65"}}`)
// 	expectedOutput := `user not found`
// 	result := s.tester.Advance(admin, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// /////////////////////////// Withdraw ////////////////////////////

// func (s *AppSuite) TestItWithdrawVolt() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc")

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), sender, big.NewInt(10000), []byte(""))

// 	withdrawInput := []byte(`{"path":"withdrawVolt"}`)

// 	expectedNoticePayload := `withdrawn VOLT and 10000 from 0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc with voucher index: 1`
// 	expectedWithdrawVoucherPayload := make([]byte, 0, 4+32+32)
// 	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
// 	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
// 	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, sender[:]...)
// 	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(10000).FillBytes(make([]byte, 32))...)
// 	withdrawResult := s.tester.Advance(sender, withdrawInput)
// 	s.Len(withdrawResult.Notices, 1)
// 	s.Len(withdrawResult.Vouchers, 1)
// 	s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)
// 	s.Equal(common.HexToAddress("0x0000000000000000000000000000000000000001"), withdrawResult.Vouchers[0].Destination)
// 	s.Equal(expectedNoticePayload, string(withdrawResult.Notices[0].Payload))
// }

// func (s *AppSuite) TestItWithdrawVoltWithInsuficientBalance() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc")

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	input := []byte(`{"path":"withdrawVolt"}`)
// 	expectedOutput = `no balance of VOLT to withdraw`
// 	result = s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItWithdrawStablecoin() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc")

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), sender, big.NewInt(10000), []byte(""))

// 	input := []byte(`{"path":"withdrawStablecoin"}`)
// 	expectedNoticePayload := `withdrawn STABLECOIN and 10000 from 0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc with voucher index: 1`
// 	expectedWithdrawVoucherPayload := make([]byte, 0, 4+32+32)
// 	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
// 	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
// 	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, sender[:]...)
// 	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(10000).FillBytes(make([]byte, 32))...)
// 	withdrawResult := s.tester.Advance(sender, input)
// 	s.Len(withdrawResult.Notices, 1)
// 	s.Len(withdrawResult.Vouchers, 1)
// 	s.Equal(expectedWithdrawVoucherPayload, withdrawResult.Vouchers[0].Payload)
// 	s.Equal(common.HexToAddress("0x0000000000000000000000000000000000000001"), withdrawResult.Vouchers[0].Destination)
// 	s.Equal(expectedNoticePayload, string(withdrawResult.Notices[0].Payload))
// }

// func (s *AppSuite) TestItWithdrawStablecoinWithInsuficientBalance() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc")

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	input := []byte(`{"path":"withdrawStablecoin"}`)

// 	expectedOutput = `no balance of STABLECOIN to withdraw`
// 	result = s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// ///////////////// Contract ///////////////////

// func (s *AppSuite) TestItCreateContract() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, input)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }

// func (s *AppSuite) TestItCreateContractWithoutPermissions() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	input := []byte(`{"path":"updateContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000002"}}`)
// 	result = s.tester.Advance(admin, input)
// 	s.Len(result.Notices, 1)
// 	expectedOutput = `failed to find user by address 0x1234567890AbcdEF1234567890aBcdef12345678: record not found`
// 	result = s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItCreateContractWithInvalidData() {
// 	sender := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"createContract","payload":{"symbol":"","address":""}}`)
// 	expectedOutput := `invalid contract`
// 	result := s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItDeleteContract() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	input := []byte(`{"path":"deleteContract","payload":{"symbol":"VOLT"}}`)
// 	expectedOutput = `deleted contract with - {"Symbol":"VOLT"}`
// 	result = s.tester.Advance(admin, input)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }

// func (s *AppSuite) TestItDeleteContractWithoutPermissions() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	input := []byte(`{"path":"deleteContract","payload":{"symbol":"VOLT"}}`)
// 	expectedOutput = `failed to find user by address 0x1234567890AbcdEF1234567890aBcdef12345678: record not found`
// 	result = s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItDeleteNonExistentContract() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"deleteContract","payload":{"symbol":"NONEXISTENT"}}`)
// 	expectedOutput := `contract not found`
// 	result := s.tester.Advance(admin, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItUpdateContract() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	input := []byte(`{"path":"updateContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000002"}}`)
// 	expectedOutput = fmt.Sprintf(`updated contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000002","created_at":%d,"updated_at":%d}`, time.Now().Unix(), time.Now().Unix())
// 	result = s.tester.Advance(admin, input)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }

// func (s *AppSuite) TestItUpdateContractWithoutPermissions() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	input := []byte(`{"path":"updateContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000003"}}`)
// 	expectedOutput = `record not found`
// 	result = s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItUpdateNonExistentContract() {
// 	sender := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"updateContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000002"}}`)
// 	expectedOutput := `contract not found`
// 	result := s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// ///////////////// Station ///////////////////

// func (s *AppSuite) TestItCreateStation() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"createStation","payload":{"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","consumption":100,"interest_rate":50,"latitude":40.7128,"longitude":-74.0060}}`)
// 	expectedOutput := fmt.Sprintf(`created station - {"id":1,"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","state":"active","consumption":"100","interest_rate":"50","latitude":40.7128,"longitude":-74.006,"created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, input)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }

// func (s *AppSuite) TestItCreateStationWithoutPermissions() {
// 	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
// 	input := []byte(`{"path":"createStation","payload":{"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","consumption":100,"interest_rate":50,"latitude":40.7128,"longitude":-74.0060}}`)
// 	expectedOutput := `failed to find user by address 0x1234567890AbcdEF1234567890aBcdef12345678: record not found`
// 	result := s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItCreateStationWithInvalidData() {
// 	sender := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"createStation","payload":{"owner":"","consumption":-100,"interest_rate":-50,"latitude":91.0000,"longitude":181.0000}}`)
// 	expectedOutput := `invalid station`
// 	result := s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItUpdateStation() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")

// 	createStationInput := []byte(`{"path":"createStation","payload":{"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","consumption":100,"interest_rate":50,"latitude":40.7128,"longitude":-74.0060}}`)
// 	expectedOutput := fmt.Sprintf(`created station - {"id":1,"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","state":"active","consumption":"100","interest_rate":"50","latitude":40.7128,"longitude":-74.006,"created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createStationInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	input := []byte(`{"path":"updateStation","payload":{"id":1, "owner": "0x1234567890abcdef1234567890abcdef12345678", "consumption": 150, "interest_rate": 75, "state": "active", "latitude": 34.0522, "longitude": -118.2437}}`)
// 	expectedOutput = fmt.Sprintf(`updated station - {"id":1,"consumption":"150","owner":"0x1234567890AbcdEF1234567890aBcdef12345678","interest_rate":"75","state":"active","latitude":34.0522,"longitude":-118.2437,"created_at":%d,"updated_at":%d}`, time.Now().Unix(), time.Now().Unix())
// 	result = s.tester.Advance(admin, input)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }

// func (s *AppSuite) TestItUpdateStationWithoutPermissions() {
// 	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678") // Not an admin
// 	input := []byte(`{"path":"updateStation","payload":{"id":1, "owner": "0x1234567890abcdef1234567890abcdef12345678", "consumption": 150, "interest_rate": 75, "latitude": 34.0522, "longitude": -118.2437}}`)
// 	expectedOutput := `failed to find user by address 0x1234567890AbcdEF1234567890aBcdef12345678: record not found`
// 	result := s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItUpdateNonExistentStation() {
// 	sender := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"updateStation","payload":{"id":10, "owner": "0x1234567890abcdef1234567890abcdef12345678", "consumption": 150, "interest_rate": 75, "latitude": 34.0522, "longitude": -118.2437}}`)
// 	expectedOutput := `station not found`
// 	result := s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItDeleteStation() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")

// 	createStationInput := []byte(`{"path":"createStation","payload":{"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","consumption":100,"interest_rate":50,"latitude":40.7128,"longitude":-74.0060}}`)
// 	expectedOutput := fmt.Sprintf(`created station - {"id":1,"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","state":"active","consumption":"100","interest_rate":"50","latitude":40.7128,"longitude":-74.006,"created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createStationInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	input := []byte(`{"path":"deleteStation","payload":{"id":1}}`)
// 	expectedOutput = `deleted station with - {"id":1}`
// 	result = s.tester.Advance(admin, input)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }

// func (s *AppSuite) TestItDeleteStationWithoutPermissions() {
// 	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678") // Not an admin
// 	input := []byte(`{"path":"deleteStation","payload":{"id":1}}`)
// 	expectedOutput := `failed to find user by address 0x1234567890AbcdEF1234567890aBcdef12345678: record not found`
// 	result := s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// func (s *AppSuite) TestItDeleteNonExistentStation() {
// 	sender := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	input := []byte(`{"path":"deleteStation","payload":{"id":10}}`)
// 	expectedOutput := `station not found`
// 	result := s.tester.Advance(sender, input)
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// ///////////////// Order ///////////////////

// func (s *AppSuite) TestItCreateOrder() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65")

// 	appAddressResult := s.tester.RelayAppAddress(common.HexToAddress("0xdadadadadadadadadadadadadadadadadadadada"))
// 	s.Nil(appAddressResult.Err)

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000002"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000002","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createStationInput := []byte(`{"path":"createStation","payload":{"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","consumption":100,"interest_rate":50,"latitude":40.7128,"longitude":-74.0060}}`)
// 	expectedOutput = fmt.Sprintf(`created station - {"id":1,"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","state":"active","consumption":"100","interest_rate":"50","latitude":40.7128,"longitude":-74.006,"created_at":%d}`, time.Now().Unix())
// 	result = s.tester.Advance(admin, createStationInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderPayload := []byte(`{"path":"createOrder","payload":{"station_id":1}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(10000), createOrderPayload)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":1,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"200","station_id":1,"station_owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","interest_rate":"50","created_at":%d}`, time.Now().Unix())
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }

// func (s *AppSuite) TestItCreateOrderWithInvalidData() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65")

// 	appAddressResult := s.tester.RelayAppAddress(common.HexToAddress("0xdadadadadadadadadadadadadadadadadadadada"))
// 	s.Nil(appAddressResult.Err)

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000002"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000002","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createStationInput := []byte(`{"path":"createStation","payload":{"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","consumption":100,"interest_rate":50,"latitude":40.7128,"longitude":-74.0060}}`)
// 	expectedOutput = fmt.Sprintf(`created station - {"id":1,"owner":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","state":"active","consumption":"100","interest_rate":"50","latitude":40.7128,"longitude":-74.006,"created_at":%d}`, time.Now().Unix())
// 	result = s.tester.Advance(admin, createStationInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderPayload := []byte(`{"path":"createOrder","payload":{"station_id":1}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(0), createOrderPayload)
// 	expectedOutput = "invalid order"
// 	s.ErrorContains(result.Err, expectedOutput)
// }

// ////////////// Auction //////////////////

// func (s *AppSuite) TestItCreateAuctionAndFinishAuctionWithoutPartialSelling() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65")
// 	stationOwner01 := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
// 	stationOwner02 := common.HexToAddress("0x3C44CdDdB6a9008054Ea72dBa55dDb1D8EDd895D")
// 	bidder01 := common.HexToAddress("0x0000000000000000000000000000000000000003")
// 	bidder02 := common.HexToAddress("0x0000000000000000000000000000000000000004")
// 	bidder03 := common.HexToAddress("0x0000000000000000000000000000000000000005")
// 	bidder04 := common.HexToAddress("0x0000000000000000000000000000000000000006")

// 	appAddressResult := s.tester.RelayAppAddress(common.HexToAddress("0xdadadadadadadadadadadadadadadadadadadada"))
// 	s.Nil(appAddressResult.Err)

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createContractInput = []byte(`{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000002"}}`)
// 	expectedOutput = fmt.Sprintf(`created contract - {"id":2,"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000002","created_at":%d}`, time.Now().Unix())
// 	result = s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createStationInput := []byte(fmt.Sprintf(`{"path":"createStation","payload":{"owner":"%v","consumption":100,"interest_rate":50,"latitude":40.7128,"longitude":-74.0060}}`, stationOwner01))
// 	expectedOutput = fmt.Sprintf(`created station - {"id":1,"owner":"%v","state":"active","consumption":"100","interest_rate":"50","latitude":40.7128,"longitude":-74.006,"created_at":%d}`, stationOwner01, time.Now().Unix())
// 	result = s.tester.Advance(admin, createStationInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createStationInput = []byte(fmt.Sprintf(`{"path":"createStation","payload":{"owner":"%v","consumption":100,"interest_rate":10,"latitude":40.7128,"longitude":-74.0060}}`, stationOwner02))
// 	expectedOutput = fmt.Sprintf(`created station - {"id":2,"owner":"%v","state":"active","consumption":"100","interest_rate":"10","latitude":40.7128,"longitude":-74.006,"created_at":%d}`, stationOwner02, time.Now().Unix())
// 	result = s.tester.Advance(admin, createStationInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderInput := []byte(`{"path":"createOrder","payload":{"station_id":1}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(10000), createOrderInput)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":1,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"200","station_id":1,"station_owner":"%v","interest_rate":"50","created_at":%d}`, stationOwner01, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderInput = []byte(`{"path":"createOrder","payload":{"station_id":2}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(100), createOrderInput)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":2,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"10","station_id":2,"station_owner":"%v","interest_rate":"10","created_at":%d}`, stationOwner02, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderInput = []byte(`{"path":"createOrder","payload":{"station_id":1}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(50000), createOrderInput)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":3,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"1000","station_id":1,"station_owner":"%v","interest_rate":"50","created_at":%d}`, stationOwner01, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderInput = []byte(`{"path":"createOrder","payload":{"station_id":2}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(100), createOrderInput)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":4,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"10","station_id":2,"station_owner":"%v","interest_rate":"10","created_at":%d}`, stationOwner02, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderInput = []byte(`{"path":"createOrder","payload":{"station_id":1}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(20000), createOrderInput)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":5,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"400","station_id":1,"station_owner":"%v","interest_rate":"50","created_at":%d}`, stationOwner01, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createAuctionInput := []byte(fmt.Sprintf(`{"path":"createAuction","payload":{"interest_rate":"1000","expires_at":%d,"debt_issued":%d}}`, time.Now().Add(5*time.Second).Unix(), 2))
// 	result = s.tester.Advance(admin, createAuctionInput)
// 	expectedOutput = fmt.Sprintf(`created auction - {"id":1,"debt_issued":"1620","interest_rate":"1000","state":"ongoing","expires_at":%d,"created_at":%d}`, time.Now().Add(5*time.Second).Unix(), time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createBidInput := []byte(`{"path": "createBid", "payload": {"interest_rate":"100"}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), bidder01, big.NewInt(600), createBidInput)
// 	expectedOutput = fmt.Sprintf(`created bid - {"id":1,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000003","amount":"600","interest_rate":"100","state":"pending","created_at":%d}`, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"500"}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), bidder02, big.NewInt(520), createBidInput)
// 	expectedOutput = fmt.Sprintf(`created bid - {"id":2,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000004","amount":"520","interest_rate":"500","state":"pending","created_at":%d}`, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"200"}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), bidder03, big.NewInt(200), createBidInput)
// 	expectedOutput = fmt.Sprintf(`created bid - {"id":3,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000005","amount":"200","interest_rate":"200","state":"pending","created_at":%d}`, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"300"}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), bidder04, big.NewInt(300), createBidInput)
// 	expectedOutput = fmt.Sprintf(`created bid - {"id":4,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000006","amount":"300","interest_rate":"300","state":"pending","created_at":%d}`, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	time.Sleep(5 * time.Second)

// 	finishAuctionInput := []byte(`{"path":"finishAuction"}`)
// 	result = s.tester.Advance(admin, finishAuctionInput)
// 	expectedOutput = `finished auction with - id: 1, required amount: 1620 and price limit per credit: 1000`
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	offSetStationConsumptionInput := []byte(`{"path":"offSetStationConsumption", "payload": {"id": 1, "amount_to_be_offset": 1600}}`)
// 	result = s.tester.Advance(stationOwner01, offSetStationConsumptionInput)
// 	expectedOutput = fmt.Sprintf(`offSet Amount from station: 1 by msg_sender: %v`, stationOwner01)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	offSetStationConsumptionInput = []byte(`{"path":"offSetStationConsumption", "payload": {"id": 2, "amount_to_be_offset": 20}}`)
// 	result = s.tester.Advance(stationOwner02, offSetStationConsumptionInput)
// 	expectedOutput = fmt.Sprintf(`offSet Amount from station: 2 by msg_sender: %v`, stationOwner02)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }

// func (s *AppSuite) TestItCreateAuctionAndFinishAuctionWithPartialSelling() {
// 	admin := common.HexToAddress("0x0142f501EE21f4446009C3505c51d0043feC5c68")
// 	sender := common.HexToAddress("0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65")
// 	stationOwner01 := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
// 	stationOwner02 := common.HexToAddress("0x3C44CdDdB6a9008054Ea72dBa55dDb1D8EDd895D")
// 	bidder01 := common.HexToAddress("0x0000000000000000000000000000000000000003")
// 	bidder02 := common.HexToAddress("0x0000000000000000000000000000000000000004")
// 	bidder03 := common.HexToAddress("0x0000000000000000000000000000000000000005")
// 	bidder04 := common.HexToAddress("0x0000000000000000000000000000000000000006")
// 	bidder05 := common.HexToAddress("0x0000000000000000000000000000000000000007")

// 	appAddressResult := s.tester.RelayAppAddress(common.HexToAddress("0xdadadadadadadadadadadadadadadadadadadada"))
// 	s.Nil(appAddressResult.Err)

// 	createContractInput := []byte(`{"path":"createContract","payload":{"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001"}}`)
// 	expectedOutput := fmt.Sprintf(`created contract - {"id":1,"symbol":"VOLT","address":"0x0000000000000000000000000000000000000001","created_at":%d}`, time.Now().Unix())
// 	result := s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createContractInput = []byte(`{"path":"createContract","payload":{"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000002"}}`)
// 	expectedOutput = fmt.Sprintf(`created contract - {"id":2,"symbol":"STABLECOIN","address":"0x0000000000000000000000000000000000000002","created_at":%d}`, time.Now().Unix())
// 	result = s.tester.Advance(admin, createContractInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createStationInput := []byte(fmt.Sprintf(`{"path":"createStation","payload":{"owner":"%v","consumption":100,"interest_rate":50,"latitude":40.7128,"longitude":-74.0060}}`, stationOwner01))
// 	expectedOutput = fmt.Sprintf(`created station - {"id":1,"owner":"%v","state":"active","consumption":"100","interest_rate":"50","latitude":40.7128,"longitude":-74.006,"created_at":%d}`, stationOwner01, time.Now().Unix())
// 	result = s.tester.Advance(admin, createStationInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createStationInput = []byte(fmt.Sprintf(`{"path":"createStation","payload":{"owner":"%v","consumption":100,"interest_rate":10,"latitude":40.7128,"longitude":-74.0060}}`, stationOwner02))
// 	expectedOutput = fmt.Sprintf(`created station - {"id":2,"owner":"%v","state":"active","consumption":"100","interest_rate":"10","latitude":40.7128,"longitude":-74.006,"created_at":%d}`, stationOwner02, time.Now().Unix())
// 	result = s.tester.Advance(admin, createStationInput)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderInput := []byte(`{"path":"createOrder","payload":{"station_id":1}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(10000), createOrderInput)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":1,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"200","station_id":1,"station_owner":"%v","interest_rate":"50","created_at":%d}`, stationOwner01, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderInput = []byte(`{"path":"createOrder","payload":{"station_id":2}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(100), createOrderInput)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":2,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"10","station_id":2,"station_owner":"%v","interest_rate":"10","created_at":%d}`, stationOwner02, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderInput = []byte(`{"path":"createOrder","payload":{"station_id":1}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(50000), createOrderInput)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":3,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"1000","station_id":1,"station_owner":"%v","interest_rate":"50","created_at":%d}`, stationOwner01, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderInput = []byte(`{"path":"createOrder","payload":{"station_id":2}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(100), createOrderInput)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":4,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"10","station_id":2,"station_owner":"%v","interest_rate":"10","created_at":%d}`, stationOwner02, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createOrderInput = []byte(`{"path":"createOrder","payload":{"station_id":1}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000002"), sender, big.NewInt(20000), createOrderInput)
// 	expectedOutput = fmt.Sprintf(`created order - {"id":5,"buyer":"0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65","amount":"400","station_id":1,"station_owner":"%v","interest_rate":"50","created_at":%d}`, stationOwner01, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createAuctionInput := []byte(fmt.Sprintf(`{"path":"createAuction","payload":{"interest_rate":"1000","expires_at":%d,"debt_issued":%d}}`, time.Now().Add(5*time.Second).Unix(), 2))
// 	result = s.tester.Advance(admin, createAuctionInput)
// 	expectedOutput = fmt.Sprintf(`created auction - {"id":1,"debt_issued":"1620","interest_rate":"1000","state":"ongoing","expires_at":%d,"created_at":%d}`, time.Now().Add(5*time.Second).Unix(), time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createBidInput := []byte(`{"path": "createBid", "payload": {"interest_rate":"100"}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), bidder01, big.NewInt(600), createBidInput)
// 	expectedOutput = fmt.Sprintf(`created bid - {"id":1,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000003","amount":"600","interest_rate":"100","state":"pending","created_at":%d}`, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"500"}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), bidder02, big.NewInt(520), createBidInput)
// 	expectedOutput = fmt.Sprintf(`created bid - {"id":2,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000004","amount":"520","interest_rate":"500","state":"pending","created_at":%d}`, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"200"}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), bidder03, big.NewInt(150), createBidInput)
// 	expectedOutput = fmt.Sprintf(`created bid - {"id":3,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000005","amount":"150","interest_rate":"200","state":"pending","created_at":%d}`, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"300"}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), bidder04, big.NewInt(300), createBidInput)
// 	expectedOutput = fmt.Sprintf(`created bid - {"id":4,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000006","amount":"300","interest_rate":"300","state":"pending","created_at":%d}`, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	createBidInput = []byte(`{"path": "createBid", "payload": {"interest_rate":"200"}}`)
// 	result = s.tester.DepositERC20(common.HexToAddress("0x0000000000000000000000000000000000000001"), bidder05, big.NewInt(150), createBidInput)
// 	expectedOutput = fmt.Sprintf(`created bid - {"id":5,"auction_id":1,"bidder":"0x0000000000000000000000000000000000000007","amount":"150","interest_rate":"200","state":"pending","created_at":%d}`, time.Now().Unix())
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	time.Sleep(5 * time.Second)

// 	finishAuctionInput := []byte(`{"path":"finishAuction"}`)
// 	result = s.tester.Advance(admin, finishAuctionInput)
// 	expectedOutput = `finished auction with - id: 1, required amount: 1620 and price limit per credit: 1000`
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	offSetStationConsumptionInput := []byte(`{"path":"offSetStationConsumption", "payload": {"id": 1, "amount_to_be_offset": 1600}}`)
// 	result = s.tester.Advance(stationOwner01, offSetStationConsumptionInput)
// 	expectedOutput = fmt.Sprintf(`offSet Amount from station: 1 by msg_sender: %v`, stationOwner01)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))

// 	offSetStationConsumptionInput = []byte(`{"path":"offSetStationConsumption", "payload": {"id": 2, "amount_to_be_offset": 20}}`)
// 	result = s.tester.Advance(stationOwner02, offSetStationConsumptionInput)
// 	expectedOutput = fmt.Sprintf(`offSet Amount from station: 2 by msg_sender: %v`, stationOwner02)
// 	s.Len(result.Notices, 1)
// 	s.Equal(expectedOutput, string(result.Notices[0].Payload))
// }
