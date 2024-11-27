package go_gcp_service

import (
	"context"

	"cloud.google.com/go/billing/budgets/apiv1"
	"google.golang.org/api/option"
)

type budgetClientLogger interface {
	Panic(args ...interface{})
}

type BudgetClientConfig struct {
	Logger       budgetClientLogger
	ClientOption *option.ClientOption
}

type BudgetClient struct {
	*budgets.BudgetClient
}

func NewGCPBudgetClient(clientConfig BudgetClientConfig) BudgetClient {
	budgetClient, err := budgets.NewBudgetClient(context.Background(), *clientConfig.ClientOption)

	if err != nil {
		clientConfig.Logger.Panic("Failed to create cloud budget api client: %v \n", err)
	}
	return BudgetClient{
		budgetClient,
	}
}
