package database

func init() {
	Register("clickhouse", &Clickhouse{})
}

type Clickhouse struct {
	Driver
}
