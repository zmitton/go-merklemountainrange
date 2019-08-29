# MerkleMountainRange
Golang version of [this](https://github.com/zmitton/merkle-mountain-range). Please use the readme over there to understand how this all works. It's almost exactly the same API. This go implimentation is finally feature complete and well tested.

As of now, the only function missing here is the ability to add more db nodes to an existing proof (sparse) mmr. You can already serialize an mmr and create an mmr from serilized data, so I'm not sure when you would need to use that other feature.

Both packages support a memoryBasedDb or a fileBaseddb. The fileBased format is identical (`.mmr` can be opened by either implimentation). I would like to add leveldb based db support (mostly for use with ethereum), but I havent yet because the file format I designed is about half the space/cost. The `.mmr` file is essentially treated as a random access array. Here is the format.

```
[[wordsize](8) [leafLength](8)](wordsize) [node0](wordsize) [node1](wordsize)...
```

I might make a rust version, any help would be greatly appreciated (I dont know rust at all). This is also my first golang project so feel free to correct any small or large mistakes.


Every operation benchmarked thus far has been almost _exactly_ 20x faster than its JS version.

```
memoryBased GetUnverified:           800ns
memoryBased Get:                   3.8µs
memoryBased Append:                6.0µs

fileBased GetUnverified:           2.0µs
fileBased Get:                    38.9µs
fileBased Append               2.5ms
```

This fileBased append takes much longer unfortunately because it calls fs.sync(). If anyone knows a better way to (essentially) _save_ the db after the `append`, please let me know.

<!-- 
notes
/*
make a reverse getNodePosition function (getLeafIndex?), and in the test, do a loop to
100,000 testing each result against its inverse function (actually is this possible? consider the fact that some nodes dont have a cooresponding leaf).
name change: targetIndex -> targetNodeIndex (in mountainpositions function)
 - remember to move metadata in `.mmr` in js implimentation (this is major version bump)
 - add `serialize()` method to db api and add `fromSerialized()` to membased db
 - add `getUnverified()` method to js api (note: has to check leaflength)
 - add pretections for file-based in case it already exists (doesnt overwrite)
 - add persistent leveldb support with NAMESPACE feature
*/

 -->
