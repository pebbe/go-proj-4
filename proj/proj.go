/*
Package proj provides an interface to the Cartographic Projections Library PROJ.4.

See: http://trac.osgeo.org/proj/

Example usage:

    merc, err := proj.NewProj("+proj=merc +ellps=clrk66 +lat_ts=33")
    defer merc.Close() // if omitted, this will be called on garbage collection
    if err != nil {
        log.Fatal(err)
    }

    ll, err := proj.NewProj("+proj=latlong")
    defer ll.Close()
    if err != nil {
        log.Fatal(err)
    }

    x, y, err := proj.Transform2(ll, merc, 1, 1, proj.DegToRad(-16), proj.DegToRad(20.25))
    if err != nil {
        log.Fatal(err)
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
	"runtime"
	"sync"
	"unsafe"
)

var (
	mu sync.Mutex
)

type Proj struct {
	pj     _Ctype_projPJ
	opened bool
}

func NewProj(definition string) (*Proj, error) {
	cs := C.CString(definition)
	defer C.free(unsafe.Pointer(cs))
	proj := Proj{opened: false}

	mu.Lock()
	proj.pj = C.pj_init_plus(cs)
	errstring := C.GoString(C.get_err())
	mu.Unlock()

	var err error = nil
	if errstring == "" {
		proj.opened = true
	} else {
		err = errors.New(errstring)
	}

	runtime.SetFinalizer(&proj, (*Proj).Close)

	return &proj, err
}

func (p *Proj) Close() {
	if p.opened {
		C.pj_free(p.pj)
		p.opened = false
	}
}

func transform(srcpj, dstpj *Proj, point_count int64, point_offset int, x, y, z float64, has_z bool) (float64, float64, float64, error) {
	if !(srcpj.opened && dstpj.opened) {
		return math.NaN(), math.NaN(), math.NaN(), errors.New("projection is closed")
	}

	x1 := C.double(x)
	y1 := C.double(y)
	z1 := C.double(z)

	var e *C.char
	if has_z {
		e = C.transform(srcpj.pj, dstpj.pj, C.long(point_count), C.int(point_offset), &x1, &y1, &z1, C.int(1))
	} else {
		e = C.transform(srcpj.pj, dstpj.pj, C.long(point_count), C.int(point_offset), &x1, &y1, &z1, C.int(0))
	}

	if e != nil {
		return math.NaN(), math.NaN(), math.NaN(), errors.New(C.GoString(e))
	}

	return float64(x1), float64(y1), float64(z1), nil
}

func Transform2(srcpj, dstpj *Proj, point_count int64, point_offset int, x, y float64) (float64, float64, error) {
	xx, yy, _, err := transform(srcpj, dstpj, point_count, point_offset, x, y, 0, false)
	return xx, yy, err
}

func Transform3(srcpj, dstpj *Proj, point_count int64, point_offset int, x, y, z float64) (float64, float64, float64, error) {
	return transform(srcpj, dstpj, point_count, point_offset, x, y, z, true)
}

// Longitude and latitude in degrees
func Fwd(proj *Proj, long, lat float64) (x float64, y float64, err error) {
	if !proj.opened {
		return math.NaN(), math.NaN(), errors.New("projection is closed")
	}
	x1 := C.double(long)
	y1 := C.double(lat)
	e := C.fwd(proj.pj, &x1, &y1)
	if e != nil {
		return math.NaN(), math.NaN(), errors.New(C.GoString(e))
	}
	return float64(x1), float64(y1), nil
}

// Longitude and latitude in degrees
func Inv(proj *Proj, x, y float64) (lat float64, long float64, err error) {
	if !proj.opened {
		return math.NaN(), math.NaN(), errors.New("projection is closed")
	}
	x2 := C.double(x)
	y2 := C.double(y)
	e := C.inv(proj.pj, &x2, &y2)
	if e != nil {
		return math.NaN(), math.NaN(), errors.New(C.GoString(e))
	}
	return float64(x2), float64(y2), nil
}

func DegToRad(deg float64) float64 {
	return deg / 180.0 * math.Pi
}

func RadToDeg(rad float64) float64 {
	return rad / math.Pi * 180.0
}
