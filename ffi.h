#include <stdlib.h>
#include <string.h>
#include "libquickjs/quickjs.h"
#include "libquickjs/quickjs-libc.h"

static inline JS_BOOL JS_IsInt(JSValueConst val) { return JS_VALUE_GET_TAG(val) == JS_TAG_INT; }

static inline JSValue JS_Null() { return JS_NULL; }

extern JSValue proxyCall(JSContext *ctx, JSValueConst this, int argc, JSValueConst *argv);
