package simulation

import (
	"math"
	"math/rand"
	"sort"

	goStats "github.com/GaryBoone/GoStats/stats"
	goMatrix "github.com/skelterjohn/go.matrix"
)

/*
	All price/asset matrices are of the following form:

	rows: 		# periods
	columns: 	# assets

		ASSET1	ASSET2	ASSET3 ...
	T1	x		y		z
	T2	x		y		z
	T3 	x		y		z
	..
	..
*/

///////////
// Types //
///////////

type AssetPerformanceResults struct {
	RealEstatePerformance ReturnsList
	InflationPerformance  ReturnsList
	PortfolioPerformance  ReturnsList
}

type ReturnResultsByAsset map[string]ReturnsList

type ReturnsList []float64

/////////////////////
// 'Class Methods' //
/////////////////////

// NumberOfMonths determines the number of months the simulation must cover,
// based on the user's ages.
// Params: none
// Returns: integer
func (s *SimulationData) NumberOfMonthsToSimulate() int {
	male := s.Parameters.MaleAge
	female := s.Parameters.FemaleAge

	var yearsToRun int
	if male == 0 {
		yearsToRun = 120 - female
	} else if female == 0 {
		yearsToRun = 120 - male
	} else {
		ages := []float64{float64(s.Parameters.MaleAge), float64(s.Parameters.FemaleAge)}
		yearsToRun = 120 - int(goStats.StatsMin(ages))
	}

	return 12 * yearsToRun
}

// GenerateAssetPerformance generates performance results for real estate,
// inflation and the overall portfolio, and returns as a struct of float arrays.
// Receiver: SimulationData
// Params: numberOfPeriods int -- number of periods to model
// Returns: AssetPerformanceResults
func (s *SimulationData) GenerateAssetPerformance(numberOfPeriods int) AssetPerformanceResults {
	return AssetPerformanceResults{
		RealEstatePerformance: s.RealEstateRandoms(numberOfPeriods),
		InflationPerformance:  s.InflationRandoms(numberOfPeriods),
		PortfolioPerformance:  s.GeneratePortfolioPerformance(numberOfPeriods),
	}
}

// GeneratePortfolioPerformance Consolidates the asset-level data into a single
// value for the user's selected portfolio.
// Receiver: SimulationData
// Params: numberOfPeriods int -- number of periods to model
// Returns: ReturnsList ([]float64)
func (s *SimulationData) GeneratePortfolioPerformance(numberOfPeriods int) ReturnsList {
	assetPerformance := s.GenerateReturns(numberOfPeriods)
	/* assetPerformance is of the following form:
		{
	    	"CDN-REALESTATE": {0.00046232370381282806, 0.0003000276461659901, 0.00039978092717385394},
	    	"INTL-BOND":      {0.0002785337404710855, 0.0003000405280607996, 0.00039872608322449193},
	    	"US-REALESTATE":  {0.0007126042852496324, 0.0003000101077287372, 0.00040144669102363373},
		}
	*/

	portfolioWeights := s.SelectedPortfolioWeights
	/* portfolioWeights is of the following form:
	{
		"INTL-BOND":0.65,
		"US-REALESTATE":0.3,
		"CDN-REALESTATE":0.05,
	}
	*/

	// weightedReturns will be a map[string][]float64 where the keys are the
	// security ID's, and the values will be an array of the weighted returns
	// in each period.
	weightedReturns := ReturnResultsByAsset{}
	for securityId, portfolioWeightOfAsset := range portfolioWeights {
		for _, periodReturn := range assetPerformance[securityId] {
			weightedReturns[securityId] = append(weightedReturns[securityId], periodReturn*portfolioWeightOfAsset)
		}
	}

	// Combine the weightedReturns into a portfolio return in each period this
	// is a straight sum as we have already weighted the asset returns by
	// portfolio weight.
	portfolioReturns := make(ReturnsList, numberOfPeriods)
	for periodIndex := 0; periodIndex < numberOfPeriods; periodIndex++ {
		returnSum := 0.0
		for _, periodReturns := range weightedReturns {
			returnSum += periodReturns[periodIndex]
		}
		portfolioReturns[periodIndex] = returnSum
	}

	return portfolioReturns
}

