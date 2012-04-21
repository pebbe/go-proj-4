/*
Package proj provides an interface to the Cartographic Projections Library PROJ.4.

See: http://trac.osgeo.org/proj/

Example usage:

    merc, err := NewProj("+proj=merc +ellps=clrk66 +lat_ts=33")
    defer merc.Close()
    if err != nil {
        log.Fatal(err)
    }

    ll, err := NewProj("+proj=latlong +ellps=clrk66")
    defer ll.Close()
    if err != nil {
        log.Fatal(err)
    }

    x, y, err := Transform2(ll, merc, 1, 1, DegToRad(-16), DegToRad(20.25))
    if err != nil {
        t.Fatal(err)
    }
    fmt.Printf("%.2f %.2f", x, y)  // should print: -1495284.21 1920596.79
*/
package proj

/*
#cgo LDFLAGS: -lproj
#include "proj.h"
*/
import "C"

import (
	"errors"
    "math"
	"unsafe"
)

type Proj struct {
	pj     _Ctype_projPJ
	opened bool
}

func NewProj(definition string) (*Proj, error) {
	cs := C.CString(definition)
	defer C.free(unsafe.Pointer(cs))
	proj := Proj{opened: false}
	proj.pj = C.pj_init_plus(cs)

	err := getError()
	if err == nil {
		proj.opened = true
	}
	return &proj, err
}

func (p *Proj) Close() {
    if p.opened {
        C.pj_free(p.pj)
        p.opened = false
    }
}

func Transform2(srcpj, dstpj *Proj, point_count int64, point_offset int, x, y float64) (float64, float64, error) {
    if !(srcpj.opened && dstpj.opened) {
        return 0.0, 0.0, errors.New("projection is closed")
    }
	triple := C.transform2(srcpj.pj, dstpj.pj, C.long(point_count), C.int(point_offset), C.double(x), C.double(y))
	if e := C.GoString(C.triple_err(triple)); e != "" {
        return 0.0, 0.0, errors.New(e)
	}
    return float64(C.triple_x(triple)), float64(C.triple_y(triple)), nil
}

func Transform3(srcpj, dstpj *Proj, point_count int64, point_offset int, x, y, z float64) (float64, float64, float64, error) {
    if !(srcpj.opened && dstpj.opened) {
        return 0.0, 0.0, 0.0, errors.New("projection is closed")
    }
    triple := C.transform3(srcpj.pj, dstpj.pj, C.long(point_count), C.int(point_offset), C.double(x), C.double(y), C.double(z))
    if e := C.GoString(C.triple_err(triple)); e != "" {
        return 0.0, 0.0, 0.0, errors.New(e)
    }
    return float64(C.triple_x(triple)), float64(C.triple_y(triple)), float64(C.triple_z(triple)), nil
}

func getError() error {
	s := C.GoString(C.get_err())
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func DegToRad(deg float64) float64 {
    return deg / 180.0 * math.Pi
}

func RadToDeg(rad float64) float64 {
    return rad / math.Pi * 180.0
}