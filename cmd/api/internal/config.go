package internal

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Env            string
	Version        string
	Port           int
	DisplayVersion bool

	Db struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
}

func (c *Config) Parse() {
	flag.StringVar(&c.Env, "env", c.defaultEnv(), "Environment (development|staging|production)\nDotenv variable: ENV\n")
	flag.IntVar(&c.Port, "port", c.defaultPort(), "API server port\nDotenv variable: PORT\n")
	flag.BoolVar(&c.DisplayVersion, "version", false, "Display version and build time")

	flag.StringVar(&c.Db.Dsn, "db-dsn", c.defaultDbDsn(), "PostgreSQL DSN\nDotenv variable: DB_DSN\n")
	flag.IntVar(&c.Db.MaxOpenConns, "db-max-open-conns", c.defaultDbMaxOpenConns(), "PostgreSQL maximum number of open connections\nDotenv variable: DB_MAX_OPEN_CONNS\n")
	flag.IntVar(&c.Db.MaxIdleConns, "db-max-idle-conns", c.defaultDbMaxIdleConns(), "PostgreSQL maximum number of idle connections\nDotenv variable: DB_MAX_IDLE_CONNS\n")
	flag.StringVar(&c.Db.MaxIdleTime, "db-max-idle-time", c.defaultDbMaxIdleTime(), "PostgreSQL maximumn idle time\nDotenv variable: DB_MAX_IDLE_TIME\n")

	flag.Parse()
}

// Validate ensures required flags or environment variables are present
func (c *Config) Validate() error {
	//if c.Db.Dsn == "" {
	//	return errors.New("the 'db-dsn' flag is required")
	//}
	return nil
}

func (c *Config) defaultEnv() string {
	const defaultEnv = "development"

	if env, exists := os.LookupEnv("ENV"); exists {
		return env
	}
	return defaultEnv
}

func (c *Config) defaultPort() int {
	const defaultPort = 4000

	if portEnv, exists := os.LookupEnv("PORT"); exists {
		port, err := strconv.Atoi(portEnv)
		if err == nil {
			return port
		}
	}
	return defaultPort
}

func (c *Config) defaultDbDsn() string {
	const defaultDsn = ""

	if dsn, exists := os.LookupEnv("DB_DSN"); exists {
		return dsn
	}
	return defaultDsn
}

func (c *Config) defaultDbMaxOpenConns() int {
	const defMaxOpenConns = 25

	if maxOpenConnsEnv, exists := os.LookupEnv("DB_MAX_OPEN_CONNS"); exists {
		maxOpenConns, err := strconv.Atoi(maxOpenConnsEnv)
		if err == nil {
			return maxOpenConns
		}
	}
	return defMaxOpenConns
}

func (c *Config) defaultDbMaxIdleConns() int {
	const defMaxIdleConns = 25

	if maxIdleConnsEnv, exists := os.LookupEnv("DB_MAX_IDLE_CONNS"); exists {
		maxIdleConns, err := strconv.Atoi(maxIdleConnsEnv)
		if err == nil {
			return maxIdleConns
		}
	}
	return defMaxIdleConns
}

func (c *Config) defaultDbMaxIdleTime() string {
	const defMaxIdleTime = "15m"

	if maxIdleTime, exists := os.LookupEnv("DB_MAX_IDLE_TIME"); exists {
		return maxIdleTime
	}
	return defMaxIdleTime
}
