#ifndef MLR_DSL_CST_H
#define MLR_DSL_CST_H

#include "containers/mlr_dsl_ast.h"
#include "rval_evaluators.h"

// ================================================================
// Concrete syntax tree (CST) derived from an abstract syntax tree (AST).
//
// At present (January 2016) the grammar for 'mlr put' is flat: there is no looping or branching (other than the ternary
// operator, and the gate statement).  Statements are of the form:
//
// * Assignment to oosvar (out-of-stream variables, prefixed with @ sigil)
//
// * Assignment to srec (in-stream records, with field names prefixed with $ sigil)
//
// * filter statements: if false, do not pass the current record from the input stream to the output stream
//
// * gate statements: if false, stop executing further statements
//
// * bare-boolean statements: no-ops unless they have side effects: namely, the matches/does-not-match
//   operators =~ and !=~ setting regex captures \1, \2, etc.
//
// * emit statements: these place oosvar key-value pairs into the output stream.  These can be of the form 'emit @a;
//   emit @b' which produce separate records such as a=3 and b=4, or of the form 'emit @a, @b' which produce records
//   such as a=3,b=4.
//
// This means that as of the present, the CSTs are just list of left-hand sides and right-hand sides.  The left-hand
// side is an output field name, with an LHS-type flag (srec or oosvar). For filter/gate/bare-boolean these are unused;
// for assignments and emits, the output field names are used. The right-hand sides are lrec-evaluators which take an
// input record and oosvar-map, and produce a mlrval which can be assigned. Additionally, for all but emit, there is a
// single LHS name and single RHS lrec-evaluator per statement.
//
// Further, these statements are organized into three groups:
//
// * begin: executed once, before the first input record is read.
// * main : executed for each input record.
// * end :  executed once, after the last input record is read.
// ================================================================

typedef enum _mlr_dsl_cst_lhs_type_t {
	MLR_DSL_CST_LHS_TYPE_NONE,
	MLR_DSL_CST_LHS_TYPE_SREC,
	MLR_DSL_CST_LHS_TYPE_OOSVAR
} mlr_dsl_cst_lhs_type_t;

typedef struct _mlr_dsl_cst_statement_item_t {
	// LHS:
	// xxx replace with function-pointer:
	mlr_dsl_cst_lhs_type_t lhs_type;
	char* output_field_name;
	sllv_t* poosvar_lhs_keylist_evaluators;

	// RHS:
	rval_evaluator_t* prhs_evaluator;
	//rval_evaluator_t* prhs_or_cond_evaluator; <-- rename
	//sllv_t*           pcond_statements;
} mlr_dsl_cst_statement_item_t;

typedef struct _mlr_dsl_cst_statement_t {
	mlr_dsl_ast_node_type_t ast_node_type;
	sllv_t* pitems;
} mlr_dsl_cst_statement_t;

typedef struct _mlr_dsl_cst_t {
	sllv_t* pbegin_statements;
	sllv_t* pmain_statements;
	sllv_t* pend_statements;
} mlr_dsl_cst_t;

// ----------------------------------------------------------------
mlr_dsl_cst_t* mlr_dsl_cst_alloc(mlr_dsl_ast_t* past, int type_inferencing);
void mlr_dsl_cst_free(mlr_dsl_cst_t* pcst);

void mlr_dsl_cst_evaluate(
	sllv_t*          pcst_statements, // begin/main/end
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pemit_rec,
	sllv_t*          poutrecs);

int mlr_dsl_cst_node_evaluate(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pemit_rec,
	sllv_t*          poutrecs);

#endif // MLR_DSL_CST_H
