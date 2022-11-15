package eth

import (
	"context"
	"math/big"
	"time"

	gethethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

const CtxTimeout = 10 * time.Second
const TryGetGasPriceTimeInterval = 10 * time.Second

type Eth struct {
	contractCfgGroup *ContractCfgGroup
	contracts        *contractGroup
	bindOpts         *bindOpts

	slot           int64
	maxGasPrice    *big.Int
	tipCoefficient float64

	ethClient  *gethethclient.Client
	gethCli    *gethclient.Client
	gethRpcCli *gethrpc.Client
}

func NewEth(config *ChainConfig) (*Eth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()
	rpcClient, err := gethrpc.DialContext(ctx, config.ChainURI)
	if err != nil {
		return nil, err
	}

	ethClient := gethethclient.NewClient(rpcClient)
	gethCli := gethclient.New(rpcClient)

	contractGroup, err := newContractGroup(ethClient, config.ContractCfgGroup)
	if err != nil {
		return nil, err
	}

	tmpBindOpts, err := newBindOpts(config.ContractBindOptsCfg)

	if err != nil {
		return nil, err
	}

	return &Eth{
		contractCfgGroup: config.ContractCfgGroup,
		ethClient:        ethClient,
		gethCli:          gethCli,
		gethRpcCli:       rpcClient,
		contracts:        contractGroup,
		bindOpts:         tmpBindOpts,
		slot:             config.Slot,
		tipCoefficient:   config.TipCoefficient,
		maxGasPrice:      new(big.Int).SetUint64(config.ContractBindOptsCfg.MaxGasPrice),
	}, nil
}
