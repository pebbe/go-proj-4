#include <proj_api.h>

typedef struct {
    char *err;
    double x, y, z;
} triple;

triple *transform(projPJ srcdefn, projPJ dstdefn, long point_count, int point_offset, double x, double y, double z, int has_z);
char *get_err();
