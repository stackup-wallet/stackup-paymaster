package contract

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Data struct {
	Paymaster    common.Address
	ValidUntil   *big.Int
	ValidAfter   *big.Int
	ERC20Token   common.Address
	ExchangeRate *big.Int
}

func NewData(pm common.Address, token common.Address, exchangeRate *big.Int) *Data {
	return &Data{
		Paymaster:    pm,
		ValidUntil:   big.NewInt(int64(time.Now().Add(time.Hour).Unix())),
		ValidAfter:   big.NewInt(0),
		ERC20Token:   token,
		ExchangeRate: exchangeRate,
	}
}
