package converter

import (
	"time"

	"github.com/shooketh/sakura/module/generator"
)

// ToTime returns the time when id was generated.
func ToTime(id int64) time.Time {
	ts := id >> (generator.DatacenterIDBits + generator.WorkerIDBits + generator.SequenceBits)
	d := time.Duration(ts * int64(time.Millisecond))
	return generator.SakuraEpoch.Add(d)
}

// Dump returns the structure of id.
func Dump(id int64) (t time.Time, datacenterID, workerID int64, sequence int64) {
	datacenterID = (id & (generator.MaxDatacenterID << generator.DatacenterIDShift)) >> generator.DatacenterIDShift
	workerID = (id & (generator.MaxWorkerID << generator.SequenceBits)) >> generator.SequenceBits
	sequence = id & generator.SequenceMask
	return ToTime(id), datacenterID, workerID, sequence
}
