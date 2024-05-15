package fsm

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

type Fsm struct {
	Db *sync.Map
}

type setPayload struct {
	Key   string
	Value string
}

func (kf *Fsm) Apply(log *raft.Log) any {
	switch log.Type {
	case raft.LogCommand:
		var sp setPayload
		err := json.Unmarshal(log.Data, &sp)
		if err != nil {
			return fmt.Errorf("could not parse payload: %s", err)
		}

		kf.Db.Store(sp.Key, sp.Value)
	default:
		return fmt.Errorf("unknown node log type: %#v", log.Type)
	}

	return nil
}

func (kf *Fsm) Restore(rc io.ReadCloser) error {
	// Must always restore from a clean state!!
	kf.Db.Range(func(key any, _ any) bool {
		kf.Db.Delete(key)
		return true
	})

	decoder := json.NewDecoder(rc)

	for decoder.More() {
		var sp setPayload
		err := decoder.Decode(&sp)
		if err != nil {
			return fmt.Errorf("could not decode payload: %s", err)
		}

		kf.Db.Store(sp.Key, sp.Value)
	}

	return rc.Close()
}

type snapshotNoop struct{}

func (sn snapshotNoop) Persist(_ raft.SnapshotSink) error { return nil }
func (sn snapshotNoop) Release()                          {}

func (kf *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	return snapshotNoop{}, nil
}
