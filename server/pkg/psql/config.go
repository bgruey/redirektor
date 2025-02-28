package psql

import (
	"regexp"
	"strconv"
)

var dbUrlRegex = regexp.MustCompile("//(.+):(.+)@(.+):(.+)/(.+)")

type Config struct {
	Url                   string
	Host                  string
	Port                  int
	Username              string
	Password              string
	Database              string
	SSLMode               string
	MaxIdleConnections    int
	MaxOpenConnections    int
	SaveSQLAfterExecution bool
	CallbackSqlKey        string
}

func (c *Config) parse() {
	if c.MaxOpenConnections == 0 {
		c.MaxOpenConnections = 50
	}

	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = 50
	}

	if c.Url != "" {
		parts := dbUrlRegex.FindStringSubmatch(c.Url)
		c.Username = parts[1]
		c.Password = parts[2]
		c.Host = parts[3]
		c.Port, _ = strconv.Atoi(parts[4])
		c.Database = parts[5]
	}
}
