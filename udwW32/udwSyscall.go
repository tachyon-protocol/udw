// +build windows

package udwW32

import (
	"github.com/tachyon-protocol/udw/udwStrconv"
	"strconv"
	"sync"
	"syscall"
)

type Dll struct {
	Name       string
	locker     sync.Mutex
	syscallDll *syscall.DLL
}

type DllProc struct {
	Dll    *Dll
	Name   string
	locker sync.Mutex
	addr   uintptr
}

func (dll *Dll) mustGetSyscallDll() *syscall.DLL {
	dll.locker.Lock()
	defer dll.locker.Unlock()
	if dll.syscallDll != nil {
		return dll.syscallDll
	}
	dll.syscallDll = syscall.MustLoadDLL(dll.Name)
	return dll.syscallDll
}
func MustLoadDLL(name string) *Dll {
	dll := &Dll{Name: name}
	dll.syscallDll = syscall.MustLoadDLL(dll.Name)
	return dll
}

func (dll *Dll) NewProc(name string) *DllProc {

	return &DllProc{
		Dll:  dll,
		Name: name,
	}
}
func (dll *Dll) MustFindProc(name string) *DllProc {

	proc := &DllProc{
		Dll:  dll,
		Name: name,
	}
	syscallDll := proc.Dll.mustGetSyscallDll()
	proc.addr = syscallDll.MustFindProc(proc.Name).Addr()
	return proc
}
func (dll *Dll) MustLoadDll(name string) {
	dll.locker.Lock()
	defer dll.locker.Unlock()
	dll.Name = name
	dll.syscallDll = syscall.MustLoadDLL(dll.Name)
}
func (dll *Dll) Release() (err error) {
	dll.locker.Lock()
	defer dll.locker.Unlock()
	if dll.syscallDll != nil {
		err = dll.syscallDll.Release()
		if err != nil {
			return err
		}
		dll.syscallDll = nil
		dll.Name = ""
		return nil
	}
	return nil
}
func (dll *Dll) IsRelease() bool {
	dll.locker.Lock()
	b := dll.Name == ""
	dll.locker.Unlock()
	return b
}

func (proc *DllProc) Release() {
	proc.locker.Lock()
	proc.addr = 0
	proc.Dll = nil
	proc.Name = ""
	proc.locker.Unlock()
}
func (proc *DllProc) IsRelease() bool {
	proc.locker.Lock()
	b := proc.Dll == nil && proc.Name == "" && proc.addr == 0
	proc.locker.Unlock()
	return b
}

func (proc *DllProc) MustFindProc(dll *Dll, name string) {
	syscallDll := dll.mustGetSyscallDll()
	addr := syscallDll.MustFindProc(name).Addr()
	proc.locker.Lock()
	proc.addr = addr
	proc.Dll = dll
	proc.Name = name
	proc.locker.Unlock()
}

func (proc *DllProc) Addr() uintptr {
	proc.locker.Lock()
	defer proc.locker.Unlock()
	if proc.addr != 0 {
		return proc.addr
	}
	syscallDll := proc.Dll.mustGetSyscallDll()
	proc.addr = syscallDll.MustFindProc(proc.Name).Addr()
	return proc.addr
}

func (proc *DllProc) Call(a ...uintptr) (r1, r2 uintptr, err error) {
	r1, r2, errorNo := proc.callErrorNo(a...)
	if errorNo != 0 {
		return r1, r2, errorNo

	}
	return r1, r2, nil
}
func (proc *DllProc) callErrorNo(a ...uintptr) (r1, r2 uintptr, errorNo syscall.Errno) {
	addr := proc.Addr()
	switch len(a) {
	case 0:
		return syscall.Syscall(addr, uintptr(len(a)), 0, 0, 0)
	case 1:
		return syscall.Syscall(addr, uintptr(len(a)), a[0], 0, 0)
	case 2:
		return syscall.Syscall(addr, uintptr(len(a)), a[0], a[1], 0)
	case 3:
		return syscall.Syscall(addr, uintptr(len(a)), a[0], a[1], a[2])
	case 4:
		return syscall.Syscall6(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], 0, 0)
	case 5:
		return syscall.Syscall6(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], 0)
	case 6:
		return syscall.Syscall6(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5])
	case 7:
		return syscall.Syscall9(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], 0, 0)
	case 8:
		return syscall.Syscall9(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], 0)
	case 9:
		return syscall.Syscall9(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8])
	case 10:
		return syscall.Syscall12(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], 0, 0)
	case 11:
		return syscall.Syscall12(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], 0)
	case 12:
		return syscall.Syscall12(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11])
	case 13:
		return syscall.Syscall15(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], 0, 0)
	case 14:
		return syscall.Syscall15(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], 0)
	case 15:
		return syscall.Syscall15(addr, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], a[14])
	default:
		panic("Call " + proc.Name + " with too many arguments " + strconv.Itoa(len(a)) + ".")
	}

}

func (proc *DllProc) CallErrorMsg(a ...uintptr) (r1, r2 uintptr, errMsg string) {
	r1, r2, errorNo := proc.callErrorNo(a...)
	if errorNo != 0 {
		return r1, r2, proc.Name + " fail W32: " + udwStrconv.FormatUint64Hex(uint64(uintptr(errorNo))) + " " + errorNo.Error()
	}
	return r1, r2, ""
}
