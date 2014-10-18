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
		fmt.Println(fmt.Sprintf("DC #%d:", i+1))

		// this should *never* error
		nnWidth := len(strconv.Itoa(len(tokenList)))

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

	data["order"] = make([]*string, 0)

	for x, v := range t {
		var dcTokens []*big.Int
		dcStr := fmt.Sprintf("dc_%d", x)
		dcList = append(dcList, &dcStr)

		data[dcStr] = make([]*big.Int, 0)

		for _, z := range v {
			dcTokens = append(dcTokens, z)
		}

		data[dcStr] = dcTokens
	}

	data["order"] = dcList

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
