package go_gcp_service

import (
	"fmt"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/type/money"
)

type billingServiceLogger interface {
	Errorf(template string, args ...interface{})
	Panic(args ...interface{})
}

// BillingService -> handles the gcp billing related functions
type BillingService struct {
	projectName       string
	billingAccountID  string
	budgetDisplayName string
	budgetAmount      int64
	logger            billingServiceLogger
	gcpBilling        BillingClient
	budgetClient      BudgetClient
}

// NewGCPBillingService for the GCPBilling struct
func NewGCPBillingService(
	configService BillingService,
) BillingService {
	return configService
}

// GetBillingInfo Get Billing info for certain date
func (s BillingService) GetBillingInfo() (*cloudbilling.ProjectBillingInfo, error) {
	projectName := fmt.Sprintf("projects/%s", s.projectName)
	billingInfo, err := s.gcpBilling.Projects.GetBillingInfo(projectName).Do()

	return billingInfo, err
}

// GetExistingBudgetList Get Billing info for certain date
func (s BillingService) GetExistingBudgetList(
	ctx context.Context,
) (*budgetspb.Budget, error) {
	var budgetList []*budgetspb.Budget
	var err error
	parentId := fmt.Sprintf("billingAccounts/%s", s.billingAccountID)
	req := budgetspb.ListBudgetsRequest{
		Parent: parentId,
	}

	budgetsIter := s.budgetClient.ListBudgets(ctx, &req)
	for {
		budget, err := budgetsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			s.logger.Errorf("failed to retrieve budget: %v", err)
		}
		budgetList = append(budgetList, budget)
	}
	if len(budgetList) > 0 {
		return budgetList[0], err
	}

	return nil, err
}

func (s BillingService) GetBudgetCreateUpdateRequest() *budgetspb.Budget {
	projectName := fmt.Sprintf("projects/%s", s.projectName)

	budget := &budgetspb.Budget{
		DisplayName: "Project Budget",
		Name:        s.budgetDisplayName,
		BudgetFilter: &budgetspb.Filter{
			CreditTypesTreatment: budgetspb.Filter_INCLUDE_ALL_CREDITS,
			Projects:             []string{projectName},
		},
		Amount: &budgetspb.BudgetAmount{
			BudgetAmount: &budgetspb.BudgetAmount_SpecifiedAmount{
				SpecifiedAmount: &money.Money{
					Units:        s.budgetAmount,
					Nanos:        0,
					CurrencyCode: "JPY",
				},
			},
		},
		ThresholdRules: []*budgetspb.ThresholdRule{
			{
				ThresholdPercent: 0.25,
				SpendBasis:       budgetspb.ThresholdRule_CURRENT_SPEND,
			},
			{
				ThresholdPercent: 0.50,
				SpendBasis:       budgetspb.ThresholdRule_CURRENT_SPEND,
			},
			{
				ThresholdPercent: 0.75,
				SpendBasis:       budgetspb.ThresholdRule_CURRENT_SPEND,
			},
			{
				ThresholdPercent: 1,
				SpendBasis:       budgetspb.ThresholdRule_CURRENT_SPEND,
			},
		},
	}
	return budget
}

func (s BillingService) CreateBudget(ctx context.Context) (*budgetspb.Budget, error) {
	parentId := fmt.Sprintf("billingAccounts/%s", s.billingAccountID)
	budget := s.GetBudgetCreateUpdateRequest()
	createRequest := budgetspb.CreateBudgetRequest{
		Parent: parentId,
		Budget: budget,
	}

	billingInfo, err := s.budgetClient.CreateBudget(ctx, &createRequest)
	if err != nil {
		s.logger.Errorf("failed to create budget: %v\n", err)
	}

	return billingInfo, err
}

func (s BillingService) CreateOrUpdateBudget(ctx context.Context) (*budgetspb.Budget, error) {
	budget, _ := s.GetExistingBudgetList(ctx)
	if budget == nil {
		return s.CreateBudget(ctx)
	} else {
		return s.EditBudget(ctx, budget)
	}

}

func (s BillingService) EditBudget(ctx context.Context, budget *budgetspb.Budget) (*budgetspb.Budget, error) {
	// Modify the budget configuration here
	envBudget := s.GetBudgetCreateUpdateRequest()

	budget.Amount.BudgetAmount = envBudget.Amount.BudgetAmount
	budget.ThresholdRules = envBudget.ThresholdRules

	editRequest := budgetspb.UpdateBudgetRequest{
		Budget: budget,
	}

	billingInfo, err := s.budgetClient.UpdateBudget(ctx, &editRequest)
	if err != nil {
		s.logger.Errorf("failed to retrieve budget: %v\n", err)
	}

	return billingInfo, err
}