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
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

// PrintTokens prints the generated tokens in the same format as the Cassandra token generator
func PrintTokens(t [][]*big.Int, w int) {
	for i, tokenList := range t {
		// print the header
		fmt.Println(fmt.Sprintf("DC #%d:", i+1))

		// get the width of the largest number to properly space the column
		nnWidth := len(strconv.Itoa(len(tokenList)))

		// print each node in the datacenter
		for ni, nt := range tokenList {
			fmt.Println(fmt.Sprintf("  Node #%*d: % *d", nnWidth, ni+1, w+1, nt))
		}
	}
}

// PrintJSON prints the results of the token generation in a JSON format
func PrintJSON(t [][]*big.Int, prettyPrint bool) {
	var dcList []*string
	var jsonBytes []byte
	var err error

	data := make(map[string]interface{})

	data["keys"] = make([]*string, 0)

	for x, v := range t {
		// set the key name for this datacenter
		dcStr := fmt.Sprintf("dc%d", x+1)

		// append the key name to the array for specifying datacenter order
		dcList = append(dcList, &dcStr)

		// create the entry in the map for this datacenter
		data[dcStr] = v
	}

	data["keys"] = dcList

	if prettyPrint {
		jsonBytes, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(data)
	}

	if err != nil {
		fmt.Println("error printing json:", err)
		return
	}

	fmt.Println(string(jsonBytes))
}
