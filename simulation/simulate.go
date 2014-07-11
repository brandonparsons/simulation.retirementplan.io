package simulation

import (
	"math/rand"

	goStats "github.com/GaryBoone/GoStats/stats"
)

type simulationResponse struct {
	Results []IndividualRunResults `json:"results"`
}

type IndividualRunResults struct {
	// TimeStep []LifeResult
	HappinessLevel int    `json:"happiness_level"`
	Exclamation    string `json:"exclamation"`
}

// type LifeResult struct {
//  Assets   float64 `json:"assets"`
//  Income   float64 `json:"income"`
//  Expenses float64 `json:"expenses"`
//  jsTime   int     `json:"js_time"`
// }

// simulate is the main call for simulations - it runs all of the trials etc.
// Receiver: SimulationData
// Params: none
// Returns: simulationResponse
func (s *SimulationData) simulate() simulationResponse {
	n := s.NumberOfTrials
	results := make([]IndividualRunResults, n)

	type empty struct{}
	notifier := make(chan empty, n)

	for trial := uint(0); trial < n; trial++ {
		go func(i uint) {
			results[i] = s.RunIndividualSimulation()
			notifier <- empty{}
		}(trial)
	}

	// Wait for goroutines to finish
	for i := uint(0); i < n; i++ {
		<-notifier
	}

	return simulationResponse{Results: results}
}

// NumberOfMonthsToSimulate determines the number of months the simulation must cover,
// based on the user's ages.
// Params: none
// Returns: integer
func (s *SimulationData) NumberOfMonthsToSimulate() int {
	male := s.Parameters.MaleAge
	female := s.Parameters.FemaleAge

	var yearsToRun int
	if male == 0 {
		yearsToRun = 120 - female
	} else if female == 0 {
		yearsToRun = 120 - male
	} else {
		ages := []float64{float64(s.Parameters.MaleAge), float64(s.Parameters.FemaleAge)}
		yearsToRun = 120 - int(goStats.StatsMin(ages))
	}

	return 12 * yearsToRun
}

// RunIndividualSimulation is a single loop through the simulation. It is called
// by the `simulate` function
// Receiver: SimulationData
// Params: none
// Returns: IndividualRunResults
func (s *SimulationData) RunIndividualSimulation() IndividualRunResults {
	// perf := s.GenerateAssetPerformance(s.NumberOfMonthsToSimulate())
	// log.Println(perf)

	return IndividualRunResults{
		HappinessLevel: rand.Intn(100),
		Exclamation:    "Booyah!",
	}
}
