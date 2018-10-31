package scheduler


import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestRange(t *testing.T) {
	ranges := []struct {
		expr     string
		min, max uint
		expected uint64
	}{
		{"5", 0, 7, 1 << 5},
		{"0", 0, 7, 1 << 0},
		{"7", 0, 7, 1 << 7},

		{"5-5", 0, 7, 1 << 5},
		{"5-6", 0, 7, 1<<5 | 1<<6},
		{"5-7", 0, 7, 1<<5 | 1<<6 | 1<<7},

		{"5-6/2", 0, 7, 1 << 5},
		{"5-7/2", 0, 7, 1<<5 | 1<<7},
		{"5-7/1", 0, 7, 1<<5 | 1<<6 | 1<<7},

		{"*", 1, 3, 1<<1 | 1<<2 | 1<<3 | starBit},
		{"*/2", 1, 3, 1<<1 | 1<<3 | starBit},
	}

	for _, c := range ranges {
		actual := getRange(c.expr, bounds{c.min, c.max, nil})
		if actual != c.expected {
			t.Errorf("%s => (expected) %d != %d (actual)", c.expr, c.expected, actual)
		}
	}
}

func TestField(t *testing.T) {
	fields := []struct {
		expr     string
		min, max uint
		expected uint64
	}{
		{"5", 1, 7, 1 << 5},
		{"5,6", 1, 7, 1<<5 | 1<<6},
		{"5,6,7", 1, 7, 1<<5 | 1<<6 | 1<<7},
		{"1,5-7/2,3", 1, 7, 1<<1 | 1<<5 | 1<<7 | 1<<3},
	}

	for _, c := range fields {
		actual := getField(c.expr, bounds{c.min, c.max, nil})
		if actual != c.expected {
			t.Errorf("%s => (expected) %d != %d (actual)", c.expr, c.expected, actual)
		}
	}
}

func TestBits(t *testing.T) {
	allBits := []struct {
		r        bounds
		expected uint64
	}{
		{minutes, 0xfffffffffffffff}, // 0-59: 60 ones
		{hours, 0xffffff},            // 0-23: 24 ones
		{dom, 0xfffffffe},            // 1-31: 31 ones, 1 zero
		{months, 0x1ffe},             // 1-12: 12 ones, 1 zero
		{dow, 0x7f},                  // 0-6: 7 ones
	}

	for _, c := range allBits {
		actual := all(c.r) // all() adds the starBit, so compensate for that..
		if c.expected|starBit != actual {
			t.Errorf("%d-%d/%d => (expected) %b != %b (actual)",
				c.r.min, c.r.max, 1, c.expected|starBit, actual)
		}
	}

	bits := []struct {
		min, max, step uint
		expected       uint64
	}{

		{0, 0, 1, 0x1},
		{1, 1, 1, 0x2},
		{1, 5, 2, 0x2a}, // 101010
		{1, 4, 2, 0xa},  // 1010
	}

	for _, c := range bits {
		actual := getBits(c.min, c.max, c.step)
		if c.expected != actual {
			t.Errorf("%d-%d/%d => (expected) %b != %b (actual)",
				c.min, c.max, c.step, c.expected, actual)
		}
	}
}

func TestSpecSchedule(t *testing.T) {
	entries := []struct {
		expr     string
		expected *SpecSchedule
	}{
		{"5 * * * *", &SpecSchedule{1 << 1, 1 << 5, all(hours), all(dom), all(months), all(dow)}},
		{"@every 5m", &SpecSchedule{1 << 1, 1 << 5, all(hours), all(dom), all(months), all(dow)}},
	}

	for _, c := range entries {
		actual, err := Parse(c.expr)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("%s => (expected) %b != %b (actual)", c.expr, c.expected, actual)
		}
	}

	actual, err := Parse("5 */8 2-6 1,3,5,7,9 *")
	if err != nil {
		t.Error(err)
	}

	start := time.Date(2009, 1, 1, 0, 0, 1, 0, time.Local)
	fmt.Println(start.Format("2006-01-02 15:04"))
	end := time.Date(2010, 1, 1, 0, 0, 0, 0, time.Local)
	for ; end.After(start); start = start.Add(1 * time.Minute) {
		if actual.IsTimeMatches(start) {
			fmt.Println(start.Format("2006-01-02 15:04"))
		}
	}
}
