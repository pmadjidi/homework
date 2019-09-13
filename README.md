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



Performance data...

➜  pedometers git:(master) ✗ ab -n 10000 -c 100 "http://localhost:8080/get/step/payam" 
This is ApacheBench, Version 2.3 <$Revision: 1826891 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 1000 requests
Completed 2000 requests
Completed 3000 requests
Completed 4000 requests
Completed 5000 requests
Completed 6000 requests
Completed 7000 requests
Completed 8000 requests
Completed 9000 requests
Completed 10000 requests
Finished 10000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /get/step/payam
Document Length:        27 bytes

Concurrency Level:      100
Time taken for tests:   9.629 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1440000 bytes
HTML transferred:       270000 bytes
Requests per second:    1038.53 [#/sec] (mean)
Time per request:       96.290 [ms] (mean)
Time per request:       0.963 [ms] (mean, across all concurrent requests)
Transfer rate:          146.04 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.4      0       9
Processing:     8   96  29.4     92     261
Waiting:        4   81  23.7     80     219
Total:         10   96  29.4     92     261

Percentage of the requests served within a certain time (ms)
  50%     92
  66%    104
  75%    112
  80%    119
  90%    134
  95%    149
  98%    168
  99%    180
 100%    261 (longest request)