// GenerateReturns generates performance specifically for assets, and returns
// separately by asset class as a map
// Receiver: SimulationData
// Params: numberOfPeriods int -- number of periods to model
// Returns: ReturnResultsByAsset
func (s *SimulationData) GenerateReturns(numberOfPeriods int) ReturnResultsByAsset {

	assetPerformanceData := s.AssetPerformanceData // map[string]Distribution
	assetClassIds := s.AssetClassIds()             // []string
	numberOfAssets := len(assetClassIds)

	choleskyApplied := s.ApplyCholeskyDecomposition(numberOfPeriods)

	prices := goMatrix.Zeros(numberOfPeriods, numberOfAssets)

	for row := 0; row < prices.Rows(); row++ {
		for column := 0; column < prices.Cols(); column++ {
			var startingValue float64
			if row == 0 {
				startingValue = 0.0
			} else {
				startingValue = prices.Get((row - 1), column)
			}

			assetId := assetClassIds[column] // These are sorted alphabetically
			assetStats := assetPerformanceData[assetId]
			assetMeanReturn := assetStats.Mean
			assetStdDev := assetStats.StdDev

			b := assetMeanReturn - 0.5*math.Pow(assetStdDev, 2)
			c := assetStdDev * choleskyApplied.Get(row, column)

			// Can't do the exp(x) in the same step as you need to use the previous value as a starting price!
			prices.Set(row, column, (startingValue + b + c))
		}

	}

	for row := 0; row < prices.Rows(); row++ {
		for column := 0; column < prices.Cols(); column++ {
			basePrice := prices.Get(row, column)
			expPrice := math.Exp(basePrice)
			prices.Set(row, column, expPrice)
		}
	}

	/* Convert prices to relative returns. */
	// Add T=0, price=1.0
	initialPrices := goMatrix.Ones(1, numberOfAssets)
	augmentedPrices, err := initialPrices.Stack(prices)
	panicIf(err)

	// Create base matrix
	asRelativeReturns := goMatrix.Zeros(numberOfPeriods, numberOfAssets)

	// Each row in the prices matrix is a list of asset prices in each year - NOT
	// the progression of a single asset. We'll need to grab columns for that.
	for column := 0; column < augmentedPrices.Cols(); column++ {
		priceValues := augmentedPrices.GetColVector(column).Array() // []float64
		for periodIndex, periodPrice := range priceValues {
			if periodIndex == 0 {
				// Nothing to do on the first column
				continue
			}
			lastPrice := priceValues[periodIndex-1]
			pctReturn := (periodPrice - lastPrice) / lastPrice
			asRelativeReturns.Set((periodIndex - 1), column, pctReturn)
		}
	}
	/* */

	results := asRelativeReturns.Arrays()    // [][]float64
	resultsByAsset := ReturnResultsByAsset{} // map[string][]float64
	for _, assetClassId := range assetClassIds {
		resultsByAsset[assetClassId] = make([]float64, numberOfPeriods)
	}

	for periodIndex, resultSet := range results {
		for assetIndex, returnResult := range resultSet {
			assetClassId := assetClassIds[assetIndex]
			resultsByAsset[assetClassId][periodIndex] = returnResult
		}
	}

	return resultsByAsset
}

// CholeskyMatrix takes the array of floats provided by the JSON data, and
// converts it to a matrix.
// Receiver: SimulationData
// Params: none
// Returns: *goMatrix.DenseMatrix
func (s *SimulationData) CholeskyMatrix() *goMatrix.DenseMatrix {
	vals := s.CholeskyDecomposition
	noOfVals := float64(len(vals))
	noRows := int(math.Pow(noOfVals, 0.5))
	return goMatrix.MakeDenseMatrix(vals, noRows, noRows)
}

