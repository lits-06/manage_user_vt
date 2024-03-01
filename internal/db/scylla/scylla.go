package scylla

import (
	"log"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func ConnectDB() {
    var cluster = gocql.NewCluster("127.0.0.1:9042")
    cluster.Keyspace = "system"
    
    s, err := cluster.CreateSession()
    if err != nil {
        log.Printf("Failed to connect to cluster")
    }

    createKeyspace(s)
    cluster.Keyspace = "history"

    s, err = cluster.CreateSession()
    if err != nil {
        log.Printf("Failed to connect to cluster")
    }
    createTable(s)

    Session = s
}

func createKeyspace(s *gocql.Session) {
    if err := s.Query("CREATE KEYSPACE IF NOT EXISTS history WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}").Exec(); err != nil {
        log.Printf("Failed to create keyspace")
    }
}
    
func createTable(s *gocql.Session) {
    if err := s.Query("CREATE TABLE IF NOT EXISTS info (id UUID PRIMARY KEY, email TEXT, info TEXT)").Exec(); err != nil {
        log.Printf("Failed to create table")
    }
}