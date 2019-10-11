package sqlite3

import "C"

type ErrNo int

const ErrNoMask C.int = 0xff

type ErrNoExtended int

type Error struct {
	Code         ErrNo
	ExtendedCode ErrNoExtended
	err          string
}

var (
	ErrError      = ErrNo(1)
	ErrInternal   = ErrNo(2)
	ErrPerm       = ErrNo(3)
	ErrAbort      = ErrNo(4)
	ErrBusy       = ErrNo(5)
	ErrLocked     = ErrNo(6)
	ErrNomem      = ErrNo(7)
	ErrReadonly   = ErrNo(8)
	ErrInterrupt  = ErrNo(9)
	ErrIoErr      = ErrNo(10)
	ErrCorrupt    = ErrNo(11)
	ErrNotFound   = ErrNo(12)
	ErrFull       = ErrNo(13)
	ErrCantOpen   = ErrNo(14)
	ErrProtocol   = ErrNo(15)
	ErrEmpty      = ErrNo(16)
	ErrSchema     = ErrNo(17)
	ErrTooBig     = ErrNo(18)
	ErrConstraint = ErrNo(19)
	ErrMismatch   = ErrNo(20)
	ErrMisuse     = ErrNo(21)
	ErrNoLFS      = ErrNo(22)
	ErrAuth       = ErrNo(23)
	ErrFormat     = ErrNo(24)
	ErrRange      = ErrNo(25)
	ErrNotADB     = ErrNo(26)
	ErrNotice     = ErrNo(27)
	ErrWarning    = ErrNo(28)
)

func (err ErrNo) Error() string {
	return Error{Code: err}.Error()
}

func (err ErrNo) Extend(by int) ErrNoExtended {
	return ErrNoExtended(int(err) | (by << 8))
}

func (err ErrNoExtended) Error() string {
	return Error{Code: ErrNo(C.int(err) & ErrNoMask), ExtendedCode: err}.Error()
}

func (err Error) Error() string {
	if err.err != "" {
		return err.err
	}
	return errorString(err)
}
