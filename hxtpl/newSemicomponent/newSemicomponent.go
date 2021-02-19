// Copyright 07-Nov-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package newSemicomponent

import (
	"github.com/dedeme/golib/file"
	"github.com/dedeme/hxtpl/zenity"
  gpath "path"
  "strings"
  "time"
)

func code(class string) string {
	r := `// Copyright $DATE$ ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

import dm.Domo;
import dm.Ui;
import dm.Ui.Q;

/// $CLASS$ component.
class $CLASS$ {
  /// Object -------------------------------------------------------------------

  final wg = Q("div");

  public function new () {
  }

  /// DOM ----------------------------------------------------------------------

  public function show (parent: Domo): Void {
    parent.removeAll().add(wg);
    update();
  }

  public function update (): Void {
    wg.removeAll().add(Q("p").text(Std.string("$CLASS$")));
  }
}
`
  date := time.Now().Format("02-Jan-2006")
	return strings.ReplaceAll(
    strings.ReplaceAll(r, "$CLASS$", class),
    "$DATE$", date,
  )
}

func create(path string) bool {
	if path[0] == '/' {
		zenity.Msgbox("File path '" + path + "' is not a relative one.")
		return false
	}

	pathHx := path + ".hx"
	if file.Exists(pathHx) {
		zenity.Msgbox("File '" + pathHx + "' already exists.")
		return false
	}

	file.WriteAll(pathHx, code(strings.Title(gpath.Base(path))))
	return true
}

func Show() {
	path := zenity.Input(
		"New Semicomponent",
		"\n\nFile relative path:\n(e.g.: view)\n(e.g.: src/view)",
	)
	switch path {
	case "":
		zenity.Msgbox("File path is missing.")
	default:
		if create(path) {
			zenity.Msgbox("File '" + path + ".hx' created.")
			return
		}
	}
}
