package eth

import (
	gethethclient "github.com/ethereum/go-ethereum/ethclient"
)

// ==================================================================================================================
// contract client group
type contractGroup struct {
}

func newContractGroup(ethClient *gethethclient.Client, cfgGroup *ContractCfgGroup) (*contractGroup, error) {
	return nil, nil
}

// ==================================================================================================================
// contract bind opts
type bindOpts struct {
}

func newBindOpts(cfg *ContractBindOptsCfg) (*bindOpts, error) {

	return nil, nil
}
