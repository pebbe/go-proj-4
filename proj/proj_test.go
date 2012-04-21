package proj

import (
	"fmt"
	"testing"
)

func TestMercToLat(t *testing.T) {

	merc, err := NewProj("+proj=merc +ellps=clrk66 +lat_ts=33")
	defer merc.Close()
	if err != nil {
		t.Error(err)
	}

	ll, err := NewProj("+proj=latlong +ellps=clrk66")
	defer ll.Close()
	if err != nil {
		t.Error(err)
	}

	x, y, err := Transform2(ll, merc, 1, 1, DegToRad(-16), DegToRad(20.25))
	if err != nil {
		t.Error(err)
	} else {
	    s := fmt.Sprintf("%.2f %.2f", x, y)
    	s1 := "-1495284.21 1920596.79"
    	if s != s1 {
    		t.Errorf("MercToLat = %v, want %v", s, s1)
    	}
    }
}
