package simulation

import (
	"testing"
)

func TestIsOnetime(t *testing.T) {
	t.Skip("Pending...")
}
func TestIsWeekly(t *testing.T) {
	t.Skip("Pending....")
}
func TestIsMonthly(t *testing.T) {
	t.Skip("Pending....")
}
func TestIsAnnual(t *testing.T) {
	t.Skip("Pending....")
}
func TestEnds(t *testing.T) {
	t.Skip("Pending....")
}
func TestHasEnded(t *testing.T) {
	t.Skip("Pending....")
}

func TestIsRelevantOnetimeDate(t *testing.T) {
	jul := dateToInt(moveDateToEndOfMonth(dateToTime(1406851199)))
	aug := dateToInt(moveDateToEndOfMonth(dateToTime(1409529599)))
	sep := dateToInt(moveDateToEndOfMonth(dateToTime(1412207999)))

	e := &expense{Amount: 25000, Frequency: "onetime", OneTimeOn: aug, Ends: 0} // Aug 31, 2014

	var res bool

	res = e.isRelevantOnetimeDate(jul)
	if res {
		t.Error(
			"Failed for:", dateToTime(jul),
			"Got:", res,
		)
	}

	res = e.isRelevantOnetimeDate(aug)
	if !res {
		t.Error(
			"Failed for:", dateToTime(aug),
			"Got:", res,
		)
	}

	res = e.isRelevantOnetimeDate(sep)
	if res {
		t.Error(
			"Failed for:", dateToTime(sep),
			"Got:", res,
		)
	}

}

func TestFilterExpenses(t *testing.T) {
	expenses := []expense{
		expense{Amount: 100, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
		expense{Amount: 25, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
		expense{Amount: 45, Frequency: "weekly", OneTimeOn: 0, Ends: 1406851199}, // Jul 31, 2014
		expense{Amount: 50, Frequency: "weekly", OneTimeOn: 0, Ends: 0},
		expense{Amount: 300, Frequency: "monthly", OneTimeOn: 0, Ends: 0},
		expense{Amount: 3000, Frequency: "annual", OneTimeOn: 0, Ends: 1472687999}, // Aug 31, 2016
		expense{Amount: 5000, Frequency: "annual", OneTimeOn: 0, Ends: 0},
		expense{Amount: 25000, Frequency: "onetime", OneTimeOn: 1409529599, Ends: 0}, // Aug 31, 2014
	}

	ret := filterExpenses(expenses)

	if len(ret["weekly"]) != 4 {
		t.Error("Incorrect mapping weekly!")
	}

	if len(ret["monthly"]) != 1 {
		t.Error("Incorrect mapping weekly!")
	}

	if len(ret["annual"]) != 2 {
		t.Error("Incorrect mapping weekly!")
	}

	if len(ret["onetime"]) != 1 {
		t.Error("Incorrect mapping weekly!")
	}

	if ret["onetime"][0].Amount != 25000 {
		t.Error("Invalid mapping!")
	}
}
