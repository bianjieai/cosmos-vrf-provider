package eth

import (
	"context"
	"math/big"
	"time"

	"gitlab.bianjie.ai/avata/contracts/vrf-provider/internal/app/client/repostitory/eth/contracts"

	"github.com/ethereum/go-ethereum"

	gethcmn "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/irisnet/core-sdk-go/types"

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
		tipCoefficient:   config.TipCoefficient,
		maxGasPrice:      new(big.Int).SetUint64(config.ContractBindOptsCfg.MaxGasPrice),
	}, nil
}

func (eth Eth) GetLatestHeight() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CtxTimeout)
	defer cancel()
	return eth.ethClient.BlockNumber(ctx)
}

func (eth Eth) GetRandomWordsRequestedEvent(height uint64) ([]*contracts.VrfRandomWordsRequested, error) {
	address := gethcmn.HexToAddress(eth.contractCfgGroup.VRF.Addr)
	topic := eth.contractCfgGroup.VRF.Topic
	logs, err := eth.getLogs(address, topic, height, height)
	if err != nil {
		return nil, err
	}
	var vrfRqs []*contracts.VrfRandomWordsRequested
	for _, log := range logs {
		vrfRq, err := eth.contracts.VRF.ParseRandomWordsRequested(log)
		if err != nil {
			return nil, err
		}
		vrfRqs = append(vrfRqs, vrfRq)
	}

	return vrfRqs, nil
}

func (eth Eth) SubscriptionConsumerAddedEvent(height uint64) ([]*contracts.VrfSubscriptionConsumerAdded, error) {
	address := gethcmn.HexToAddress(eth.contractCfgGroup.VRF.Addr)
	topic := "SubscriptionConsumerAdded(uint64,address)"
	logs, err := eth.getLogs(address, topic, height, height)
	if err != nil {
		return nil, err
	}
	var vrfRqs []*contracts.VrfSubscriptionConsumerAdded
	for _, log := range logs {
		vrfRq, err := eth.contracts.VRF.ParseSubscriptionConsumerAdded(log)
		if err != nil {
			return nil, err
		}
		vrfRqs = append(vrfRqs, vrfRq)
	}

	return vrfRqs, nil
}

func (eth Eth) FulfillRandomWords(msgs types.Msgs) types.Error {

	return nil
}

func (eth Eth) RegisterProvider() {
}

func (eth Eth) ServiceName() string {
	return "vrf"
}

func (eth Eth) Contracts() *contractGroup {
	return eth.contracts
}

func (eth *Eth) getLogs(address gethcmn.Address, topic string, fromBlock, toBlock uint64) ([]gethtypes.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: []gethcmn.Address{address},
		Topics:    [][]gethcmn.Hash{{gethcrypto.Keccak256Hash([]byte(topic))}},
	}
	return eth.ethClient.FilterLogs(context.Background(), query)
}

func (eth *Eth) pair(x, y *big.Int) [2]*big.Int { return [2]*big.Int{x, y} }
