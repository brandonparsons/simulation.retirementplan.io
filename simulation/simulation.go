// Package simulation contains all of the business logic required to simulate
// user's retirement. At the time of writing, this ingests JSON data and outputs
// JSON data.
package simulation

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"

	goStats "github.com/GaryBoone/GoStats/stats"
	"github.com/kr/pretty"
)

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
	resp := Simulate(&simulationData)

	return ApiResponse{
		Response: map[string]interface{}{
			"success":   true,
			"timesteps": resp,
		},
		StatusCode: http.StatusOK,
	}
}

type simulationResponse []summarizedTimeStep

// Simulate is the main call for simulations - it runs all of the trials, munges
// data, etc.
// Receiver: None
// Params: s *SimulationData
// Returns: simulationResponse ([]summarizedTimeStep)
func Simulate(s *SimulationData) simulationResponse {
	detailedResults := runSimulations(s)
	summarizedResults := summarizeResults(detailedResults)
	return summarizedResults
}

type simulationTimeStep struct {
	assets        float64
	income        float64
	expenses      float64
	jsTime        int
	maleAge       int
	femaleAge     int
	maleAlive     bool
	femaleAlive   bool
	maleRetired   bool
	femaleRetired bool
}

// runSimulations Gathers invididual trial results, as passes detailed data up
// to be summarized.
// Receiver: None
// Params: s *SimulationData
// Returns: [][]simulationTimeStep
func runSimulations(s *SimulationData) [][]simulationTimeStep {
	n := s.NumberOfTrials
	results := make([][]simulationTimeStep, n)

	// This does not change trial-to-trial, do only once.
	timeSteps := s.applyExpenses(s.numberOfMonthsToSimulate())

	type empty struct{}
	notifier := make(chan empty, n)
	for trial := 0; trial < n; trial++ {
		go func(i int) {
			results[i] = s.runIndividualSimulation(timeSteps)
			notifier <- empty{}
		}(trial)
	}

	// Wait for goroutines to finish
	for i := 0; i < n; i++ {
		<-notifier
	}

	return results
}

type summarizedTimeStep struct {
	AssetsMean   float64 `json:"assets_mean"`
	AssetsCILow  float64 `json:"assets_ci_low"`
	AssetsCIHigh float64 `json:"assets_ci_high"`

	IncomeMean   float64 `json:"income_mean"`
	IncomeCILow  float64 `json:"income_ci_low"`
	IncomeCIHigh float64 `json:"income_ci_high"`

	ExpensesMean   float64 `json:"expenses_mean"`
	ExpensesCILow  float64 `json:"expenses_ci_low"`
	ExpensesCIHigh float64 `json:"expenses_ci_high"`

	OutOfMoneyPercentage float64 `json:"out_of_money_percentage"`
	JsTime               int     `json:"js_time"`
}

func summarizeResults(detailedData [][]simulationTimeStep) []summarizedTimeStep {
	numberOfTrials := len(detailedData)
	numberOfPeriods := len(detailedData[0])

	// Prep the slice
	summarizedResults := make([]summarizedTimeStep, numberOfPeriods)

	for period := 0; period < numberOfPeriods; period++ {

		outOfMoneyOccurences := 0.0              // initialize value
		jsTime := detailedData[0][period].jsTime // same in every trial

		/* 	Transpose the arrays so that we have a list of asset/income/expense
		results by period, instead of by trial. */

		periodAssetResults := make([]float64, numberOfTrials)
		periodIncomeResults := make([]float64, numberOfTrials)
		periodExpensesResults := make([]float64, numberOfTrials)

		for trialIndex, arrayOfTrialResults := range detailedData {
			periodAssetResults[trialIndex] = arrayOfTrialResults[period].assets
			periodIncomeResults[trialIndex] = arrayOfTrialResults[period].income
			periodExpensesResults[trialIndex] = arrayOfTrialResults[period].expenses

			if arrayOfTrialResults[period].assets < 0 {
				outOfMoneyOccurences++
			}
		}

		/* */

		/* Generate descriptive statistics for each period */

		outOfMoneyPercentage := outOfMoneyOccurences / float64(numberOfTrials)

		assetsMean := goStats.StatsMean(periodAssetResults)
		assetsStdDev := goStats.StatsSampleStandardDeviation(periodAssetResults)
		assetsCIFactor := 1.96 * assetsStdDev / math.Pow(float64(numberOfTrials), 0.5)

		incomeMean := goStats.StatsMean(periodIncomeResults)
		incomeStdDev := goStats.StatsSampleStandardDeviation(periodIncomeResults)
		incomeCIFactor := 1.96 * incomeStdDev / math.Pow(float64(numberOfTrials), 0.5)

		expensesMean := goStats.StatsMean(periodExpensesResults)
		expensesStdDev := goStats.StatsSampleStandardDeviation(periodExpensesResults)
		expensesCIFactor := 1.96 * expensesStdDev / math.Pow(float64(numberOfTrials), 0.5)

		/* */

		summarizedResults[period] = summarizedTimeStep{
			AssetsMean:   assetsMean,
			AssetsCILow:  assetsMean - assetsCIFactor,
			AssetsCIHigh: assetsMean + assetsCIFactor,

			IncomeMean:   incomeMean,
			IncomeCILow:  incomeMean - incomeCIFactor,
			IncomeCIHigh: incomeMean + incomeCIFactor,

			ExpensesMean:   expensesMean,
			ExpensesCILow:  expensesMean - expensesCIFactor,
			ExpensesCIHigh: expensesMean + expensesCIFactor,

			OutOfMoneyPercentage: outOfMoneyPercentage,
			JsTime:               jsTime,
		}
	}

	return summarizedResults
}
