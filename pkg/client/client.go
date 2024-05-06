package client

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-logr/logr"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"github.com/stackup-wallet/stackup-paymaster/pkg/handlers"
	"github.com/stackup-wallet/stackup-paymaster/pkg/handlers/payg"
)

type Client struct {
	rpc          *rpc.Client
	eth          *ethclient.Client
	chainID      *big.Int
	ov           *gas.Overhead
	ep2pms       map[common.Address][]common.Address
	paygHandler  *payg.Handler
	logger       logr.Logger
	nativeTracer bool
}

func New(
	signer *signer.EOA,
	rpc *rpc.Client,
	eth *ethclient.Client,
	chain *big.Int,
	ov *gas.Overhead,
	ep2pms map[common.Address][]common.Address,
	l logr.Logger,
	nt bool,
) *Client {
	return &Client{
		rpc:          rpc,
		eth:          eth,
		chainID:      chain,
		ov:           ov,
		ep2pms:       ep2pms,
		paygHandler:  payg.New(signer, rpc, eth, chain, ov, nt),
		logger:       l,
		nativeTracer: nt,
	}
}

func (c *Client) Accounts(ep string) ([]string, error) {
	l := c.logger.WithName("pm_accounts")

	epAddr := common.HexToAddress(ep)
	pmAddr, ok := c.ep2pms[epAddr]
	if !ok {
		err := errors.New("entryPoint: Implementation not supported")
		l.Error(err, "pm_accounts error")
		return nil, err
	}
	l = l.WithValues("entrypoint", epAddr.String()).
		WithValues("paymasters", pmAddr).
		WithValues("chain_id", c.chainID.String())

	res := []string{}
	for _, pm := range pmAddr {
		res = append(res, pm.Hex())
	}
	l.Info("pm_accounts ok")

	return res, nil
}

func (c *Client) SponsorUserOperation(
	op map[string]any,
	ep string,
	ctx map[string]any,
) (*handlers.SponsorUserOperationResponse, error) {
	l := c.logger.WithName("pm_sponsorUserOperation")

	epAddr := common.HexToAddress(ep)
	pmAddrs, ok := c.ep2pms[epAddr]
	if !ok {
		err := errors.New("entryPoint: Implementation not supported")
		l.Error(err, "pm_sponsorUserOperation error")
		return nil, err
	}
	l = l.WithValues("entrypoint", epAddr.String()).
		WithValues("paymasters", pmAddrs).
		WithValues("chain_id", c.chainID.String())

	userOp, err := userop.New(op)
	if err != nil {
		err = fmt.Errorf("bad userOp: %s", err)
		l.Error(err, "pm_sponsorUserOperation error")
		return nil, err
	}

	ct, err := handlers.NewContextType(ctx)
	if err != nil {
		err = fmt.Errorf("bad context: %s", err)
		l.Error(err, "pm_sponsorUserOperation error")
		return nil, err
	}
	l = l.WithValues("type", ct.Type)

	switch ct.Type {
	case "payg":
		res, err := c.paygHandler.Run(userOp, epAddr, pmAddrs[0])
		if err != nil {
			l.Error(err, "pm_sponsorUserOperation error")
			return nil, err
		}

		l.Info("pm_sponsorUserOperation ok")
		return res, nil
	default:
		err := fmt.Errorf("type: %s not recognized", ct.Type)
		l.Error(err, "pm_sponsorUserOperation error")
		return nil, err
	}
}
