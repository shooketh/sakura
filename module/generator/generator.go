// +--------+----------------+-----------------------------------------+-----------------+
// | unused | timestamp (ms) | machine id (10bit)                      | sequence number |
// +--------+----------------+----------------------+------------------+-----------------+
// |  (1bit)|         (41bit)| datacenter id (5bit) | worker id (5bit) |          (12bit)|
// +--------+----------------+----------------------+------------------+-----------------+

package generator

import (
	"sync"
	"time"

	"github.com/shooketh/sakura/common/errors"
	"github.com/shooketh/sakura/common/location"
)

const (
	TimeUnit          = int64(time.Millisecond)
	WorkerIDBits      = 5
	DatacenterIDBits  = 5
	SequenceBits      = 12
	MaxDatacenterID   = 1<<DatacenterIDBits - 1
	MaxWorkerID       = 1<<WorkerIDBits - 1
	SequenceMask      = 1<<SequenceBits - 1
	TimestampShift    = SequenceBits + WorkerIDBits + DatacenterIDBits
	DatacenterIDShift = SequenceBits + WorkerIDBits
)

var SakuraEpoch = time.Date(2022, 1, 1, 0, 0, 0, 0, location.UTC())

func init() {
	time.Local = location.UTC()
}

type Generator struct {
	mu           sync.Mutex
	epoch        time.Time
	time         int64
	datacenterID int64
	workerID     int64
	sequence     int64
}

func New(datacenterID, workerID, lastTime int64) (*Generator, error) {
	if datacenterID < 0 || datacenterID > MaxDatacenterID {
		return nil, errors.InvalidDatacenterID
	}

	if workerID < 0 || workerID > MaxWorkerID {
		return nil, errors.InvalidWorkerID
	}

	now := time.Now()
	if now.UnixMilli() <= lastTime {
		return nil, errors.ClockRollback
	}

	g := Generator{}
	g.datacenterID = datacenterID
	g.workerID = workerID
	g.epoch = now.Add(SakuraEpoch.Sub(now))

	return &g, nil
}

func (g *Generator) Generate() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Since(g.epoch).Nanoseconds() / TimeUnit

	if now == g.time {
		g.sequence = (g.sequence + 1) & SequenceMask

		if g.sequence == 0 {
			for now <= g.time {
				now = time.Since(g.epoch).Nanoseconds() / TimeUnit
			}
		}
	} else {
		g.sequence = 0
	}

	g.time = now

	return now<<TimestampShift |
		g.datacenterID<<DatacenterIDShift |
		g.workerID<<SequenceBits |
		g.sequence
}

func (g *Generator) LastGenerateTimeUnixMilli() int64 {
	t := SakuraEpoch.Add(time.Duration(g.time * TimeUnit))
	return t.UnixMilli()
}
