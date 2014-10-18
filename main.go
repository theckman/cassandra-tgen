package main

import (
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/theckman/cassandra-tgen/formatter"
	"github.com/theckman/cassandra-tgen/options"
	"github.com/theckman/cassandra-tgen/ring"
)

func interactiveInput() (counts []*big.Int) {
	var line string

	counts = make([]*big.Int, 0)

	// incrementer for i++ operations
	increment := big.NewInt(1)

	fmt.Println("Token Generator Interactive Mode")
	fmt.Println("--------------------------------")
	fmt.Println("")

	// retry label for when the number of datacenters cannot be converted to a big.Int
CountRetry:
	fmt.Print(" How many datacenters will participate in this Cassandra cluster? ")

	// get the user input
	fmt.Scanln(&line)

	// try to convert the user input to a big.Int
	numDcs, ok := new(big.Int).SetString(line, 10)

	// if the user's input cannot be converted inform them over the error and try again
	if !ok {
		fmt.Println(fmt.Sprintf("Oops, '%v' can't be converted to a big.Int\n", line))
		goto CountRetry
	}

	// build the array of node counts per datacenter in a loop
	for i := big.NewInt(0); i.Cmp(numDcs) == -1; i.Add(i, increment) {
		// retry label for when the number of nodes in this datacenter cannot be converted to a big.Int
	InputRetry:
		dcNum := new(big.Int).Set(i)
		dcNum.Add(dcNum, increment)

		fmt.Print(fmt.Sprintf(" How many nodes are in datacenter #%d? ", dcNum))

		// get the user input
		fmt.Scanln(&line)

		// try to convert the user input to a big.Int
		dcCount, ok := new(big.Int).SetString(line, 10)

		// if the input couldnt be converted print an error and try again
		if !ok {
			fmt.Println(fmt.Sprintf("Oops, '%v' can't be converted to to a big.Int\n", line))
			goto InputRetry
		}

		// append the value to the count array
		counts = append(counts, dcCount)
	}

	fmt.Println("")

	return
}

func main() {
	var tokenResults [][]*big.Int
	var tokenRing *ring.TokenRing

	// build a new option struct and parse it
	opts := options.New()
	_, err := opts.Parse()

	// if parsing bombed...
	if err != nil {
		// and it's not the help message...
		if !strings.Contains(fmt.Sprintf("%v", err), "Usage") {
			// print the error and exit indicating failure
			fmt.Println("error:", err)
			os.Exit(1)
		}
		// just return as it was the help message
		return
	}

	// if the count string is the default value we should use the interactive input mode
	// otherwise use the command-line options
	if opts.DcCountStr == "0" {
		nodeCounts := interactiveInput()
		tokenRing = ring.NewRing(nodeCounts, opts.RingRange)
	} else {
		tokenRing = ring.NewRing(opts.NodeCounts, opts.RingRange)
	}

	// call the proper function based on NTS or ONTS option
	if opts.Nts {
		tokenResults = tokenRing.CalcOffsetTokensNTS()
	} else {
		tokenResults = tokenRing.CalcOffsetTokensONTS()
	}

	// if we are printing in JSON format determine which method we're using
	// if not JSON just print the table format
	if opts.JSON {
		if opts.Pretty {
			formatter.PrintJSON(tokenResults, true)
		} else {
			formatter.PrintJSON(tokenResults, false)
		}
	} else {
		longestTokenLen := len(opts.RingRangeStr)
		formatter.PrintTokens(tokenResults, longestTokenLen)
	}
}
