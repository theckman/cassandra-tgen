// Copyright (c) 2014-2015 Tim Heckman
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

package ctgformatter_test

import (
	"math/big"
	"testing"

	"github.com/theckman/cassandra-tgen/ctgformatter"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {
	tests  [][][]*big.Int
	rrsLen int
}

var _ = Suite(&TestSuite{})

var formatTokenOutputs = []string{
	`DC #1:
  Node #1:                                        0
`,
	`DC #1:
  Node #1:                                        0
  Node #2:   56713727820156410577229101238628035242
  Node #3:  113427455640312821154458202477256070484
`,
	`DC #1:
  Node #1:                                        0
  Node #2:   56713727820156410577229101238628035242
  Node #3:  113427455640312821154458202477256070484
DC #2:
  Node #1:  169417178424467235000914166253263322299
  Node #2:   41811290829115311202148688466350243003
  Node #3:   84346586694232619135070514395321269435
  Node #4:  126881882559349927067992340324292295867
`,
}

var formatJSONUgly = []string{
	"[[0]]",
	"[[0,56713727820156410577229101238628035242,113427455640312821154458202477256070484]]",
	"[[0,56713727820156410577229101238628035242,113427455640312821154458202477256070484],[169417178424467235000914166253263322299,41811290829115311202148688466350243003,84346586694232619135070514395321269435,126881882559349927067992340324292295867]]",
}

var formatJSONPretty = []string{
	`[
  [
    0
  ]
]`,
	`[
  [
    0,
    56713727820156410577229101238628035242,
    113427455640312821154458202477256070484
  ]
]`,
	`[
  [
    0,
    56713727820156410577229101238628035242,
    113427455640312821154458202477256070484
  ],
  [
    169417178424467235000914166253263322299,
    41811290829115311202148688466350243003,
    84346586694232619135070514395321269435,
    126881882559349927067992340324292295867
  ]
]`,
}

func (s *TestSuite) SetUpSuite(c *C) {
	s.rrsLen = len("170141183460469231731687303715884105728")

	// build out tests
	w := make([]*big.Int, 0, 1)
	d := make([][]*big.Int, 0, 1)

	bI, ok := new(big.Int).SetString("0", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)
	s.tests = [][][]*big.Int{append(d, w)}

	w = make([]*big.Int, 0, 3)
	d = make([][]*big.Int, 0, 1)

	bI, ok = new(big.Int).SetString("0", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)

	bI, ok = new(big.Int).SetString("56713727820156410577229101238628035242", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)

	bI, ok = new(big.Int).SetString("113427455640312821154458202477256070484", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)
	s.tests = append(s.tests, append(d, w))

	w = make([]*big.Int, 0, 3)
	d = make([][]*big.Int, 0, 2)

	bI, ok = new(big.Int).SetString("0", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)

	bI, ok = new(big.Int).SetString("56713727820156410577229101238628035242", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)

	bI, ok = new(big.Int).SetString("113427455640312821154458202477256070484", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)
	d = append(d, w)

	w = make([]*big.Int, 0, 4)

	bI, ok = new(big.Int).SetString("169417178424467235000914166253263322299", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)

	bI, ok = new(big.Int).SetString("41811290829115311202148688466350243003", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)

	bI, ok = new(big.Int).SetString("84346586694232619135070514395321269435", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)

	bI, ok = new(big.Int).SetString("126881882559349927067992340324292295867", 0)
	c.Assert(ok, Equals, true)

	w = append(w, bI)

	d = append(d, w)
	s.tests = append(s.tests, d)
}

func (s *TestSuite) TearDownSuite(c *C) {
	s.tests = nil
}

func (s *TestSuite) TestFormatTokens(c *C) {
	for i := range s.tests {
		val := ctgformatter.FormatTokens(s.tests[i], 39)
		c.Check(string(val), Equals, formatTokenOutputs[i])
	}
}

func (s *TestSuite) TestFormatJSON(c *C) {
	// test without pretty print
	for i := range s.tests {
		val := ctgformatter.FormatJSON(s.tests[i], false)
		c.Check(string(val), Equals, formatJSONUgly[i])
	}

	// test with pretty print
	for i := range s.tests {
		val := ctgformatter.FormatJSON(s.tests[i], true)
		c.Check(string(val), Equals, formatJSONPretty[i])
	}
}
