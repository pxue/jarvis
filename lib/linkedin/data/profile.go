package data

import "github.com/upper/bond"

type Profile struct{}

type ProfileStore struct {
	bond.Store
}

func (store ProfileStore) CollectionName() string {
	return `linkedin_profiles`
}
