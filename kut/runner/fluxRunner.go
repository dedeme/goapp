// Copyright 02-Feb-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package runner

import (
	"fmt"
	"github.com/dedeme/kut/expression"
	"github.com/dedeme/kut/heap"
	"github.com/dedeme/kut/heap0"
	"github.com/dedeme/kut/iterator"
	"github.com/dedeme/kut/runner/fail"
	"github.com/dedeme/kut/statement"
)

func runIf(
	stackTrace []*statement.T,
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T,
) (
	withReturn bool, withBreak bool, withContinue bool,
	ret *expression.T, err error,
	stackT []*statement.T,
) {
	ps := st.Value.([]interface{})
	cond := ps[0].(*expression.T)
	var ex *expression.T
	ex, err = Solve(imports, hp0, hps, cond, stackTrace)
	if err == nil {
		switch v := (ex.Value).(type) {
		case bool:
			if v {
				withReturn, withBreak, withContinue, ret, err, stackT =
					RunStat(stackTrace, imports, hp0, hps, ps[1].(*statement.T))
			} else if ps[2] != nil {
				withReturn, withBreak, withContinue, ret, err, stackT =
					RunStat(stackTrace, imports, hp0, hps, ps[2].(*statement.T))
			}
		default:
			err = fail.Type(ex, stackTrace, "bool")
		}
	}
	return
}

func runWhile(
	stackTrace []*statement.T,
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T,
) (withReturn bool, ret *expression.T, err error) {
	ret = expression.MkEmpty()
	ps := st.Value.([]interface{})
	cond := ps[0].(*expression.T)
	wst := ps[1].(*statement.T)
	for {
		if cond != nil {
			var ex *expression.T
			ex, err = Solve(imports, hp0, hps, cond, stackTrace)
			if err != nil {
				break
			}
			stop := false
			switch v := (ex.Value).(type) {
			case bool:
				stop = !v
			default:
				err = fail.Type(ex, stackTrace, "bool")
				stop = true
			}
			if stop {
				break
			}
		}
		var withBreak bool
		withReturn, withBreak, _, ret, err, _ =
			RunStat(stackTrace, imports, hp0, hps, wst)
		if withBreak || withReturn || err != nil {
			break
		}
	}
	return
}

func runFor(
	stackTrace []*statement.T,
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T,
) (withReturn bool, ret *expression.T, err error) {
	ret = expression.MkEmpty()
	ps := st.Value.([]interface{})
	varName := ps[0].(string)
	var a *expression.T
	a, err = Solve(imports, hp0, hps, ps[1].(*expression.T), stackTrace)
	if err != nil {
		return
	}
	fst := ps[2].(*statement.T)
	sts := []*statement.T{statement.NewBuilt(statement.Empty, nil)}
	if fst.Type == statement.Block {
		sts = append(sts, fst.Value.([]*statement.T)...)
	} else {
		sts = append(sts, fst)
	}
	switch v := (a.Value).(type) {
	case []*expression.T:
		for _, e := range v {
			var ex *expression.T
			ex, err = Solve(imports, hp0, hps, e, stackTrace)
			if err != nil {
				break
			}
			sts[0] = statement.NewBuilt(statement.Assign, []*expression.T{
				expression.New(expression.Sym, varName), ex})
			fst2 := statement.NewBuilt(statement.Block, sts)
			var withBreak bool
			withReturn, withBreak, _, ret, err, _ =
				RunStat(stackTrace, imports, hp0, hps, fst2)
			if withBreak || withReturn || err != nil {
				break
			}
		}
	case *iterator.T:
		for v.HasNext() {
			e := v.Next()
			var ex *expression.T
			ex, err = Solve(imports, hp0, hps, e, stackTrace)
			if err != nil {
				break
			}
			sts[0] = statement.NewBuilt(statement.Assign, []*expression.T{
				expression.New(expression.Sym, varName), ex})
			fst2 := statement.NewBuilt(statement.Block, sts)
			var withBreak bool
			withReturn, withBreak, _, ret, err, _ =
				RunStat(stackTrace, imports, hp0, hps, fst2)
			if withBreak || withReturn || err != nil {
				break
			}
		}
	default:
		err = fail.Type(a, stackTrace, "array or <iterator>")
	}
	return
}

