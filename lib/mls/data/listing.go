package data

import (
	"fmt"
	"time"

	db "upper.io/db.v2"

	"github.com/pxue/jarvis/lib/gmaps"
	"github.com/upper/bond"
)

type Listing struct {
	ID      int64  `db:"id,omitempty,pk" json:"id"`
	MLS     string `db:"mls_id" json:"mlsId"`
	Address string `db:"address" json:"address"`
	Unit    int    `db:"unit" json:"unit"`
	Price   int    `db:"price" json:"price"`

	Description  string          `db:"description" json:"description"`
	HasLocker    bool            `db:"has_locker" json:"hasLocker"`
	HasParking   bool            `db:"has_parking" json:"hasParking"`
	HasBeanField bool            `db:"has_beanfield" json:"hasBeanField"`
	Size         ListSize        `db:"apt_size" json:"aptSize"`
	Exposure     string          `db:"exposure" json:"exposure"`
	Distance     *gmaps.Distance `db:"distance,jsonb" json:"distance"`

	CreatedAt *time.Time `db:"created_at,omitempty" json:"createdAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

type ListSize uint32

const (
	_ ListSize = iota
	ListSizeMini
	ListSizeSmall
	ListSizeMedium
	ListSizeLarge
	ListSizeXLarge
)

var (
	listSizes = []string{
		"-",
		"0-499",
		"500-599",
		"600-699",
		"700-799",
		"1000-1999",
	}
)

type ListingStore struct {
	bond.Store
}

func (l *Listing) CollectionName() string {
	return `mls_listings`
}

func (store ListingStore) FindByMLS(mlsID string) (*Listing, error) {
	var listing *Listing
	if err := store.Find(db.Cond{"mls_id": mlsID}).One(&listing); err != nil {
		return nil, err
	}
	return listing, nil
}

func (s *ListSize) UnmarshalText(text string) error {
	enum := string(text)
	for i := 0; i < len(listSizes); i++ {
		if enum == listSizes[i] {
			*s = ListSize(i)
			return nil
		}
	}
	return fmt.Errorf("unknown size %s", enum)
}

// String returns the string value of the status.
func (s ListSize) String() string {
	return listSizes[s]
}
