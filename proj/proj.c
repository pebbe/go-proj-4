#include "proj.h"

triple *transform(projPJ srcdefn, projPJ dstdefn, long point_count, int point_offset, double x, double y, double z, int has_z) {
    triple
      *trip;
    int
	err;

    trip = (triple *) malloc (sizeof (triple));
    if (! trip)
	return NULL;

    trip->x = x;
    trip->y = y;
    trip->z = z;
    if (has_z)
	err = pj_transform(srcdefn, dstdefn, point_count, point_offset, &(trip->x), &(trip->y), &(trip->z));
    else
	err = pj_transform(srcdefn, dstdefn, point_count, point_offset, &(trip->x), &(trip->y), NULL);
    trip->err = err ? pj_strerrno(err) : "";
    return trip;
}

void fwd(projPJ pj, double *x, double *y) {
    projUV
	p;

    p.u = *x * DEG_TO_RAD;
    p.v = *y * DEG_TO_RAD;
    p = pj_fwd(p, pj);

    *x = p.u;
    *y = p.v;
}

void inv(projPJ pj, double *x, double *y) {
    projUV
	p;

    p.u = *x;
    p.v = *y;
    p = pj_inv(p, pj);

    *x = p.u / DEG_TO_RAD;
    *y = p.v / DEG_TO_RAD;
}

char *get_err()
{
    int
        *i;
    i = pj_get_errno_ref();
    if (*i)
        return pj_strerrno(*i);
    else
        return "";
}
