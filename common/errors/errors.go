package errors

import "errors"

var (
	InvalidDatacenterID = errors.New("invalid datacenter id")
	InvalidWorkerID     = errors.New("invalid worker id")
	DuplicationWorkerID = errors.New("duplication worker id")
	OverflowWorkerID    = errors.New("no more workerId available in the datacenter")
	ClockRollback       = errors.New("system clock was rollbacked")
)
