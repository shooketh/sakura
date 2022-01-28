package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/shooketh/sakura/common/config"
	"github.com/shooketh/sakura/common/location"
	"github.com/shooketh/sakura/common/log"
	"github.com/shooketh/sakura/module/app"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var exitCode = 0

func init() {
	time.Local = location.UTC()

	if err := config.Init(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := log.Init(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	app, err := app.New(ctx)
	if err != nil {
		log.Logger.Error().Err(err).Msg("new sakura app failed")
		os.Exit(1)
	}

	go signalHandler(ctx, stop, app)

	serverListen(ctx, app)

	wg.Add(1)
	go serverServe(ctx, &wg, app)

	key, leaseID := serviceRegister(ctx, app)

	go serviceKeepAlive(ctx, key, leaseID, app)

	wg.Wait()

	if err := app.UnregisterWorker(); err != nil {
		log.Logger.Error().Err(err).Msg("")
		exitCode = 1
	}

	log.Logger.Info().Msg("shutdown completed")

	os.Exit(exitCode)
}

func signalHandler(ctx context.Context, stop context.CancelFunc, app *app.App) {
	<-ctx.Done()

	stop()
	log.Logger.Info().Msg("shutting down, press Ctrl+C again to force")
	if app.Server == nil {
		log.Logger.Info().Msg("shutdown completed")
		if err := app.UnregisterWorker(); err != nil {
			log.Logger.Error().Err(err).Msg("")
		}
		os.Exit(0)
	}
	app.Server.Stop()
}

func serverListen(ctx context.Context, app *app.App) {
	if err := app.Listen(ctx); err != nil {
		log.Logger.Error().Err(err).Msg("")
		if err := app.UnregisterWorker(); err != nil {
			log.Logger.Error().Err(err).Msg("")
		}
		os.Exit(1)
	}
	log.Logger.Info().Msgf("server listening at %v", app.Listener.Addr())
}

func serverServe(ctx context.Context, wg *sync.WaitGroup, app *app.App) {
	defer wg.Done()
	if err := app.Serve(ctx); err != nil {
		log.Logger.Error().Err(err).Msg("")
		if err := app.UnregisterWorker(); err != nil {
			log.Logger.Error().Err(err).Msg("")
		}
		os.Exit(1)
	}
}

func serviceRegister(ctx context.Context, app *app.App) (key string, leaseID clientv3.LeaseID) {
	leaseResp, err := app.EtcdHandler.Grant(ctx, config.Config.Etcd.ServiceLeaseTTL)
	if err != nil {
		log.Logger.Error().Err(err).Msg("")
		if err := app.UnregisterWorker(); err != nil {
			log.Logger.Error().Err(err).Msg("")
		}
		os.Exit(1)
	}

	key = fmt.Sprintf("%s/%d/%d", config.Config.Etcd.ServicePrefix, app.DatacenterID, app.WorkerID)

	if err := app.EtcdHandler.Add(ctx, config.Config.Etcd.ServicePrefix, key, app.Addr, clientv3.WithLease(leaseResp.ID)); err != nil {
		log.Logger.Error().Err(err).Msg("")
		if err := app.UnregisterWorker(); err != nil {
			log.Logger.Error().Err(err).Msg("")
		}
		os.Exit(1)
	}
	log.Logger.Info().Msgf("registered the service[key=%s] in etcd", key)

	return key, leaseResp.ID
}

func serviceKeepAlive(ctx context.Context, key string, leaseID clientv3.LeaseID, app *app.App) {
	if err := app.EtcdHandler.KeepAlive(ctx, leaseID, key); err != nil {
		log.Logger.Error().Err(err).Msg("")
		app.Server.Stop()
		exitCode = 1
	}
}
