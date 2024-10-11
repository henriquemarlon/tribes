package custom_type

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
)

type BigInt struct {
	*big.Int
}

func NewBigInt(value *big.Int) BigInt {
	return BigInt{Int: value}
}

func (b *BigInt) Scan(value interface{}) error {
	if value == nil {
		b.Int = big.NewInt(0)
		return nil
	}
	switch v := value.(type) {
	case int64:
		b.Int = big.NewInt(v)
	case []byte:
		b.Int = new(big.Int)
		_, ok := b.Int.SetString(string(v), 10)
		if !ok {
			return fmt.Errorf("failed to parse BigInt from bytes: %s", v)
		}
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

func (b BigInt) Value() (driver.Value, error) {
	return b.String(), nil
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	var strValue string
	var number json.Number
	if err := json.Unmarshal(data, &number); err == nil {
		strValue = number.String()
	} else {
		if err := json.Unmarshal(data, &strValue); err != nil {
			return err
		}
	}

	b.Int = new(big.Int)
	_, ok := b.Int.SetString(strValue, 10)
	if !ok {
		return fmt.Errorf("failed to parse BigInt from string: %s", strValue)
	}
	return nil
}
