package scylla

import (
	"fmt"
	"os"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func ConnectDB() {
    var cluster = gocql.NewCluster(os.Getenv("NODE_0"), os.Getenv("NODE_1"), os.Getenv("NODE_2"))
    cluster.Authenticator = gocql.PasswordAuthenticator{Username: os.Getenv("USERNAME_SCYLLA"), Password: os.Getenv("PASSWORD_SCYLLA")}
    cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy("REGION")
    cluster.Keyspace = os.Getenv("KEYSPACE_SCYLLA")

    Session, err := cluster.CreateSession()
    if err != nil {
        panic("Failed to connect to cluster")
    }

    defer Session.Close()

    var query = Session.Query("SELECT * FROM system.clients")

    if rows, err := query.Iter().SliceMap(); err == nil {
        for _, row := range rows {
            fmt.Printf("%v\n", row)
        }
    } else {
        panic("Query error: " + err.Error())
    }
}
    