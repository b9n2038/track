package rating

import (
	"fmt"
	"time"
)

type Rating int

const (
	Bad Rating = iota + 1
	Poor
	Fair
	Good
	Awesome
)

func (r Rating) String() string {
	labels := map[Rating]string{
		Bad:     "Bad",
		Poor:    "Poor",
		Fair:    "Fair",
		Good:    "Good",
		Awesome: "Awesome",
	}
	if label, exists := labels[r]; exists {
		return label
	}
	return "Invalid Rating"
}

func (r Rating) Emoji() string {
	emojis := map[Rating]string{
		Bad:     "💩",
		Poor:    "😠",
		Fair:    "😐",
		Good:    "😊",
		Awesome: "🤩",
	}
	if emoji, exists := emojis[r]; exists {
		return emoji
	}
	return "❌"
}

func (r Rating) IsValid() bool {
	return r >= Bad && r <= Awesome
}

type DayRating struct {
	ID     string
	Date   time.Time
	Rating Rating
}

func (dr DayRating) Label() string {
	yr := dr.Date.Format("06")   // Last two digits of year
	_, week := dr.Date.ISOWeek() // Get ISO week number
	weekday := dr.Date.Weekday() // Get day of week (0-6)
	return fmt.Sprintf("%sw%02d-%d", yr, week, weekday)
}

func (dr DayRating) String() string {
	return fmt.Sprintf("%s: %s %s", dr.Label(), dr.Rating.String(), dr.Rating.Emoji())
}

// NewDayRating creates a new DayRating with validation
func NewDayRating(date time.Time, rating Rating) (DayRating, error) {
	if !rating.IsValid() {
		return DayRating{}, fmt.Errorf("invalid rating value: %d", rating)
	}
	return DayRating{
		Date:   date,
		Rating: rating,
	}, nil
}

// NewRating creates a new Rating with validation
func NewRating(value int) (Rating, error) {
	rating := Rating(value)
	if !rating.IsValid() {
		return 0, fmt.Errorf("rating must be between %d and %d, got %d", Bad, Awesome, value)
	}
	return rating, nil
}
