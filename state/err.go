package state

import "errors"

var (
	ErrOidNotFound    = errors.New("oid not found") //
	ErrInvalidUid     = errors.New("invalid uid")   // 无效的客服ID或者访客ID 一般是通过ID长度检测出了不正常
	ErrUpdateNoAffect = errors.New("update no affect")
)
