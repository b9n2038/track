package cli

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestParseDayID(t *testing.T) {
	got, err := parseDayID("25w08-1")
	if err != nil {
		t.Errorf("failed to parse")
	}
	want := time.Date(2025, time.Month(2), 17, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, got.YearDay(), want.YearDay())
	//t.Errorf("Abs(-1) = %d; want 1")
}
