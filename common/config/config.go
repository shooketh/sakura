package config

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Config = conf{}

type conf struct {
	App  app  `mapstructure:"app"`
	Log  log  `mapstructure:"log"`
	GRPC grpc `mapstructure:"grpc"`
	Etcd etcd `mapstructure:"etcd"`
}

type app struct {
	Env            string        `mapstructure:"env" env:"SAKURA_ENV" flag:"env"`
	Timeout        time.Duration `mapstructure:"timeout" env:"SAKURA_TIMEOUT" flag:"timeout"`
	DatacenterID   int64         `mapstructure:"datacenterID" env:"SAKURA_DETACENTER_ID" flag:"datacenter-id"`
	WorkerID       int64         `mapstructure:"workerID" env:"SAKURA_WORKER_ID" flag:"worker-id"`
	WorkerPrefix   string        `mapstructure:"workerPrefix" env:"SAKURA_WORKER_PREFIX" flag:"worker-prefix"`
	LastTimePrefix string        `mapstructure:"lastTimePrefix" env:"SAKURA_LAST_TIME_PREFIX" flag:"last-time-prifix"`
}

type log struct {
	Level    string `mapstructure:"level" env:"SAKURA_LOG_LEVEL" flag:"log-level"`
	Path     string `mapstructure:"path" env:"SAKURA_LOG_PATH" flag:"log-path"`
	MaxSize  int    `mapstructure:"maxSize" env:"SAKURA_LOG_MAX_SIZE" flag:"log-max-size"`
	MaxAge   int    `mapstructure:"maxAge" env:"SAKURA_LOG_MAX_AGE" flag:"log-max-age"`
	Compress bool   `mapstructure:"compress" env:"SAKURA_LOG_COMPRESS" flag:"log-compress"`
}

type grpc struct {
	IP   string `mapstructure:"ip" env:"SAKURA_GRPC_IP" flag:"ip"`
	Port string `mapstructure:"port" env:"SAKURA_GRPC_PORT" flag:"port"`
}

type etcd struct {
	Endpoints       []string      `mapstructure:"endpoints" env:"SAKURA_ETCD_ENDPOINTS" flag:"etcd-endpoints"`
	Timeout         time.Duration `mapstructure:"timeout" env:"SAKURA_ETCD_TIMEOUT" flag:"etcd-timeout"`
	Username        string        `mapstructure:"username" env:"SAKURA_ETCD_USERNAME" flag:"etcd-username"`
	Password        string        `mapstructure:"password" env:"SAKURA_ETCD_PASSWORD" flag:"etcd-password"`
	ServicePrefix   string        `mapstructure:"servicePrefix" env:"SAKURA_ETCD_SERVICE_PREFIX" flag:"etcd-service-prefix"`
	ServiceLeaseTTL int64         `mapstructure:"serviceLeaseTTL" env:"SAKURA_ETCD_SERVICE_LEASE_TTL" flag:"etcd-service-lease-ttl"`
}

func Init() error {
	Config.App.Env = os.Getenv("SAKURA_ENV")

	parseFlags()

	if flag.Lookup("env").Value.String() == "" {
		return fmt.Errorf("failed to set SAKURA_ENV or --env flag")
	}

	systemPath := "/etc/sakura/"
	currentPath := "../../config/"
	configType := "yaml"

	viper.SetConfigName(Config.App.Env)
	viper.AddConfigPath(systemPath)
	viper.AddConfigPath(currentPath)
	viper.SetConfigType(configType)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config info failed: %w", err)
	}

	bindEnvs(Config)

	bindFlags(Config)

	if viper.GetString("grpc.ip") == "" {
		ip, err := localIP()
		if err != nil {
			return err
		}
		viper.Set("grpc.ip", ip)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("unmarshal error config file: %w", err)
	}

	return nil
}

func parseFlags() {
	flag.StringVar(&Config.App.Env, "env", Config.App.Env, "sakura application environment")
	flag.Duration("timeout", Config.App.Timeout, "")
	flag.Int64("datacenter-id", Config.App.DatacenterID, "datacenter id (0 ~ 31)")
	flag.Int64("worker-id", Config.App.WorkerID, "worker id (0 ~ 31)")
	flag.String("worker-prefix", Config.App.WorkerPrefix, "")
	flag.String("last-time-prefix", Config.App.LastTimePrefix, "")
	flag.String("log-level", Config.Log.Level, "log level (panic, fatal, error, warn, info)")
	flag.String("log-path", Config.Log.Path, "")
	flag.Int("log-max-size", Config.Log.MaxSize, "")
	flag.Int("log-max-age", Config.Log.MaxAge, "")
	flag.Bool("log-compress", Config.Log.Compress, "")
	flag.String("ip", Config.GRPC.IP, "ip address to listen")
	flag.String("port", Config.GRPC.Port, "port to listen")
	flag.StringSlice("etcd-endpoints", Config.Etcd.Endpoints, "endpoints for connecting etcd")
	flag.Duration("etcd-timeout", Config.Etcd.Timeout, "etcd operation timeout(s)")
	flag.String("etcd-username", Config.Etcd.Username, "username for connecting etcd")
	flag.String("etcd-password", Config.Etcd.Password, "password for connecting etcd")
	flag.String("etcd-service-prefix", Config.Etcd.ServicePrefix, "")
	flag.Int64("etcd-service-lease-ttl", Config.Etcd.ServiceLeaseTTL, "")
	flag.Parse()
}

func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		e, _ := t.Tag.Lookup("env")
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		case reflect.Slice:
			ev := os.Getenv(e)
			if ev == "" {
				break
			}
			s := strings.Split(ev, ",")
			viper.Set(strings.Join(append(parts, tv), "."), s)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."), e)
		}
	}
}

func bindFlags(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		f, _ := t.Tag.Lookup("flag")
		switch v.Kind() {
		case reflect.Struct:
			bindFlags(v.Interface(), append(parts, tv)...)
		default:
			viper.BindPFlag(strings.Join(append(parts, tv), "."), flag.Lookup(f))
		}
	}
}

func localIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
			return ip.IP.String(), nil
		}
	}

	return "", fmt.Errorf("can not find the server ip address")
}
