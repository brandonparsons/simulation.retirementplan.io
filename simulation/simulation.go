// Package simulation contains all of the business logic required to simulate
// user's retirement. At the time of writing, this ingests JSON data and outputs
// JSON data.
package simulation

import goStats "github.com/GaryBoone/GoStats/stats"

type simulationResponse []trialResult
type trialResult []simulationTimeStep

type simulationTimeStep struct {
	Assets        float64 `json:"assets"`
	Income        float64 `json:"income"`
	Expenses      float64 `json:"expenses"`
	JsTime        int     `json:"js_time"`
	maleAge       int
	femaleAge     int
	maleAlive     bool
	femaleAlive   bool
	maleRetired   bool
	femaleRetired bool
}

// numberOfMonthsToSimulate determines the number of months the simulation must
// cover, based on the user's ages.
// Params: none
// Returns: integer
func (s *SimulationData) numberOfMonthsToSimulate() int {
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

	return yearsToRun * 12
}

// simulate is the main call for simulations - it runs all of the trials etc.
// Receiver: SimulationData
// Params: none
// Returns: simulationResponse ([]trialResult)
func (s *SimulationData) Simulate() simulationResponse {
	n := s.NumberOfTrials
	trialResults := make([]trialResult, n)

	// This does not change trial-to-trial, do only once.
	timeSteps := s.applyExpenses(s.numberOfMonthsToSimulate())

	type empty struct{}
	notifier := make(chan empty, n)
	for trial := 0; trial < n; trial++ {
		go func(i int) {
			trialResults[i] = s.runIndividualSimulation(timeSteps)
			notifier <- empty{}
		}(trial)
	}

	// Wait for goroutines to finish
	for i := 0; i < n; i++ {
		<-notifier
	}

	return trialResults
}

