package channels

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/contracts/vrf-provider/internal/app/client/domain"
	"gitlab.bianjie.ai/avata/contracts/vrf-provider/internal/app/client/repostitory"
	typeserr "gitlab.bianjie.ai/avata/contracts/vrf-provider/internal/pkg/types/errors"
)

type IChannel interface {
	Relay() error
	IsNotRelay() bool
}

type Channel struct {
	source repostitory.IChain

	context *domain.Context

	logger *log.Logger
}

func NewChannel(source repostitory.IChain, startHeight uint64, logger *log.Logger) (IChannel, error) {
	return &Channel{
		logger:  logger,
		source:  source,
		context: domain.NewContext(startHeight, source.ServiceName()),
	}, nil
}

func (channel *Channel) Relay() error {
	return channel.relay()
}

func (channel *Channel) IsNotRelay() bool {
	curHeight := channel.Context().Height()
	latestHeight, err := channel.source.GetLatestHeight()
	if err != nil {
		return false
	}

	if curHeight < latestHeight {
		return true
	}

	return false
}

func (channel *Channel) Context() *domain.Context {
	return channel.context
}

func (channel *Channel) relay() error {
	logger := channel.logger.WithFields(log.Fields{
		"start_height": channel.Context().Height(),
		"option":       "relay",
	})
	latestHeight, err := channel.source.GetLatestHeight()
	if err != nil {
		logger.Error("failed to get latest height")
		return typeserr.ErrGetLatestHeight
	}
	if latestHeight <= channel.Context().Height() {
		logger.Info("the current height cannot be relayed yet")
		return typeserr.ErrNotProduced
	}
	time.Sleep(5 * time.Second)
	logger.Info("testing.....")
	// todo
	channel.Context().IncrHeight()
	return nil
}
