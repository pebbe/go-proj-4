#include "proj_go.h"

const char *get_err() {
    int
        i;
    i = proj_errno(0);
    if (i)
        return proj_errno_string(i);
    else
        return NULL;
}

const char *get_proj_err(PJ *pj) {
    int
        i;
    i = proj_errno(pj);
    if (i)
        return proj_errno_string(i);
    else
        return NULL;
}

PJ *create(const char *definition) {
    return proj_create(0, definition);
}

const char *fwd(PJ *pj, double *x, double *y) {
    PJ_COORD
	in,
	out;

    in.uv.u = *x * DEG_TO_RAD;
    in.uv.v = *y * DEG_TO_RAD;

    out = proj_trans(pj, PJ_FWD, in);

    *x = out.uv.u;
    *y = out.uv.v;

    return get_proj_err(pj);
}

const char *inv(PJ *pj, double *x, double *y) {
    PJ_COORD
	in,
	out;

    in.uv.u = *x;
    in.uv.v = *y;

    out = proj_trans(pj, PJ_INV, in);

    *x = out.uv.u / DEG_TO_RAD;
    *y = out.uv.v / DEG_TO_RAD;

    return get_proj_err(pj);
}

const char *transform(PJ *srcdefn, PJ *dstdefn, long point_count, double *x, double *y, double *z) {
    size_t
	err;
    long
	i;
    PJ_COORD
	*coord;
    
    coord = (PJ_COORD *) malloc(point_count * sizeof(PJ_COORD));
    if (!coord)
	return "allocation failed";
    for (i = 0; i < point_count; i++) {
	coord[i].xyz.x = x[i];
	coord[i].xyz.y = y[i];
	if (z)
	    coord[i].xyz.z = z[i];
    }
    
    err = proj_trans_array(srcdefn, PJ_INV, point_count, coord);
    if (!err)
	err = proj_trans_array(dstdefn, PJ_FWD, point_count, coord);
    if (err) {
	free(coord);
	return proj_errno_string(err);
    }

    for (i = 0; i < point_count; i++) {
	x[i] = coord[i].xyz.x;
	y[i] = coord[i].xyz.y;
	if (z)
	    z[i] = coord[i].xyz.z;
    }

    free(coord);
    return NULL;
}

