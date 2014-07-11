// Package simulation contains all of the business logic required to simulate
// user's retirement. At the time of writing, this ingests JSON data and outputs
// JSON data.
package simulation

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Response   map[string]interface{}
	StatusCode int
}

type simulationData struct {
	InTodaysDollars          bool                    `json:"in_todays_dollars"`
	NumberOfTrials           int                     `json:"number_of_trials"`
	CholeskyDecomposition    []float64               `json:"cholesky_decomposition"`
	Inflation                distribution            `json:"inflation"`
	RealEstate               distribution            `json:"real_estate"`
	AssetPerformanceData     map[string]distribution `json:"asset_performance_data"`
	Parameters               parameters              `json:"simulation_parameters"`
	Expenses                 []expense               `json:"expenses"`
	SelectedPortfolioWeights portfolioWeights        `json:"selected_portfolio_weights"`
}

type portfolioWeights map[string]float64

type parameters struct {
	Male                   bool    `json:"male"`
	Married                bool    `json:"married"`
	Retired                bool    `json:"retired"`
	MaleAge                int     `json:"male_age"`
	RetirementAgeMale      int     `json:"retirement_age_male"`
	FemaleAge              int     `json:"female_age"`
	RetirementAgeFemale    int     `json:"retirement_age_female"`
	ExpensesMultiplier     float64 `json:"expenses_multiplier"`
	FractionSingleIncome   float64 `json:"fraction_single_income"`
	StartingAssets         float64 `json:"starting_assets"`
	Income                 float64 `json:"income"`
	CurrentTax             float64 `json:"current_tax"`
	SalaryIncrease         float64 `json:"salary_increase"`
	IncomeInflationIndex   float64 `json:"income_inflation_index"`
	ExpensesInflationIndex float64 `json:"expenses_inflation_index"`
	RetirementIncome       float64 `json:"retirement_income"`
	RetirementExpenses     float64 `json:"retirement_expenses"`
	RetirementTax          float64 `json:"retirement_tax"`
	LifeInsurance          float64 `json:"life_insurance"`
	IncludeHome            bool    `json:"include_home"`
	HomeValue              float64 `json:"home_value"`
	SellHouseIn            int     `json:"sell_house_in"`
	NewHomeRelVal          float64 `json:"new_home_relative_value"`
}

type distribution struct {
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
}

// ValidateAndHandleJsonInput is the main entry point into this package given
// a POST'ed JSON body.
// Receiver: None
// Params: r *http.Request
// Returns: ApiResponse {Response/StatusCode}
func ValidateAndHandleJsonInput(r *http.Request) ApiResponse {
	decoder := json.NewDecoder(r.Body)

	var simulationData simulationData

	err := decoder.Decode(&simulationData)
	if err != nil {
		return ApiResponse{
			Response: map[string]interface{}{
				"success": false,
				"message": "Invalid JSON structure.",
			},
			StatusCode: http.StatusBadRequest,
		}
	}

	prettyPrint(simulationData)
	resp := simulationData.simulate()

	return ApiResponse{
		Response: map[string]interface{}{
			"success":       true,
			"trial_results": resp,
		},
		StatusCode: http.StatusOK,
	}
}
