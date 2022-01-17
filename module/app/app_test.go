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
	defaultIP              = "192.168.1.1"
	defaultPort            = "50051"
	defaultEndpoints       = []string{"127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"}
	defaultTimeout         = 5 * time.Second
	workerLeaseTTL   int64 = 100
	workerPrefix           = "/sakura/worker"
	logger           zerolog.Logger
)

func init() {
	time.Local = location.UTC()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	o := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
	logger = zerolog.New(o).With().Timestamp().Caller().Logger()
}

func TestDuplicatedWorkerID(t *testing.T) {
	var datacenterID int64 = 1
	var workerID int64 = 0

	handler, err := handler.New(&handler.Option{
		Logger: logger,
		EtcdOption: &handler.EtcdOption{
			Endpoints: defaultEndpoints,
			Timeout:   defaultTimeout,
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

	err = app.GetWorkerID(context.TODO())
	assert.Nil(t, err)
	err = app.GetWorkerID(context.TODO())
	assert.Equal(t, err, errors.DuplicationWorkerID)
	_, err = app.EtcdHandler.Delete(context.TODO(), fmt.Sprintf("%s/%d/%d", workerPrefix, datacenterID, workerID))
	assert.Nil(t, err)
}

func TestGenerateWorkerID(t *testing.T) {
	var datacenterID int64 = 1
	var workerID int64 = -1

	handler, err := handler.New(&handler.Option{
		Logger: logger,
		EtcdOption: &handler.EtcdOption{
			Endpoints: defaultEndpoints,
			Timeout:   defaultTimeout,
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

	err = app.GetWorkerID(context.TODO())
	t.Logf("workerID is: %d", workerID)
	assert.Nil(t, err)

	_, err = handler.Delete(context.TODO(), fmt.Sprintf("%s/%d/%d", workerPrefix, datacenterID, workerID))
	assert.Nil(t, err)
}

func TestWorkerIDOverflow(t *testing.T) {
	var datacenterID int64 = 1
	var workerID int64 = -1

	handler, err := handler.New(&handler.Option{
		Logger: logger,
		EtcdOption: &handler.EtcdOption{
			Endpoints: defaultEndpoints,
			Timeout:   defaultTimeout,
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
		err = app.GetWorkerID(context.TODO())
		assert.Nil(t, err)
		if err != nil {
			break
		}
		t.Logf("current workerID is: %d", workerID)
		app.WorkerID = -1
	}

	err = app.GetWorkerID(context.TODO())
	assert.Equal(t, err, errors.OverflowWorkerID)

	withRange := clientv3.WithRange(fmt.Sprintf("%s/%d/%s", workerPrefix, datacenterID, "\\0"))
	_, err = handler.Delete(context.TODO(), fmt.Sprintf("%s/%d/%d", workerPrefix, datacenterID, 0), withRange)
	assert.Nil(t, err)
}
