// Code generated by "stringer -type=Opcode"; DO NOT EDIT.

package isa

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpcodeMask-63]
	_ = x[Add-0]
	_ = x[Arg-1]
	_ = x[ArgLen-2]
	_ = x[Bind-3]
	_ = x[BindRef-4]
	_ = x[Call-5]
	_ = x[Call0-6]
	_ = x[Call1-7]
	_ = x[Call2-8]
	_ = x[Call3-9]
	_ = x[CallWith-10]
	_ = x[Car-11]
	_ = x[Cdr-12]
	_ = x[Closure-13]
	_ = x[CondJump-14]
	_ = x[Cons-15]
	_ = x[Const-16]
	_ = x[Deref-17]
	_ = x[Div-18]
	_ = x[Dup-19]
	_ = x[Empty-20]
	_ = x[Eq-21]
	_ = x[False-22]
	_ = x[Jump-23]
	_ = x[Label-24]
	_ = x[Load-25]
	_ = x[Mod-26]
	_ = x[Mul-27]
	_ = x[Neg-28]
	_ = x[NegInt-29]
	_ = x[NewRef-30]
	_ = x[NoOp-31]
	_ = x[Not-32]
	_ = x[Null-33]
	_ = x[NumEq-34]
	_ = x[NumGt-35]
	_ = x[NumGte-36]
	_ = x[NumLt-37]
	_ = x[NumLte-38]
	_ = x[Panic-39]
	_ = x[Pop-40]
	_ = x[PopArgs-41]
	_ = x[PosInt-42]
	_ = x[Private-43]
	_ = x[Public-44]
	_ = x[PushArgs-45]
	_ = x[Resolve-46]
	_ = x[RestArg-47]
	_ = x[RetFalse-48]
	_ = x[RetNull-49]
	_ = x[RetTrue-50]
	_ = x[Return-51]
	_ = x[Store-52]
	_ = x[Sub-53]
	_ = x[TailCall-54]
	_ = x[True-55]
	_ = x[Vector-56]
	_ = x[Zero-57]
}

const (
	_Opcode_name_0 = "AddArgArgLenBindBindRefCallCall0Call1Call2Call3CallWithCarCdrClosureCondJumpConsConstDerefDivDupEmptyEqFalseJumpLabelLoadModMulNegNegIntNewRefNoOpNotNullNumEqNumGtNumGteNumLtNumLtePanicPopPopArgsPosIntPrivatePublicPushArgsResolveRestArgRetFalseRetNullRetTrueReturnStoreSubTailCallTrueVectorZero"
	_Opcode_name_1 = "OpcodeMask"
)

var (
	_Opcode_index_0 = [...]uint16{0, 3, 6, 12, 16, 23, 27, 32, 37, 42, 47, 55, 58, 61, 68, 76, 80, 85, 90, 93, 96, 101, 103, 108, 112, 117, 121, 124, 127, 130, 136, 142, 146, 149, 153, 158, 163, 169, 174, 180, 185, 188, 195, 201, 208, 214, 222, 229, 236, 244, 251, 258, 264, 269, 272, 280, 284, 290, 294}
)

func (i Opcode) String() string {
	switch {
	case i <= 57:
		return _Opcode_name_0[_Opcode_index_0[i]:_Opcode_index_0[i+1]]
	case i == 63:
		return _Opcode_name_1
	default:
		return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
