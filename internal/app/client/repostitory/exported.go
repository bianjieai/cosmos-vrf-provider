package repostitory

import (
	"github.com/irisnet/core-sdk-go/types"
	"gitlab.bianjie.ai/avata/contracts/vrf-provider/internal/app/client/repostitory/eth/contracts"
)

type IChain interface {
	GetLatestHeight() (uint64, error)
	GetRandomWordsRequestedEvent(height uint64) ([]*contracts.VrfRandomWordsRequested, error)

	FulfillRandomWords(msgs types.Msgs) types.Error

	ServiceName() string
}
