package accrual

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/config"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/accrual"
	"net/http"
)

type WebClient struct {
	accrualClient *resty.Client
}

type accrualResponse struct {
	OrderNumber string           `json:"number"`
	Status      string           `json:"status"`
	Accrual     *decimal.Decimal `json:"accrual,omitempty"`
}

func NewWebClient(cfg config.Config) accrual.Repository {
	client := resty.New()
	log.Info().Str("accrualBaseURL", cfg.AccrualURL).Msg("Setting base url")
	client.SetBaseURL(cfg.AccrualURL)
	client.SetRetryCount(5)

	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	// Registering Request Middleware
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		log.Debug().Str("url", req.URL).Msg("trying to get data from external accrual system")
		return nil // if its success otherwise return error
	})

	return &WebClient{accrualClient: client}
}

func (w *WebClient) Get(ctx context.Context, orderNumber entity.OrderNumber) (*entity.AccrualOrder, error) {
	var result accrualResponse
	resp, err := w.
		accrualClient.
		R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/orders/%s", orderNumber))
	if err != nil {
		log.Error().Err(err).Msg("error while getting accrual from external service")
		return nil, err
	}
	log.Info().Str("status", resp.Status()).Str("response", resp.String()).Msg("got response")
	if resp.StatusCode() != http.StatusOK {
		return nil, nil
	}
	status, err := entity.StringToAccrualOrderStatus(result.Status)
	if err != nil {
		return nil, err
	}

	return &entity.AccrualOrder{
		OrderNumber: orderNumber,
		Status:      status,
		Accrual:     result.Accrual,
	}, nil

}
