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

package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

// PrintTokens prints the generated tokens in the same format as the Cassandra token generator
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

func jsonMarshal(v map[string]interface{}, pp bool) ([]byte, error) {
	if pp {
		return json.MarshalIndent(v, "", "  ")
	} else {
		return json.Marshal(v)
	}
}

// PrintJSON prints the results of the token generation in a JSON format
func FormatJSON(t [][]*big.Int, prettyPrint bool) []byte {
	data := make(map[string]interface{})

	data["keys"] = make([]*string, 0)

	var dcList []*string

	for x, v := range t {
		// set the key name for this datacenter
		dcStr := fmt.Sprintf("dc%d", x+1)

		// append the key name to the array for specifying datacenter order
		dcList = append(dcList, &dcStr)

		// create the entry in the map for this datacenter
		data[dcStr] = v
	}

	data["keys"] = dcList

	jsonBytes, err := jsonMarshal(data, prettyPrint)

	if err != nil {
		m := make(map[string]interface{})
		m["error"] = err.Error()

		j, err2 := jsonMarshal(m, prettyPrint)

		// if this is not nilwe've hit a very unexpected and very
		// unavoidable error this really should never be hit unless
		// something really bizarre happens...
		if err2 != nil {
			panic(fmt.Sprintf("unavoidable error; '%v' when JSON Marshaling error for data map. Original: '%v'\n", err2.Error(), err.Error()))
		}

		return j
	}

	return jsonBytes
}
