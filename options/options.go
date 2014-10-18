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
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/jessevdk/go-flags"
)

// ArgOpts is a struct for specifying command line options for this utility
type ArgOpts struct {
	// The -j or --json flag specifies whether or not to use the JSON output format
	JSON bool `short:"j" long:"json" description:"use the JSON output format" default:"false"`

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

	// These two options are not use by go-flags -- some of the values parsed from string => big.Int are stored here
	DcCount, RingRange *big.Int
	NodeCounts         []*big.Int
}

func badConversion(o *ArgOpts, v interface{}) (*ArgOpts, error) {
	return o, fmt.Errorf("failed to convert '%v' to big.Int", v)
}

// New returns a pointer to a new ArgOpts struct
func New() *ArgOpts {
	return &ArgOpts{}
}

// Parse parses the command line arguments for the option struct.
// It sets the values directly on the struct, as well as returns a pointer to the struct as
// well as an error
func (opts *ArgOpts) Parse() (*ArgOpts, error) {
	var nc []*big.Int
	var lastVal *string

	// try to parse the arguments
	_, err := flags.Parse(opts)

	// if there was an error parsing the args returns a new ArgOpts
	// as well as the error
	if err != nil {
		opts = &ArgOpts{}
		return opts, err
	}

	// if --onts unset -nts which is enabled by default
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
		return badConversion(opts, opts.DcCountStr)
	}

	ringRange, ok := new(big.Int).SetString(opts.RingRangeStr, 10)

	if !ok {
		return badConversion(opts, opts.RingRangeStr)
	}

	opts.RingRange = ringRange

	// if the specific DC count doesn't match up with the node counts provided
	// return an error because things are whacked
	if count.Cmp(big.NewInt(int64(len(nodeCounts)))) != 0 && opts.DcCountStr != "0" {
		err := "the datacenter count (-d) must be equivalent to count of items in the node count (-c) array"
		return opts, errors.New(err)
	}

	// assign the DC count to opts.DcCount
	opts.DcCount = count

	// this iterates over the nodeCounts and converts each item to
	// a *big.Int while appending it to the 'nc' slice
	for _, v := range nodeCounts {
		lastVal = &v
		bi, ok := new(big.Int).SetString(*lastVal, 10)
		if !ok {
			break
		}
		nc = append(nc, bi)
	}

	// if the length of the two slices fails to match up we aborted the above loop
	// because we were unable to convert some value
	if len(nc) != len(nodeCounts) {
		return badConversion(opts, *lastVal)
		// return opts, fmt.Errorf("failed to convert '%v' to a big.Int", *lastVal)
	}

	// assign the []*big.Int to opts.NodeCounts
	opts.NodeCounts = nc

	// return the options and an error of nil
	return opts, nil
}
