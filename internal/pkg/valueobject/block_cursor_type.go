package valueobject

import (
	"database/sql/driver"
	"fmt"
)

type BlockCursorType string

const (
	BlockCursorTypeScanner   BlockCursorType = "scanner"
	BlockCursorTypeFinalizer BlockCursorType = "finalizer"
)

// Scan scan value into Jsonb, implements sql.Scanner interface
func (t *BlockCursorType) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value (%v)", value)
	}
	*t = BlockCursorType(s)
	return nil
}

// Value return BlockCursorType value, implement driver.Valuer interface
func (t BlockCursorType) Value() (driver.Value, error) {
	return string(t), nil
}
