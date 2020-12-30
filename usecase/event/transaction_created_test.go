package event_test

import (
	event_entity "github.com/crypto-com/chain-indexing/entity/event"
	"github.com/crypto-com/chain-indexing/internal/primptr"
	"github.com/crypto-com/chain-indexing/test/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chain-indexing/usecase/coin"
	event_usecase "github.com/crypto-com/chain-indexing/usecase/event"
	"github.com/crypto-com/chain-indexing/usecase/model"
)

var _ = Describe("Event", func() {
	registry := event_entity.NewRegistry()
	event_usecase.RegisterEvents(registry)

	Describe("En/DecodeTransactionCreated", func() {
		It("should able to encode and decode to the same Event", func() {
			anyTxHash := factory.RandomTxHash()
			anyHeight := int64(1000)
			anyParams := model.CreateTransactionParams{
				TxHash:   anyTxHash,
				Code:     0,
				Log:      "{\"events\":[]}",
				MsgCount: 1,
				Signers: []model.TransactionSigner{
					{
						Type:            "/cosmos.crypto.secp256k1.PubKey",
						Pubkeys:         []string{"tcro1fmprm0sjy6lz9llv7rltn0v2azzwcwzvk2lsyn"},
						AccountSequence: uint64(1),
					},
					{
						Type: "/cosmos.crypto.multisig.LegacyAminoPubKey",
						Pubkeys: []string{
							"tcro1fmprm0sjy6lz9llv7rltn0v2azzwcwzvk2lsyn",
							"tcro1fmprm0sjy6lz9llv7rltn0v2azzwcwzvk2lsyn",
							"tcro1fmprm0sjy6lz9llv7rltn0v2azzwcwzvk2lsyn",
						},
						MaybeThreshold:  primptr.Int(2),
						AccountSequence: uint64(1),
					},
				},
				Fee:           coin.MustNewCoinFromString("1000"),
				FeePayer:      "tcro1fmprm0sjy6lz9llv7rltn0v2azzwcwzvk2lsyn",
				FeeGranter:    "tcro1fmprm0sjy6lz9llv7rltn0v2azzwcwzvk2lsyn",
				GasWanted:     200000,
				GasUsed:       10000,
				Memo:          "Test memo",
				TimeoutHeight: int64(10),
			}
			event := event_usecase.NewTransactionCreated(anyHeight, anyParams)

			encoded, err := event.ToJSON()
			Expect(err).To(BeNil())

			decodedEvent, err := registry.DecodeByType(
				event_usecase.TRANSACTION_CREATED, 1, []byte(encoded),
			)
			Expect(err).To(BeNil())
			Expect(decodedEvent).To(Equal(event))
			typedEvent, _ := decodedEvent.(*event_usecase.TransactionCreated)
			Expect(typedEvent.Name()).To(Equal(event_usecase.TRANSACTION_CREATED))
			Expect(typedEvent.Version()).To(Equal(1))

			Expect(typedEvent.TxHash).To(Equal(anyTxHash))
		})
	})
})
