package doltdb

import (
	"context"
	"errors"

	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/hash"
	"github.com/attic-labs/noms/go/types"
	"github.com/liquidata-inc/ld/dolt/go/libraries/utils/pantoerr"
)

const (
	metaField      = "meta"
	parentsField   = "parents"
	rootValueField = "value"
)

// Commit contains information on a commit that was written to noms
type Commit struct {
	vrw      types.ValueReadWriter
	commitSt types.Struct
}

// HashOf returns the hash of the commit
func (c *Commit) HashOf() hash.Hash {
	return c.commitSt.Hash()
}

// GetCommitMeta gets the metadata associated with the commit
func (c *Commit) GetCommitMeta() *CommitMeta {
	metaVal := c.commitSt.Get(metaField)

	if metaVal != nil {
		if metaSt, ok := metaVal.(types.Struct); ok {
			cm, err := commitMetaFromNomsSt(metaSt)

			if err == nil {
				return cm
			}
		}
	}

	panic(c.HashOf().String() + " is a commit without the required metadata.")
}

func (c *Commit) ParentHashes() []hash.Hash {
	parentSet := c.getParents()

	hashes := make([]hash.Hash, 0, parentSet.Len())
	parentSet.IterAll(context.TODO(), func(parentVal types.Value) {
		parentRef := parentVal.(types.Ref)
		parentHash := parentRef.TargetHash()
		hashes = append(hashes, parentHash)
	})

	return hashes
}

func (c *Commit) getParents() types.Set {
	if parVal := c.commitSt.Get(parentsField); parVal != nil {
		return parVal.(types.Set)
	}

	return types.EmptySet
}

// NumParents gets the number of parents a commit has.
func (c *Commit) NumParents() int {
	parents := c.getParents()
	return int(parents.Len())
}

func (c *Commit) getParent(ctx context.Context, idx int) *types.Struct {
	parentSet := c.getParents()

	itr := parentSet.IteratorAt(ctx, uint64(idx))
	parentVal := itr.Next(ctx)

	if parentVal == nil {
		return nil
	}

	parentRef := parentVal.(types.Ref)
	parentSt := parentRef.TargetValue(ctx, c.vrw).(types.Struct)
	return &parentSt
}

// GetRootValue gets the RootValue of the commit.
func (c *Commit) GetRootValue() *RootValue {
	rootVal := c.commitSt.Get(rootValueField)

	if rootVal != nil {
		if rootSt, ok := rootVal.(types.Struct); ok {
			return newRootValue(c.vrw, rootSt)
		}
	}

	panic(c.HashOf().String() + " is a commit without a value.")
}

var ErrNoCommonAnscestor = errors.New("no common anscestor")

func GetCommitAnscestor(cm1, cm2 *Commit) (*Commit, error) {
	ref1, ref2 := types.NewRef(cm1.commitSt), types.NewRef(cm2.commitSt)

	var ancestorSt types.Struct
	err := pantoerr.PanicToErrorInstance(ErrNomsIO, func() error {
		ref, err := getCommitAncestorRef(context.TODO(), ref1, ref2, cm1.vrw)

		if err != nil {
			return err
		}

		ancestorSt, _ = ref.TargetValue(context.TODO(), cm1.vrw).(types.Struct)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Commit{cm1.vrw, ancestorSt}, nil
}

func getCommitAncestorRef(ctx context.Context, ref1, ref2 types.Ref, vrw types.ValueReadWriter) (types.Ref, error) {
	ancestorRef, ok := datas.FindCommonAncestor(ctx, ref1, ref2, vrw)

	if !ok {
		return types.Ref{}, ErrNoCommonAnscestor
	}

	return ancestorRef, nil
}

func (c *Commit) CanFastForwardTo(new *Commit) (bool, error) {
	ancestor, err := GetCommitAnscestor(c, new)

	if err != nil {
		return false, err
	} else if ancestor == nil {
		return false, errors.New("cannot perform fast forward merge; commits have no common ancestor")
	} else if ancestor.commitSt.Equals(c.commitSt) {
		if ancestor.commitSt.Equals(new.commitSt) {
			return true, ErrUpToDate
		}
		return true, nil
	} else if ancestor.commitSt.Equals(new.commitSt) {
		return false, ErrIsAhead
	}

	return false, nil
}

func (c *Commit) CanFastReverseTo(new *Commit) (bool, error) {
	ancestor, err := GetCommitAnscestor(c, new)

	if err != nil {
		return false, err
	} else if ancestor == nil {
		return false, errors.New("cannot perform fast forward merge; commits have no common ancestor")
	} else if ancestor.commitSt.Equals(new.commitSt) {
		if ancestor.commitSt.Equals(c.commitSt) {
			return true, ErrUpToDate
		}
		return true, nil
	} else if ancestor.commitSt.Equals(c.commitSt) {
		return false, ErrIsBehind
	}

	return false, nil
}
