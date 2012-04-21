#include "proj.h"

triple *transform2(projPJ srcdefn, projPJ dstdefn, long point_count, int point_offset, double x, double y) {
    static triple trip;
    trip.x = x;
    trip.y = y;
    trip.z = 0;
    trip.err = pj_transform(srcdefn, dstdefn, point_count, point_offset, &trip.x, &trip.y, NULL);
    return &trip;
}

triple *transform3(projPJ srcdefn, projPJ dstdefn, long point_count, int point_offset, double x, double y, double z) {
    static triple trip;
    trip.x = x;
    trip.y = y;
    trip.z = z;
    trip.err = pj_transform(srcdefn, dstdefn, point_count, point_offset, &trip.x, &trip.y, &trip.z);
    return &trip;
}

double triple_x(triple *t)
{
    return t->x;
}

double triple_y(triple *t)
{
    return t->y;
}

double triple_z(triple *t)
{
    return t->z;
}

char *triple_err(triple *t)
{
    if (t->err) {
        return pj_strerrno(t->err);
    }
    return "";
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