// runIndividualSimulation is a single loop through the simulation. It is called
// by the `simulate` function
// Receiver: SimulationData
// Params: timeSteps []*timeStep -- prebuilt date steps with expenses applied
// Returns: trialResult ([]simulationTimeStep)
func (s *SimulationData) runIndividualSimulation(timeSteps []*timeStep) trialResult {

	// Copy in data from timeSteps (includes date and expenses)
	trialResult := make([]simulationTimeStep, len(timeSteps))
	for i, v := range timeSteps {
		trialResult[i] = simulationTimeStep{
			Expenses: v.expenses,
			JsTime:   v.date * 1000,
		}
	}

	/* Gather starting data */

	married := s.Parameters.Married
	male := s.Parameters.Male
	maleAge := s.Parameters.MaleAge
	femaleAge := s.Parameters.FemaleAge
	retirementExpenseFactor := s.Parameters.RetirementExpenses

	oneHasAlreadyDied := false // Outside of loop -- using as flag

	assetPerformance := s.generateAssetPerformance(s.numberOfMonthsToSimulate())

	var maleAlive bool
	var femaleAlive bool
	if married {
		maleAlive = true
		femaleAlive = true
	} else {
		if male {
			maleAlive = true
			femaleAlive = false
		} else {
			maleAlive = false
			femaleAlive = true
		}
	}

	trialResult[0].Income = s.Parameters.Income / 12.0 // Monthly

	/* */

	/* Do most of the calculation in one loop */

	for monthIndex := range trialResult {

		// Mortality results: check alive, retirement & dead
		if monthIndex != 0 && monthIndex%12 == 0 {
			// Mortality is tied to age, which only changes every 12 months
			if maleAlive {
				maleAge++
				maleAlive = !maleDiesAt(maleAge)
			}
			if femaleAlive {
				femaleAge++
				femaleAlive = !femaleDiesAt(femaleAge)
			}
		}

		// Handle male age / retired
		trialResult[monthIndex].maleAlive = maleAlive
		trialResult[monthIndex].maleAge = maleAge
		if maleAge >= s.Parameters.RetirementAgeMale {
			trialResult[monthIndex].maleRetired = true
		} else {
			trialResult[monthIndex].maleRetired = false
		}

		// Handle female age / retired
		trialResult[monthIndex].femaleAlive = femaleAlive
		trialResult[monthIndex].femaleAge = femaleAge
		if femaleAge >= s.Parameters.RetirementAgeFemale {
			trialResult[monthIndex].femaleRetired = true
		} else {
			trialResult[monthIndex].femaleRetired = false
		}

		// Apply the retirement expense reduction. If married, and both are
		// retired, or if single, and retired, cut expenses down by the provided
		// factor.
		var applyRetirementExpenseReduction bool
		if married {
			applyRetirementExpenseReduction = trialResult[monthIndex].maleRetired && trialResult[monthIndex].femaleRetired
		} else {
			if male {
				applyRetirementExpenseReduction = trialResult[monthIndex].maleRetired
			} else {
				applyRetirementExpenseReduction = trialResult[monthIndex].femaleRetired
			}
		}
		if applyRetirementExpenseReduction {
			trialResult[monthIndex].Expenses = trialResult[monthIndex].Expenses * (retirementExpenseFactor / 100)
		}

		// Handle death.  Add life insurance if just died, and any timestep
		// where one parter is dead apply expenses reduction. Life insurance
		// portion is relevant if married or single. Expenses reduction is only
		// relevant if married.
		var someoneHasDied bool
		if married {
			someoneHasDied = !(trialResult[monthIndex].maleAlive && trialResult[monthIndex].femaleAlive)
		} else {
			someoneHasDied = (male && !trialResult[monthIndex].maleAlive) || (!male && !trialResult[monthIndex].femaleAlive)
		}

		if someoneHasDied {
			if !oneHasAlreadyDied {
				oneHasAlreadyDied = true
				trialResult[monthIndex].Income += s.Parameters.LifeInsurance
			}

			if s.Parameters.ExpensesMultiplier != 0.0 {
				trialResult[monthIndex].Expenses = trialResult[monthIndex].Expenses / s.Parameters.ExpensesMultiplier
			}
		}
	}

	/* */

	/* Apply income events (retirement changes, salary increases) */

	applyFractionForSingleIncome := false
	retired := false
	haveNotAppliedSingleIncomeFraction := true

	if married && (s.Parameters.MaleAge < s.Parameters.RetirementAgeMale) && (s.Parameters.FemaleAge < s.Parameters.RetirementAgeFemale) && (s.Parameters.FractionSingleIncome != 0) {
		applyFractionForSingleIncome = true
	}

	for monthIndex := range trialResult {
		if monthIndex == 0 {
			continue
		}

		// Bring forward last week's income. The zero-index case was handled
		// above by dividing provided income by 12.
		trialResult[monthIndex].Income = trialResult[monthIndex-1].Income

		// Apply fraction for single income if someone is retired, and we
		// haven't already. `applyFractionForSingleIncome` implies married.
		if applyFractionForSingleIncome && haveNotAppliedSingleIncomeFraction {
			if trialResult[monthIndex].maleRetired || trialResult[monthIndex].femaleRetired {
				haveNotAppliedSingleIncomeFraction = false
				trialResult[monthIndex].Income = trialResult[monthIndex].Income * (s.Parameters.FractionSingleIncome / 100)
			}
		}

		// Evaluate if we are in a fully-retired state - i.e. both retired if
		// married, otherwise person is retired.
		if married {
			if trialResult[monthIndex].maleRetired && trialResult[monthIndex].femaleRetired {
				retired = true
			}
		} else {
			if (male && trialResult[monthIndex].maleRetired) || (!male && trialResult[monthIndex].femaleRetired) {
				retired = true
			}
		}

		if retired {
			trialResult[monthIndex].Income = s.Parameters.RetirementIncome / 12.0
		} else if monthIndex%12 == 0 {
			// Apply salary increase if first month of year
			trialResult[monthIndex].Income = trialResult[monthIndex].Income * (1 + s.Parameters.SalaryIncrease/100)
		}
	}

	/* */

	/* Handle inflation projections */

	// Inflation data comes in as monthly values. Convert to a cumulative basis
	// so it can be cleanly mapped to an array of income/expense values.
	monthlyInflationFactors := make([]float64, len(assetPerformance.inflationPerformance))
	currentCumulativeValue := 1.0
	for monthIndex, monthlyInflation := range assetPerformance.inflationPerformance {
		appliedInflation := currentCumulativeValue * (1 + monthlyInflation)
		monthlyInflationFactors[monthIndex] = appliedInflation
		currentCumulativeValue = appliedInflation
	}

	// Apply inflation to income and expenses.
	for monthIndex := range trialResult {
		// Apply inflation to expenses on a monthly basis (it's not tied to pay
		//raises etc.)
		expensesInflationFactor := (monthlyInflationFactors[monthIndex]-1)*(s.Parameters.ExpensesInflationIndex/100) + 1
		trialResult[monthIndex].Expenses = trialResult[monthIndex].Expenses * expensesInflationFactor

		// Apply inflation to income only on a yearly basis -- assume it is tied
		// to a portion of your raise, rather than your income increased every
		// month.
		incomeInflationFactor := 1.0
		if monthIndex != 0 && monthIndex%12 == 0 {
			incomeInflationFactor = (monthlyInflationFactors[monthIndex]-1)*(s.Parameters.IncomeInflationIndex/100) + 1
		}
		trialResult[monthIndex].Income = trialResult[monthIndex].Income * incomeInflationFactor
	}

	/* */

	// Apply taxes to income. Include varying tax rates during employment, and
	// during retirement.
	for monthIndex := range trialResult {
		applyRetirementTax := false
		if married {
			if trialResult[monthIndex].maleRetired && trialResult[monthIndex].femaleRetired {
				applyRetirementTax = true
			}
		} else {
			if (male && trialResult[monthIndex].maleRetired) || (!male && trialResult[monthIndex].femaleRetired) {
				applyRetirementTax = true
			}
		}

		if applyRetirementTax {
			trialResult[monthIndex].Income = trialResult[monthIndex].Income * (1 - s.Parameters.RetirementTax/100)
		} else {
			trialResult[monthIndex].Income = trialResult[monthIndex].Income * (1 - s.Parameters.CurrentTax/100)
		}
	}

	// If including the home value in the simulation, apply downsize income to
	// the appropriate time step.
	if s.Parameters.IncludeHome {
		houseSaleMonth := s.Parameters.SellHouseIn * 12 // Sell house in provided as year
		relevantRealEstateReturnData := assetPerformance.realEstatePerformance[0:houseSaleMonth]
		futureValueFactor := 1.0
		for _, v := range relevantRealEstateReturnData {
			futureValueFactor = futureValueFactor * (1 + v)
		}
		futureHomeValue := s.Parameters.HomeValue * futureValueFactor
		trialResult[houseSaleMonth].Assets += futureHomeValue * (1 - s.Parameters.NewHomeRelVal/100)
	}

	// If everyone has died, reduce the income and expenses to zero.
	for monthIndex := range trialResult {
		if !trialResult[monthIndex].maleAlive && !trialResult[monthIndex].femaleAlive {
			trialResult[monthIndex].Income = 0
			trialResult[monthIndex].Expenses = 0
		}
	}

	// Run through the timeSteps, and adjust the asset balance based on income
	// shortfall or excess.
	lastPeriodEndingAssets := s.Parameters.StartingAssets
	for monthIndex := range trialResult {
		trialResult[monthIndex].Assets = lastPeriodEndingAssets

		thisMonthPortfolioReturns := lastPeriodEndingAssets * assetPerformance.portfolioPerformance[monthIndex]
		thisMonthIncomeShortfall := trialResult[monthIndex].Expenses - trialResult[monthIndex].Income
		thisMonthAssetImpact := thisMonthPortfolioReturns - thisMonthIncomeShortfall

		lastPeriodEndingAssets += thisMonthAssetImpact
	}

	return trialResult
}
