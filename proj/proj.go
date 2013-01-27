/*
Package proj provides an interface to the Cartographic Projections Library PROJ.4.

See: http://trac.osgeo.org/proj/

Example usage:

    merc, err := proj.NewProj("+proj=merc +ellps=clrk66 +lat_ts=33")
    defer merc.Close() // if omitted, this will be called on garbage collection
    if err != nil {
        log.Fatal(err)
    }

    ll, err := proj.NewProj("+proj=latlong +ellps=clrk66")
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

	var err error = nil
	s := C.GoString(C.get_err())
	if s == "" {
		proj.opened = true
	} else {
		err = errors.New(s)
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
		return 0.0, 0.0, 0.0, errors.New("projection is closed")
	}

	var triple *C.triple
	if has_z {
		triple = C.transform(srcpj.pj, dstpj.pj, C.long(point_count), C.int(point_offset), C.double(x), C.double(y), C.double(z), C.int(0))
	} else {
		triple = C.transform(srcpj.pj, dstpj.pj, C.long(point_count), C.int(point_offset), C.double(x), C.double(y), C.double(0), C.int(0))
	}

	if triple == nil {
		return 0.0, 0.0, 0.0, errors.New("transform malloc failed")
	}
	defer C.free(unsafe.Pointer(triple))

	if e := C.GoString(triple.err); e != "" {
		return 0.0, 0.0, 0.0, errors.New(e)
	}

	return float64(triple.x), float64(triple.y), float64(triple.z), nil
}

func Transform2(srcpj, dstpj *Proj, point_count int64, point_offset int, x, y float64) (float64, float64, error) {
	xx, yy, _, err := transform(srcpj, dstpj, point_count, point_offset, x, y, 0, false)
	return xx, yy, err
}

func Transform3(srcpj, dstpj *Proj, point_count int64, point_offset int, x, y, z float64) (float64, float64, float64, error) {
	return transform(srcpj, dstpj, point_count, point_offset, x, y, z, true)
}

func DegToRad(deg float64) float64 {
	return deg / 180.0 * math.Pi
}

func RadToDeg(rad float64) float64 {
	return rad / math.Pi * 180.0
}
