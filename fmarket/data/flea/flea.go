// Copyright 05-Dic-2022 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Flea data
package flea

import (
	"github.com/dedeme/fmarket/data/cts"
	"github.com/dedeme/fmarket/data/model"
	"github.com/dedeme/fmarket/data/quotes"
	"github.com/dedeme/ktlib/arr"
	"github.com/dedeme/ktlib/js"
	"github.com/dedeme/ktlib/math"
)

type FmodelT struct {
	// Model identifier
	modelId string
	// Parameter values
	params []float64
}

func NewFmodel(modelId string, params []float64) *FmodelT {
	return &FmodelT{modelId, params}
}

func (fm *FmodelT) fmodelEq(fm2 *FmodelT) bool {
	return fm.modelId == fm2.modelId &&
		arr.Eqf(fm.params, fm2.params, func(f1, f2 float64) bool {
			return math.Eq(f1, f2, 0.000001)
		})
}

func fmodelToJs(md *FmodelT) string {
	return js.Wa([]string{
		js.Ws(md.modelId),
		js.Wa(arr.Map(md.params, func(value float64) string {
			return js.WdDec(value, 4)
		})),
	})
}

func fmodelFromJs(j string) *FmodelT {
	a := js.Ra(j)
	return NewFmodel(
		js.Rs(a[0]),
		arr.Map(js.Ra(a[1]), js.Rd),
	)
}

// Flea data
type T struct {
	// Identifier
	id int64
	// Creation cycle
	cycle int
	// 'true' if flea is male
	isMale bool
	// Models data(three different models)
	models []*FmodelT
	// Assets after evaluation of models
	assets float64
}

func New(
	id int64, cycle int, isMale bool, models []*FmodelT, assets float64,
) *T {
	if models[0].modelId < models[1].modelId {
		if models[0].modelId < models[2].modelId {
			if models[2].modelId < models[1].modelId {
				models[1], models[2] = models[2], models[1]
			}
		} else {
			models[0], models[1], models[2] = models[2], models[0], models[1]
		}
	} else {
		if models[1].modelId < models[2].modelId {
			if models[0].modelId < models[2].modelId {
				models[0], models[1], models[2] = models[1], models[0], models[2]
			} else {
				models[0], models[1], models[2] = models[1], models[2], models[0]
			}
		} else {
			models[0], models[1], models[2] = models[2], models[1], models[0]
		}
	}

	return &T{id, cycle, isMale, models, assets}
}

func (f *T) IsMale() bool {
	return f.isMale
}

func (f *T) Id() int64 {
	return f.id
}

func (f *T) EqModel(f2 *T) bool {
	return arr.Eqf(f.models, f2.models, func(m1, m2 *FmodelT) bool {
		return m1.fmodelEq(m2)
	})
}

func (f *T) HasModel(modelId string) bool {
	return f.models[0].modelId == modelId ||
		f.models[1].modelId == modelId ||
		f.models[2].modelId == modelId
}

func (f *T) ModelIds() []string {
	mds := f.models
	return []string{mds[0].modelId, mds[1].modelId, mds[2].modelId}
}

func (f *T) Assets() float64 {
	return f.assets
}

// Modifies f.assets after a new evaluation.
func (f *T) Update(qs *quotes.T) {
	assets := 0.0
	for mdIx := 0; mdIx < 3; mdIx++ {
		md := f.models[mdIx]
		assets += model.FromId(md.modelId).Simulation(qs, md.params)
	}
	f.assets += (assets - f.assets) * cts.InheritanceRatio
}

func ToJs(f *T) string {
	return js.Wa([]string{
		js.Wl(f.id),
		js.Wi(f.cycle),
		js.Wb(f.isMale),
		js.Wa(arr.Map(f.models, fmodelToJs)),
		js.WdDec(f.assets, 2),
	})
}

func FromJs(j string) *T {
	a := js.Ra(j)
	return &T{
		js.Rl(a[0]),
		js.Ri(a[1]),
		js.Rb(a[2]),
		arr.Map(js.Ra(a[3]), fmodelFromJs),
		js.Rd(a[4]),
	}
}

// Generates a new flea from 'f1' and 'f2'.
func Generate(id int64, cycle int, f1, f2 *T) *T {
	mds := arr.Copy(f1.models)
	for _, md := range f2.models {
		if !arr.Anyf(mds, func(e *FmodelT) bool {
			return e.modelId == md.modelId
		}) {
			arr.Push(&mds, md)
		}
	}
	arr.Shuffle(mds)
	newMds := arr.Copy(mds[:3])

	// Mutations
	for i, md := range newMds {
		if math.Rndi(200) == 120 {
			models := model.List()
			arr.Shuffle(models)
			for _, m2 := range models {
				if m2.Id() != newMds[0].modelId &&
					m2.Id() != newMds[1].modelId &&
					m2.Id() != newMds[2].modelId {
					newMds[i] = NewFmodel(m2.Id(), m2.Mutation())
					break
				}
			}
		} else if math.Rndi(20) == 12 {
			m2 := model.FromId(md.modelId)
			newMds[i] = NewFmodel(m2.Id(), m2.Mutation())
		}
	}

	assets := (f1.assets + f2.assets) / 2.0
	return New(id, cycle, math.Rndi(2) == 0, newMds, assets)
}

// Sorts 'fleas' in place from greater to lesser
func Sort(fleas []*T) {
	arr.Sort(fleas, func(f1, f2 *T) bool {
		return f1.assets > f2.assets // greater to lesser
	})
}

// Ordered fleas table
type TbT struct {
	// Next flea indentifier
	nextId int64
	// Next generation cycle
	nextCycle int
	// Current active fleas
	fleas []*T
}

// Create a new fleas table from the sorted array 'fleas'
func NewTbFromSorted(nextId int64, nextCycle int, fleas []*T) *TbT {
	return &TbT{nextId, nextCycle, fleas}
}

// Next flea indentifier
func (tb *TbT) NextId() int64 {
	return tb.nextId
}

// Next generation cycle
func (tb *TbT) NextCycle() int {
	return tb.nextCycle
}

// Returns a sorted array of fleas
func (tb *TbT) Fleas() []*T {
	return tb.fleas
}

func TbToJs(tb *TbT) string {
	return js.Wa([]string{
		js.Wl(tb.nextId),
		js.Wi(tb.nextCycle),
		js.Wa(arr.Map(tb.fleas, ToJs)),
	})
}

func TbFromJs(j string) *TbT {
	a := js.Ra(j)

	return NewTbFromSorted(
		js.Rl(a[0]),
		js.Ri(a[1]),
		arr.Map(js.Ra(a[2]), FromJs),
	)
}
