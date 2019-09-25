# homework

Background:

The design is influenced by Hoare's work on CSP.

Three implementations provided. The branch master provides an implementation with a flat map structure and
the branch sharded_map provides an implementation with sharded map structure to avoid serialization of the
requests due to serialization of the map data structure guarded by channel.
Sharded map scales better but the code is more complex.
Branch syncMap is the least complex code and build on Golang sync.Map but the caveat with this implementation is that 
sync.Map data structure is build upon key value of type interface{} and thus the type safety of Golang is out of the 
window. Access to the map needs explicit type casts which can result in run time panics...


Caveats branch sharded_map, the implementation with shard map structure needs more testing
Range check on all applications parameters should be implemented




To run locally:

checkout the branch master,sharded_map or syncMap


./make

./homework   


To run in docker container

checkout the branch master or sharded_map

./build

./run

to call use api:

localhost:8080/apiUrl 

charded...

localhost:8090/apiUrl 

and syncMap

localhost:8050/apiUrl 



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
