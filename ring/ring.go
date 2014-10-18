package ring

import (
	"math/big"
	"sort"
)

// MinDcOffsetDivider is minimum token offset divider
const MinDcOffsetDivider int64 = 235

// OffsetSpacer is the default divider for the token offset
const OffsetSpacer int64 = 2

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
func NewRing(d []*big.Int, r *big.Int) (t *TokenRing) {
	t = &TokenRing{
		DcCounts:  d,
		RingRange: r,
	}
	return
}

// BestPerDcOffset is something that somethings
func (r *TokenRing) BestPerDcOffset() (iOffset *big.Int) {
	var iNumDcs, iDivider, iMdod *big.Int

	iOffset = big.NewInt(0)

	iMostNodes := big.NewInt(1)
	iLowestDivision := big.NewInt(0)

	// set the total number of datacenters
	iNumDcs = big.NewInt(int64(len(r.DcCounts)))

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

	iMdod = big.NewInt(MinDcOffsetDivider)

	if iLowestDivision.Cmp(iMdod) == 1 {
		iDivider = iLowestDivision
	} else {
		iDivider = iMdod
	}

	iOffset.Div(iOffset.Neg(r.RingRange), iDivider)

	return
}

// CalcOffsetTokensNTS is something that somethings
func (r *TokenRing) CalcOffsetTokensNTS() (dcList [][]*big.Int) {
	var dcOffset, wOffset, wArcSize, wToken *big.Int
	var wDcTokens []*big.Int

	// variable initializers
	dcList = make([][]*big.Int, 0)
	dcOffset = r.BestPerDcOffset()
	wOffset = big.NewInt(0)
	wArcSize = big.NewInt(0)

	// loop over the definition for each datacenter
	for i, v := range r.DcCounts {
		// instatiate a new temporary array of big.Int pointers
		wDcTokens = make([]*big.Int, 0)

		// the offset is r.BestPerDcOffset() * loop index
		wOffset.Mul(big.NewInt(int64(i)), dcOffset)

		// the arc is the biggest token / numNodes in this dc
		wArcSize.Div(r.RingRange, v)

		// we need to build the array of *big.Int that we need to append to dcList
		for x := big.NewInt(0); x.Cmp(v) == -1; x.Add(x, biIncrementer) {
			// instantiate a new *big.Int for wToken
			wToken = big.NewInt(0)

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
	return
}

// CalcOffsetTokensONTS is something that somethings
func (r *TokenRing) CalcOffsetTokensONTS() (dcList [][]*big.Int) {
	var dcsByCount [][]*big.Int
	var wCounts, nodeMap, layoutMap []*big.Int
	var biggestDcCount *big.Int

	dcsByCount = make([][]*big.Int, 0)
	nodeMap = make([]*big.Int, 0)
	layoutMap = make([]*big.Int, 0)

	for i, v := range r.DcCounts {
		wCounts = make([]*big.Int, 2)
		wCounts[0] = big.NewInt(int64(i))
		wCounts[1] = v

		dcsByCount = append(dcsByCount, wCounts)
	}

	sortBigInts(dcsByCount)

	biggestDcCount = dcsByCount[0][1]

	for _, v := range dcsByCount {
		dcNum, nodeCount := v[0], v[1]
		for i := big.NewInt(0); i.Cmp(nodeCount) == -1; i.Add(biIncrementer, i) {
			nodeMap = append(nodeMap, dcNum)
		}
	}

	// kind of backed myself in to a corner here
	// if a DC will ever have N nodes where N is in BigInt range, this will break
	// hopefully that doesn't happen in my lifetime
	for i := 0; i < int(biggestDcCount.Uint64()); i++ {
		for j := int(i); j < len(nodeMap); j += int(biggestDcCount.Uint64()) {
			layoutMap = append(layoutMap, nodeMap[j])
		}
	}

	dcList = make([][]*big.Int, len(dcsByCount))

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
	return
}