func runForIx(
	stackTrace []*statement.T,
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T,
) (withReturn bool, ret *expression.T, err error) {
	ret = expression.MkEmpty()
	ps := st.Value.([]interface{})
	varName := ps[0].(string)
	ixName := ps[1].(string)
	var a *expression.T
	a, err = Solve(imports, hp0, hps, ps[2].(*expression.T), stackTrace)
	if err != nil {
		return
	}
	fst := ps[3].(*statement.T)
	sts := []*statement.T{
		statement.NewBuilt(statement.Empty, nil),
		statement.NewBuilt(statement.Empty, nil),
	}
	if fst.Type == statement.Block {
		sts = append(sts, fst.Value.([]*statement.T)...)
	} else {
		sts = append(sts, fst)
	}
	switch v := (a.Value).(type) {
	case []*expression.T:
		for i, e := range v {
			var ex *expression.T
			ex, err = Solve(imports, hp0, hps, e, stackTrace)
			if err != nil {
				break
			}
			sts[0] = statement.NewBuilt(statement.Assign, []*expression.T{
				expression.New(expression.Sym, varName), ex})
			sts[1] = statement.NewBuilt(statement.Assign, []*expression.T{
				expression.New(expression.Sym, ixName),
				expression.MkFinal(int64(i))})
			fst2 := statement.NewBuilt(statement.Block, sts)
			var withBreak bool
			withReturn, withBreak, _, ret, err, _ =
				RunStat(stackTrace, imports, hp0, hps, fst2)
			if withBreak || withReturn || err != nil {
				break
			}
		}
	case *iterator.T:
		i := 0
		for v.HasNext() {
			e := v.Next()
			var ex *expression.T
			ex, err = Solve(imports, hp0, hps, e, stackTrace)
			if err != nil {
				break
			}
			sts[0] = statement.NewBuilt(statement.Assign, []*expression.T{
				expression.New(expression.Sym, varName), ex})
			sts[1] = statement.NewBuilt(statement.Assign, []*expression.T{
				expression.New(expression.Sym, ixName),
				expression.MkFinal(int64(i))})
			fst2 := statement.NewBuilt(statement.Block, sts)
			var withBreak bool
			withReturn, withBreak, _, ret, err, _ =
				RunStat(stackTrace, imports, hp0, hps, fst2)
			if withBreak || withReturn || err != nil {
				break
			}
			i++
		}
	default:
		err = fail.Type(a, stackTrace, "array or <iterator>")
	}
	return
}

func runForR(
	stackTrace []*statement.T,
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T,
) (withReturn bool, ret *expression.T, err error) {
	ret = expression.MkEmpty()
	ps := st.Value.([]interface{})
	varName := ps[0].(string)
	var start *expression.T
	start, err = Solve(imports, hp0, hps, ps[1].(*expression.T), stackTrace)
	if err != nil {
		return
	}
	var end *expression.T
	end, err = Solve(imports, hp0, hps, ps[2].(*expression.T), stackTrace)
	if err != nil {
		return
	}
	fst := ps[3].(*statement.T)
	sts := []*statement.T{statement.NewBuilt(statement.Empty, nil)}
	if fst.Type == statement.Block {
		sts = append(sts, fst.Value.([]*statement.T)...)
	} else {
		sts = append(sts, fst)
	}
	switch vstart := (start.Value).(type) {
	case int64:
		switch vend := (end.Value).(type) {
		case int64:
			if vend >= vstart {
				for i := vstart; i < vend; i++ {
					sts[0] = statement.NewBuilt(statement.Assign, []*expression.T{
						expression.New(expression.Sym, varName),
						expression.MkFinal(i)})
					fst2 := statement.NewBuilt(statement.Block, sts)
					var withBreak bool
					withReturn, withBreak, _, ret, err, _ =
						RunStat(stackTrace, imports, hp0, hps, fst2)
					if withBreak || withReturn || err != nil {
						break
					}
				}
			} else {
				for i := vstart; i > vend; i-- {
					sts[0] = statement.NewBuilt(statement.Assign, []*expression.T{
						expression.New(expression.Sym, varName),
						expression.MkFinal(i)})
					fst2 := statement.NewBuilt(statement.Block, sts)
					var withBreak bool
					withReturn, withBreak, _, ret, err, _ =
						RunStat(stackTrace, imports, hp0, hps, fst2)
					if withBreak || withReturn || err != nil {
						break
					}
				}
			}
		default:
			err = fail.Type(end, stackTrace, "int")
		}
	default:
		err = fail.Type(start, stackTrace, "int")
	}
	return
}

