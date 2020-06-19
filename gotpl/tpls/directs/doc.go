// Copyright 08-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

//  Go templates.
//
//  Direct forms
//
//  Templates have the following form:
//    ·<key> · <par1> · <par2> ...
//  For example:
//    ·for·i·len(x)
//
//  Directives
//
//    for each map mapTo mapFrom filter sort yes not eqs eqi copyright
//
//  for
//    ·for·end
//      ·for·len(x)
//    ----------
//    ·for·var·end
//      ·for·i·len(x)
//    ----------
//    ·for·var·start·end
//      ·for·i·3·len(x)
//    ----------
//    ·for·var·start·end·step
//      ·for·i·3·len(x)·4
//
//  each
//    ·each·collection
//      ·each·ls
//    ----------
//    ·each·elem·collection
//      ·each·e·ls
//    ----------
//    ·each·index·elem·collection
//      ·for·i·e·ls
//
//  map
//    ·map·source·taget·targetType·function
//      ·map·src·tg·string·fcopy(e + 4)
//
//  mapTo
//    ·mapTo·source·fn
//      ·mapTo·src·json.Ws(e)
//
//  mapFrom
//    ·mapFrom·index·targetType·function
//      ·mapFrom·2·string·e.Rs()
//
//  filter
//    ·filter·source·taget·type·function
//      ·filter·src·tg·string·e > 33
//
//  sort
//    ·sort·iname·type
//      ·sort·sortInt·string
//
//  ·yes
//    ·yes·cond
//    --------
//    if !(cond) {
//      t.Fatal(fail)
//    }
//
//  not
//    ·not·cond
//    --------
//    if cond {
//      t.Fatal(failNot)
//    }
//
//  eqs
//    ·eqs·str1·"abc"
//    --------
//    if r := eqs(str1, "abc"); r != "" {
//      t.Fatal(r)
//    }
//
//  eqi
//    ·eqi·n1·33
//    --------
//    if r := eqi(n1, 33); r != "" {
//      t.Fatal(r)
//    }
//
//  copyright
//    ·copyright
//      ·copyright
package directs
