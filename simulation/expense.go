package simulation

type expense struct {
	Amount    float64 `json:"amount"`
	Frequency string  `json:"frequency"`
	OneTimeOn int     `json:"onetime_on"`
	Ends      int     `json:"ends"`
}

// isOnetime Determines if an Expense has one-time frequency
// Receiver: Expense
// Params: None
// Returns: bool
func (e *expense) isOnetime() bool {
	return e.Frequency == "onetime"
}

// isWeekly Determines if an Expense has weekly frequency
// Receiver: Expense
// Params: None
// Returns: bool
func (e *expense) isWeekly() bool {
	return e.Frequency == "weekly"
}

// isMonthly Determines if an Expense has monthly frequency
// Receiver: Expense
// Params: None
// Returns: bool
func (e *expense) isMonthly() bool {
	return e.Frequency == "monthly"
}

// isAnnual Determines if an Expense has annual frequency
// Receiver: Expense
// Params: None
// Returns: bool
func (e *expense) isAnnual() bool {
	return e.Frequency == "annual"
}

// ends Determines if an Expense ends
// Receiver: Expense
// Params: None
// Returns: bool
func (e *expense) ends() bool {
	return e.Ends != 0
}

// hasEnded Determines if an Expense has already ended given a date
// Receiver: Expense
// Params: currentDate -- int (UTC)
// Returns: bool
func (e *expense) hasEnded(currentDate int) bool {
	endDate := dateToInt(moveDateToEndOfMonth(dateToTime(e.Ends)))
	current := dateToInt(moveDateToEndOfMonth(dateToTime(currentDate)))
	return endDate < current
}

// isRelevantOnetimeDate Determines if a one-time Expense is triggered in a
// given month.
// Receiver: Expense
// Params: currentDate -- int (UTC)
// Returns: bool
func (e *expense) isRelevantOnetimeDate(currentDate int) bool {
	oneTime := dateToInt(moveDateToEndOfMonth(dateToTime(e.OneTimeOn)))
	current := dateToInt(moveDateToEndOfMonth(dateToTime(currentDate)))
	return oneTime == current
}

// filterExpenses splits expenses into buckets by frequency
// Params: expenses -- []expense
// Returns: map[string][]expense -- keys are frequencies
func filterExpenses(expenses []expense) map[string][]expense {
	weeklyExpenses := make([]expense, 0)
	monthlyExpenses := make([]expense, 0)
	annualExpenses := make([]expense, 0)
	onetimeExpenses := make([]expense, 0)

	for _, expense := range expenses {
		if expense.isWeekly() {
			weeklyExpenses = append(weeklyExpenses, expense)
		} else if expense.isMonthly() {
			monthlyExpenses = append(monthlyExpenses, expense)
		} else if expense.isAnnual() {
			annualExpenses = append(annualExpenses, expense)
		} else if expense.isOnetime() {
			onetimeExpenses = append(onetimeExpenses, expense)
		} else {
			panic("Invalid expense frequency.")
		}
	}

	return map[string][]expense{
		"weekly":  weeklyExpenses,
		"monthly": monthlyExpenses,
		"annual":  annualExpenses,
		"onetime": onetimeExpenses,
	}
}
