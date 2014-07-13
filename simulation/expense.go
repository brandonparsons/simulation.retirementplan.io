package simulation

type Expense struct {
	Amount    float64 `json:"amount"`
	Frequency string  `json:"frequency"`
	OneTimeOn int     `json:"onetime_on"`
	Ends      int     `json:"ends"`
}

// isOnetime Determines if an Expense has one-time frequency
// Receiver: Expense
// Params: None
// Returns: bool
func (e *Expense) isOnetime() bool {
	return e.Frequency == "onetime"
}

// isWeekly Determines if an Expense has weekly frequency
// Receiver: Expense
// Params: None
// Returns: bool
func (e *Expense) isWeekly() bool {
	return e.Frequency == "weekly"
}

// isMonthly Determines if an Expense has monthly frequency
// Receiver: Expense
// Params: None
// Returns: bool
func (e *Expense) isMonthly() bool {
	return e.Frequency == "monthly"
}

// isAnnual Determines if an Expense has annual frequency
// Receiver: Expense
// Params: None
// Returns: bool
func (e *Expense) isAnnual() bool {
	return e.Frequency == "annual"
}

// ends Determines if an Expense ends
// Receiver: Expense
// Params: None
// Returns: bool
func (e *Expense) ends() bool {
	return e.Ends != 0
}

// hasEnded Determines if an Expense has already ended given a date
// Receiver: Expense
// Params: currentDate -- int (UTC)
// Returns: bool
func (e *Expense) hasEnded(currentDate int) bool {
	endDate := dateToInt(moveDateToEndOfMonth(dateToTime(e.Ends)))
	current := dateToInt(moveDateToEndOfMonth(dateToTime(currentDate)))
	return endDate < current
}

// isRelevantOnetimeDate Determines if a one-time Expense is triggered in a
// given month.
// Receiver: Expense
// Params: currentDate -- int (UTC)
// Returns: bool
func (e *Expense) isRelevantOnetimeDate(currentDate int) bool {
	oneTime := dateToInt(moveDateToEndOfMonth(dateToTime(e.OneTimeOn)))
	current := dateToInt(moveDateToEndOfMonth(dateToTime(currentDate)))
	return oneTime == current
}

// filterExpenses splits expenses into buckets by frequency
// Params: expenses -- []Expense
// Returns: map[string][]Expense -- keys are frequencies
func filterExpenses(expenses []Expense) map[string][]Expense {
	weeklyExpenses := make([]Expense, 0)
	monthlyExpenses := make([]Expense, 0)
	annualExpenses := make([]Expense, 0)
	onetimeExpenses := make([]Expense, 0)

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

	return map[string][]Expense{
		"weekly":  weeklyExpenses,
		"monthly": monthlyExpenses,
		"annual":  annualExpenses,
		"onetime": onetimeExpenses,
	}
}
