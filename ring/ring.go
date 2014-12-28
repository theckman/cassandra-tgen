// Copyright (c) 2014 Tim Heckman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ring

import (
	"math/big"
	"sort"
)

const (
	// MinDcOffsetDivider is minimum token offset divider
	MinDcOffsetDivider int64 = 235

	// OffsetSpacer is the default divider for the token offset
	OffsetSpacer int64 = 2
)

// *big.Int -- used for i++ operations on big.Int objects
var biIncrementer = big.NewInt(1)

type bigIntSliceSlice [][]*big.Int

func (s bigIntSliceSlice) Swap(x, y int) {
	s[x], s[y] = s[y], s[x]
}

func (s bigIntSliceSlice) Less(x, y int) bool {
	return s[x][1].Cmp(s[y][1]) == -1
}

func (s bigIntSliceSlice) Len() int {
	return len(s)
}

func (s bigIntSliceSlice) SortReverse() {
	sort.Sort(sort.Reverse(s))
}

func sortBigInts(b [][]*big.Int) {
	sort.Sort(sort.Reverse(bigIntSliceSlice(b)))
}

// TokenRing is the configuration of the cluster's ring for calculating tokens
type TokenRing struct {
	DcCounts  []*big.Int
	RingRange *big.Int
}

// NewRing returns a new instance of TokenRing
func New(d []*big.Int, r *big.Int) *TokenRing {
	return &TokenRing{
		DcCounts:  d,
		RingRange: r,
	}
}

// BestPerDcOffset is something that somethings
func (r *TokenRing) BestPerDcOffset() *big.Int {
	iOffset := big.NewInt(0)

	iMostNodes := big.NewInt(1)
	iLowestDivision := big.NewInt(0)

	// set the total number of datacenters
	iNumDcs := big.NewInt(int64(len(r.DcCounts)))

	// figure out the one with the most nodes
	for _, v := range r.DcCounts {
		if v.Cmp(iMostNodes) == 1 {
			iMostNodes = v
		}
	}

	// multiply the number of datacenters, the total number of nodes in the datacnter with the most
	// and the offset spacer and assign the result to iLowestDivision
	iLowestDivision.Mul(iNumDcs, iMostNodes)
	iLowestDivision.Mul(iLowestDivision, big.NewInt(OffsetSpacer))

	iMdod := big.NewInt(MinDcOffsetDivider)

	var iDivider *big.Int

	if iLowestDivision.Cmp(iMdod) == 1 {
		iDivider = iLowestDivision
	} else {
		iDivider = iMdod
	}

	iOffset.Div(iOffset.Neg(r.RingRange), iDivider)

	return iOffset
}

// CalcOffsetTokensNTS is something that somethings
func (r *TokenRing) CalcOffsetTokensNTS() [][]*big.Int {
	var (
		dcOffset = r.BestPerDcOffset()
		wOffset  = big.NewInt(0)
		wArcSize = big.NewInt(0)
		dcList   = make([][]*big.Int, 0)
	)

	// loop over the definition for each datacenter
	for i, v := range r.DcCounts {
		// instatiate a new temporary array of big.Int pointers
		wDcTokens := make([]*big.Int, 0)

		// the offset is r.BestPerDcOffset() * loop index
		wOffset.Mul(big.NewInt(int64(i)), dcOffset)

		// the arc is the biggest token / numNodes in this dc
		wArcSize.Div(r.RingRange, v)

		// we need to build the array of *big.Int that we need to append to dcList
		for x := big.NewInt(0); x.Cmp(v) == -1; x.Add(x, biIncrementer) {
			// instantiate a new *big.Int for wToken
			wToken := big.NewInt(0)

			// the generated token is (loop index * wArcSize + wOffset) % r.RingRange
			wToken.Mul(x, wArcSize)
			wToken.Add(wToken, wOffset)
			wToken.Mod(wToken, r.RingRange)

			// append the generated token to the working Token array (wDcTokens)
			wDcTokens = append(wDcTokens, wToken)
		}
		// append wDcTokens to dcList
		dcList = append(dcList, wDcTokens)
	}

	return dcList
}

// CalcOffsetTokensONTS is something that somethings
func (r *TokenRing) CalcOffsetTokensONTS() [][]*big.Int {
	dcsByCount := make([][]*big.Int, 0)

	for i, v := range r.DcCounts {
		wCounts := make([]*big.Int, 2)

		wCounts[0] = big.NewInt(int64(i))
		wCounts[1] = v

		dcsByCount = append(dcsByCount, wCounts)
	}

	sortBigInts(dcsByCount)

	nodeMap := make([]*big.Int, 0)
	for _, v := range dcsByCount {
		dcNum, nodeCount := v[0], v[1]
		for i := big.NewInt(0); i.Cmp(nodeCount) == -1; i.Add(biIncrementer, i) {
			nodeMap = append(nodeMap, dcNum)
		}
	}

	// kind of backed myself in to a corner here
	// if a DC will ever have N nodes where N is in BigInt range, this will break
	// hopefully that doesn't happen in my lifetime
	layoutMap := make([]*big.Int, 0)
	biggestDcCount := dcsByCount[0][1]
	for i := 0; i < int(biggestDcCount.Uint64()); i++ {
		for j := int(i); j < len(nodeMap); j += int(biggestDcCount.Uint64()) {
			layoutMap = append(layoutMap, nodeMap[j])
		}
	}

	dcList := make([][]*big.Int, len(dcsByCount))

	for i := range dcsByCount {
		dcList[i] = make([]*big.Int, 0)
	}

	for i, v := range layoutMap {
		index := int(v.Uint64())

		val := big.NewInt(0)
		val.Mul(r.RingRange, big.NewInt(int64(i)))
		val.Div(val, big.NewInt(int64(len(layoutMap))))

		dcList[index] = append(dcList[index], val)
	}

	return dcList
}
