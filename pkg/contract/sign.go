package contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
)

func Sign(message []byte, signer *signer.EOA) ([]byte, error) {
	digest := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(digest))
	sig, err := crypto.Sign(hash.Bytes(), signer.PrivateKey)
	if err != nil {
		return nil, err
	}
	sig[64] += 27

	return sig, nil
}
