package block

import (
	"fmt"

	view2 "github.com/crypto-com/chain-indexing/appinterface/projection/block/view"

	entity_projection "github.com/crypto-com/chain-indexing/entity/projection"

	"github.com/crypto-com/chain-indexing/appinterface/projection/rdbbase"
	"github.com/crypto-com/chain-indexing/appinterface/rdb"
	event_entity "github.com/crypto-com/chain-indexing/entity/event"
	applogger "github.com/crypto-com/chain-indexing/internal/logger"
	event_usecase "github.com/crypto-com/chain-indexing/usecase/event"
)

var _ entity_projection.Projection = &Block{}

// TODO: Listen to council node related events and project council node
type Block struct {
	*rdbbase.RDbBase

	rdbConn rdb.Conn
	logger  applogger.Logger
}

func NewBlock(logger applogger.Logger, rdbConn rdb.Conn) *Block {
	return &Block{
		rdbbase.NewRDbBase(rdbConn.ToHandle(), "Block"),

		rdbConn,
		logger,
	}
}

func (_ *Block) GetEventsToListen() []string {
	return []string{event_usecase.BLOCK_CREATED}
}

func (projection *Block) OnInit() error {
	return nil
}

func (projection *Block) HandleEvents(height int64, events []event_entity.Event) error {
	rdbTx, err := projection.rdbConn.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = rdbTx.Rollback()
		}
	}()

	rdbTxHandle := rdbTx.ToHandle()
	blocksView := view2.NewBlocks(rdbTxHandle)

	for _, event := range events {
		if blockCreatedEvent, ok := event.(*event_usecase.BlockCreated); ok {
			if handleErr := projection.handleBlockCreatedEvent(blocksView, blockCreatedEvent); handleErr != nil {
				return fmt.Errorf("error handling BlockCreatedEvent: %v", handleErr)
			}
		} else {
			return fmt.Errorf("received unexpected event %sV%d(%s)", event.Name(), event.Version(), event.UUID())
		}
	}
	if err = projection.UpdateLastHandledEventHeight(rdbTxHandle, height); err != nil {
		return fmt.Errorf("error updating last handled event height: %v", err)
	}

	if err = rdbTx.Commit(); err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}
	committed = true
	return nil
}

func (projection *Block) handleBlockCreatedEvent(blocksView *view2.Blocks, event *event_usecase.BlockCreated) error {
	committedCouncilNodes := make([]view2.BlockCommittedCouncilNode, 0)
	for _, signature := range event.Block.Signatures {
		committedCouncilNodes = append(committedCouncilNodes, view2.BlockCommittedCouncilNode{
			Address:    signature.ValidatorAddress,
			Time:       signature.Timestamp,
			Signature:  signature.Signature,
			IsProposer: event.Block.ProposerAddress == signature.ValidatorAddress,
		})
	}

	if err := blocksView.Insert(&view2.Block{
		Height:                event.Block.Height,
		Hash:                  event.Block.Hash,
		Time:                  event.Block.Time,
		AppHash:               event.Block.AppHash,
		TransactionCount:      len(event.Block.Txs),
		CommittedCouncilNodes: committedCouncilNodes,
	}); err != nil {
		return err
	}

	return nil
}