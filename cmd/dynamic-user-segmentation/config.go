package main

import (
	"errors"
	"fmt"
	zap_logger "github.com/romandnk/dynamic-user-segmentation-service/internal/logger/zap"
	http_server "github.com/romandnk/dynamic-user-segmentation-service/internal/server/http"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/storage/postgres"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

var (
	ErrZapLoggerInvalidEncoding       = errors.New("invalid encoding (json, console)")
	ErrZapLoggerEmptyOutputPath       = errors.New("empty output path")
	ErrZapLoggerEmptyErrorOutputPath  = errors.New("empty error output path")
	ErrZapLoggerInvalidLevel          = errors.New("invalid level (debug, info, warn, error, dpanic, panic, fatal)")
	ErrPostgresParseMaxConnLifetime   = errors.New("invalid max conn lifetime (format 1h2m3s)")
	ErrPostgresParseMaxConnIdleTime   = errors.New("invalid max conn idle time (format 1h2m3s)")
	ErrPostgresEmptyHost              = errors.New("empty host")
	ErrPostgresInvalidPort            = errors.New("invalid port (from 0 to 65535 inclusively)")
	ErrPostgresEmptyUsername          = errors.New("empty username")
	ErrPostgresEmptyPassword          = errors.New("empty password")
	ErrPostgresEmptyDBName            = errors.New("empty database name")
	ErrPostgresInvalidSSLMode         = errors.New("invalid ssl mode (disable, allow, prefer, require, verify-ca, verify-full")
	ErrPostgresInvalidMaxConns        = errors.New("max conns cannot be less than zero")
	ErrPostgresInvalidMinConns        = errors.New("min conns cannot be less than zero")
	ErrPostgresInvalidMaxConnLifetime = errors.New("max conn lifetime cannot be less than zero")
	ErrPostgresInvalidMaxConnIdleTime = errors.New("max conn idle time cannot be less than zero")
	ErrServerParseReadTimeout         = errors.New("invalid read timeout (format 1h2m3s)")
	ErrServerParseWriteTimeout        = errors.New("invalid write timeout (format 1h2m3s)")
	ErrServerEmptyHost                = errors.New("empty host")
	ErrServerInvalidPort              = errors.New("invalid port (from 0 to 65535 inclusively)")
	ErrServerInvalidReadTimeout       = errors.New("invalid read timeout (must be only positive)")
	ErrServerInvalidWriteTimeout      = errors.New("invalid write timeout (must be only positive)")
	ErrParseTicker                    = errors.New("invalid auto add ticker (format 1h2m3s)")
)

type Config struct {
	ZapLogger     zap_logger.Config
	Postgres      postgres.Config
	Server        http_server.Config
	Ticker        time.Duration
	PathToReports string
}

func NewConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("DUS") // DUS stands for "dynamic user segmentation" service
	viper.AutomaticEnv()

	zapLoggerConfig, err := newZapLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("zap logger: %w", err)
	}

	postgresConfig, err := newPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	serverConfig, err := newServerConfig()
	if err != nil {
		return nil, fmt.Errorf("server: %w", err)
	}

	tickerStr := viper.GetString("auto_add_ticker")
	ticker, err := time.ParseDuration(tickerStr)
	if err != nil {
		return nil, fmt.Errorf("auto add ticker: %w", ErrParseTicker)
	}

	pathToReports := viper.GetString("path_to_reports")

	config := &Config{
		ZapLogger:     zapLoggerConfig,
		Postgres:      postgresConfig,
		Server:        serverConfig,
		Ticker:        ticker,
		PathToReports: pathToReports,
	}

	return config, nil
}

func newZapLoggerConfig() (zap_logger.Config, error) {
	levelStr := viper.GetString("zap_logger.level")
	level, err := stringToZapLogLevel(levelStr)
	if err != nil {
		return zap_logger.Config{}, err
	}

	encoding := viper.GetString("zap_logger.encoding")
	outputPath := viper.GetStringSlice("zap_logger.output_path")
	errorOutputPath := viper.GetStringSlice("zap_logger.error_output_path")

	cfg := zap_logger.Config{
		Level:           level,
		Encoding:        encoding,
		OutputPath:      outputPath,
		ErrorOutputPath: errorOutputPath,
	}

	err = validateZapLoggerConfig(cfg)
	if err != nil {
		return zap_logger.Config{}, err
	}

	return cfg, nil
}

func validateZapLoggerConfig(cfg zap_logger.Config) error {
	encodings := map[string]struct{}{
		"json":    {},
		"console": {},
	}
	if _, ok := encodings[cfg.Encoding]; !ok {
		return fmt.Errorf("encoding: %w", ErrZapLoggerInvalidEncoding)
	}
	if len(cfg.OutputPath) == 0 {
		return fmt.Errorf("output path: %w", ErrZapLoggerEmptyOutputPath)
	}
	if len(cfg.ErrorOutputPath) == 0 {
		return fmt.Errorf("error output path: %w", ErrZapLoggerEmptyErrorOutputPath)
	}

	return nil
}

func stringToZapLogLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "dpanic":
		return zapcore.DPanicLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InvalidLevel, fmt.Errorf("level: %w", ErrZapLoggerInvalidLevel)
	}
}

func newPostgresConfig() (postgres.Config, error) {
	host := viper.GetString("postgres_storage.host")
	port := viper.GetInt("postgres_storage.port")
	username := viper.GetString("POSTGRES_USERNAME")
	password := viper.GetString("POSTGRES_PASSWORD")
	DBName := viper.GetString("postgres_storage.db_name")
	sslMode := viper.GetString("postgres_storage.sslmode")
	maxConns := viper.GetInt("postgres_storage.max_conns")
	minConns := viper.GetInt("postgres_storage.min_conns")
	maxConnLifetime := viper.GetString("postgres_storage.max_conn_lifetime")
	parsedMaxConnLifetime, err := time.ParseDuration(maxConnLifetime)
	if err != nil {
		return postgres.Config{}, fmt.Errorf("max conn lifetime: %w", ErrPostgresParseMaxConnLifetime)
	}

	maxConnIdleTime := viper.GetString("postgres_storage.max_conn_idle_time")
	parsedMaxConnIdleTime, err := time.ParseDuration(maxConnIdleTime)
	if err != nil {
		return postgres.Config{}, fmt.Errorf("max conn idle time: %w", ErrPostgresParseMaxConnIdleTime)
	}

	cfg := postgres.Config{
		Host:            host,
		Port:            port,
		Username:        username,
		Password:        password,
		DBName:          DBName,
		SSLMode:         sslMode,
		MaxConns:        maxConns,
		MinConns:        minConns,
		MaxConnLifetime: parsedMaxConnLifetime,
		MaxConnIdleTime: parsedMaxConnIdleTime,
	}

	err = validatePostgresConfig(cfg)
	if err != nil {
		return postgres.Config{}, err
	}

	return cfg, nil
}

func validatePostgresConfig(cfg postgres.Config) error {
	if cfg.Host == "" {
		return fmt.Errorf("host: %w", ErrPostgresEmptyHost)
	}
	if cfg.Port < 0 || cfg.Port > 65535 {
		return fmt.Errorf("port: %w", ErrPostgresInvalidPort)
	}
	if cfg.Username == "" {
		return fmt.Errorf("username: %w", ErrPostgresEmptyUsername)
	}
	if cfg.Password == "" {
		return fmt.Errorf("password: %w", ErrPostgresEmptyPassword)
	}
	if cfg.DBName == "" {
		return fmt.Errorf("database name: %w", ErrPostgresEmptyDBName)
	}
	sslModes := map[string]struct{}{
		"disable":     {},
		"allow":       {},
		"prefer":      {},
		"require":     {},
		"verify-ca":   {},
		"verify-full": {},
	}
	if _, ok := sslModes[cfg.SSLMode]; !ok {
		return fmt.Errorf("ssl mode: %w", ErrPostgresInvalidSSLMode)
	}
	if cfg.MaxConns <= 0 {
		return fmt.Errorf("max conns: %w", ErrPostgresInvalidMaxConns)
	}
	if cfg.MinConns <= 0 {
		return fmt.Errorf("min conns: %w", ErrPostgresInvalidMinConns)
	}
	if cfg.MaxConnLifetime <= 0 {
		return fmt.Errorf("max conn lifetime: %w", ErrPostgresInvalidMaxConnLifetime)
	}
	if cfg.MaxConnIdleTime <= 0 {
		return fmt.Errorf("max conn idle time: %w", ErrPostgresInvalidMaxConnIdleTime)
	}

	return nil
}

func newServerConfig() (http_server.Config, error) {
	host := viper.GetString("server.host")
	port := viper.GetInt("server.port")
	readTimeout := viper.GetString("server.read_timeout")
	parsedReadTimeout, err := time.ParseDuration(readTimeout)
	if err != nil {
		return http_server.Config{}, fmt.Errorf("read timeout: %w", ErrServerParseReadTimeout)
	}

	writeTimeout := viper.GetString("server.write_timeout")
	parsedWriteTimeout, err := time.ParseDuration(writeTimeout)
	if err != nil {
		return http_server.Config{}, fmt.Errorf("write timeout: %w", ErrServerParseWriteTimeout)
	}

	cfg := http_server.Config{
		Host:         host,
		Port:         port,
		ReadTimeout:  parsedReadTimeout,
		WriteTimeout: parsedWriteTimeout,
	}

	err = validateServerConfig(cfg)
	if err != nil {
		return http_server.Config{}, err
	}

	return cfg, nil
}

func validateServerConfig(cfg http_server.Config) error {
	if cfg.Host == "" {
		return fmt.Errorf("host: %w", ErrServerEmptyHost)
	}
	if cfg.Port < 0 || cfg.Port > 65535 {
		return fmt.Errorf("port: %w", ErrServerInvalidPort)
	}
	if cfg.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout: %w", ErrServerInvalidReadTimeout)
	}
	if cfg.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout: %w", ErrServerInvalidWriteTimeout)
	}

	return nil
}
