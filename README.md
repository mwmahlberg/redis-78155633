Code for [my answer][answer] to the question [go-redis always get dial tcp i/o timeout][question] on StackOverflow
====

Prerequisites
-------------

* A running docker
* docker compose


Run
-------------

```none
$ docker compose build
...
=> => naming to docker.io/library/myapp:latest                             0.0s
$ docker compose up app
[+] Running 15/15
 ✔ Network redis-78155633_default        Created                           0.0s 
 ✔ Volume "redis-78155633_redis-data-4"  Created                           0.0s 
 ✔ Volume "redis-78155633_redis-data-5"  Created                           0.0s 
 ✔ Volume "redis-78155633_redis-data-6"  Created                           0.0s 
 ✔ Volume "redis-78155633_redis-data-1"  Created                           0.0s 
 ✔ Volume "redis-78155633_redis-data-2"  Created                           0.0s 
 ✔ Volume "redis-78155633_redis-data-3"  Created                           0.0s 
 ✔ Container redis-2                     Created                           0.1s 
 ✔ Container redis-3                     Created                           0.1s 
 ✔ Container redis-1                     Created                           0.1s 
 ✔ Container redis-5                     Created                           0.1s 
 ✔ Container redis-6                     Created                           0.1s 
 ✔ Container redis-4                     Created                           0.1s 
 ✔ Container redis-cluster-creator       Created                           0.0s 
 ✔ Container redis-78155633-app-1        Created                           0.0s 
Attaching to app-1
app-1  | myapp: Cluster is healhty and has 3 shards
app-1  | myapp: * shard with keyrange 0-5460 has 2 nodes and 1 slots
app-1  | myapp:     master:      172.26.0.6 (online) 
app-1  | myapp:    replica:      172.26.0.4 (loading) 
app-1  | myapp: * shard with keyrange 10923-16383 has 2 nodes and 1 slots
app-1  | myapp:     master:      172.26.0.3 (online) 
app-1  | myapp:    replica:      172.26.0.7 (loading) 
app-1  | myapp: * shard with keyrange 5461-10922 has 2 nodes and 1 slots
app-1  | myapp:     master:      172.26.0.5 (online) 
app-1  | myapp:    replica:      172.26.0.2 (loading) 
app-1  | myapp: 
app-1  |        Pinging each shard member
app-1  |        =========================
app-1  | myapp: 172.26.0.5:6379 - ping: PONG
app-1  | myapp: 172.26.0.6:6379 - ping: PONG
app-1  | myapp: 172.26.0.3:6379 - ping: PONG
app-1  | myapp: my_counter: 1
app-1  |        
```

[question]: https://stackoverflow.com/questions/78155633/go-redis-always-get-dial-tcp-i-o-timeout/
[answer]: https://stackoverflow.com/a/78162786/1296707