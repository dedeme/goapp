// Copyright 01-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Program constants.
package cts

const (
	// GENERAL

	// Application name
	AppName = "FMarket"

	// Data directory of FMarket
	FMarketPath = "/dm/dmWeb/app/KtWeb/dmcgi/" + AppName

	// Log path
	LogPath = FMarketPath + "/log.txt"

	// data path from KtMMarket
	MMarketPath = "/home/deme/.dmGoApp/KtMMarket/data"

	// MODEL

	// Investors initial capital for each cycle
	InitialCapital = 100000.0

	// Bet
	Bet = 10000.0

	// Minimun cash to bet
	MinToBet = 11000.0

	// No lost multiplicator
	NoLostMultiplicator = 1.05

	// FLEA

	// Number of fleas per model
	FleasPerModel = 10000

	// Number of fleas in ranking
	FleasInRanking = 40

	// Number of single rankings in rankings data base
	RankingsInDatabase = 10

	// Number of model groups in model statitics
	GroupsInModelRanking = 10

	//
	InheritanceRatio = 1.0 / 201.0
)
