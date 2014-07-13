package simulation

import (
	"time"

	goMoment "github.com/jinzhu/now"
)

type timeList []time.Time

type timeStep struct {
	date     int
	expenses float64
}

// dateToTime Converts an integer time (UTC) to a time.Time
// Receiver: None
// Params: intTime -- int (UTC)
// Returns: time.Time (UTC)
func dateToTime(intTime int) time.Time {
	return time.Unix(int64(intTime), 0).UTC()
}

// dateToInt Converts a UTC time.Time to an integer
// Receiver: None
// Params: regularTimeType (time.Time) -- UTC
// Returns: int
func dateToInt(regularTimeType time.Time) int {
	return int(regularTimeType.Unix())
}

// generateMonthsList returns a list of months (end of month), length per arg
// Receiver: None
// Params: numberOfMonths -- integer
// Returns: timeList ([]time.Time)
func generateMonthsList(numberOfMonths int) timeList {
	goMoment.FirstDayMonday = true

	// Do everything in UTC, otherwise you get messed up EndOfMonth's on days
	// where daylight savings time switches
	endofThisMonth := goMoment.New(time.Now().UTC()).EndOfMonth()

	monthsList := make(timeList, numberOfMonths)
	onMonth := endofThisMonth
	monthsList[0] = onMonth
	for i := 1; i < numberOfMonths; i++ {
		onMonth = moveDateToEndOfMonth(onMonth.AddDate(0, 0, 1))
		monthsList[i] = onMonth
	}

	return monthsList
}

// moveDateToEndOfMonth moves any date to the end of its respective month
// Receiver: None
// Params: date time.Time -- original date
// Returns: time.Time -- end of the month
func moveDateToEndOfMonth(date time.Time) time.Time {
	return goMoment.New(date).EndOfMonth()
}

// isYearEnd Determines if a given time, when moved to end of month, is Dec-31
// Receiver: None
// Params: date time.Time -- date of interest
// Returns: bool
func isYearEnd(date time.Time) bool {
	endOfMonth := moveDateToEndOfMonth(date)
	if endOfMonth.Day() == 31 && endOfMonth.Month().String() == "December" {
		return true
	} else {
		return false
	}
}

// applyWeeklyExpenses Takes an array of timesteps (struct{expenses, date}) and
// an array of weekly expenses, and adjusts each steps's expenses value to
// reflect all expenses. Similar to what was passed in - returns array of
// pointers to timesteps so we can keep mapping various expenses to same set
// of timeSteps.
// Receiver: None
// Params: timeSteps -- []*timeStep, weeklyExpenses -- []Expense
// Returns: []*timeStep
func applyWeeklyExpenses(timeSteps []*timeStep, weeklyExpenses []Expense) []*timeStep {
	for _, expense := range weeklyExpenses {
		monthlyAmount := expense.Amount * (52.0 / 12)
		for _, timeStep := range timeSteps {
			if expense.ends() && expense.hasEnded(timeStep.date) {
				break
			}
			timeStep.expenses += monthlyAmount
		}
	}
	return timeSteps
}

// applyMonthlyExpenses Takes an array of timesteps (struct{expenses, date}) and
// an array of monthly expenses, and adjusts each steps's expenses value to
// reflect all expenses. Similar to what was passed in - returns array of
// pointers to timesteps so we can keep mapping various expenses to same set
// of timeSteps.
// Receiver: None
// Params: timeSteps -- []*timeStep, monthlyExpenses -- []Expense
// Returns: []*timeStep
func applyMonthlyExpenses(timeSteps []*timeStep, monthlyExpenses []Expense) []*timeStep {
	for _, expense := range monthlyExpenses {
		for _, timeStep := range timeSteps {
			if expense.ends() && expense.hasEnded(timeStep.date) {
				break
			}
			timeStep.expenses += expense.Amount
		}
	}
	return timeSteps
}

// applyAnnualExpenses Takes an array of timesteps (struct{expenses, date}) and
// an array of annual expenses, and adjusts each steps's expenses value to
// reflect all expenses. Similar to what was passed in - returns array of
// pointers to timesteps so we can keep mapping various expenses to same set
// of timeSteps.
// Receiver: None
// Params: timeSteps -- []*timeStep, annualExpenses -- []Expense
// Returns: []*timeStep
func applyAnnualExpenses(timeSteps []*timeStep, annualExpenses []Expense) []*timeStep {
	for _, expense := range annualExpenses {
		for _, timeStep := range timeSteps {
			if expense.ends() && expense.hasEnded(timeStep.date) {
				break
			}
			if !isYearEnd(dateToTime(timeStep.date)) {
				continue
			}
			timeStep.expenses += expense.Amount
		}
	}
	return timeSteps
}

// applyOnetimeExpenses Takes an array of timesteps (struct{expenses, date}) and
// an array of one-time expenses, and adjusts each steps's expenses value to
// reflect all expenses. Similar to what was passed in - returns array of
// pointers to timesteps so we can keep mapping various expenses to same set
// of timeSteps.
// Receiver: None
// Params: timeSteps -- []*timeStep, onetimeExpenses -- []Expense
// Returns: []*timeStep
func applyOnetimeExpenses(timeSteps []*timeStep, onetimeExpenses []Expense) []*timeStep {
	for _, expense := range onetimeExpenses {
		for _, timeStep := range timeSteps {
			if expense.isRelevantOnetimeDate(timeStep.date) {
				timeStep.expenses += expense.Amount
				break
			}
		}
	}
	return timeSteps
}

// applyExpenses pulls the array of expenses present in the simulationData struct
// and builds out the simulation timeSteps, applying weekly/monthly/annual/onetime
// expenses as appropriate.
// Params: numberOfMonths int -- how many months to simulate
// Returns: []timeStep
func (s *SimulationData) applyExpenses(numberOfMonths int) []*timeStep {
	// This is called ONCE at the beginning of a set of simulation trials. Do
	// **not** do any run-specific calculations here.

	// Initialize the timesteps
	months := generateMonthsList(numberOfMonths)
	timeSteps := make([]*timeStep, numberOfMonths)
	for monthIndex, month := range months {
		step := &timeStep{
			date:     dateToInt(month),
			expenses: 0.0,
		}
		timeSteps[monthIndex] = step
	}

	// Split expenses into buckets
	arrangedExpenses := filterExpenses(s.Expenses)

	// Apply expenses to the timesteps.
	// Doing it this way instead of looping over expenses and applying expenses
	// as there are always going to be more timesteps than expenses.
	timeSteps = applyWeeklyExpenses(timeSteps, arrangedExpenses["weekly"])
	timeSteps = applyMonthlyExpenses(timeSteps, arrangedExpenses["monthly"])
	timeSteps = applyAnnualExpenses(timeSteps, arrangedExpenses["annual"])
	timeSteps = applyOnetimeExpenses(timeSteps, arrangedExpenses["onetime"])

	return timeSteps
}
