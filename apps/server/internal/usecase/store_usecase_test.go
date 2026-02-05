package usecase

import (
	"testing"

	"github.com/price-comparison/server/internal/domain"
	"github.com/price-comparison/server/internal/query"
)

type storeRepoStub struct {
	lastLimit     int
	lastOffset    int
	lastSortField string
	lastSortOrder string
}

func (s *storeRepoStub) FindNearby(lat, lon float64, radiusMeters int, limit, offset int) ([]domain.Store, error) {
	return nil, nil
}

func (s *storeRepoStub) FindAll(filters query.StoreFilters, limit, offset int, sortField, sortOrder string) ([]domain.Store, error) {
	s.lastLimit = limit
	s.lastOffset = offset
	s.lastSortField = sortField
	s.lastSortOrder = sortOrder
	return []domain.Store{}, nil
}

func (s *storeRepoStub) FindByID(id int) (*domain.Store, error) {
	return nil, nil
}

func TestStoreListDefaults(t *testing.T) {
	stub := &storeRepoStub{}
	uc := NewStoreUsecase(stub, nil, 0)

	if _, err := uc.List(StoreListOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stub.lastLimit != DefaultLimit {
		t.Fatalf("expected limit %d, got %d", DefaultLimit, stub.lastLimit)
	}
	if stub.lastOffset != 0 {
		t.Fatalf("expected offset 0, got %d", stub.lastOffset)
	}
	if stub.lastSortField != "name" {
		t.Fatalf("expected sort field name, got %s", stub.lastSortField)
	}
	if stub.lastSortOrder != "ASC" {
		t.Fatalf("expected sort order ASC, got %s", stub.lastSortOrder)
	}
}
