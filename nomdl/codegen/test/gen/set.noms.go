// This file was generated by nomdl/codegen.

package gen

import (
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

// SetOfBool

type SetOfBool struct {
	s   types.Set
	ref *ref.Ref
}

func NewSetOfBool() SetOfBool {
	return SetOfBool{types.NewSet(), &ref.Ref{}}
}

type SetOfBoolDef map[bool]bool

func (def SetOfBoolDef) New() SetOfBool {
	l := make([]types.Value, len(def))
	i := 0
	for d, _ := range def {
		l[i] = types.Bool(d)
		i++
	}
	return SetOfBool{types.NewSet(l...), &ref.Ref{}}
}

func (s SetOfBool) Def() SetOfBoolDef {
	def := make(map[bool]bool, s.Len())
	s.s.Iter(func(v types.Value) bool {
		def[bool(v.(types.Bool))] = true
		return false
	})
	return def
}

func (s SetOfBool) Equals(other types.Value) bool {
	return other != nil && __typeRefForSetOfBool.Equals(other.Type()) && s.Ref() == other.Ref()
}

func (s SetOfBool) Ref() ref.Ref {
	return types.EnsureRef(s.ref, s)
}

func (s SetOfBool) Chunks() (chunks []ref.Ref) {
	chunks = append(chunks, s.Type().Chunks()...)
	chunks = append(chunks, s.s.Chunks()...)
	return
}

func (s SetOfBool) ChildValues() []types.Value {
	return append([]types.Value{}, s.s.ChildValues()...)
}

// A Noms Value that describes SetOfBool.
var __typeRefForSetOfBool types.Type

func (m SetOfBool) Type() types.Type {
	return __typeRefForSetOfBool
}

func init() {
	__typeRefForSetOfBool = types.MakeCompoundTypeRef(types.SetKind, types.MakePrimitiveTypeRef(types.BoolKind))
	types.RegisterValue(__typeRefForSetOfBool, builderForSetOfBool, readerForSetOfBool)
}

func builderForSetOfBool(v types.Value) types.Value {
	return SetOfBool{v.(types.Set), &ref.Ref{}}
}

func readerForSetOfBool(v types.Value) types.Value {
	return v.(SetOfBool).s
}

func (s SetOfBool) Empty() bool {
	return s.s.Empty()
}

func (s SetOfBool) Len() uint64 {
	return s.s.Len()
}

func (s SetOfBool) Has(p bool) bool {
	return s.s.Has(types.Bool(p))
}

type SetOfBoolIterCallback func(p bool) (stop bool)

func (s SetOfBool) Iter(cb SetOfBoolIterCallback) {
	s.s.Iter(func(v types.Value) bool {
		return cb(bool(v.(types.Bool)))
	})
}

type SetOfBoolIterAllCallback func(p bool)

func (s SetOfBool) IterAll(cb SetOfBoolIterAllCallback) {
	s.s.IterAll(func(v types.Value) {
		cb(bool(v.(types.Bool)))
	})
}

func (s SetOfBool) IterAllP(concurrency int, cb SetOfBoolIterAllCallback) {
	s.s.IterAllP(concurrency, func(v types.Value) {
		cb(bool(v.(types.Bool)))
	})
}

type SetOfBoolFilterCallback func(p bool) (keep bool)

func (s SetOfBool) Filter(cb SetOfBoolFilterCallback) SetOfBool {
	out := s.s.Filter(func(v types.Value) bool {
		return cb(bool(v.(types.Bool)))
	})
	return SetOfBool{out, &ref.Ref{}}
}

func (s SetOfBool) Insert(p ...bool) SetOfBool {
	return SetOfBool{s.s.Insert(s.fromElemSlice(p)...), &ref.Ref{}}
}

func (s SetOfBool) Remove(p ...bool) SetOfBool {
	return SetOfBool{s.s.Remove(s.fromElemSlice(p)...), &ref.Ref{}}
}

func (s SetOfBool) Union(others ...SetOfBool) SetOfBool {
	return SetOfBool{s.s.Union(s.fromStructSlice(others)...), &ref.Ref{}}
}

func (s SetOfBool) Subtract(others ...SetOfBool) SetOfBool {
	return SetOfBool{s.s.Subtract(s.fromStructSlice(others)...), &ref.Ref{}}
}

func (s SetOfBool) Any() bool {
	return bool(s.s.Any().(types.Bool))
}

func (s SetOfBool) fromStructSlice(p []SetOfBool) []types.Set {
	r := make([]types.Set, len(p))
	for i, v := range p {
		r[i] = v.s
	}
	return r
}

func (s SetOfBool) fromElemSlice(p []bool) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = types.Bool(v)
	}
	return r
}
