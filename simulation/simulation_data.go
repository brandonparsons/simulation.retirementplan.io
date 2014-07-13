package simulation

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/kr/pretty"
)

type SimulationData struct {
	NumberOfTrials           int                     `json:"number_of_trials"`
	CholeskyDecomposition    []float64               `json:"cholesky_decomposition"`
	Inflation                Distribution            `json:"inflation"`
	RealEstate               Distribution            `json:"real_estate"`
	AssetPerformanceData     map[string]Distribution `json:"asset_performance_data"`
	Parameters               Parameters              `json:"simulation_parameters"`
	Expenses                 []Expense               `json:"expenses"`
	SelectedPortfolioWeights map[string]float64      `json:"selected_portfolio_weights"`
}

type Parameters struct {
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

type Distribution struct {
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
}

type ApiResponse struct {
	Response   map[string]interface{}
	StatusCode int
}

// ValidateAndHandleJsonInput is the main entry point into this package given
// a POST'ed JSON body.
// Receiver: None
// Params: j io.ReadCloser (via r.Body)
// Returns: ApiResponse {Response/StatusCode}
func ValidateAndHandleJsonInput(j io.ReadCloser) ApiResponse {
	decoder := json.NewDecoder(j)

	var simulationData SimulationData

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

	log.Printf("%# v", pretty.Formatter(simulationData))
	resp := simulationData.Simulate()

	return ApiResponse{
		Response: map[string]interface{}{
			"success":       true,
			"trial_results": resp,
		},
		StatusCode: http.StatusOK,
	}
}
