package simulation

import (
	"math/rand"
	"testing"
)

func TestReturnsSameValuesMale(t *testing.T) {
	rand.Seed(42)

	first := maleDiesAt(90)
	second := maleDiesAt(90)
	third := maleDiesAt(90)

	if first {
		t.Error("Was expecting false")
	}

	if !second {
		t.Error("Was expecting true")
	}

	if third {
		t.Error("Was expecting false")
	}
}

func TestReturnsSameValuesFemale(t *testing.T) {
	rand.Seed(42)

	first := femaleDiesAt(90)
	second := femaleDiesAt(90)
	third := femaleDiesAt(90)

	if first {
		t.Error("Was expecting false")
	}

	if !second {
		t.Error("Was expecting true")
	}

	if third {
		t.Error("Was expecting false")
	}
}

func TestAlwaysTrueAt120ForMale(t *testing.T) {
	allTrue := true
	for i := 0; i < 100; i++ {
		maleDies := maleDiesAt(121)
		if !maleDies {
			allTrue = false
		}
	}

	if !allTrue {
		t.Error("Male should have died every time")
	}
}

func TestAlwaysTrueAt120ForFemale(t *testing.T) {
	allTrue := true
	for i := 0; i < 100; i++ {
		femaleDies := femaleDiesAt(121)
		if !femaleDies {
			allTrue = false
		}
	}

	if !allTrue {
		t.Error("Female should have died every time")
	}
}
