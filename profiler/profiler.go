package main

import (
	"fmt"

	"bitbucket.org/retirementplanio/go-simulation/simulation"
	"github.com/davecheney/profile"
)

//////////
// Main //
//////////

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	s := simulation.SimulationData{
		NumberOfTrials:        5000,
		CholeskyDecomposition: []float64{0.0094794922, 0, 0, -7.36e-05, 0.0055677999, 0, 0.0050681903, -0.0004821709, 0.013367741},
		Inflation:             simulation.Distribution{Mean: 0.00046346514957523, StdDev: 0.00024792742828969},
		RealEstate:            simulation.Distribution{Mean: 0.0029064094738571, StdDev: 0.014660011854061},
		AssetPerformanceData: map[string]simulation.Distribution{
			"INTL-BOND":      simulation.Distribution{Mean: 0.0003, StdDev: 0.0002},
			"US-REALESTATE":  simulation.Distribution{Mean: 0.0004, StdDev: 0.00025},
			"CDN-REALESTATE": simulation.Distribution{Mean: 0.0005, StdDev: 0.00021},
		},
		Parameters: simulation.Parameters{Male: true, Married: true, Retired: false, MaleAge: 29, RetirementAgeMale: 62, FemaleAge: 30, RetirementAgeFemale: 35, ExpensesMultiplier: 1.6, FractionSingleIncome: 65, StartingAssets: 125000, Income: 120000, CurrentTax: 35, SalaryIncrease: 3, IncomeInflationIndex: 20, ExpensesInflationIndex: 100, RetirementIncome: 12000, RetirementExpenses: 80, RetirementTax: 25, LifeInsurance: 250000, IncludeHome: true, HomeValue: 550000, SellHouseIn: 25, NewHomeRelVal: 65},
		Expenses: []simulation.Expense{
			simulation.Expense{Amount: 100, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 25, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 50, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 45, Frequency: "weekly", OneTimeOn: 0, Ends: 1420095599},
			simulation.Expense{Amount: 300, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 3000, Frequency: "annual", OneTimeOn: 0, Ends: 1422773999},
			simulation.Expense{Amount: 5000, Frequency: "annual", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 25000, Frequency: "onetime", OneTimeOn: 1409551199, Ends: 0},
		},
		SelectedPortfolioWeights: map[string]float64{"INTL-BOND": 0.65, "US-REALESTATE": 0.3, "CDN-REALESTATE": 0.05},
	}

	results := s.Simulate()
	fmt.Printf("Length of results: %d\n", len(results))
}
