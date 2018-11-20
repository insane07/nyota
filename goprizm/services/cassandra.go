package services

import (
	"github.com/gocql/gocql"
	"goprizm/sysutils"
	"strings"
)

// Cassandra holds gocql session
type Cassandra struct {
	Session *gocql.Session
}

// GetCassandra returns Cassandra object with live session
func GetCassandra() (Cassandra, error) {

	// Get details from environment variable
	hosts := sysutils.Getenv("CASSANDRA_HOSTS", "localhost")
	port := sysutils.GetenvInt("CASSANDRA_PORT", 9042)
	keyspace := sysutils.Getenv("CASSANDRA_KEYSPACE", "prizm")

	// ACP has cassandra with username/password but	in our cluster
	// its probably none.
	username := sysutils.Getenv("CASSANDRA_USER", "")
	password := sysutils.Getenv("CASSANDRA_PASSWORD", "")

	cluster := gocql.NewCluster(strings.Split(hosts, ",")...)
	cluster.Keyspace = keyspace
	cluster.Port = port

	if len(username) != 0 || len(password) != 0 {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}
	}

	sess, err := cluster.CreateSession()

	if err != nil {
		return Cassandra{}, err
	}

	cassy := Cassandra{
		Session: sess,
	}

	return cassy, nil
}
