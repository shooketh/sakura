package app_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/shooketh/etcd-handler"
	"github.com/shooketh/sakura/common/errors"
	"github.com/shooketh/sakura/common/location"
	"github.com/shooketh/sakura/module/app"
	"github.com/shooketh/sakura/module/generator"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	defaultIP        = "192.168.1.1"
	defaultPort      = "50051"
	defaultEndpoints = []string{"127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"}
	defaultTimeout   = 5 * time.Second
	workerPrefix     = "/sakura/worker"
	logger           zerolog.Logger
)

func init() {
	time.Local = location.UTC()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	o := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
	logger = zerolog.New(o).With().Timestamp().Caller().Logger()
}

func TestDuplicatedWorkerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var datacenterID int64 = 1
	var workerID int64 = 0

	handler, err := handler.New(ctx, &handler.Option{
		Logger: logger,
		Config: &clientv3.Config{
			Endpoints:   defaultEndpoints,
			DialTimeout: defaultTimeout,
		},
	})
	assert.Nil(t, err)

	app := app.App{
		Addr:         defaultIP + ":" + defaultPort,
		DatacenterID: datacenterID,
		WorkerID:     workerID,
		WorkerPrefix: workerPrefix,
		EtcdHandler:  handler,
		Logger:       logger,
	}

	err = app.GetWorkerID(ctx)
	assert.Nil(t, err)
	err = app.GetWorkerID(ctx)
	assert.Equal(t, err, errors.DuplicationWorkerID)
	_, err = app.EtcdHandler.Delete(ctx, fmt.Sprintf("%s/%d/%d", workerPrefix, app.DatacenterID, app.WorkerID))
	assert.Nil(t, err)
}

func TestGenerateWorkerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var datacenterID int64 = 1
	var workerID int64 = -1

	handler, err := handler.New(ctx, &handler.Option{
		Logger: logger,
		Config: &clientv3.Config{
			Endpoints:   defaultEndpoints,
			DialTimeout: defaultTimeout,
		},
	})
	assert.Nil(t, err)

	app := app.App{
		Addr:         defaultIP + ":" + defaultPort,
		DatacenterID: datacenterID,
		WorkerID:     workerID,
		WorkerPrefix: workerPrefix,
		EtcdHandler:  handler,
		Logger:       logger,
	}

	err = app.GetWorkerID(ctx)
	t.Logf("workerID is: %d", app.WorkerID)
	assert.Nil(t, err)

	_, err = handler.Delete(ctx, fmt.Sprintf("%s/%d/%d", workerPrefix, app.DatacenterID, app.WorkerID))
	assert.Nil(t, err)
}

func TestWorkerIDOverflow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var datacenterID int64 = 1
	var workerID int64 = -1

	handler, err := handler.New(ctx, &handler.Option{
		Logger: logger,
		Config: &clientv3.Config{
			Endpoints:   defaultEndpoints,
			DialTimeout: defaultTimeout,
		},
	})
	assert.Nil(t, err)

	app := app.App{
		Addr:         defaultIP + ":" + defaultPort,
		DatacenterID: datacenterID,
		WorkerID:     workerID,
		WorkerPrefix: workerPrefix,
		EtcdHandler:  handler,
		Logger:       logger,
	}

	for i := 0; i <= generator.MaxWorkerID; i++ {
		err = app.GetWorkerID(ctx)
		assert.Nil(t, err)
		if err != nil {
			break
		}
		t.Logf("current workerID is: %d", app.WorkerID)
		app.WorkerID = -1
	}

	err = app.GetWorkerID(ctx)
	assert.Equal(t, err, errors.OverflowWorkerID)

	withRange := clientv3.WithRange(fmt.Sprintf("%s/%d/%s", workerPrefix, app.DatacenterID, "\\0"))
	_, err = handler.Delete(ctx, fmt.Sprintf("%s/%d/%d", workerPrefix, app.DatacenterID, 0), withRange)
	assert.Nil(t, err)
}
