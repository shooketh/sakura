package converter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/shooketh/sakura/common/location"
	"github.com/shooketh/sakura/module/converter"
	"github.com/stretchr/testify/assert"
)

func TestDump(t *testing.T) {
	testCases := []struct {
		id  int64
		ts  time.Time
		did int64
		wid int64
		seq int64
	}{
		{
			2708845965213696,
			time.Date(2022, 01, 8, 11, 23, 59, 206000000, location.UTC()),
			1,
			0,
			0,
		},
		{
			2708846166540288,
			time.Date(2022, 01, 8, 11, 23, 59, 254000000, location.UTC()),
			1,
			0,
			0,
		},
	}
	for _, tc := range testCases {
		ts, did, wid, seq := converter.Dump(tc.id)
		assert.Equal(t, ts, tc.ts, fmt.Sprintf("%d timestamp is not expected. got %s expected %s", tc.id, ts, tc.ts))
		assert.Equal(t, did, tc.did, fmt.Sprintf("%d datacenterID is not expected. got %d expected %d", tc.id, did, tc.did))
		assert.Equal(t, wid, tc.wid, fmt.Sprintf("%d workerID is not expected. got %d expected %d", tc.id, wid, tc.wid))
		assert.Equal(t, seq, tc.seq, fmt.Sprintf("%d sequence is not expected. got %d expected %d", tc.id, seq, tc.seq))
	}
}
