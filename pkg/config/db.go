package config

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// DB holds the DB configuration
type DB struct {
	Host            string
	Port            int
	SslMode         string
	Name            string
	User            string
	Password        string
	Debug           bool
	MaxOpenConn     int
	MaxIdleConn     int
	MaxConnLifetime time.Duration
}

var db = &DB{}

func DBCfg() *DB { return db }

// LoadDBCfg loads DB configuration from (in order):
//  1. DATABASE_URL (e.g. postgres://user:pass@host:5432/db?sslmode=disable)
//  2. DB_* vars (your current scheme)
//  3. POSTGRES_* / PG* fallbacks
func LoadDBCfg() {
	// Prefer DATABASE_URL if present
	if raw := strings.TrimSpace(os.Getenv("DATABASE_URL")); raw != "" {
		if parseDatabaseURL(raw, db) {
			fillDBDefaults(db)
			return
		}
		// fall through to envs on parse error
	}

	// DB_* first (your existing layout)
	db.Host = firstNonEmpty(os.Getenv("DB_HOST"), os.Getenv("POSTGRES_HOST"), os.Getenv("PGHOST"))
	db.Port = firstInt(5432, os.Getenv("DB_PORT"), os.Getenv("POSTGRES_PORT"), os.Getenv("PGPORT"))
	db.User = firstNonEmpty(os.Getenv("DB_USER"), os.Getenv("POSTGRES_USER"), os.Getenv("PGUSER"))
	db.Password = firstNonEmpty(os.Getenv("DB_PASSWORD"), os.Getenv("DB_PASS"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("PGPASSWORD"))
	db.Name = firstNonEmpty(os.Getenv("DB_NAME"), os.Getenv("POSTGRES_DB"), os.Getenv("PGDATABASE"))
	db.SslMode = strings.ToLower(strings.TrimSpace(firstNonEmpty(
		os.Getenv("DB_SSL_MODE"), // your current key
		os.Getenv("DB_SSLMODE"),  // common variant
		os.Getenv("DB_SSLMODE"),  // (yes, keep it if you want)
		os.Getenv("PGSSLMODE"),
		os.Getenv("SSLMODE"),
	)))

	db.Debug = firstBool(false, os.Getenv("DB_DEBUG"))
	db.MaxOpenConn = firstInt(10, os.Getenv("DB_MAX_OPEN_CONNECTIONS"))
	db.MaxIdleConn = firstInt(2, os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	db.MaxConnLifetime = time.Duration(firstInt(600, os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))) * time.Second

	fillDBDefaults(db)
}

func fillDBDefaults(c *DB) {
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.SslMode == "" {
		// Supabase Postgres requires TLS
		c.SslMode = "require"
	}
	// Hard fail early with a clear message if critical fields are missing.
	if c.Host == "" || c.User == "" || c.Name == "" {
		log.Fatalf("database env incomplete: host=%q user=%q dbname=%q (set DB_* or POSTGRES_* or DATABASE_URL)",
			c.Host, c.User, c.Name)
	}
}

// BuildPostgresDSN returns a lib/pq style DSN.
func BuildPostgresDSN(c *DB) string {
	return "host=" + c.Host +
		" port=" + strconv.Itoa(c.Port) +
		" sslmode=" + c.SslMode +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.Name
}

// ---- helpers ----

func parseDatabaseURL(raw string, out *DB) bool {
	u, err := url.Parse(raw)
	if err != nil {
		log.Printf("WARN: bad DATABASE_URL (%v); ignoring", err)
		return false
	}
	pw, _ := u.User.Password()
	port := firstInt(5432, u.Port())
	out.Host = u.Hostname()
	out.Port = port
	out.User = u.User.Username()
	out.Password = pw
	out.Name = strings.TrimPrefix(u.Path, "/")
	out.SslMode = strings.ToLower(strings.TrimSpace(u.Query().Get("sslmode")))
	return true
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

func firstInt(def int, vals ...string) int {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			if n, err := strconv.Atoi(s); err == nil {
				return n
			}
		}
	}
	return def
}

func firstBool(def bool, vals ...string) bool {
	for _, v := range vals {
		if s := strings.TrimSpace(strings.ToLower(v)); s != "" {
			return s == "1" || s == "true" || s == "yes" || s == "y"
		}
	}
	return def
}
