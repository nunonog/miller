package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// This is for Lvalues, i.e. things on the left-hand-side of an assignment
// statement.

// ----------------------------------------------------------------
func BuildAssignableNode(
	astNode *dsl.ASTNode,
) (IAssignable, error) {

	switch astNode.Type {

	case dsl.NodeTypeDirectFieldValue:
		return BuildDirectFieldValueLvalueNode(astNode)
		break
	case dsl.NodeTypeIndirectFieldValue:
		return BuildIndirectFieldValueLvalueNode(astNode)
		break
	case dsl.NodeTypeFullSrec:
		return BuildFullSrecLvalueNode(astNode)
		break

	case dsl.NodeTypeDirectOosvarValue:
		return BuildDirectOosvarValueLvalueNode(astNode)
		break
	case dsl.NodeTypeIndirectOosvarValue:
		return BuildIndirectOosvarValueLvalueNode(astNode)
		break
	case dsl.NodeTypeFullOosvar:
		return BuildFullOosvarLvalueNode(astNode)
		break
	case dsl.NodeTypeLocalVariable:
		return BuildLocalVariableLvalueNode(astNode)
		break

	case dsl.NodeTypeArrayOrMapIndexAccess:
		return BuildIndexedLvalueNode(astNode)
		break
	}

	return nil, errors.New(
		"CST BuildAssignableNode: unhandled AST node " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
type DirectFieldValueLvalueNode struct {
	lhsFieldName *types.Mlrval
}

func BuildDirectFieldValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDirectFieldValue)

	lhsFieldName := types.MlrvalFromString(string(astNode.Token.Lit))
	return NewDirectFieldValueLvalueNode(&lhsFieldName), nil
}

func NewDirectFieldValueLvalueNode(lhsFieldName *types.Mlrval) *DirectFieldValueLvalueNode {
	return &DirectFieldValueLvalueNode{
		lhsFieldName: lhsFieldName,
	}
}

func (this *DirectFieldValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *DirectFieldValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		err := state.Inrec.PutCopyWithMlrvalIndex(this.lhsFieldName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Inrec.PutIndexed(
			append([]*types.Mlrval{this.lhsFieldName}, indices...),
			rvalue,
		)
	}
}

// ----------------------------------------------------------------
type IndirectFieldValueLvalueNode struct {
	lhsFieldNameExpression IEvaluable
}

func BuildIndirectFieldValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectFieldValue)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	lhsFieldNameExpression, err := BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return NewIndirectFieldValueLvalueNode(lhsFieldNameExpression), nil
}

func NewIndirectFieldValueLvalueNode(
	lhsFieldNameExpression IEvaluable,
) *IndirectFieldValueLvalueNode {
	return &IndirectFieldValueLvalueNode{
		lhsFieldNameExpression: lhsFieldNameExpression,
	}
}

func (this *IndirectFieldValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *IndirectFieldValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	lhsFieldName := this.lhsFieldNameExpression.Evaluate(state)

	if indices == nil {
		err := state.Inrec.PutCopyWithMlrvalIndex(&lhsFieldName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Inrec.PutIndexed(
			append([]*types.Mlrval{&lhsFieldName}, indices...),
			rvalue,
		)
	}
}

// ----------------------------------------------------------------
type FullSrecLvalueNode struct {
}

func BuildFullSrecLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFullSrec)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(astNode.Children != nil)
	return NewFullSrecLvalueNode(), nil
}

func NewFullSrecLvalueNode() *FullSrecLvalueNode {
	return &FullSrecLvalueNode{}
}

func (this *FullSrecLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *FullSrecLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// The input record is a *Mlrmap so just invoke its PutIndexed.
	err := state.Inrec.PutIndexed(indices, rvalue)
	if err != nil {
		return err
	}
	return nil
}

// ----------------------------------------------------------------
type DirectOosvarValueLvalueNode struct {
	lhsOosvarName *types.Mlrval
}

func BuildDirectOosvarValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDirectOosvarValue)

	lhsOosvarName := types.MlrvalFromString(string(astNode.Token.Lit))
	return NewDirectOosvarValueLvalueNode(&lhsOosvarName), nil
}

func NewDirectOosvarValueLvalueNode(lhsOosvarName *types.Mlrval) *DirectOosvarValueLvalueNode {
	return &DirectOosvarValueLvalueNode{
		lhsOosvarName: lhsOosvarName,
	}
}

func (this *DirectOosvarValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *DirectOosvarValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		err := state.Oosvars.PutCopyWithMlrvalIndex(this.lhsOosvarName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Oosvars.PutIndexed(
			append([]*types.Mlrval{this.lhsOosvarName}, indices...),
			rvalue,
		)
	}
}

// ----------------------------------------------------------------
type IndirectOosvarValueLvalueNode struct {
	lhsOosvarNameExpression IEvaluable
}

func BuildIndirectOosvarValueLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIndirectOosvarValue)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	lhsOosvarNameExpression, err := BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return NewIndirectOosvarValueLvalueNode(lhsOosvarNameExpression), nil
}

func NewIndirectOosvarValueLvalueNode(
	lhsOosvarNameExpression IEvaluable,
) *IndirectOosvarValueLvalueNode {
	return &IndirectOosvarValueLvalueNode{
		lhsOosvarNameExpression: lhsOosvarNameExpression,
	}
}

