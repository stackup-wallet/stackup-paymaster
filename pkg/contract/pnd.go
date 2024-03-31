package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func getAbiArgs() abi.Arguments {
	return abi.Arguments{
		{Name: "validUntil", Type: uint48},
		{Name: "validAfter", Type: uint48},
		{Name: "erc20Token", Type: address},
		{Name: "exchangeRate", Type: uint256},
	}
}

func EncodePaymasterAndData(
	data *Data,
	signature []byte,
) ([]byte, error) {
	args := getAbiArgs()
	packed, err := args.Pack(data.ValidUntil, data.ValidAfter, data.ERC20Token, data.ExchangeRate)
	if err != nil {
		return nil, err
	}

	concat := data.Paymaster.Bytes()
	concat = append(concat, packed...)
	concat = append(concat, signature...)
	return concat, nil
}
