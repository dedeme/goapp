// Copyright 07-Nov-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"github.com/dedeme/hxtpl/zenity"
	"github.com/dedeme/hxtpl/newSemicomponent"
)

func main() {

	option := zenity.Menu()
	switch option {
	case "New Semicomponent":
		newSemicomponent.Show()
	default:
		zenity.Msgbox("No option selected.")
	}
}
