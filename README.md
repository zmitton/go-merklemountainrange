# MerkleMountainRange
Golang version of [this](https://github.com/zmitton/merkle-mountain-range)

The plan is to have the exact same API. Both packages use a fileBaseddb adapter that can interface with the same database format (the database is just a file).

This golang version is still being made however (Javascript one works now).

The format for `.mmr` files is being changed to hold the  `wordsize` information.
```
[[wordsize]8 [leafLength]8 ]wordsize [leaf0] wordsize [leaf1] wordsize...
```

Every operation benchmarked thus far has been almost _exactly_ 20x faster than its JS version.

memoryBased Get (unverified):        800ns
memoryBased GetVerified:           3.8µs
memoryBased Append:                6.0µs

fileBased Get (unverified):        2.0µs
fileBased GetVerified:            38.9µs
fileBased Append               2.5ms




/*
make a reverse getNodePosition function (getLeafIndex?), and in the test, do a loop to
100,000 testing each result against its inverse function
targetIndex -> targetNodeIndex (in mountainpositions function)
*/


