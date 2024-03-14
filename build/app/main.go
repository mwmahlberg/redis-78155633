package main

import (
	"context"
	"math"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/redis/go-redis/v9"
)

var cfg struct {
	RedisCache struct {
		Addresses             []string      `kong:"default='redis-node-1:6379',help='A list of Redis addresses to connect to.'"`
		Password              string        `kong:"help='The password to use when connecting to Redis.'"`
		PoolSize              int           `kong:"default='10',help='The size of the Redis connection pool.'"`
		MaxRetries            int           `kong:"default='3',help='The maximum number of retries before giving up on a Redis command.'"`
		InitialClusterTimeout time.Duration `kong:"default='10s',help='The maximum time to wait for the cluster to be ready.'"`
	} `kong:"embed,prefix=redis-"`
}

func main() {
	appCtx := kong.Parse(&cfg,
		kong.DefaultEnvars("MY_APP"), // Optional: set up environment variable prefixing.
	)

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:      cfg.RedisCache.Addresses,
		Password:   cfg.RedisCache.Password,
		PoolSize:   cfg.RedisCache.PoolSize,
		MaxRetries: cfg.RedisCache.MaxRetries,

		// Have at least one tenth of the max pool size as min idle connections,
		// rounded up.
		MinIdleConns: int(math.Ceil(float64(cfg.RedisCache.PoolSize) / 10.0)),
	})
	defer rdb.Close()

	/*
		This is a simple example of how to wait for the cluster to be ready.
	*/

	// We set a timeout
	end := time.After(cfg.RedisCache.InitialClusterTimeout)

	// We create a ticker to check the cluster health every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

cluster_health:

	// We loop until the cluster is ready or the timeout is reached
	for {
		select {
		case <-end:
			appCtx.Printf("Timeout reached, the cluster is not ready yet.")
			return
		case <-ticker.C:
			v := rdb.ClusterInfo(context.Background()).Val()
			if strings.Contains(v, "cluster_state:ok") {
				// The label is used to break the outer loop.
				// Without it, the select would be broken, but the for loop would
				// continue.
				break cluster_health
			}
		}
	}

	_, err := rdb.Ping(context.Background()).Result()
	appCtx.FatalIfErrorf(err, "pinging cluster: %#v", err)

	// We have several ways to examine the shards.

	// We can get the shards directly...
	shards, err := rdb.ClusterShards(signalCtx).Result()
	appCtx.FatalIfErrorf(err, "getting cluster shards: %#v", err)
	appCtx.Printf("Cluster is healhty and has %d shards", len(shards))

	for _, shard := range shards {
		// ...and then get the nodes and slots for each shard.
		appCtx.Printf(
			"* shard with keyrange %d-%d has %d nodes and %d slots",
			shard.Slots[0].Start,
			shard.Slots[len(shard.Slots)-1].End,
			len(shard.Nodes),
			len(shard.Slots))
		for _, node := range shard.Nodes {
			appCtx.Printf("%10s: %15s (%s) ", node.Role, node.Endpoint, node.Health)
		}
	}

	appCtx.Printf("\nPinging each shard member\n=========================")
	if err := rdb.ForEachShard(context.Background(), func(ctx context.Context, client *redis.Client) error {
		r := client.Ping(ctx)
		appCtx.Printf("%s - %s", client.Options().Addr, r)
		return r.Err()
	}); err != nil {
		appCtx.Printf("pinging cluster shards: %#v", err)
	}

	doSomethingWithRedis(appCtx, rdb)
	<-signalCtx.Done()
	appCtx.Printf("Shutting down gracefully\n")
	stop()
}

func doSomethingWithRedis(appCtx *kong.Context, rdb *redis.ClusterClient) {
	// Use the Redis client.
	execCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rdb.Incr(execCtx, "my_counter")
	select {
	case <-execCtx.Done():
		d, _ := execCtx.Deadline()
		appCtx.Printf("timeout exceeded at %s", d)
	default:
		v, err := rdb.Get(execCtx, "my_counter").Int()
		if err != nil {
			appCtx.Printf("getting my_counter: %s", err)
			return
		}
		appCtx.Printf("my_counter: %d\n", v)
	}

}
