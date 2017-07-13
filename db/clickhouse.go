package db

func init() {
	Register("clickhouse", &Clickhouse{})
}

type Clickhouse struct {
	Driver
}
