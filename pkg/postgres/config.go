package postgres

type Config struct {
	Host         string `default:"localhost"`
	Port         int    `default:"5432"`
	DBName       string `default:"poseidon"`
	User         string `default:"poseidon"`
	Password     string `default:"123456"`
	ConnLifeTime int    `default:"300"`
	MaxIdleConns int    `default:"10"`
	MaxOpenConns int    `default:"80"`
	LogLevel     int    `default:"1"`
}