// InflationRandoms generates random inflation performance of a given length
// based on the statistics in the SimulationData struct
// Receiver: SimulationData
// Params: numberOfPeriods int -- number of periods to generate performance for
// Returns: ReturnsList
func (s *SimulationData) InflationRandoms(numberOfPeriods int) ReturnsList {
	return generateRandomsFromDistribution(s.Inflation, numberOfPeriods)
}

// RealEstateRandoms generates random real estate performance of a given length
// based on the statistics in the SimulationData struct
// Receiver: SimulationData
// Params: numberOfPeriods int -- number of periods to generate performance for
// Returns: ReturnsList
func (s *SimulationData) RealEstateRandoms(numberOfPeriods int) ReturnsList {
	return generateRandomsFromDistribution(s.RealEstate, numberOfPeriods)
}

// ApplyCholeskyDecomposition returns a matrix with an applied cholesky
// decomposition - i.e. it creates the random normal matrix, and applies the
// cholesky matrix. The number of assets is implied from the cholesky
// decomposition matrix size.
// Receiver: SimulationData
// Params: numberOfPeriods int -- number of periods to generate performance for
// Returns: *goMatrix.DenseMatrix
func (s *SimulationData) ApplyCholeskyDecomposition(numberOfPeriods int) *goMatrix.DenseMatrix {
	choleskyDecomposition := s.CholeskyMatrix()
	numberOfAssets := choleskyDecomposition.Cols()
	randomValueMatrix := randomNormalsMatrix(numberOfPeriods, numberOfAssets)
	choleskyApplied := zerosMatrix(numberOfPeriods, numberOfAssets)

	for row := 0; row < choleskyApplied.Rows(); row++ {
		for column := 0; column < choleskyApplied.Cols(); column++ {
			answer := 0.0
			if column == 0 {
				answer = randomValueMatrix.Get(row, 0)
			} else {
				for i := 0; i < column; i++ {
					answer += randomValueMatrix.Get(row, i) * choleskyDecomposition.Get(column, i)
				}
			}
			choleskyApplied.Set(row, column, answer)
		}
	}
	return choleskyApplied
}

// AssetClassIds returns the asset class IDs of interest - those that the user
// has invested in.
// Receiver: SimulationData
// Params: None
// Returns: []string
func (s *SimulationData) AssetClassIds() []string {
	mapWithAssetClasses := s.SelectedPortfolioWeights
	assetClassIds := make([]string, len(mapWithAssetClasses))
	i := 0
	for k, _ := range mapWithAssetClasses {
		assetClassIds[i] = k
		i++
	}
	sort.Strings(assetClassIds)

	return assetClassIds
}

///////////////
// Utilities //
///////////////

// generateRandomsFromDistribution is a utility method that will generate a set
// of random values from a given normal distribution
// Params: distribution Distribution -- contains stats
// Params: numberOfPeriods int -- number of periods to generate randoms for
// Returns: []float64
func generateRandomsFromDistribution(distribution Distribution, numberOfPeriods int) []float64 {
	results := make([]float64, numberOfPeriods)
	for i := range results {
		sample := rand.NormFloat64()*distribution.StdDev + distribution.Mean
		results[i] = sample
	}
	return results
}

// randomNormalsMatrix returns a matrix filled with random float64's of a given size
// Params: rows int -- number of rows to fill
// Params: cols int -- number of cols to fill
// Returns: *goMatrix.DenseMatrix
func randomNormalsMatrix(rows, cols int) *goMatrix.DenseMatrix {
	return goMatrix.Normals(rows, cols)
}

// zerosMatrix returns a matrix filled with zeroes of a given size
// Params: rows int -- number of rows to fill
// Params: cols int -- number of cols to fill
// Returns: *goMatrix.DenseMatrix
func zerosMatrix(rows, cols int) *goMatrix.DenseMatrix {
	return goMatrix.Zeros(rows, cols)
}
