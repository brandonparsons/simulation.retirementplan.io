package simulation

import (
	"fmt"
	"time"

	goMoment "github.com/jinzhu/now"
)

func TestTimeStuff() {
	goMoment.FirstDayMonday = true
	goMoment.BeginningOfWeek()
	fmt.Println("Beginning of this week:", goMoment.BeginningOfWeek())
	fmt.Println("Beginning of this month:", goMoment.BeginningOfMonth())

	end := goMoment.EndOfMonth()
	next := goMoment.New(end.AddDate(0, 0, 1)).BeginningOfMonth()
	fmt.Println("Beginning of next month:", next)
	fmt.Println("Should be the same as:", time.Date(2014, 7, 1, 0, 0, 0, 0, time.Now().Location()))
	fmt.Println("Unix:", next.Unix())
}
