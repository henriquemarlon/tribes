package custom_type

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type Address common.Address

func HexToAddress(hex string) Address {
	return Address(common.HexToAddress(hex))
}

func (a *Address) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*a = Address(common.HexToAddress(v))
		return nil
	case []byte:
		*a = Address(common.HexToAddress(string(v)))
		return nil
	default:
		return fmt.Errorf("unsupported type for address scan: %T", value)
	}
}

func (a Address) Value() (driver.Value, error) {
	return common.Address(a).Hex(), nil
}

// MarshalJSON serializes the Address into a JSON string.
func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(common.Address(a).Hex())
}

// UnmarshalJSON deserializes a JSON string into the Address.
func (a *Address) UnmarshalJSON(data []byte) error {
	var hex string
	if err := json.Unmarshal(data, &hex); err != nil {
		return fmt.Errorf("failed to unmarshal Address: %v", err)
	}
	if !common.IsHexAddress(hex) {
		return fmt.Errorf("invalid hex address: %s", hex)
	}
	*a = Address(common.HexToAddress(hex))
	return nil
}
