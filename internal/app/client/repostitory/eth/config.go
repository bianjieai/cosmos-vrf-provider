package eth

func NewChainConfig() *ChainConfig {
	return &ChainConfig{}
}

type ChainConfig struct {
	ChainURI string
	ChainID  uint64

	Slot           int64
	TipCoefficient float64

	ContractCfgGroup    *ContractCfgGroup
	ContractBindOptsCfg *ContractBindOptsCfg
}

func NewContractCfgGroup() *ContractCfgGroup {
	return &ContractCfgGroup{}
}

type ContractCfgGroup struct {
	Erc1155Bank ContractCfg
}

type ContractCfg struct {
	Addr       string
	Topic      string
	OptPrivKey string
}

func NewContractBindOptsCfg() *ContractBindOptsCfg {
	return &ContractBindOptsCfg{}
}

type ContractBindOptsCfg struct {
	CountPrivKey string
	GasLimit     uint64
	MaxGasPrice  uint64
	ChainID      uint64
}
