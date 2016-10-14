#ifndef LOCAL_STACK_H
#define LOCAL_STACK_H

#include "containers/mlrval.h"

// ================================================================
// Bound & scoped variables for use in for-loops, function bodies, and
// subroutine bodies. Indices of local variables, and max-depth for top-level
// statement blocks, are compted by the stack-allocator which marks up the AST
// before the CST is built from it.
//
// A convention shared between the stack-allocator and this data structure is
// that slot 0 is an absent-null which is used for reads of undefined (or
// as-yet-undefined) local variables.
// ================================================================

// ----------------------------------------------------------------
typedef struct _local_stack_t {
	int in_use;
	int size;
	int frame_base;
	mv_t* pvars;
	// xxx has-absent-read flag ...
} local_stack_t;

// ----------------------------------------------------------------
// A stack is allocated for a top-level statement block: begin, end, or main, or
// user-defined function/subroutine. (The latter two may be called recursively
// in which case the in_use flag notes the need to allocate a new stack.)

local_stack_t* local_stack_alloc(int size);
void local_stack_free(local_stack_t* pstack);

// ----------------------------------------------------------------
// Sets/clear the in-use flag for top-level statement blocks, and verifies
// the contract for absent-null at slot 0.

static inline void local_stack_enter(local_stack_t* pstack) {
	pstack->in_use = TRUE;
}
static inline void local_stack_exit (local_stack_t* pstack) {
	pstack->in_use = FALSE;
	if (!mv_is_absent(&pstack->pvars[0])) {
		fprintf(stderr, "%s: Internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
}

// ----------------------------------------------------------------
// Frames are entered/exited for each curly-braced statement block, including
// the top-level block itself as well as ifs/fors/whiles.

static inline void local_stack_frame_enter(local_stack_t* pstack, int count) {
	// xxx try to avoid with absent-read flag at stack-allocator ...
	mv_t* pframe = &pstack->pvars[pstack->frame_base];
	for (int i = 0; i < count; i++) {
		pframe[i] = mv_absent();
	}
	pstack->frame_base += count;
}
static inline void local_stack_frame_exit(local_stack_t* pstack, int count) {
	pstack->frame_base -= count;
}

#endif // LOCAL_STACK_H