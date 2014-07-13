package simulation

import (
	"testing"
	"time"
)

func TestDateToTime(t *testing.T) {
	asInt := 1388534400
	toDate := dateToTime(asInt)
	if toDate.Year() != 2014 || (toDate.Month().String() != "January") || (toDate.Day() != 1) {
		t.Error(
			"For", asInt,
			"got", toDate,
		)
	}
}

func TestDateToInt(t *testing.T) {
	now, _ := time.Parse("02/01/2006", "11/07/2014")
	toInt := dateToInt(now)
	if toInt != 1405036800 {
		t.Error(
			"For", now,
			"expected", 1405036800,
			"got", toInt,
		)
	}
}

func TestGenerateMonthsList(t *testing.T) {
	t.Skip("Pending....")
}
func TestMoveDateToEndOfMonth(t *testing.T) {
	t.Skip("Pending....")
}

func TestApplyWeeklyExpenses(t *testing.T) {
	weeklyExpenses := []Expense{
		Expense{Amount: 100, Frequency: "weekly", OneTimeOn: 0, Ends: 1409529590}, // Just before Aug-31-2014. Moves up to end of aug, breaks in sept (doesn't apply to sept)
		Expense{Amount: 25, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
		Expense{Amount: 45, Frequency: "weekly", OneTimeOn: 0, Ends: 1412207999}, // On Sep 31, 2014. Breaks in oct (doesn't apply oct)
		Expense{Amount: 50, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
	}

	timeSteps := []*timeStep{
		&timeStep{date: 1406851199, expenses: 0.0}, // Jul-31-2014
		&timeStep{date: 1409529599, expenses: 0.0}, // Aug-31-2014
		&timeStep{date: 1412207999, expenses: 0.0}, // Sep-31-2014
		&timeStep{date: 1414821599, expenses: 0.0},
		&timeStep{date: 1417417199, expenses: 0.0},
		&timeStep{date: 1420095599, expenses: 0.0},
		&timeStep{date: 1422773999, expenses: 0.0},
	}

	applied := applyWeeklyExpenses(timeSteps, weeklyExpenses)
	factor := (52.0 / 12)

	if applied[0].expenses != (100*factor + 25*factor + 45*factor + 50*factor) {
		t.Error("Messed up first bucket!")
	}

	if applied[1].expenses != (100*factor + 25*factor + 45*factor + 50*factor) {
		t.Error("Messed up second bucket!")
	}

	if applied[2].expenses != (25*factor + 45*factor + 50*factor) {
		t.Error("Messed up third bucket!")
	}

	if applied[3].expenses != (25*factor + 50*factor) {
		t.Error("Messed up fourth bucket!")
	}
}

func TestApplyMonthlyExpenses(t *testing.T) {
	monthlyExpenses := []Expense{
		Expense{Amount: 300, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
		Expense{Amount: 100, Frequency: "monthly", OneTimeOn: 0, Ends: 1409529590}, // Just before Aug-31-2014. Moves up to end of aug, breaks in sept (doesn't apply to sept)
		Expense{Amount: 45, Frequency: "monthly", OneTimeOn: 0, Ends: 1412207999},  // On Sep 31, 2014. Breaks in oct (doesn't apply oct)
	}

	timeSteps := []*timeStep{
		&timeStep{date: 1406851199, expenses: 0.0}, // Jul-31-2014
		&timeStep{date: 1409529599, expenses: 0.0}, // Aug-31-2014
		&timeStep{date: 1412207999, expenses: 0.0}, // Sep-31-2014
		&timeStep{date: 1414821599, expenses: 0.0},
		&timeStep{date: 1417417199, expenses: 0.0},
		&timeStep{date: 1420095599, expenses: 0.0},
		&timeStep{date: 1422773999, expenses: 0.0},
	}

	applied := applyMonthlyExpenses(timeSteps, monthlyExpenses)

	if applied[0].expenses != (300 + 100 + 45) {
		t.Error("Messed up first bucket!")
	}

	if applied[1].expenses != (300 + 100 + 45) {
		t.Error("Messed up second bucket!")
	}

	if applied[2].expenses != (300 + 45) {
		t.Error("Messed up third bucket!")
	}

	if applied[3].expenses != 300 {
		t.Error("Messed up fourth bucket!")
	}
}

func TestApplyAnnualExpenses(t *testing.T) {
	annualExpenses := []Expense{
		Expense{Amount: 1000, Frequency: "annual", OneTimeOn: 0, Ends: 1412207999}, // On Sep 31, 2014
		Expense{Amount: 3000, Frequency: "annual", OneTimeOn: 0, Ends: 1472687999}, // Aug 31, 2016
		Expense{Amount: 5000, Frequency: "annual", OneTimeOn: 0, Ends: 0},
	}

	timeSteps := []*timeStep{
		&timeStep{date: 1409529599, expenses: 0.0}, // Aug-31-2014
		&timeStep{date: 1412207999, expenses: 0.0}, // Sep-31-2014
		&timeStep{date: 1414713599, expenses: 0.0}, // Oct-31-2014
		&timeStep{date: 1417391999, expenses: 0.0}, // Nov-30-2014
		&timeStep{date: 1420070399, expenses: 0.0}, // Dec-31-2014
		&timeStep{date: 1422748799, expenses: 0.0}, // Jan-31-2015
		&timeStep{date: 1451606399, expenses: 0.0}, // Dec-31-2015
		&timeStep{date: 1483228799, expenses: 0.0}, // Dec-31-2016
		&timeStep{date: 1514764799, expenses: 0.0}, // Dec-31-2017
	}

	applied := applyAnnualExpenses(timeSteps, annualExpenses)

	if applied[0].expenses != 0 && applied[1].expenses != 0 && applied[2].expenses != 0 && applied[3].expenses != 0 {
		t.Error("Messed up early buckets!")
	}

	if applied[4].expenses != 8000 {
		t.Error("Messed up first!")
	}

	if applied[5].expenses != 0 {
		t.Error("Messed mid bucket!")
	}

	if applied[6].expenses != 8000 {
		t.Error("Messed up 2015!")
	}

	if applied[7].expenses != 5000 {
		t.Error("Messed up 2016!")
	}

	if applied[8].expenses != 5000 {
		t.Error("Messed up 2017!")
	}
}

func TestApplyOnetimeExpenses(t *testing.T) {
	onetimeExpenses := []Expense{
		Expense{Amount: 100, Frequency: "onetime", OneTimeOn: 1406851199, Ends: 0}, // Jul 31, 2014
		Expense{Amount: 210, Frequency: "onetime", OneTimeOn: 1409529599, Ends: 0}, // Aug 31, 2014
		Expense{Amount: 320, Frequency: "onetime", OneTimeOn: 1409529590, Ends: 0}, // Just before Aug-31-2014. Moves up to end of aug.
		Expense{Amount: 430, Frequency: "onetime", OneTimeOn: 1412207999, Ends: 0}, // On Sep 31, 2014
	}

	timeSteps := []*timeStep{
		&timeStep{date: 1406851199, expenses: 0.0}, // Jul-31-2014
		&timeStep{date: 1409529599, expenses: 0.0}, // Aug-31-2014
		&timeStep{date: 1412207999, expenses: 0.0}, // Sep-31-2014
		&timeStep{date: 1414713599, expenses: 0.0}, // Oct-31-2014
		&timeStep{date: 1417391999, expenses: 0.0}, // Nov-30-2014
		&timeStep{date: 1420070399, expenses: 0.0}, // Dec-31-2014
		&timeStep{date: 1422748799, expenses: 0.0}, // Jan-31-2015
	}

	applied := applyOnetimeExpenses(timeSteps, onetimeExpenses)

	if applied[0].expenses != 100 {
		t.Error("First messed up")
	}

	if applied[1].expenses != 530 {
		t.Error("Second messed up")
	}

	if applied[2].expenses != 430 {
		t.Error(
			"Third messed up",
		)
	}

	if applied[3].expenses != 0 && applied[4].expenses != 0 && applied[5].expenses != 0 && applied[6].expenses != 0 {
		t.Error("Rest messed up")
	}
}

func TestApplyExpenses(t *testing.T) {
	t.Skip("Pending....")
}
