package store

import (
	"fmt"
	"sync"

	vinyl "github.com/platform-edn/vinyl/internal"
)

type RecordMap map[string]vinyl.Record

type Memory struct {
	Records RecordMap
	mutex   sync.RWMutex
}

func NewMemory(records ...vinyl.Record) *Memory {
	rmap := RecordMap{}

	for _, r := range records {
		rmap[r.Domain] = r
	}

	return &Memory{
		Records: rmap,
		mutex:   sync.RWMutex{},
	}
}

func (store *Memory) GetRecord(domain string) (*vinyl.Record, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	record, exist := store.Records[domain]
	if !exist {
		return nil, fmt.Errorf("GetRecord: %w", &MissingRecordError{
			Domain: domain,
		})
	}

	return &record, nil
}

func (store *Memory) RemoveRecord(domain string) (*vinyl.Record, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

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
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	records := []vinyl.Record{}
	for _, record := range store.Records {
		records = append(records, record)
	}

	return records, nil
}

func (store *Memory) CreateRecord(domain string, address string, ttl uint32) (*vinyl.Record, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

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
