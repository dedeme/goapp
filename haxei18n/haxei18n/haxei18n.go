// Copyright 28-Mar-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/file"
	"github.com/dedeme/ktlib/path"
	"github.com/dedeme/ktlib/str"
	"github.com/dedeme/ktlib/sys"
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
	sys.Print(str.Fmt("%v\n"+
		"Use:\n"+
		"  haxei18n <languages> <roots> <target>\n"+
		"For example:\n"+
		"  haxei18n \\\n"+
		"    \"en:es_ES\" \\\n"+
		"    \"/deme/hx18/src/client/bet:/deme/hx18/src/lib/basic\" \\\n"+
		"    \"/deme/hx18/src/client/bet\"\n"+
		"Note that different languages and roots are separated by ':'\n\n",
		msg,
	))
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
	arr.Each(str.Split(file.Read(source), "\n"), func(l string) {
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
						sys.Print(str.Fmt("Line %d: Sign '=' is not allowed in keys\n", nl))
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
	sys.Println(source)
	if file.IsDirectory(source) {
		for _, fname := range file.Dir(source) {
			extract(keys, path.Cat(source, fname))
		}
		return
	}
	if file.Exists(source) {
		if path.Extension(source) == ".hx" {
			extractFile(keys, source)
		}
		return
	}
	panic(str.Fmt("'%s' not found", source))
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

	dicDir := path.Cat(target, "i18n")
	if !file.Exists(dicDir) {
		file.Mkdir(dicDir)
	}

	if !file.IsDirectory(dicDir) {
		panic(str.Fmt("'%s' is not a directory", dicDir))
	}

	dicPath := path.Cat(dicDir, lang+".txt")
	if !file.Exists(dicPath) {
		file.Write(dicPath, "")
	}

	oldDic := map[string]string{}
	text := file.Read(dicPath)
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
				done.WriteString(str.Fmt("# %s: %d\n", p.file, p.line))
			}
			done.WriteString(str.Fmt("%s = %s\n\n", k, v))
			r = append(r, str.Fmt("\"%s\" => \"%s\"", k, v))
			delete(oldDic, k)
		} else {
			todo.WriteString("# TO DO\n")
			for _, p := range pss {
				todo.WriteString(str.Fmt("# %s: %d\n", p.file, p.line))
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
		orphan.WriteString(str.Fmt(
			"# ORPHAN\n%s = %s\n\n", k, oldDic[k],
		))
	}

	var newText strings.Builder
	newText.WriteString("# File generated by haxei18n.\n\n")
	newText.WriteString(orphan.String())
	newText.WriteString(todo.String())
	newText.WriteString(done.String())
	file.Write(dicPath, newText.String())

	return r
}

func main() {
	if len(sys.Args()) != 4 {
		help("Wrong number of parameters calling haxei18n")
		return
	}

	langs := strings.Split(sys.Args()[1], ":")
	sources := strings.Split(sys.Args()[2], ":")
	target := sys.Args()[3]
	hxtarget := path.Cat(target, "src")

	if !file.IsDirectory(target) {
		panic(str.Fmt("'%s' is not a directory", target))
	}

	if !file.IsDirectory(hxtarget) {
		panic(str.Fmt("'%s' is not a directory", hxtarget))
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
		code.WriteString(str.Fmt("  static var %sDic = [\n    ", lang))
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

	file.Write(path.Cat(hxtarget, "I18n.hx"), code.String())

}
