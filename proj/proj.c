#include "proj.h"

char *get_err()
{
    int
        *i;
    i = pj_get_errno_ref();
    if (*i)
        return pj_strerrno(*i);
    else
        return NULL;
}

char *fwd(projPJ pj, double *x, double *y) {
    projUV
	p;

    p.u = *x * DEG_TO_RAD;
    p.v = *y * DEG_TO_RAD;
    p = pj_fwd(p, pj);

    *x = p.u;
    *y = p.v;

    return get_err();
}

char *inv(projPJ pj, double *x, double *y) {
    projUV
	p;

    p.u = *x;
    p.v = *y;
    p = pj_inv(p, pj);

    *x = p.u / DEG_TO_RAD;
    *y = p.v / DEG_TO_RAD;

    return get_err();
}

char *transform(projPJ srcdefn, projPJ dstdefn, long point_count, int point_offset, double *x, double *y, double *z, int has_z) {
    int
	err;
    if (has_z)
	err = pj_transform(srcdefn, dstdefn, point_count, point_offset, x, y, z);
    else
	err = pj_transform(srcdefn, dstdefn, point_count, point_offset, x, y, NULL);
    return err ? pj_strerrno(err) : NULL;
}

