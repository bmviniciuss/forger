package generators

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type UUIDType string

var (
	UUIDTypeV4   UUIDType = "v4"
	UUIDTypeULID UUIDType = "ulid"
)

func NewUUIDType(uuidType string) UUIDType {
	switch uuidType {
	case "v4":
		return UUIDTypeV4
	case "ulid":
		return UUIDTypeULID
	default:
		return UUIDTypeV4
	}
}

func UUID(ctx context.Context, options ...interface{}) (string, error) {
	// Signature UUID(ctx)
	if len(options) <= 0 {
		return uuid.NewString(), nil
	}

	// Signature UUID(ctx, type)
	tRaw, ok := options[0].(string)
	if !ok {
		return "", errors.New("invalid type for uuid function")
	}
	t := NewUUIDType(tRaw)

	switch t {
	case UUIDTypeV4:
		return uuid.NewString(), nil
	case UUIDTypeULID:
		uu, _ := uuid.FromBytes(ulid.Make().Bytes())
		return uu.String(), nil
	default:
		return uuid.NewString(), nil
	}
}
