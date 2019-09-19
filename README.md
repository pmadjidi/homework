# homework

Background

The design is influenced by Hoare's work on CSP.

Two implementations provided. The branch master provides an implementation with a flat map structure and
the branch sharded_map provides an implementation with sharded map structure to avoid serialization of the
requests due to serialization of the map data structure guarded by channel.
Sharded map scales better but the code is more complex.

Caveats branch sharded_map, the implementation with shard map structure needs more testing
Range check on all app parameters should be implemented




Some performance data is gathered, please refer to perfdata.txt
API hosted on http://payam.ite.kth.se:8080
for the distributed version please refer to  http://payam.ite.kth.se:9090


To run locally:

checkout the branch master or sharded_map


./make

./homework   


To run in docker container

checkout the branch master or sharded_map

./build

./run

to call use api:

localhost:8080/apiUrl


API
Method GET

/add/step/{person}

/inc/{person}/{steps}

/get/step/{person}

/get/allsteps

/add/group/{name}

/extend/{group}/{person}

/get/group/{name}

/get/allgroups
