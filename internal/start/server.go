package start

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/jsonrpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"github.com/stackup-wallet/stackup-paymaster/internal/config"
	"github.com/stackup-wallet/stackup-paymaster/internal/logger"
	"github.com/stackup-wallet/stackup-paymaster/internal/o11y"
	"github.com/stackup-wallet/stackup-paymaster/pkg/client"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Server() {
	conf := config.GetValues()

	logr := logger.NewZeroLogr().WithName("stackup_paymaster")

	signer, err := signer.New(conf.SigningKey)
	if err != nil {
		log.Fatal(err)
	}

	rpc, err := rpc.Dial(conf.EthClientUrl)
	if err != nil {
		log.Fatal(err)
	}
	eth := ethclient.NewClient(rpc)

	chain, err := eth.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if o11y.IsEnabled(conf.OTELServiceName) {
		o11yOpts := &o11y.Opts{
			ServiceName:     conf.OTELServiceName,
			CollectorHeader: conf.OTELCollectorHeaders,
			CollectorUrl:    conf.OTELCollectorUrl,
			InsecureMode:    conf.OTELInsecureMode,

			ChainID:       chain,
			SignerAddress: signer.Address,
		}

		tracerCleanup := o11y.InitTracer(o11yOpts)
		defer tracerCleanup()

		metricsCleanup := o11y.InitMetrics(o11yOpts)
		defer metricsCleanup()
	}

	ov := gas.NewDefaultOverhead()
	if conf.IsArbStackNetwork || config.ArbStackChains.Contains(chain.Uint64()) {
		ov.SetCalcPreVerificationGasFunc(gas.CalcArbitrumPVGWithEthClient(rpc, conf.DefaultEntryPoint))
		ov.SetPreVerificationGasBufferFactor(16)
	}

	if conf.IsOpStackNetwork || config.OpStackChains.Contains(chain.Uint64()) {
		ov.SetCalcPreVerificationGasFunc(
			gas.CalcOptimismPVGWithEthClient(rpc, chain, conf.DefaultEntryPoint),
		)
		ov.SetPreVerificationGasBufferFactor(1)
	}

	c := client.New(signer, rpc, eth, chain, ov, conf.EntryPointToPaymasters, logr)

	gin.SetMode(conf.GinMode)
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}
	if o11y.IsEnabled(conf.OTELServiceName) {
		r.Use(otelgin.Middleware(conf.OTELServiceName))
	}
	r.Use(
		cors.Default(),
		logger.WithLogr(logr),
		gin.Recovery(),
	)
	r.GET("/ping", func(g *gin.Context) {
		g.Status(http.StatusOK)
	})
	handlers := []gin.HandlerFunc{
		jsonrpc.Controller(client.NewRpcAdapter(c)),
		jsonrpc.WithOTELTracerAttributes(),
	}
	r.POST("/", handlers...)
	r.POST("/rpc", handlers...)

	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		log.Fatal(err)
	}
}
