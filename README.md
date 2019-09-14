# homework

Background

The design is influenced by Hoare's work on CSP.
The homework could have been implemented equally well with mutex guarded maps
or  Golangs sync.Map for concurrent safe access to map data structure.
The former would probably perform slightly better but would not be as "Goish" and the later 
violates type safety by use of interface{} and type cast could cause runtime panics... and thus is considered by the author as premature optimization
for this task.
All thought Golangs sync.Map could be a solution to scale a cross massive multiprocessing platforms...
To scale beyond single host solutions, distributed event logs and event processors is a way to proceed.

Some performance data is gathered, please refer to perfdata.txt
API hosted on http://payam.ite.kth.se:8080


To run locally:

./make

./homework   


To run in docker container

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

