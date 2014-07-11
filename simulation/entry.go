// Package simulation contains all of the business logic required to simulate
// user's retirement. At the time of writing, this ingests JSON data and outputs
// JSON data.
package simulation

import (
	"encoding/json"
	"log"
	"net/http"
)

type SimulationData struct {
	InTodaysDollars          bool                    `json:"in_todays_dollars"`
	NumberOfTrials           uint                    `json:"number_of_trials"`
	CholeskyDecomposition    []float64               `json:"cholesky_decomposition"`
	Inflation                Distribution            `json:"inflation"`
	RealEstate               Distribution            `json:"real_estate"`
	AssetPerformanceData     map[string]Distribution `json:"asset_performance_data"`
	Parameters               Parameters              `json:"simulation_parameters"`
	Expenses                 []Expense               `json:"expenses"`
	SelectedPortfolioWeights PortfolioWeights        `json:"selected_portfolio_weights"`
}

type PortfolioWeights map[string]float64

type AssetPerformance map[string]Distribution

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

type Expense struct {
	Amount    int    `json:"amount"`
	Frequency string `json:"frequency"`
	OneTimeOn int    `json:"onetime_on"`
	Ends      int    `json:"ends"`
}

type Distribution struct {
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
}

func ValidateAndHandleJsonInput(r *http.Request) (response, int) {
	decoder := json.NewDecoder(r.Body)

	var simulationData SimulationData

	err := decoder.Decode(&simulationData)
	if err != nil {
		log.Println("Invalid JSON structure.")
		return response{
			"success": false,
			"message": "Invalid JSON structure.",
		}, http.StatusBadRequest
	}

	prettyPrint(simulationData)
	resp := simulationData.simulate()

	return response{
		"success":            true,
		"simulation_results": resp,
	}, http.StatusOK
}
