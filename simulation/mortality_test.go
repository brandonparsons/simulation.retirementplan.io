package simulation_test

import (
	"math/rand"
	"testing"
	. "../simulation"
)

func TestReturnsSameValuesMale(t *testing.T) {
	rand.Seed(42)

	first := MaleDiesAt(90)
	second := MaleDiesAt(90)
	third := MaleDiesAt(90)

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

	first := FemaleDiesAt(90)
	second := FemaleDiesAt(90)
	third := FemaleDiesAt(90)

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
		maleDies := MaleDiesAt(121)
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
		femaleDies := FemaleDiesAt(121)
		if !femaleDies {
			allTrue = false
		}
	}

	if !allTrue {
		t.Error("Female should have died every time")
	}
}
