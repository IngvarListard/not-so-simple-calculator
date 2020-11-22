package sqlstore_test

import (
	"encoding/json"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/models"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/store/fixtures"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHistoryRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t)
	defer teardown("history")

	s := sqlstore.New(db)
	h := models.TestHistory(t)
	assert.NoError(t, s.History().Create(h))
	assert.NotNil(t, h.ID)
}

func TestHistoryRepository_GetAllHistory(t *testing.T) {
	db, teardown := sqlstore.TestDB(t)
	defer teardown("history")

	expectedJSON, err := fixtures.PopulateHistory(db)
	if err != nil {
		t.Fatal(err)
	}

	s := sqlstore.New(db)
	h, err := s.History().GetAllHistory()
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(h)

	assert.Equal(t, expectedJSON, b)
}

func TestHistoryRepository_GetHistoryByTimeRange(t *testing.T) {
	db, teardown := sqlstore.TestDB(t)
	defer teardown("history")

	expectedJSON, err := fixtures.PopulateHistory(db)
	if err != nil {
		t.Fatal(err)
	}

	s := sqlstore.New(db)
	h, err := s.History().GetHistoryByTimeRange("2020-11-22T14:00:00+03:00", "2020-11-22T14:40:00+03:00")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(h)

	assert.Equal(t, expectedJSON, b)
}
