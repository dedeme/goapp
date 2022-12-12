// Copyright 03-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Investor model.
package model

import (
	"github.com/dedeme/fmarket/data/broker"
	"github.com/dedeme/fmarket/data/cts"
	"github.com/dedeme/fmarket/data/quotes"
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/js"
	"github.com/dedeme/ktlib/math"
)

type T struct {
	// Short name
	id string
	// Parameter names
	paramNames []string
	// Paramteres decimals to convert to string
	paramDecs []int
	// Minimum values of parameters (inclusive)
	paramMins []float64
	// Maximum values of parameters (exclusive)
	paramMaxs []float64
	// Function to calculate operations.
	//    closes: Normalized closes in matrix 'dates x cos' ordered from before to after.
	//    params: Values to calculate.
	//    action: Function called after calculate 'refs'
	//            When closes[i] >= refs[i], position is bought.
	//            When closes[i] < refs[i], position is sold.
	//            Params:
	//              closes: Last closes. One for each company.
	//              refs  : Last references. One for each company.
	calc func(
		closes [][]float64,
		params []float64,
		action func(closes []float64, refs []float64),
	)
}

func (m *T) Id() string {
	return m.id
}

// Returns random parameters of a model
func (m *T) Mutation() []float64 {
	r := make([]float64, len(m.paramMins))
	for i := range r {
		r[i] = m.paramMins[i] + (m.paramMaxs[i]-m.paramMins[i])*math.Rnd()
	}
	return r
}

// Returns assets resulting of simulation.
//  qs: Quotes for simulation.
//  params: Model parameters.
func (md *T) Simulation(qs *quotes.T, params []float64) float64 {
	nks := qs.Cos
	nCos := len(qs.Cos)
	opens := qs.Opens
	closes := qs.Closes
	maxs := qs.Maxs
	cashIn := cts.InitialCapital
	withdrawal := 0.0

	prfCashs := make([]float64, nCos)
	toSells := make([]bool, nCos)
	for i := 0; i < nCos; i++ {
		toSells[i] = true
		prfCashs[i] = cts.Bet
	}
	toDos := make([]bool, nCos)
	stockss := make([]int, nCos)
	prfStockss := make([]int, nCos)
	prices := make([]float64, nCos)
	prfPrices := make([]float64, nCos)
	coughts := make([]bool, nCos)
	prfCoughts := make([]bool, nCos)
	ix := 0

	md.calc(closes, params, func(cs, rs []float64) {
		os := opens[ix]
		mxs := maxs[ix]

		assets := 0.0
		for i := range nks {
			op := os[i]
			cl := cs[i]
			rf := rs[i]
			toSell := toSells[i]
			toDo := toDos[i]

			if toDo {
				if toSell { // there is buy order.
					if !prfCoughts[i] {
						prfCash := prfCashs[i]
						stocks := int((prfCash - broker.Fees(prfCash)) / op)
						cost := broker.Buy(stocks, op)
						for cost > prfCash {
							stocks--
							cost = broker.Buy(stocks, op)
						}
						prfStockss[i] = stocks
						prfCashs[i] -= cost
						prfPrices[i] = op
					}
					if cashIn > cts.MinToBet && !coughts[i] {
						stocks := int(cts.Bet / op)
						stockss[i] = stocks
						cashIn -= broker.Buy(stocks, op)
						prices[i] = op
					}
				} else {
					stocks := stockss[i]
					if stocks > 0 && !coughts[i] {
						if op > prices[i]*cts.NoLostMultiplicator {
							cashIn += broker.Sell(stocks, op)
							stockss[i] = 0
						} else {
							coughts[i] = true
						}
					}
					stocks = prfStockss[i]
					if stocks > 0 && !prfCoughts[i] {
						if op > prfPrices[i]*cts.NoLostMultiplicator {
							prfCashs[i] += broker.Sell(stocks, op)
							prfStockss[i] = 0
						} else {
							prfCoughts[i] = true
						}
					}
				}

				toDos[i] = false
			}

			if coughts[i] {
				price := prices[i] * cts.NoLostMultiplicator
				if mxs[i] > price {
					stocks := stockss[i]
					cashIn += broker.Sell(stocks, price)
					stockss[i] = 0
					coughts[i] = false
				}
			}

			if prfCoughts[i] {
				price := prfPrices[i] * cts.NoLostMultiplicator
				if mxs[i] > price {
					stocks := prfStockss[i]
					prfCashs[i] += broker.Sell(stocks, price)
					prfStockss[i] = 0
					prfCoughts[i] = false
				}
			}

			stks := stockss[i]
			if stks > 0 {
				assets += broker.Sell(stks, cl)
			}

			if toSell {
				if rf > cl {
					toDos[i] = true
					toSells[i] = false
				}
			} else if rf < cl {
				toDos[i] = true
				toSells[i] = true
			}
		}

		total := cashIn + assets
		if total > cts.InitialCapital+cts.Bet+cts.Bet {
			dif := total - cts.InitialCapital - cts.Bet
			securAmount := cts.MinToBet - cts.Bet
			withdraw := -1.0
			if cashIn > dif+securAmount {
				withdraw = dif
			} else if cashIn > cts.MinToBet {
				withdraw = math.Floor((cashIn-securAmount)/cts.Bet) * cts.Bet
			}
			if withdraw > 0 {
				withdrawal += withdraw
				cashIn -= withdraw
			}
		}

		ix++
	})

	cash := cashIn + withdrawal

	lastCloses := closes[len(closes)-1]
	stocksValue := 0.0
	for i := range lastCloses {
		stocks := stockss[i]
		if stocks > 0 {
			stocksValue += float64(stocks) * lastCloses[i]
		}
	}

	return cash + stocksValue
}

// NOTE: Parameter m.Calc is not serialized.
func ToJs(m *T) string {
	return js.Wa([]string{
		js.Ws(m.id),
		js.Wa(arr.Map(m.paramNames, js.Ws)),
		js.Wa(arr.Map(m.paramDecs, js.Wi)),
		js.Wa(arr.Map(m.paramMins, js.Wd)),
		js.Wa(arr.Map(m.paramMaxs, js.Wd)),
	})
}

// Returns the model with identifier 'id'.
//
// If no model is found, panic is raised.
func FromId(id string) *T {
	md, ok := arr.Find(List(), func(m *T) bool {
		return m.id == id
	})
	if !ok {
		panic("Model '" + id + "' not found")
	}
	return md
}

// Returns list of model identifiers
func List() []*T {
	return []*T{
		ApprNew(),
		EaNew(),
		Ea2New(),
		MaNew(),
		MmNew(),
		QfixNew(),
		QmobNew(),
	}
}
