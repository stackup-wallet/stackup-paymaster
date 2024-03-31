package client

import "github.com/stackup-wallet/stackup-paymaster/pkg/handlers"

type RpcAdapter struct {
	client *Client
}

func NewRpcAdapter(c *Client) *RpcAdapter {
	return &RpcAdapter{
		client: c,
	}
}

func (r *RpcAdapter) Pm_accounts(ep string) ([]string, error) {
	return r.client.Accounts(ep)
}

func (r *RpcAdapter) Pm_sponsorUserOperation(op map[string]any,
	ep string,
	ctx map[string]any) (*handlers.SponsorUserOperationResponse, error) {
	return r.client.SponsorUserOperation(op, ep, ctx)
}
