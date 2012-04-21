#include <proj_api.h>

typedef struct {
    int err;
    double x, y, z;
} triple;

triple *transform2(projPJ srcdefn, projPJ dstdefn, long point_count, int point_offset, double x, double y);
triple *transform3(projPJ srcdefn, projPJ dstdefn, long point_count, int point_offset, double x, double y, double z);
double triple_x(triple *t);
double triple_y(triple *t);
double triple_z(triple *t);
char *triple_err(triple *t);
char *get_err();
