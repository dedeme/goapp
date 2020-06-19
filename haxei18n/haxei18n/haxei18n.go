// Copyright 28-Mar-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"github.com/dedeme/golib/file"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

type pos = struct {
	file string
	line int
}

func newPos(file string, line int) *pos {
	return &pos{file, line}
}

func help(msg string) {
	fmt.Printf("%v\n"+
		"Use:\n"+
		"  haxei18n <languages> <roots> <target>\n"+
		"For example:\n"+
		"  haxei18n \\\n"+
		"    \"en:es_ES\" \\\n"+
		"    \"/deme/hx18/src/client/bet:/deme/hx18/src/lib/basic\" \\\n"+
		"    \"/deme/hx18/src/client/bet\"\n"+
		"Note that different languages and roots are separated by ':'\n\n",
		msg,
	)
}

func extractFile(keys map[string][]*pos, source string) {
	const (
		CODE = iota
		PRECOMMENT
		LCOMMENT
		BCOMMENT
		BCOMMENT2
		QUOTE
		QUOTE2
		SQUOTE
		SQUOTE2
		ENTRY1
		ENTRY2
		EQUOTE
		EQUOTE2
	)

	codeState := func(ch byte) int {
		switch ch {
		case '/':
			return PRECOMMENT
		case '"':
			return QUOTE
		case '\'':
			return SQUOTE
		case '_':
			return ENTRY1
		default:
			return CODE
		}
	}

	var bf strings.Builder
	state := CODE
	nl := 0
	file.Lines(source, func(l string) {
		l = l + "\n"
		nl++
		length := len(l)
		ix := 0
		for ix < length {
			ch := l[ix]
			ix++
			switch state {
			case PRECOMMENT:
				switch ch {
				case '/':
					state = LCOMMENT
				case '*':
					state = BCOMMENT
				default:
					state = codeState(ch)
				}
			case LCOMMENT:
				if ch == '\n' {
					state = CODE
				}
			case BCOMMENT:
				if ch == '*' {
					state = BCOMMENT2
				}
			case BCOMMENT2:
				if ch == '/' {
					state = CODE
				} else {
					if ch == '*' {
						state = BCOMMENT2
					} else {
						state = BCOMMENT
					}
				}
			case QUOTE:
				switch ch {
				case '"':
					state = CODE
				case '\n':
					state = CODE
				case '\\':
					state = QUOTE2
				}
			case SQUOTE:
				switch ch {
				case '\'':
					state = CODE
				case '\n':
					state = CODE
				case '\\':
					state = SQUOTE2
				}
			case QUOTE2:
				state = QUOTE
			case SQUOTE2:
				state = SQUOTE
			case ENTRY1:
				if ch == '(' {
					state = ENTRY2
				} else {
					state = codeState(ch)
				}
			case ENTRY2:
				if ch == '"' {
					state = EQUOTE
				} else {
					state = codeState(ch)
				}
			case EQUOTE:
				switch ch {
				case '\n':
					bf.Reset()
					state = CODE
				case '"':
					var key = bf.String()
					if strings.IndexByte(key, '=') == -1 {
						var ps = newPos(source, nl)
						if pss, ok := keys[key]; ok {
							keys[key] = append(pss, ps)
						} else {
							keys[key] = []*pos{ps}
						}
					} else {
						fmt.Printf("Line %d: Sign '=' is not allowed in keys", nl)
					}
					bf.Reset()
					state = CODE
				case '\\':
					bf.WriteByte(ch)
					state = EQUOTE2
				default:
					bf.WriteByte(ch)
				}
			case EQUOTE2:
				bf.WriteByte(ch)
				state = EQUOTE
			default:
				state = codeState(ch)
			}
		}
	})
}

func extract(keys map[string][]*pos, source string) {
	fmt.Println(source)
	if file.IsDirectory(source) {
		for _, info := range file.List(source) {
			extract(keys, path.Join(source, info.Name()))
		}
		return
	}
	if file.Exists(source) {
		if path.Ext(source) == ".hx" {
			extractFile(keys, source)
		}
		return
	}
	log.Fatalf("'%s' not found", source)
}

