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

package options

import (
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

// ArgOpts is a struct for specifying command line options for this utility
type ArgOpts struct {
	// The -j or --json flag specifies whether or not to use the JSON output format
	JSON bool `short:"j" long:"json" default:"false" description:"use the JSON output format"`

	// The -J or --pretty-json option specifies whether or not to use the formatter JSON output
	Pretty bool `short:"J" long:"pretty-json" default:"false" description:"print the output in a prettier JSON format assumes '-j'"`

	// The -n or --nts flag specifies that your cluster will be using the new NetworkTopologyStrategy
	Nts bool `short:"n" long:"nts" default:"true" description:"Optimize multi-cluster distribution for NetworkTopologyStrategy [default]"`

	// The -o or --onts flag specifies that your cluster will be using the antiquated OldNetworkTopologyStrategy
	Onts bool `short:"o" long:"onts" default:"false" description:"Optimize multi-cluster distribution for OldNetworkTopologyStrategy"`

	// The -r or --ringrange flag specifies the maximum size of the ring
	// This value is a string so that it can be converted to a big.Int value later (too large for uint64)
	RingRangeStr string `short:"r" long:"ringrange" default:"170141183460469231731687303715884105728" description:"Specify a numeric maximum token value for your ring, different from the default value of 2^127"`

	// The -d or --dc-count flag specifies the total number of datacenters in the cluster
	DcCountStr string `short:"d" long:"dc-count" default:"0" description:"The number of datacenters to calculate token ranges for"`

	// The -c or --node-count flag is a comma-separated list node counts per DC, For example:
	// ./cassandra-tgen -d 3 -c 3,2,1
	NodeCountStr string `short:"c" long:"node-count" default:"1" description:"Comma-delimited list of datacenter node counts: --count '3,2,1'"`

	// These two options are not use by go-flags, they are some of the
	// values parsed from string => big.Int that are needed later
	DcCount, RingRange *big.Int
	NodeCounts         []*big.Int
}

func badConversion(v interface{}) error {
	return fmt.Errorf("failed to convert '%v' to big.Int", v)
}

// New returns a pointer to a new ArgOpts struct
func New() *ArgOpts {
	return &ArgOpts{}
}

// Parse parses the command line arguments for the option struct.
// It sets the values directly on the struct, as well as returns a pointer to the struct as
// well as an error
func (opts *ArgOpts) Parse() error {
	var nc []*big.Int

	parser := flags.NewParser(opts, flags.HelpFlag|flags.PassDoubleDash)

	_, err := parser.Parse()

	// if parsing bombed and it's not the help message
	// we need to print the error and bail out with exit code 1
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
			os.Exit(1)
		} else {
			fmt.Printf("%v\n", err.Error())
			os.Exit(0)
		}
	}

	// if Onts unset Nts which is enabled by default
	if opts.Onts {
		opts.Nts = false
	}

	// --pretty-json infers that you are using the JSON formatter
	// so if --pretty-json was provided but not --json, make sure --json is set to true
	if opts.Pretty && !opts.JSON {
		opts.JSON = true
	}

	// split the node counts in to a slice of strings
	nodeCounts := strings.Split(opts.NodeCountStr, ",")

	// convert the total DC numbers to a *big.Int
	count, ok := new(big.Int).SetString(opts.DcCountStr, 10)

	// if converting opts.DcCountStr => *big.Int failed return the failure
	if !ok {
		return badConversion(opts.DcCountStr)
	}

	ringRange, ok := new(big.Int).SetString(opts.RingRangeStr, 10)

	if !ok {
		return badConversion(opts.RingRangeStr)
	}

	opts.RingRange = ringRange

	// if the specific DC count doesn't match up with the node counts provided
	// return an error because things are whacked
	if count.Cmp(big.NewInt(int64(len(nodeCounts)))) != 0 && opts.DcCountStr != "0" {
		return fmt.Errorf("the datacenter count (-d) must be equivalent to count of items in the node count (-c) array")
	}

	// assign the DC count to opts.DcCount
	opts.DcCount = count

	// this iterates over the nodeCounts and converts each item to
	// a *big.Int while appending it to the 'nc' slice
	for _, v := range nodeCounts {
		bi, ok := new(big.Int).SetString(v, 10)

		if !ok {
			return badConversion(v)
		}

		nc = append(nc, bi)
	}

	// assign the []*big.Int to opts.NodeCounts
	opts.NodeCounts = nc

	// return the options and an error of nil
	return nil
}
