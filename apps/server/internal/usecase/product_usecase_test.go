package usecase

import (
	"testing"

	"github.com/price-comparison/server/internal/domain"
)

type productRepoStub struct {
	lastLimit     int
	lastOffset    int
	lastSortField string
	lastSortOrder string
}

func (p *productRepoStub) FindAll(limit, offset int, sortField, sortOrder string) ([]domain.Product, error) {
	p.lastLimit = limit
	p.lastOffset = offset
	p.lastSortField = sortField
	p.lastSortOrder = sortOrder
	return []domain.Product{}, nil
}

func (p *productRepoStub) FindByID(id int) (*domain.Product, error) {
	return nil, nil
}

func (p *productRepoStub) Search(keyword string, limit, offset int, sortField, sortOrder string) ([]domain.Product, error) {
	p.lastLimit = limit
	p.lastOffset = offset
	p.lastSortField = sortField
	p.lastSortOrder = sortOrder
	return []domain.Product{}, nil
}

func (p *productRepoStub) ListCategories() ([]string, error) {
	return []string{}, nil
}

func TestProductSearchRequiresKeyword(t *testing.T) {
	stub := &productRepoStub{}
	uc := NewProductUsecase(stub, nil, 0)

	if _, err := uc.Search(ProductSearchOptions{}); err == nil {
		t.Fatalf("expected error for empty keyword")
	}
}

func TestProductListDefaults(t *testing.T) {
	stub := &productRepoStub{}
	uc := NewProductUsecase(stub, nil, 0)

	if _, err := uc.List(ProductListOptions{}); err != nil {
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
