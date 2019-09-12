# homework

Background

The design is influenced by Hoare's work on CSP
The homework could have been implemented equally well with mutex guarded maps
or  Golangs sync.Map for concurrent safe access.
The former would probably perform slightly better but would not be "Goish" and the later 
violates type  safety and thus is considered by the author as premature optimization
for this task
All thought Golangs sync.Map could be a solution to scale a cross massive multiprocessing platforms...
To scale beyond singe host solutions, distributed event logs and event processors is a way to proceed.





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