func runForRS(
	stackTrace []*statement.T,
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T,
) (withReturn bool, ret *expression.T, err error) {
	ret = expression.MkEmpty()
	ps := st.Value.([]interface{})
	varName := ps[0].(string)
	var start *expression.T
	start, err = Solve(imports, hp0, hps, ps[1].(*expression.T), stackTrace)
	if err != nil {
		return
	}
	var end *expression.T
	end, err = Solve(imports, hp0, hps, ps[2].(*expression.T), stackTrace)
	if err != nil {
		return
	}
	var step *expression.T
	step, err = Solve(imports, hp0, hps, ps[3].(*expression.T), stackTrace)
	if err != nil {
		return
	}
	fst := ps[4].(*statement.T)
	sts := []*statement.T{statement.NewBuilt(statement.Empty, nil)}
	if fst.Type == statement.Block {
		sts = append(sts, fst.Value.([]*statement.T)...)
	} else {
		sts = append(sts, fst)
	}
	switch vstart := (start.Value).(type) {
	case int64:
		switch vend := (end.Value).(type) {
		case int64:
			switch step := (step.Value).(type) {
			case int64:
				if step > 0 {
					for i := vstart; i <= vend; i += step {
						sts[0] = statement.NewBuilt(statement.Assign, []*expression.T{
							expression.New(expression.Sym, varName),
							expression.MkFinal(i)})
						fst2 := statement.NewBuilt(statement.Block, sts)
						var withBreak bool
						withReturn, withBreak, _, ret, err, _ =
							RunStat(stackTrace, imports, hp0, hps, fst2)
						if withBreak || withReturn || err != nil {
							break
						}
					}
				} else if step < 0 {
					for i := vstart; i >= vend; i += step {
						sts[0] = statement.NewBuilt(statement.Assign, []*expression.T{
							expression.New(expression.Sym, varName),
							expression.MkFinal(i)})
						fst2 := statement.NewBuilt(statement.Block, sts)
						var withBreak bool
						withReturn, withBreak, _, ret, err, _ =
							RunStat(stackTrace, imports, hp0, hps, fst2)
						if withBreak || withReturn || err != nil {
							break
						}
					}
				} else {
					err = fail.Mk("'step' paramter of 'for' can not be 0", stackTrace)
				}
			default:
				err = fail.Type(end, stackTrace, "int")
			}
		default:
			err = fail.Type(end, stackTrace, "int")
		}
	default:
		err = fail.Type(start, stackTrace, "int")
	}
	return
}

func runSwitch(
	stackTrace []*statement.T,
	imports map[string]int, hp0 heap0.T, hps []heap.T,
	st *statement.T,
) (
	withReturn bool, withBreak bool, withContinue bool,
	ret *expression.T, err error,
	stackT []*statement.T,
) {
	ps0 := st.Value.([]interface{})
	var ex0, ex1 *expression.T
	ex0, err = Solve(imports, hp0, hps, ps0[0].(*expression.T), stackTrace)

	if err == nil {
		ps1 := ps0[1].([][]interface{})
		var r *statement.T
		for _, ps := range ps1 {
			ex := ps[0].(*expression.T)
			if ex.Type == expression.Sym && (ex.Value).(string) == "default" {
				r = ps[1].(*statement.T)
				break
			}
			ex1, err = Solve(
				imports, hp0, hps,
				expression.New(expression.Eq, []*expression.T{ex0, ex}),
				stackTrace)
			if err != nil {
				break
			}
			switch v := (ex1.Value).(type) {
			case bool:
				if v {
					r = ps[1].(*statement.T)
				}
			default:
				err = fail.Type(ex1, stackTrace, "bool")
			}

			if r != nil || err != nil {
				break
			}
		}

		if err == nil {
			if r != nil {
				withReturn, withBreak, withContinue, ret, err, stackT =
					RunStat(stackTrace, imports, hp0, hps, r)
			} else {
				err = fail.Mk(fmt.Sprintf("switch did not catch '%v'", ex0), stackTrace)
			}
		}
	}
	return
}
