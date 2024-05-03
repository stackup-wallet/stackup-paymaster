package config

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Values struct {
	Port                   int
	DefaultEntryPoint      common.Address
	SigningKey             string
	EntryPointToPaymasters map[common.Address][]common.Address
	EthClientUrl           string

	// Observability variables.
	OTELServiceName      string
	OTELCollectorHeaders map[string]string
	OTELCollectorUrl     string
	OTELInsecureMode     bool

	// Rollup related variables.
	IsOpStackNetwork  bool
	IsArbStackNetwork bool

	// Undocumented variables.
	GinMode string
}

func envKeyValStringToMap(s string) map[string]string {
	out := map[string]string{}
	for _, pair := range strings.Split(s, "&") {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			break
		}
		out[kv[0]] = kv[1]
	}
	return out
}

func envArrayToAddressSlice(s string) []common.Address {
	env := strings.Split(s, ",")
	slc := []common.Address{}
	for _, ep := range env {
		slc = append(slc, common.HexToAddress(strings.TrimSpace(ep)))
	}

	return slc
}

func envKeyValAddressToAddressSlice(s string) map[common.Address][]common.Address {
	out := map[common.Address][]common.Address{}
	for _, pair := range strings.Split(s, "&") {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			break
		}
		out[common.HexToAddress(kv[0])] = envArrayToAddressSlice(kv[1])
	}
	return out
}

func variableNotSetOrIsNil(env string) bool {
	return !viper.IsSet(env) || viper.GetString(env) == ""
}

func GetValues() *Values {
	// Default variables
	viper.SetDefault("erc4337_paymaster_port", 43371)
	viper.SetDefault("erc4337_paymaster_default_entrypoint", "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")
	viper.SetDefault("erc4337_paymaster_otel_insecure_mode", false)
	viper.SetDefault("erc4337_paymaster_is_op_stack_network", false)
	viper.SetDefault("erc4337_paymaster_gin_mode", gin.ReleaseMode)

	// Read in from .env file if available
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			// Can ignore
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	// Read in from environment variables
	_ = viper.BindEnv("erc4337_paymaster_port")
	_ = viper.BindEnv("erc4337_paymaster_default_entrypoint")
	_ = viper.BindEnv("erc4337_paymaster_signing_key")
	_ = viper.BindEnv("erc4337_paymaster_entrypoint_to_paymasters")
	_ = viper.BindEnv("erc4337_paymaster_eth_client_url")
	_ = viper.BindEnv("erc4337_paymaster_otel_service_name")
	_ = viper.BindEnv("erc4337_paymaster_otel_collector_headers")
	_ = viper.BindEnv("erc4337_paymaster_otel_collector_url")
	_ = viper.BindEnv("erc4337_paymaster_otel_insecure_mode")
	_ = viper.BindEnv("erc4337_paymaster_is_op_stack_network")
	_ = viper.BindEnv("erc4337_paymaster_is_arb_stack_network")
	_ = viper.BindEnv("erc4337_paymaster_gin_mode")

	// Validate required variables
	if variableNotSetOrIsNil("erc4337_paymaster_signing_key") {
		panic("Fatal config error: erc4337_paymaster_signing_key not set")
	}

	if variableNotSetOrIsNil("erc4337_paymaster_entrypoint_to_paymasters") {
		panic("Fatal config error: erc4337_paymaster_entrypoint_to_paymasters not set")
	}

	if variableNotSetOrIsNil("erc4337_paymaster_eth_client_url") {
		panic("Fatal config error: erc4337_paymaster_eth_client_url not set")
	}

	// Validate O11Y variables
	if viper.IsSet("erc4337_paymaster_otel_service_name") &&
		variableNotSetOrIsNil("erc4337_paymaster_otel_collector_url") {
		panic("Fatal config error: erc4337_paymaster_otel_service_name is set without a collector URL")
	}

	// Return values
	port := viper.GetInt("erc4337_paymaster_port")
	defaultEntryPoint := common.HexToAddress(viper.GetString("erc4337_paymaster_default_entrypoint"))
	signingKey := viper.GetString("erc4337_paymaster_signing_key")
	entryPointToPaymasters := envKeyValAddressToAddressSlice(
		viper.GetString("erc4337_paymaster_entrypoint_to_paymasters"),
	)
	ethClientUrl := viper.GetString("erc4337_paymaster_eth_client_url")
	otelServiceName := viper.GetString("erc4337_paymaster_otel_service_name")
	otelCollectorHeader := envKeyValStringToMap(viper.GetString("erc4337_paymaster_otel_collector_headers"))
	otelCollectorUrl := viper.GetString("erc4337_paymaster_otel_collector_url")
	otelInsecureMode := viper.GetBool("erc4337_paymaster_otel_insecure_mode")
	isOpStackNetwork := viper.GetBool("erc4337_paymaster_is_op_stack_network")
	isArbStackNetwork := viper.GetBool("erc4337_paymaster_is_arb_stack_network")

	ginMode := viper.GetString("erc4337_paymaster_gin_mode")
	return &Values{
		Port:                   port,
		DefaultEntryPoint:      defaultEntryPoint,
		SigningKey:             signingKey,
		EntryPointToPaymasters: entryPointToPaymasters,
		EthClientUrl:           ethClientUrl,
		OTELServiceName:        otelServiceName,
		OTELCollectorHeaders:   otelCollectorHeader,
		OTELCollectorUrl:       otelCollectorUrl,
		OTELInsecureMode:       otelInsecureMode,
		IsOpStackNetwork:       isOpStackNetwork,
		IsArbStackNetwork:      isArbStackNetwork,
		GinMode:                ginMode,
	}
}
