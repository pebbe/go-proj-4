package proj

import (
	"fmt"
	"testing"
)

func TestMercToLatlong(t *testing.T) {
	merc, err := NewProj("+proj=merc +ellps=clrk66 +lat_ts=33")
	defer merc.Close()
	if err != nil {
		t.Error(err)
	}

	ll, err := NewProj("+proj=latlong")
	defer ll.Close()
	if err != nil {
		t.Error(err)
	}

	x, y, err := Transform2(ll, merc, DegToRad(-16), DegToRad(20.25))
	if err != nil {
		t.Error(err)
	} else {
		s := fmt.Sprintf("%.2f %.2f", x, y)
		s1 := "-1495284.21 1920596.79"
		if s != s1 {
			t.Errorf("MercToLatlong = %v, want %v", s, s1)
		}
	}
}

func TestInvalidErrorProblem(t *testing.T) {
	merc, err := NewProj("+proj=merc +ellps=clrk66 +lat_ts=33")
	defer merc.Close()
	if err != nil {
		t.Error(err)
	}

	ll, err := NewProj("+proj=latlong")
	defer ll.Close()
	if err != nil {
		t.Error(err)
	}

	_, _, err = Transform2(ll, merc, DegToRad(3000), DegToRad(500))
	if err == nil {
		t.Error("err should not be nil")
	}

	// Try create a new projection after an error
	merc2, err := NewProj("+proj=merc +ellps=clrk66 +lat_ts=33")
	defer merc2.Close()
	if err != nil {
		t.Error(err)
	}
}
