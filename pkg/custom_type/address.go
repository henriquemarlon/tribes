package custom_type

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type Address struct {
	Address common.Address
}

func NewAddress(value common.Address) Address {
	return Address{Address: value}
}

func (a *Address) Scan(value interface{}) error {
	if value == nil {
		a.Address = common.Address{}
		return nil
	}

	switch v := value.(type) {
	case string:
		a.Address = common.HexToAddress(v)
	case []byte:
		a.Address = common.BytesToAddress(v)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

func (a Address) Value() (driver.Value, error) {
	return a.Address.Hex(), nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Address.Hex())
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var hexAddress string
	if err := json.Unmarshal(data, &hexAddress); err != nil {
		return err
	}
	a.Address = common.HexToAddress(hexAddress)
	return nil
}
