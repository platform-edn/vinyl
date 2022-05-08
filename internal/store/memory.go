package store

import (
	"fmt"

	vinyl "github.com/platform-edn/vinyl/internal"
)

type RecordMap map[string]vinyl.Record

type Memory struct {
	Records RecordMap
}

func NewMemory(records ...vinyl.Record) *Memory {
	rmap := RecordMap{}

	for _, r := range records {
		rmap[r.Domain] = r
	}

	return &Memory{
		Records: rmap,
	}
}

func (store *Memory) GetRecord(domain string) (*vinyl.Record, error) {
	record, exist := store.Records[domain]
	if !exist {
		return nil, fmt.Errorf("GetRecord: %w", &MissingRecordError{
			Domain: domain,
		})
	}

	return &record, nil
}

func (store *Memory) RemoveRecord(domain string) (*vinyl.Record, error) {
	record, exist := store.Records[domain]
	if !exist {
		return nil, fmt.Errorf("RemoveRecord: %w", &MissingRecordError{
			Domain: domain,
		})
	}

	delete(store.Records, domain)

	return &record, nil
}

func (store *Memory) ListRecords() ([]vinyl.Record, error) {
	records := []vinyl.Record{}
	for _, record := range store.Records {
		records = append(records, record)
	}

	return records, nil
}

func (store *Memory) CreateRecord(domain string, address string, ttl uint32) (*vinyl.Record, error) {
	_, exist := store.Records[domain]
	if exist {
		return nil, fmt.Errorf("CreateRecord: %w", &ExistingRecordError{
			Domain: domain,
		})
	}

	record, err := vinyl.NewRecord(domain, address, ttl)
	if err != nil {
		return nil, fmt.Errorf("CreateRecord: %w", err)
	}

	store.Records[domain] = *record
	return record, nil
}
