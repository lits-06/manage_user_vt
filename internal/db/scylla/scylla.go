package scylla

import (
	"log"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func ConnectDB() {
    var cluster = gocql.NewCluster("127.0.0.1")
    cluster.Keyspace = "history"

    s, err := cluster.CreateSession()
    if err != nil {
        log.Fatal("Failed to connect to cluster")
    }

    if err := s.Query("CREATE KEYSPACE IF NOT EXISTS history WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}").Exec(); err != nil {
        log.Fatal("Failed to create keyspace")
    }

    if err := s.Query("CREATE TABLE IF NOT EXISTS info (email TEXT PRIMARY KEY, id UUID, info TEXT)").Exec(); err != nil {
        log.Fatal("Failed to create table")
    }

    Session = s
}
    