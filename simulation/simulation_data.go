package simulation

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

type simulationTimeStep struct {
	assets        float64
	income        float64
	expenses      float64
	dateInt       int
	maleAge       int
	femaleAge     int
	maleAlive     bool
	femaleAlive   bool
	maleRetired   bool
	femaleRetired bool
}

// runIndividualSimulation is a single loop through the simulation. It is called
// by the `simulate` function
// Receiver: SimulationData
// Params: timeSteps []*timeStep -- prebuilt date steps with expenses applied
// Params: numberOfMonthsToSimulate -- int
// Returns: []simulationTimeStep
func (s *SimulationData) runIndividualSimulation(timeSteps []*timeStep, numberOfMonthsToSimulate int) []simulationTimeStep {

	// Copy in data from timeSteps (includes date and expenses)
	trialResult := make([]simulationTimeStep, len(timeSteps))
	for i, v := range timeSteps {
		trialResult[i] = simulationTimeStep{
			expenses: v.expenses,
			dateInt:  v.date, // Used to be * 1000 for javascript, but we will munge this in ruby as we need string formatted dates
		}
	}

	/* Gather starting data */

	married := s.Parameters.Married
	male := s.Parameters.Male
	maleAge := s.Parameters.MaleAge
	femaleAge := s.Parameters.FemaleAge
	retirementExpenseFactor := s.Parameters.RetirementExpenses

	oneHasAlreadyDied := false // Outside of loop -- using as flag

	assetPerformance := s.generateAssetPerformance(numberOfMonthsToSimulate)

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

	trialResult[0].income = s.Parameters.Income / 12.0 // Monthly

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
			trialResult[monthIndex].expenses = trialResult[monthIndex].expenses * (retirementExpenseFactor / 100)
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
				trialResult[monthIndex].income += s.Parameters.LifeInsurance
			}

			if s.Parameters.ExpensesMultiplier != 0.0 {
				trialResult[monthIndex].expenses = trialResult[monthIndex].expenses / s.Parameters.ExpensesMultiplier
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
		trialResult[monthIndex].income = trialResult[monthIndex-1].income

		// Apply fraction for single income if someone is retired, and we
		// haven't already. `applyFractionForSingleIncome` implies married.
		if applyFractionForSingleIncome && haveNotAppliedSingleIncomeFraction {
			if trialResult[monthIndex].maleRetired || trialResult[monthIndex].femaleRetired {
				haveNotAppliedSingleIncomeFraction = false
				trialResult[monthIndex].income = trialResult[monthIndex].income * (s.Parameters.FractionSingleIncome / 100)
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
			trialResult[monthIndex].income = s.Parameters.RetirementIncome / 12.0
		} else if monthIndex%12 == 0 {
			// Apply salary increase if first month of year
			trialResult[monthIndex].income = trialResult[monthIndex].income * (1 + s.Parameters.SalaryIncrease/100)
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
		trialResult[monthIndex].expenses = trialResult[monthIndex].expenses * expensesInflationFactor

		// Apply inflation to income only on a yearly basis -- assume it is tied
		// to a portion of your raise, rather than your income increased every
		// month.
		incomeInflationFactor := 1.0
		if monthIndex != 0 && monthIndex%12 == 0 {
			incomeInflationFactor = (monthlyInflationFactors[monthIndex]-1)*(s.Parameters.IncomeInflationIndex/100) + 1
		}
		trialResult[monthIndex].income = trialResult[monthIndex].income * incomeInflationFactor
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
			trialResult[monthIndex].income = trialResult[monthIndex].income * (1 - s.Parameters.RetirementTax/100)
		} else {
			trialResult[monthIndex].income = trialResult[monthIndex].income * (1 - s.Parameters.CurrentTax/100)
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
		trialResult[houseSaleMonth].assets += futureHomeValue * (1 - s.Parameters.NewHomeRelVal/100)
	}

	// If everyone has died, reduce the income and expenses to zero.
	for monthIndex := range trialResult {
		if !trialResult[monthIndex].maleAlive && !trialResult[monthIndex].femaleAlive {
			trialResult[monthIndex].income = 0
			trialResult[monthIndex].expenses = 0
		}
	}

	// Run through the timeSteps, and adjust the asset balance based on income
	// shortfall or excess.
	lastPeriodEndingAssets := s.Parameters.StartingAssets
	for monthIndex := range trialResult {
		trialResult[monthIndex].assets = lastPeriodEndingAssets

		thisMonthPortfolioReturns := lastPeriodEndingAssets * assetPerformance.portfolioPerformance[monthIndex]
		thisMonthIncomeShortfall := trialResult[monthIndex].expenses - trialResult[monthIndex].income
		thisMonthAssetImpact := thisMonthPortfolioReturns - thisMonthIncomeShortfall

		lastPeriodEndingAssets += thisMonthAssetImpact
	}

	return trialResult
}