func makeDic(
	currentKeys map[string][]*pos, lang string, target string,
) []string {

	rdoc := func(tx string) string {
		return strings.ReplaceAll(
			strings.ReplaceAll(strings.TrimSpace(tx), "\\\"", "\""),
			"\"", "\\\"",
		)
	}

	dicDir := path.Join(target, "i18n")
	if !file.Exists(dicDir) {
		file.Mkdir(dicDir)
	}

	if !file.IsDirectory(dicDir) {
		log.Fatalf("'%s' is not a directory", dicDir)
	}

	dicPath := path.Join(dicDir, lang+".txt")
	if !file.Exists(dicPath) {
		file.WriteAll(dicPath, "")
	}

	oldDic := map[string]string{}
	text := file.ReadAll(dicPath)
	entries := strings.Split(text, "\n")
	for _, e := range entries {
		trim := strings.TrimSpace(e)
		if trim == "" || trim[0] == '#' {
			continue
		}
		kv := strings.SplitN(trim, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := rdoc(kv[0])
		value := rdoc(kv[1])
		if key != "" && value != "" {
			oldDic[key] = value
		}
	}

	var orphan strings.Builder
	var todo strings.Builder
	var done strings.Builder
	var r []string

	var currentKs []string
	for k := range currentKeys {
		currentKs = append(currentKs, k)
	}
	sort.Strings(currentKs)
	for _, k := range currentKs {
		pss := currentKeys[k]
		if v, ok := oldDic[k]; ok {
			for _, p := range pss {
				done.WriteString(fmt.Sprintf("# %s: %d\n", p.file, p.line))
			}
			done.WriteString(fmt.Sprintf("%s = %s\n\n", k, v))
			r = append(r, fmt.Sprintf("\"%s\" => \"%s\"", k, v))
      delete(oldDic, k)
		} else {
			todo.WriteString("# TO DO\n")
			for _, p := range pss {
				todo.WriteString(fmt.Sprintf("# %s: %d\n", p.file, p.line))
			}
			todo.WriteString(k + " = \n\n")
		}
	}

	var oldKs = []string{}
	for k := range oldDic {
		oldKs = append(oldKs, k)
	}
  sort.Strings(oldKs)
  for _, k := range oldKs {
    orphan.WriteString(fmt.Sprintf(
      "# ORPHAN\n%s = %s\n\n", k, oldDic[k],
    ))
  }

  var newText strings.Builder
  newText.WriteString("# File generated by jsi18n.\n\n")
  newText.WriteString(orphan.String())
  newText.WriteString(todo.String())
  newText.WriteString(done.String())
  file.WriteAll(dicPath, newText.String())

	return r
}

func main() {
	if len(os.Args) != 4 {
		help("Wrong number of parameters calling haxei18n")
		return
	}

	langs := strings.Split(os.Args[1], ":")
	sources := strings.Split(os.Args[2], ":")
	target := os.Args[3]
	hxtarget := path.Join(target, "src")

	if !file.IsDirectory(target) {
		log.Fatalf("'%s' is not a directory", target)
	}

	if !file.IsDirectory(hxtarget) {
		log.Fatalf("'%s' is not a directory", hxtarget)
	}

	var currentKeys = map[string][]*pos{}
	for _, source := range sources {
		extract(currentKeys, source)
	}

	var code strings.Builder
	code.WriteString("// Generate by hxi18n. Don't modify\n" +
		"\n" +
		"/// I18n management.\n" +
		"class I18n {\n" +
		"\n")

	for _, lang := range langs {
		code.WriteString(fmt.Sprintf("  static var %sDic = [\n    ", lang))
		dic := makeDic(currentKeys, lang, target)
		code.WriteString(strings.Join(dic, ",\n    "))
		code.WriteString("\n  ];\n\n")
	}

	code.WriteString(
		`  public static var lang(default, null) = "es";

  public static function en (): Void {
    lang = "en";
  }

  public static function es (): Void {
    lang = "es";
  }

  public static function _(key: String): String {
    final dic = lang == "en" ? enDic : esDic;
    return dic.exists(key) ? dic[key] : key;
  }

  public static function _args(key: String, args: Array<String>): String {
    var bf = "";
    final v = _(key);
    var isCode = false;
    for (i in 0...v.length) {
      final ch = v.charAt(i);
      if (isCode) {
        if (ch >= "0" && ch <= "9") bf += args[Std.parseInt(ch)];
        else bf += "%" + ch;
        isCode = false;
      } else if (ch == "%") {
        isCode = true;
      } else {
        bf += ch;
      }
    }
    return bf;
  }

}
`)

	file.WriteAll(path.Join(hxtarget, "I18n.hx"), code.String())

}
