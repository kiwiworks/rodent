package pg

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/logger"
)

type ConnectionSettings struct {
	MaxIdle     int
	MaxOpen     int
	MaxLifetime time.Duration
}

type Datasource struct {
	Host               string
	Username           string
	Password           string
	Database           string
	SslMode            string
	Port               uint16
	ConnectionSettings ConnectionSettings
}

func DatasourceFromURL(u *url.URL) (*Datasource, error) {
	log := logger.New()
	query := u.Query()
	sslmode := query.Get("sslmode")
	if sslmode == "" {
		sslmode = "disable"
	}
	host := u.Hostname()
	user := u.User
	database := strings.Trim(u.Path, "/")
	if user == nil {
		return nil, errors.Newf("a datasource cannot have empty credentials")
	}
	var portNumber uint16
	if port := u.Port(); port != "" {
		n, err := strconv.ParseUint(port, 10, 16)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid port number")
		}
		portNumber = uint16(n)
	} else {
		return nil, errors.Newf("a datasource cannot have empty port number")
	}
	password, exists := user.Password()
	if !exists {
		log.Warn("using empty password for datasource")
	}

	return &Datasource{
		Host:     host,
		Username: user.Username(),
		Password: password,
		Database: database,
		SslMode:  sslmode,
		Port:     portNumber,
		ConnectionSettings: ConnectionSettings{
			MaxIdle:     3,
			MaxOpen:     10,
			MaxLifetime: time.Second * 10,
		},
	}, nil
}

func DatasourceFromGoStyleDsn(dsn string) (*Datasource, error) {
	log := logger.New()
	fragments := strings.Split(dsn, " ")
	var datasource Datasource
	for _, fragment := range fragments {
		kv := strings.Split(fragment, "=")
		if len(kv) != 2 {
			return nil, errors.Newf("invalid dsn argument size")
		}
		key := kv[0]
		value := kv[1]
		switch key {
		case "host":
			datasource.Host = value
		case "user":
			datasource.Username = value
		case "password":
			datasource.Password = value
		case "sslmode":
			datasource.SslMode = value
		case "port":
			port, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid port number")
			}
			datasource.Port = uint16(port)
		case "database":
			if value == "" {
				log.Warn("empty database name")
			}
			datasource.Database = value
		default:
			log.Info("unrecognised datasource option", zap.String("key", key), zap.String("value", value))
		}
	}
	return &datasource, nil
}

func ParseDatasource(dsn string) (*Datasource, error) {
	u, err := url.Parse(dsn)
	if err == nil {
		return DatasourceFromURL(u)
	}
	return DatasourceFromGoStyleDsn(dsn)
}

func (p Datasource) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s port=%d",
		p.Host, p.Username, p.Password, p.Database, p.SslMode, p.Port,
	)
}
