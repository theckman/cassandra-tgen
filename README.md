# cassandra-tgen
[![TravisCI Build Status](https://img.shields.io/travis/theckman/cassandra-tgen/master.svg)](https://travis-ci.org/theckman/cassandra-tgen)
[![GoDoc](https://img.shields.io/badge/cassandra--tgen-GoDoc-blue.svg)](https://godoc.org/github.com/theckman/cassandra-tgen)

This is an implementation of the Apache Cassandra project's token generator in Go.

This implements the token generation ability, without the ability to generate the HTML graph showing the ring.

In addition, you can now output the data in a standard JSON format, or a pretty-printed JSON format. This will allow
you to use the tool to generate tokens and easily consume the results.

## LICENSE
This software is released under the [Apache 2.0](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)) license.
This work is a derivative of work done by the Apache Software Foundation (ASF) released under the Apache 2.0 license.
This work is in no way related to the works being done by the ASF.

## INSTALLATION
You'll need a properly installed Go environment ([http://golang.org/doc/install](http://golang.org/doc/install)) to use
`cassandra-tgen`

To install `cassandra-tgen`:

```
go install github.com/theckman/cassandra-tgen
```

## TODOs
It should be functionally complete. I need to write some unit tests. It has been manually tested and the output is the
same as the Apache Cassandra version.

## USAGE
Assuming you've installed Go and casandra-tgen properly the `cassandra-tgen` should be in your path.

### Interactive Mode
If you provide no flags to `cassandra-tgen` it starts in interactive mode:

```
$ cassandra-tgen
Token Generator Interactive Mode
--------------------------------

 How many datacenters will participate in this Cassandra cluster? 3
 How many nodes are in datacenter #1? 3
 How many nodes are in datacenter #2? 2
 How many nodes are in datacenter #3? 1

DC #1:
  Node #1:                                        0
  Node #2:   56713727820156410577229101238628035242
  Node #3:  113427455640312821154458202477256070484
DC #2:
  Node #1:  169417178424467235000914166253263322299
  Node #2:   84346586694232619135070514395321269435
DC #3:
  Node #1:  168693173388465238270141028790642538870
```

### Using Flags
You can also specify it via flags:

```
$ cassandra-tgen -d 3 -c3,2,1
DC #1:
  Node #1:                                        0
  Node #2:   56713727820156410577229101238628035242
  Node #3:  113427455640312821154458202477256070484
DC #2:
  Node #1:  169417178424467235000914166253263322299
  Node #2:   84346586694232619135070514395321269435
DC #3:
  Node #1:  168693173388465238270141028790642538870
```

### JSON

```
cassandra-tgen -d3 -c3,2,1 -j
{"dc1":[0,56713727820156410577229101238628035242,113427455640312821154458202477256070484],"dc2":[169417178424467235000914166253263322299,84346586694232619135070514395321269435],"dc3":[168693173388465238270141028790642538870],"keys":["dc1","dc2","dc3"]}
```

#### Pretty JSON
```
$ cassandra-tgen -d 3 -c3,2,1 -J
{
  "dc1": [
    0,
    56713727820156410577229101238628035242,
    113427455640312821154458202477256070484
  ],
  "dc2": [
    169417178424467235000914166253263322299,
    84346586694232619135070514395321269435
  ],
  "dc3": [
    168693173388465238270141028790642538870
  ],
  "keys": [
    "dc1",
    "dc2",
    "dc3"
  ]
}
```

### Help Output
```
$ cassandra-tgen -h
Usage:
  cassandra-tgen [OPTIONS]

Application Options:
  -j, --json         use the JSON output format (false)
  -J, --pretty-json  print the output in a prettier JSON format assumes '-j' (false)
  -n, --nts          Optimize multi-cluster distribution for NetworkTopologyStrategy [default] (true)
  -o, --onts         Optimize multi-cluster distribution for OldNetworkTopologyStrategy (false)
  -r, --ringrange=   Specify a numeric maximum token value for your ring, different from the default value of 2^127 (170141183460469231731687303715884105728)
  -d, --dc-count=    The number of datacenters to calculate token ranges for (0)
  -c, --node-count=  Comma-delimited list of datacenter node counts: --count '3,2,1' (1)

Help Options:
  -h, --help         Show this help message
```