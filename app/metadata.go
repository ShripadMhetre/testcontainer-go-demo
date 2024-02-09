package app

type Metadata struct {
	PostgresURL string
	RedisURL    string
}

var Connections *Metadata = &Metadata{}
