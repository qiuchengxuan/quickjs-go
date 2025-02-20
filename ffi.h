#include <stdlib.h>
#include <string.h>
#include "libquickjs/quickjs.h"
#include "libquickjs/quickjs-libc.h"

static inline JS_BOOL JS_IsInt(JSValueConst val) { return JS_VALUE_GET_TAG(val) == JS_TAG_INT; }

static inline JSValue JS_Null() { return JS_NULL; }
static inline JSValue JS_Undefined() { return JS_UNDEFINED; }

static inline void* JS_ValuePtr(JSValueConst val) { return JS_VALUE_GET_PTR(val); }
static inline int JS_ValueTag(JSValueConst val) { return JS_VALUE_GET_TAG(val); }

extern JSValue ThrowInternalError(JSContext *ctx, const char *fmt);

extern JSClassDef go_object_class;

extern JSValue proxyCall(JSContext *ctx, JSValueConst fn, JSValueConst this, int argc, JSValueConst *argv, int flags);
