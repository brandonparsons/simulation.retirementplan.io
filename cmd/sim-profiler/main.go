package main

import (
	"fmt"

	"bitbucket.org/retirementplanio/simulation.retirementplan.io/simulation"
	"github.com/davecheney/profile"
)

//////////
// Main //
//////////

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	// // LOTS OF TRIALS
	// s := simulation.SimulationData{
	// 	NumberOfTrials:        5000,
	// 	CholeskyDecomposition: []float64{0.0094794922, 0, 0, -7.36e-05, 0.0055677999, 0, 0.0050681903, -0.0004821709, 0.013367741},
	// 	Inflation:             simulation.Distribution{Mean: 0.00046346514957523, StdDev: 0.00024792742828969},
	// 	RealEstate:            simulation.Distribution{Mean: 0.0029064094738571, StdDev: 0.014660011854061},
	// 	AssetPerformanceData: map[string]simulation.Distribution{
	// 		"INTL-BOND":      simulation.Distribution{Mean: 0.0003, StdDev: 0.0002},
	// 		"US-REALESTATE":  simulation.Distribution{Mean: 0.0004, StdDev: 0.00025},
	// 		"CDN-REALESTATE": simulation.Distribution{Mean: 0.0005, StdDev: 0.00021},
	// 	},
	// 	Parameters: simulation.Parameters{Male: true, Married: true, Retired: false, MaleAge: 29, RetirementAgeMale: 62, FemaleAge: 30, RetirementAgeFemale: 35, ExpensesMultiplier: 1.6, FractionSingleIncome: 65, StartingAssets: 125000, Income: 120000, CurrentTax: 35, SalaryIncrease: 3, IncomeInflationIndex: 20, ExpensesInflationIndex: 100, RetirementIncome: 12000, RetirementExpenses: 80, RetirementTax: 25, LifeInsurance: 250000, IncludeHome: true, HomeValue: 550000, SellHouseIn: 25, NewHomeRelVal: 65},
	// 	Expenses: []simulation.Expense{
	// 		simulation.Expense{Amount: 100, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 25, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 50, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 45, Frequency: "weekly", OneTimeOn: 0, Ends: 1420095599},
	// 		simulation.Expense{Amount: 300, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 3000, Frequency: "annual", OneTimeOn: 0, Ends: 1422773999},
	// 		simulation.Expense{Amount: 5000, Frequency: "annual", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 25000, Frequency: "onetime", OneTimeOn: 1409551199, Ends: 0},
	// 	},
	// 	SelectedPortfolioWeights: map[string]float64{"INTL-BOND": 0.65, "US-REALESTATE": 0.3, "CDN-REALESTATE": 0.05},
	// }

	// // REASONABLE NUMBER OF TRIALS
	// s := simulation.SimulationData{
	// 	NumberOfTrials:        1000,
	// 	CholeskyDecomposition: []float64{0.0206140002, 0, 0, 0, 0.0058743434, 0.0107500299, 0, 0, 0.0012367294, 0.003156581, 0.0172708088, 0, 0.0086516523, 0.0062800071, 0.008154059, 0.0204417622},
	// 	Inflation:             simulation.Distribution{Mean: 0.00141579416791443, StdDev: 0.00300832469830286},
	// 	RealEstate:            simulation.Distribution{Mean: -0.00347709477116344, StdDev: 0.0097874440308587},
	// 	AssetPerformanceData: map[string]simulation.Distribution{
	// 		"INTL-REALESTATE":  simulation.Distribution{Mean: -0.0025180101, StdDev: 0.0608277917},
	// 		"US-LGCAP-STOCK":   simulation.Distribution{Mean: 0.0041452698, StdDev: 0.0441003288},
	// 		"US-LONG-GOV-BOND": simulation.Distribution{Mean: 0.0012606674, StdDev: 0.0301729388},
	// 		"US-REALESTATE":    simulation.Distribution{Mean: 0.0017808464, StdDev: 0.0687683985},
	// 	},
	// 	Parameters: simulation.Parameters{Male: true, Married: false, Retired: false, MaleAge: 30, RetirementAgeMale: 50, FemaleAge: 0, RetirementAgeFemale: 0, ExpensesMultiplier: 0, FractionSingleIncome: 0, StartingAssets: 1.565512e+06, Income: 115000, CurrentTax: 35, SalaryIncrease: 3, IncomeInflationIndex: 0, ExpensesInflationIndex: 100, RetirementIncome: 15000, RetirementExpenses: 100, RetirementTax: 35, LifeInsurance: 155000, IncludeHome: false, HomeValue: 0, SellHouseIn: 0, NewHomeRelVal: 0},
	// 	Expenses: []simulation.Expense{
	// 		simulation.Expense{Amount: 250000, Frequency: "onetime", OneTimeOn: 1420095600, Ends: 0},
	// 		simulation.Expense{Amount: 400, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 250, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 300, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 100, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 100, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 200, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 2000, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 500, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
	// 		simulation.Expense{Amount: 40, Frequency: "weekly", OneTimeOn: 0, Ends: 1408069141},
	// 	},
	// 	SelectedPortfolioWeights: map[string]float64{"INTL-REALESTATE": 0, "US-LGCAP-STOCK": 0.646172292868484, "US-LONG-GOV-BOND": 0.353827707131516, "US-REALESTATE": 0},
	// }

	// // RUNS OUT OF MONEY
	s := simulation.SimulationData{
		NumberOfTrials:        5,
		CholeskyDecomposition: []float64{0.0576625891, 0, 0, 0, 0, 0.0062658019, 0.0175730542, 0, 0, 0, 0.0056351399, 0.0051007659, 0.0031372137, 0, 0, -0.0046404388, 0.0058461108, 0.0015728539, 0.005015221, 0, 0.054254316, -0.0004589824, 0.0001774192, 0.0012348558, 0.0121795492},
		Inflation:             simulation.Distribution{Mean: 0.001387958865714121, StdDev: 0.002999574074911733},
		RealEstate:            simulation.Distribution{Mean: -0.003477094771163442, StdDev: 0.009787444030858702},
		AssetPerformanceData: map[string]simulation.Distribution{
			"US-MED-GOV-BOND":  simulation.Distribution{Mean: 0.0064879833, StdDev: 0.0136574658},
			"US-SMCAP-STOCK":   simulation.Distribution{Mean: 0.005830665, StdDev: 0.0585631555},
			"CDN-LONG-BOND":    simulation.Distribution{Mean: 0.0059680241, StdDev: 0.0576625891},
			"INTL-BOND":        simulation.Distribution{Mean: 0.006604477, StdDev: 0.0216913127},
			"US-MED-CORP-BOND": simulation.Distribution{Mean: 0.0063501862, StdDev: 0.0138674297},
		},
		Parameters: simulation.Parameters{Male: true, Married: true, Retired: false, MaleAge: 30, RetirementAgeMale: 65, FemaleAge: 30, RetirementAgeFemale: 31, ExpensesMultiplier: 1.6, FractionSingleIncome: 65, StartingAssets: 150000, Income: 200000, CurrentTax: 45, SalaryIncrease: 3, IncomeInflationIndex: 0, ExpensesInflationIndex: 100, RetirementIncome: 1000, RetirementExpenses: 100, RetirementTax: 35, LifeInsurance: 25000, IncludeHome: false, HomeValue: 0, SellHouseIn: 0, NewHomeRelVal: 0},
		Expenses: []simulation.Expense{
			{Amount: 25000, Frequency: "onetime", OneTimeOn: 1446001252, Ends: 0},
			simulation.Expense{Amount: 400, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 250, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 300, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 100, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 100, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 200, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 2000, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 500, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
			simulation.Expense{Amount: 40, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
		},
		SelectedPortfolioWeights: map[string]float64{"CDN-LONG-BOND": 0, "INTL-BOND": 0.491, "US-MED-CORP-BOND": 0.1608, "US-MED-GOV-BOND": 0.3483, "US-SMCAP-STOCK": 0},
	}

	results := simulation.Simulate(&s)
	fmt.Printf("Length of results: %d\n", len(results))
}
