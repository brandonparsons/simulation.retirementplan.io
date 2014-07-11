package simulation

import "math/rand"

// type LifeResult struct {
// 	Assets   float64 `json:"assets"`
// 	Income   float64 `json:"income"`
// 	Expenses float64 `json:"expenses"`
// 	jsTime   int     `json:"js_time"`
// }

type IndividualRunResults struct {
	// TimeStep []LifeResult
	HappinessLevel int    `json:"happiness_level"`
	Exclamation    string `json:"exclamation"`
}

func (s *SimulationData) RunIndividualSimulation() IndividualRunResults {
	return IndividualRunResults{
		HappinessLevel: rand.Intn(100),
		Exclamation:    "Booyah!",
	}
}
