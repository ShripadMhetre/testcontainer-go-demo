package app

type metadata struct {
	PostgresURL string
	RedisURL    string
}

var Connections *metadata = &metadata{}
