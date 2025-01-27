package database

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

var connectionString string = "localhost"
var clusterOptions = gocb.ClusterOptions{
	Authenticator: gocb.PasswordAuthenticator{
		Username: "class_a_rw",
		Password: "123456",
	},
}

// remove Couchbase bucket

func main() {
	cluster, err := gocb.Connect("couchbase://"+connectionString, clusterOptions)
	if err != nil {
		log.Fatal(err)
	}
	bucketMgr := cluster.Buckets()
	bucketMgr.DropBucket("market", nil)
	fmt.Print("bucket market is dropped. " + time.Now().String())
}
