package data

import "github.com/upper/bond"

type History struct{}

type HistoryStore struct {
	bond.Store
}

func (store HistoryStore) CollectionName() string {
	return `linkedin_history`
}
