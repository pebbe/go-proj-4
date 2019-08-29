#include <proj.h>
#include <stdlib.h>

#define RAD_TO_DEG    57.295779513082321
#define DEG_TO_RAD   .017453292519943296

PJ *create(const char *definition);
const char *transform(PJ *srcdefn, PJ *dstdefn, long point_count, double *x, double *y, double *z);
const char *fwd(PJ *src, double *lng, double *lat);
const char *inv(PJ *dst, double *lng, double *lat);
const char *get_err(void);
const char *get_proj_err(PJ *pj);
