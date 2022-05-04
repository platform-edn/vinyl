package discovery

import vinyl "github.com/platform-edn/vinyl/internal"

type RecordStorer interface {
	CreateRecord(string, string, uint32) (*vinyl.Record, error)
	RemoveRecord(string) (*vinyl.Record, error)
	ListRecords() ([]vinyl.Record, error)
	GetRecord(string) (*vinyl.Record, error)
}
