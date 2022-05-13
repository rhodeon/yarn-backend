package internal

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Env            string
	Version        string
	Port           int
	JwtSecret      string
	DisplayVersion bool

	Db struct {
		Name         string
		Uri          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
}

func (c *Config) Parse() {
	flag.StringVar(&c.Env, "env", c.defaultEnv(), "Environment (development|staging|production)\nDotenv variable: ENV\n")
	flag.IntVar(&c.Port, "port", c.defaultPort(), "API server port\nDotenv variable: PORT\n")
	flag.StringVar(&c.JwtSecret, "jwt-secret", c.defaultJwtSecret(), "JWT Secret Key\nDotEnv variable: JWT_SECRET\n")
	flag.BoolVar(&c.DisplayVersion, "version", false, "Display version and build time")

	flag.StringVar(&c.Db.Uri, "db-uri", c.defaultDbUri(), "MongoDB Connection String URI\nDotenv variable: DB_URI\n")
	flag.StringVar(&c.Db.Name, "db-name", c.defaultDbName(), "MongoDB Database Name\nDotenv variable: DB_NAME\n")
	flag.IntVar(&c.Db.MaxOpenConns, "db-max-open-conns", c.defaultDbMaxOpenConns(), "MongoDB maximum number of open connections\nDotenv variable: DB_MAX_OPEN_CONNS\n")
	flag.IntVar(&c.Db.MaxIdleConns, "db-max-idle-conns", c.defaultDbMaxIdleConns(), "MongoDB maximum number of idle connections\nDotenv variable: DB_MAX_IDLE_CONNS\n")
	flag.StringVar(&c.Db.MaxIdleTime, "db-max-idle-time", c.defaultDbMaxIdleTime(), "MongoDB maximumn idle time\nDotenv variable: DB_MAX_IDLE_TIME\n")

	flag.Parse()
}

// Validate ensures required flags or environment variables are present
func (c *Config) Validate() error {
	if c.Db.Uri == "" {
		return errors.New("the 'db-uri' flag is required")
	}

	if c.Db.Name == "" {
		return errors.New("the 'db-name flag is required")
	}

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

func (c *Config) defaultJwtSecret() string {
	const defaultSecret = ""

	if secret, exists := os.LookupEnv("JWT_SECRET"); exists {
		return secret
	}
	return defaultSecret
}

func (c *Config) defaultDbUri() string {
	const defaultUri = ""
	if uri, exists := os.LookupEnv("DB_URI"); exists {
		return uri
	}
	return defaultUri
}

func (c *Config) defaultDbName() string {
	const defaultName = ""

	if name, exists := os.LookupEnv("DB_NAME"); exists {
		return name
	}
	return defaultName
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