func (this *IndirectOosvarValueLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *IndirectOosvarValueLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	lhsOosvarName := this.lhsOosvarNameExpression.Evaluate(state)

	if indices == nil {
		err := state.Oosvars.PutCopyWithMlrvalIndex(&lhsOosvarName, rvalue)
		if err != nil {
			return err
		}
		return nil
	} else {
		return state.Oosvars.PutIndexed(
			append([]*types.Mlrval{&lhsOosvarName}, indices...),
			rvalue,
		)
	}
}

// ----------------------------------------------------------------
type FullOosvarLvalueNode struct {
}

func BuildFullOosvarLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFullOosvar)
	lib.InternalCodingErrorIf(astNode == nil)
	lib.InternalCodingErrorIf(astNode.Children != nil)
	return NewFullOosvarLvalueNode(), nil
}

func NewFullOosvarLvalueNode() *FullOosvarLvalueNode {
	return &FullOosvarLvalueNode{}
}

func (this *FullOosvarLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *FullOosvarLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absentness of the rvalue, so we just assign
	// whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())

	// The input record is a *Mlrmap so just invoke its PutIndexed.
	err := state.Oosvars.PutIndexed(indices, rvalue)
	if err != nil {
		return err
	}
	return nil
}

// ----------------------------------------------------------------
type LocalVariableLvalueNode struct {
	variableName string
}

func BuildLocalVariableLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeLocalVariable)

	variableName := string(astNode.Token.Lit)
	return NewLocalVariableLvalueNode(variableName), nil
}

func NewLocalVariableLvalueNode(variableName string) *LocalVariableLvalueNode {
	return &LocalVariableLvalueNode{
		variableName: variableName,
	}
}

func (this *LocalVariableLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	return this.AssignIndexed(rvalue, nil, state)
}

func (this *LocalVariableLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// AssignmentNode checks for absent, so we just assign whatever we get
	lib.InternalCodingErrorIf(rvalue.IsAbsent())
	if indices == nil {
		state.stack.SetVariable(this.variableName, rvalue)
		return nil
	} else {
		// TODO
		return errors.New("Indexed local-variable assignment has not been implemented yet.")
	}
}

// ----------------------------------------------------------------
// IndexedValueNode is a delegator to base-lvalue types.
// * The baseLvalue is some IAssignable
// * The indexEvaluables are an array of IEvaluables
// * Each needs to evaluate to int or string
// * Assignment needs to walk each level:
//   o error if ith mlrval is int and that level isn't an array
//   o error if ith mlrval is string and that level isn't a map
//   o error for any other types -- maybe absent-handling for heterogeneity ...

// ----------------------------------------------------------------
type IndexedLvalueNode struct {
	baseLvalue      IAssignable
	indexEvaluables []IEvaluable
}

func BuildIndexedLvalueNode(astNode *dsl.ASTNode) (IAssignable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapIndexAccess)
	lib.InternalCodingErrorIf(astNode == nil)

	var baseLvalue IAssignable = nil
	indexEvaluables := make([]IEvaluable, 0)
	var err error = nil

	// $ mlr -n put -v '$x[1][2]=3'
	// DSL EXPRESSION:
	// $x[1][2]=3
	// RAW AST:
	// * StatementBlock
	//     * Assignment "="
	//         * ArrayOrMapIndexAccess "[]"
	//             * ArrayOrMapIndexAccess "[]"
	//                 * DirectFieldValue "x"
	//                 * IntLiteral "1"
	//             * IntLiteral "2"
	//         * IntLiteral "3"

	// In the AST, the indices come in last-shallowest, down to first-deepest,
	// then the base Lvalue.
	walkerNode := astNode
	for {
		if walkerNode.Type == dsl.NodeTypeArrayOrMapIndexAccess {
			lib.InternalCodingErrorIf(walkerNode == nil)
			lib.InternalCodingErrorIf(len(walkerNode.Children) != 2)
			indexEvaluable, err := BuildEvaluableNode(walkerNode.Children[1])
			if err != nil {
				return nil, err
			}
			indexEvaluables = append([]IEvaluable{indexEvaluable}, indexEvaluables...)
			walkerNode = walkerNode.Children[0]
		} else {
			baseLvalue, err = BuildAssignableNode(walkerNode)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	return NewIndexedLvalueNode(baseLvalue, indexEvaluables), nil
}

func NewIndexedLvalueNode(
	baseLvalue IAssignable,
	indexEvaluables []IEvaluable,
) *IndexedLvalueNode {
	return &IndexedLvalueNode{
		baseLvalue:      baseLvalue,
		indexEvaluables: indexEvaluables,
	}
}

// Computes Lvalue indices and then delegates to the baseLvalue.  E.g. for
// '$x[1][2] = 3' or '@x[1][2] = 3', the indices are [1,2], and the baseLvalue
// is '$x' or '@x' respectively.
func (this *IndexedLvalueNode) Assign(
	rvalue *types.Mlrval,
	state *State,
) error {
	indices := make([]*types.Mlrval, len(this.indexEvaluables))

	for i, indexEvaluable := range this.indexEvaluables {
		index := indexEvaluable.Evaluate(state)
		indices[i] = &index
	}
	return this.baseLvalue.AssignIndexed(rvalue, indices, state)
}

func (this *IndexedLvalueNode) AssignIndexed(
	rvalue *types.Mlrval,
	indices []*types.Mlrval,
	state *State,
) error {
	// We are the delegator, not the delegatee
	lib.InternalCodingErrorIf(true)
	return nil // not reached
}
