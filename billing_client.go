package go_gcp_service

import (
	"context"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/option"
)

type billingClientLogger interface {
	Panic(args ...interface{})
}

type BillingClientConfig struct {
	Logger       billingClientLogger
	ClientOption *option.ClientOption
}

type BillingClient struct {
	*cloudbilling.APIService
}

// NewGCPBillingClient creates a new gcp billing api client
func NewGCPBillingClient(config BillingClientConfig) BillingClient {
	billingClient, err := cloudbilling.NewService(context.Background(), *config.ClientOption)
	if err != nil {
		config.Logger.Panic("Failed to create cloud billing api client: %v \n", err)
	}

	return BillingClient{
		billingClient,
	}
}
