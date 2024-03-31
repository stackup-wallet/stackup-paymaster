package estimator

import "math/big"

var (
	// The maximum total gas limit for the entire UserOperation.
	maxGasLimit = big.NewInt(18000000)

	// This is a placeholder paymasterAndData with the correct length and all non-zero bytes. It is required
	// to calculate an acceptable preVerificationGas prior to paymaster approval. Once preVerificationGas is
	// calculated and added to the op, we can generate the real paymasterAndData. Note that if we calculate
	// the preVerificationGas of the op again (with the real paymasterAndData), it will return a lower value
	// since the real paymasterAndData is a mix of zero and non-zero bytes. If we want the preVerificationGas
	// calculation to be the same with both dummy and real paymasterAndData, we need to know how to generate a
	// dummy with the correct ratio of zero to non-zero bytes ahead of time.
	DummyPaymasterAndDataHex = "0x010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101"
)
