package db

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/logger"
	"github.com/sean-/vpc/agent"
)

//type Pool struct {
//	Scheme             string
//	User               string
//	Host               string
//	Port               uint16
//	Database           string
//	CAPath             string
//	CertPath           string
//	KeyPath            string
//	InsecureSkipVerify bool
//
//	pool            *pgx.ConnPool
//	stdDriverConfig *stdlib.DriverConfig
//}
//
//func New(cfg *agent.AgentConfig) (*Pool, error) {
//	pool := &Pool{}
//
//	p, err := pgx.NewConnPool(cfg.DB.PoolConfig)
//	if err != nil {
//		return nil, errors.Wrap(err, "unable to create a new DB connection pool")
//	}
//	pool.pool = p
//
//	if err := pool.Ping(); err != nil {
//		return nil, errors.Wrap(err, "unable to ping database")
//	}
//
//	tlsConfig, err := cfg.TLSConfig()
//	if err != nil {
//		return nil, errors.Wrap(err, "unable to generate a TLS config")
//	}
//
//	pool.stdDriverConfig = &stdlib.DriverConfig{
//		ConnConfig: pgx.ConnConfig{
//			Database:       cfg.DB.PoolConfig.ConnConfig.Database,
//			Dial:           (&net.Dialer{Timeout: config.DefaultConnTimeout, KeepAlive: 5 * time.Minute}).Dial,
//			Host:           cfg.DB.PoolConfig.ConnConfig.Host,
//			LogLevel:       pgx.LogLevelTrace, //pgxLogLevel,
//			Logger:         logger.NewPGX(log.Logger),
//			Port:           cfg.DB.PoolConfig.ConnConfig.Port,
//			TLSConfig:      tlsConfig,
//			UseFallbackTLS: false,
//			User:           cfg.DB.PoolConfig.ConnConfig.User,
//			RuntimeParams: map[string]string{
//				"application_name": buildtime.PROGNAME,
//			},
//		},
//		AfterConnect: func(c *pgx.Conn) error {
//			return nil
//		},
//	}
//	stdlib.RegisterDriverConfig(pool.stdDriverConfig)
//
//	return pool, nil
//}
//
//func (p *Pool) Close() error {
//	p.pool.Close()
//
//	return nil
//}
//
//func (p *Pool) Ping() error {
//	pingCtx, pingCancel := context.WithTimeout(context.Background(), config.DefaultConnTimeout)
//	defer pingCancel()
//	conn, err := p.pool.Acquire()
//	if err != nil {
//		return errors.Wrap(err, "unable to acquire database connection for ping")
//	}
//	defer p.pool.Release(conn)
//
//	if err := conn.Ping(pingCtx); err != nil {
//		return errors.Wrap(err, "unable to ping database")
//	}
//
//	return nil
//}
//
//func (p *Pool) Pool() *pgx.ConnPool {
//	return p.pool
//}
//
//func (p *Pool) STDDB() (*sql.DB, error) {
//	username := p.User
//	hostname := p.Host
//	port := p.Port
//	database := p.Database
//
//	v := url.Values{}
//	if !p.InsecureSkipVerify {
//		sslMode := "require"
//
//		v.Set("sslmode", sslMode)
//		v.Set("sslrootcert", p.CAPath)
//		v.Set("sslcert", p.CertPath)
//		v.Set("sslkey", p.KeyPath)
//	}
//
//	u := url.URL{
//		Scheme:   p.Scheme,
//		User:     url.User(username),
//		Host:     net.JoinHostPort(hostname, strconv.Itoa(int(port))),
//		Path:     database,
//		RawQuery: v.Encode(),
//	}
//
//	encodedURI := p.stdDriverConfig.ConnectionString(u.String())
//	db, err := sql.Open("pgx", encodedURI)
//	if err != nil {
//		return nil, errors.Wrap(err, "unable to open a standard database connection")
//	}
//
//	return db, nil
//}
