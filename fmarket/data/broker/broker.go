// Copyright 03-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Accounting broker.
package broker

// Returns total fees of a buy or sell operation.
//    amount: Operation amount.
func Fees(amount float64) float64 {
	var brk, market float64
	if amount > 50000 {
		brk = amount * 0.001
	} else {
		brk = 9.75
	}

	market = amount * 0.00003
	market += 0.11 // Execution fee

	return brk + market
}

// Returns net cost of operation (cost + fees).
//    stocks: Stocks number.
//    price  ; Stocks price.
func Buy(stocks int, price float64) float64 {
	amount := float64(stocks) * price
	tobin := amount * 0.002

	return amount + Fees(amount) + tobin
}

// Returns net incomes of operation (incomes - fees).
//    stocks: Stocks number.
//    price : Stocks price.
func Sell(stocks int, price float64) float64 {
	amount := float64(stocks) * price
	return amount - Fees(amount)
}
