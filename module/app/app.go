package app

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/rs/zerolog"
	"github.com/shooketh/etcd-handler"
	"github.com/shooketh/sakura/common/config"
	"github.com/shooketh/sakura/common/errors"
	"github.com/shooketh/sakura/common/log"
	"github.com/shooketh/sakura/module/generator"
	"github.com/shooketh/sakura/module/grpc"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type App struct {
	Addr             string
	DatacenterID     int64
	WorkerID         int64
	WorkerPrefix     string          `json:"-"`
	EtcdHandler      handler.Handler `json:"-"`
	Logger           zerolog.Logger  `json:"-"`
	*grpc.Controller `json:"-"`
}

type datacenter struct {
	ID        int64 `json:"id"`
	Available bool  `json:"available"`
}

func New(ctx context.Context) (*App, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Config.App.Timeout*time.Second)
	defer cancel()

	var err error
	handler, err := handler.New(ctx, &handler.Option{
		Logger: log.Logger,
		Config: &clientv3.Config{
			Endpoints:   config.Config.Etcd.Endpoints,
			DialTimeout: time.Duration(config.Config.Etcd.Timeout) * time.Second,
			Username:    config.Config.Etcd.Username,
			Password:    config.Config.Etcd.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	app := App{
		Addr:         config.Config.GRPC.IP + ":" + config.Config.GRPC.Port,
		DatacenterID: config.Config.App.DatacenterID,
		WorkerID:     config.Config.App.WorkerID,
		WorkerPrefix: config.Config.App.WorkerPrefix,
		EtcdHandler:  handler,
		Logger:       log.Logger,
		Controller: &grpc.Controller{
			IP:   config.Config.GRPC.IP,
			Port: config.Config.GRPC.Port,
		},
	}

	if err := app.GetDatacenterID(ctx); err != nil {
		return nil, err
	}

	if err := app.GetWorkerID(ctx); err != nil {
		return nil, err
	}

	lastTime, err := app.GetLastGenerateTimeUnixMilli(ctx)
	if err != nil {
		return nil, err
	}

	app.Generator, err = generator.New(app.DatacenterID, app.WorkerID, lastTime)
	if err != nil {
		app.UnregisterWorker()
		return nil, err
	}

	return &app, nil
}

func (app *App) GetDatacenterID(ctx context.Context) error {
	if metadata.OnGCE() {
		z, err := metadata.Zone()
		if err != nil {
			return err
		}
		r := z[:strings.LastIndex(z, "-")]
		key := "/sakura/datacenter/gcp/region/" + r + "/zone/" + z

		resp, err := app.EtcdHandler.Get(ctx, key)
		if err != nil {
			return err
		}

		var dc datacenter
		for _, ev := range resp.Kvs {
			if err := json.Unmarshal(ev.Value, &dc); err != nil {
				return err
			}
		}
		if !dc.Available {
			return fmt.Errorf("gcp %s zone is not available", z)
		}

		app.DatacenterID = dc.ID

		return nil
	}

	if err := app.checkDatacenterID(); err != nil {
		return err
	}

	return nil
}

func (app *App) checkDatacenterID() error {
	if app.DatacenterID < 0 || app.DatacenterID > generator.MaxDatacenterID {
		return errors.InvalidDatacenterID
	}

	return nil
}

// GetWorkerID from etcd
// If the workerID is -1,then generate a unique workerID in the datacenterID range
// If the workerID is greater than -1, then register a workerID on Etcd
func (app *App) GetWorkerID(ctx context.Context) error {
	if app.WorkerID == -1 {
		err := app.generateWorkerID(ctx)
		if err != nil {
			return err
		}
	} else {
		if err := app.checkWorkerID(); err != nil {
			return err
		}

		key := fmt.Sprintf("%s/%d/%d", app.WorkerPrefix, app.DatacenterID, app.WorkerID)

		resp, err := app.EtcdHandler.PutNx(ctx, key, app.Addr)
		if err != nil {
			return err
		}

		if !resp.Succeeded {
			return errors.DuplicationWorkerID
		}
	}

	return nil
}

func (app *App) generateWorkerID(ctx context.Context) error {
	count := uint32(0)
	h := fnv.New32a()
	_, err := h.Write([]byte(app.Addr))
	if err != nil {
		return err
	}
	keyHash := h.Sum32()

	for {
		if count > generator.MaxWorkerID {
			return errors.OverflowWorkerID
		}

		app.WorkerID = int64((keyHash + count) & generator.MaxWorkerID)
		key := fmt.Sprintf("%s/%d/%d", app.WorkerPrefix, app.DatacenterID, app.WorkerID)

		resp, err := app.EtcdHandler.PutNx(ctx, key, app.Addr)
		if err != nil {
			return err
		}
		if !resp.Succeeded {
			count++
			continue
		}

		break
	}

	return nil
}

func (app *App) checkWorkerID() error {
	if app.WorkerID < 0 || app.WorkerID > generator.MaxWorkerID {
		return errors.InvalidWorkerID
	}

	return nil
}

func (app *App) GetLastGenerateTimeUnixMilli(ctx context.Context) (int64, error) {
	key := fmt.Sprintf("%s/%d/%d", config.Config.App.LastTimePrefix, app.DatacenterID, app.WorkerID)
	resp, err := app.EtcdHandler.Get(ctx, key)
	if err != nil {
		return 0, err
	}

	var lastTime int64
	for _, ev := range resp.Kvs {
		lastTime, _ = strconv.ParseInt(string(ev.Value), 10, 64)
	}

	return lastTime, nil
}

func (app *App) UnregisterWorker() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Config.App.Timeout*time.Second)
	defer cancel()

	if app.Generator != nil {
		lastTime := app.Generator.LastGenerateTimeUnixMilli()
		if lastTime != generator.SakuraEpoch.UnixMilli() {
			key := fmt.Sprintf("%s/%d/%d", config.Config.App.LastTimePrefix, app.DatacenterID, app.WorkerID)
			_, err := app.EtcdHandler.Put(ctx, key, strconv.FormatInt(lastTime, 10))
			if err != nil {
				return err
			}
		}
	}

	key := fmt.Sprintf("%s/%d/%d", config.Config.App.WorkerPrefix, app.DatacenterID, app.WorkerID)
	_, err := app.EtcdHandler.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
