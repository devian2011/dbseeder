package fake

import (
	"testing"
	"time"
)

func TestGenerateStr(t *testing.T) {
	strAny, _ := Generate("string 10")
	str, isConverted := strAny.(string)
	if !isConverted {
		t.Errorf("wrong type for string")
	}
	if len(str) != 10 {
		t.Errorf("wrong string generation, str len should be eq 10")
	}
}

func TestGenerateDate(t *testing.T) {
	dateAny, _ := Generate("date 2022-10-11 2022-11-11")
	date, isConverted := dateAny.(time.Time)
	if !isConverted {
		t.Errorf("wrong type for date generation")
	}
	begin, _ := time.Parse("2006-01-02", "2022-10-10")
	end, _ := time.Parse("2006-01-02", "2022-11-12")

	if !begin.Before(date) && !end.After(date) {
		t.Errorf("date should be betwee 2022-10-10 and 2022-11-12. Current: %s", date.Format("2006-01-02 15:03:04"))
	}
}
