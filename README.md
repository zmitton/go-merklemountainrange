# MerkleMountainRange
Golang version of [this](https://github.com/zmitton/merkle-mountain-range)

The plan is to have the exact same API. Both packages use a fileBaseddb adapter that can interface with the same database format (the database is just a file).

This golang version is still being made however (Javascript one works now).

The format for `.mmr` files is being changed to hold the  `wordsize` information.
```
[[wordsize]8 [leafLength]8 ]wordsize [leaf0] wordsize [leaf1] wordsize...
```

Every operation benchmarked thus far has been almost _exactly_ 20x faster than its JS version.

```
memoryBased GetUnverified:           800ns
memoryBased Get:                   3.8µs
memoryBased Append:                6.0µs

fileBased GetUnverified:           2.0µs
fileBased Get:                    38.9µs
fileBased Append               2.5ms
```

<!-- 
notes
/*
make a reverse getNodePosition function (getLeafIndex?), and in the test, do a loop to
100,000 testing each result against its inverse function (actually is this possible? consider the fact that some nodes dont have a cooresponding leaf).
name change: targetIndex -> targetNodeIndex (in mountainpositions function)
 - remember to move metadata in `.mmr` in js implimentation (this is major version bump)
 - add `serialize()` method to db api and add `fromSerialized()` to membased db
 - add `getUnverified()` method to js api (note: has to check leaflength)
*/

nodes (map[int64][]byte):
{ 
  30 : 0x1234567890,
  33 : 0x2143658709,
  34 : 0x1234123434
}

encodable version ([][][]byte):
[
  [12,34],
  [
    [1e],
    [12,34,56,78,90]
  ],
  [
    [21],
    [21,43,65,87,09]
  ],
  [
    [22],
    [12,34,12,34,34]
  ]
]

 -->
