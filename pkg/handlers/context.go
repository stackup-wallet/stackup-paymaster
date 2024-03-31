package handlers

import (
	"errors"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

var (
	validate = validator.New()
	onlyOnce = sync.Once{}
)

func exactFieldMatch(mapKey, fieldName string) bool {
	return mapKey == fieldName
}

func decodeCtxTypes(
	f reflect.Kind,
	t reflect.Kind,
	data interface{}) (interface{}, error) {
	// String to common.Address conversion
	if f == reflect.String && t == reflect.Array {
		return common.HexToAddress(data.(string)), nil
	}

	// String to []byte conversion
	if f == reflect.String && t == reflect.Slice {
		byteStr := data.(string)
		if len(byteStr) < 2 || byteStr[:2] != "0x" {
			return nil, errors.New("not byte string")
		}

		b, err := hexutil.Decode(byteStr)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	return data, nil
}

func validateAddressType(field reflect.Value) interface{} {
	value, ok := field.Interface().(common.Address)
	if !ok || value == common.HexToAddress("0x") {
		return nil
	}

	return field
}

func decodeMap(data map[string]any, ctx any) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook: decodeCtxTypes,
		Result:     ctx,
		ErrorUnset: true,
		MatchName:  exactFieldMatch,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

func validateStruct(ctx any) error {
	onlyOnce.Do(func() {
		validate.RegisterCustomTypeFunc(validateAddressType, common.Address{})
	})
	return validate.Struct(ctx)
}

type ContextType struct {
	Type string `json:"type" mapstructure:"type" validate:"required"`
}

func NewContextType(data map[string]any) (*ContextType, error) {
	var ctx ContextType
	if err := decodeMap(data, &ctx); err != nil {
		return nil, err
	}

	if err := validateStruct(&ctx); err != nil {
		return nil, err
	}

	return &ctx, nil
}
