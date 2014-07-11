package simulation

import (
	"encoding/json"
	"log"

	"github.com/kr/pretty"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func prettyPrint(obj interface{}) {
	log.Printf("%# v", pretty.Formatter(obj))
}

type response map[string]interface{}

func (r response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func Meanf(arr []float64) float64 {
	l := len(arr)
	var mean float64
	for i := 0; i < l; i++ {
		mean += (arr[i] - mean) / float64(i+1)
	}
	return mean
}

func Mean(arr []int) float64 {
	l := len(arr)
	var mean float64
	for i := 0; i < l; i++ {
		mean += (float64(arr[i]) - mean) / float64(i+1)
	}
	//fmt.Print(mean)
	return mean
}

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
