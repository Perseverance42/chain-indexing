package event

import (
	"bytes"

	"github.com/crypto-com/chainindex/usecase/model"

	"github.com/crypto-com/chainindex/usecase/coin"

	"github.com/luci/go-render/render"

	entity_event "github.com/crypto-com/chainindex/entity/event"
	jsoniter "github.com/json-iterator/go"
)

const TRANSACTION_CREATED_NAME = "TransactionCreated"

type TransactionCreated struct {
	entity_event.Base

	TxHash    string    `json:"txHash"`
	Code      int       `json:"code"`
	Log       string    `json:"log"`
	MsgCount  int       `json:"msgCount"`
	Fee       coin.Coin `json:"fee"`
	GasWanted string    `json:"gasWanted"`
	GasUsed   string    `json:"gasUsed"`
}

func NewTransactionCreated(blockHeight int64, params model.CreateTransactionParams) *TransactionCreated {
	return &TransactionCreated{
		Base: entity_event.NewBase(entity_event.BaseParams{
			Name:        TRANSACTION_CREATED_NAME,
			Version:     1,
			BlockHeight: blockHeight,
		}),

		TxHash:    params.TxHash,
		Code:      params.Code,
		Log:       params.Log,
		MsgCount:  params.MsgCount,
		Fee:       params.Fee,
		GasWanted: params.GasWanted,
		GasUsed:   params.GasUsed,
	}
}

func (event *TransactionCreated) ToJSON() (string, error) {
	encoded, err := jsoniter.Marshal(event)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func (event *TransactionCreated) String() string {
	return render.Render(event)
}

func DecodeTransactionCreated(encoded []byte) (entity_event.Event, error) {
	jsonDecoder := jsoniter.NewDecoder(bytes.NewReader(encoded))
	jsonDecoder.DisallowUnknownFields()

	var event *TransactionCreated
	if err := jsonDecoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}