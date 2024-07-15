package usecases

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/checkmarble/marble-backend/models"
	"github.com/checkmarble/marble-backend/usecases/executor_factory"
	"github.com/guregu/null/v5"
	"github.com/pkg/errors"
)

type convoyWebhooksRepository interface {
	GetWebhook(ctx context.Context, webhookId string) (models.Webhook, error)
	ListWebhooks(ctx context.Context, organizationId string, partnerId null.String) ([]models.Webhook, error)
	RegisterWebhook(ctx context.Context, organizationId string, partnerId null.String, input models.WebhookRegister) error
	UpdateWebhook(ctx context.Context, input models.Webhook) error
	DeleteWebhook(ctx context.Context, webhookId string) error
}

type enforceSecurityWebhook interface {
	CanCreateWebhook(ctx context.Context, organizationId string, partnerId null.String) error
	CanReadWebhook(ctx context.Context, webhook models.Webhook) error
	CanModifyWebhook(ctx context.Context, webhook models.Webhook) error
}

type WebhooksUsecase struct {
	enforceSecurity    enforceSecurityWebhook
	executorFactory    executor_factory.ExecutorFactory
	transactionFactory executor_factory.TransactionFactory
	convoyRepository   convoyWebhooksRepository
}

func NewWebhooksUsecase(
	enforceSecurity enforceSecurityWebhook,
	executorFactory executor_factory.ExecutorFactory,
	transactionFactory executor_factory.TransactionFactory,
	convoyRepository convoyWebhooksRepository,
) WebhooksUsecase {
	return WebhooksUsecase{
		enforceSecurity:    enforceSecurity,
		executorFactory:    executorFactory,
		transactionFactory: transactionFactory,
		convoyRepository:   convoyRepository,
	}
}

func (usecase WebhooksUsecase) ListWebhooks(ctx context.Context, organizationId string, partnerId null.String) ([]models.Webhook, error) {
	webhooks, err := usecase.convoyRepository.ListWebhooks(ctx, organizationId, partnerId)
	if err != nil {
		return nil, errors.Wrap(err, "error listing webhooks")
	}

	for _, webhook := range webhooks {
		if err := usecase.enforceSecurity.CanReadWebhook(ctx, webhook); err != nil {
			return nil, err
		}
	}

	return webhooks, nil
}

func (usecase WebhooksUsecase) RegisterWebhook(
	ctx context.Context,
	organizationId string,
	partnerId null.String,
	input models.WebhookRegister,
) error {
	err := usecase.enforceSecurity.CanCreateWebhook(ctx, organizationId, partnerId)
	if err != nil {
		return err
	}

	if err = input.Validate(); err != nil {
		return err
	}

	input.Secret = generateSecret()

	err = usecase.convoyRepository.RegisterWebhook(ctx, organizationId, partnerId, input)
	if err != nil {
		return errors.Wrap(err, "error registering webhook")
	}

	return nil
}

func generateSecret() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(fmt.Errorf("generateSecret: %w", err))
	}
	return hex.EncodeToString(key)
}

func (usecase WebhooksUsecase) DeleteWebhook(
	ctx context.Context, organizationId string, partnerId null.String, webhookId string,
) error {
	webhook, err := usecase.convoyRepository.GetWebhook(ctx, webhookId)
	if err != nil {
		return models.NotFoundError
	}
	if err = usecase.enforceSecurity.CanModifyWebhook(ctx, webhook); err != nil {
		return err
	}

	return usecase.convoyRepository.DeleteWebhook(ctx, webhook.Id)
}

func (usecase WebhooksUsecase) UpdateWebhook(
	ctx context.Context, organizationId string, partnerId null.String, webhookId string, input models.WebhookUpdate,
) error {
	webhook, err := usecase.convoyRepository.GetWebhook(ctx, webhookId)
	if err != nil {
		return models.NotFoundError
	}
	if err = usecase.enforceSecurity.CanModifyWebhook(ctx, webhook); err != nil {
		return err
	}

	return usecase.convoyRepository.UpdateWebhook(ctx,
		models.MergeWebhookWithUpdate(webhook, input))
}
