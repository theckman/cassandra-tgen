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

package ctgformatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

// FormatTokens prints the generated tokens in the same format as the Cassandra token generator
func FormatTokens(t [][]*big.Int, w int) []byte {
	var buf bytes.Buffer

	for i, tokenList := range t {
		// print the header
		buf.WriteString(fmt.Sprintf("DC #%d:\n", i+1))

		// get the width of the largest number to properly space the column
		nnWidth := len(strconv.Itoa(len(tokenList)))

		// print each node in the datacenter
		for ni, nt := range tokenList {
			buf.WriteString(fmt.Sprintf("  Node #%*d: % *d\n", nnWidth, ni+1, w+1, nt))
		}
	}

	return buf.Bytes()
}

// FormatJSON prints the results of the token generation in a JSON format
func FormatJSON(t [][]*big.Int, prettyPrint bool) []byte {
	var jsonBytes []byte
	var err error

	if prettyPrint {
		jsonBytes, err = json.MarshalIndent(t, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(t)
	}

	if err != nil {
		panic(fmt.Sprintf("error rendering output: %v\n\n%v", err.Error(), t))
	}

	return jsonBytes
}
