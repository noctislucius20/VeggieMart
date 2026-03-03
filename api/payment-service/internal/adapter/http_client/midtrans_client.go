package httpclient

import (
	"payment-service/config"

	"github.com/labstack/gommon/log"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransClientInterface interface {
	CreateTransaction(orderCode string, amount int64, customerName string, customerEmail string) (string, error)
}

type midtransClient struct {
	cfg    *config.Config
	logger *log.Logger
}

// CreateTransaction implements [MidtransClientInterface].
func (m *midtransClient) CreateTransaction(orderCode string, amount int64, customerName string, customerEmail string) (string, error) {
	midtrans.ServerKey = m.cfg.Midtrans.ServerKey
	midtrans.Environment = midtrans.EnvironmentType(m.cfg.Midtrans.Environment)

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderCode,
			GrossAmt: amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
	}

	snapRes, err := snap.CreateTransaction(snapReq)
	if err != nil {
		m.logger.Errorf("[MidtransClient-1] failed to create transaction: %v", err)
		return "", err
	}

	return snapRes.Token, nil
}

func NewMidtransclient(cfg *config.Config, logger *log.Logger) MidtransClientInterface {
	return &midtransClient{cfg: cfg, logger: logger}
}
