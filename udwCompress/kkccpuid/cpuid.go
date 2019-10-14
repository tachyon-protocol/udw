package kkccpuid

import "strings"

type Vendor int

const (
	Other Vendor = iota
	Intel
	AMD
	VIA
	Transmeta
	NSC
	KVM
	MSVM
	VMware
	XenHVM
)

const (
	CMOV = 1 << iota
	NX
	AMD3DNOW
	AMD3DNOWEXT
	MMX
	MMXEXT
	SSE
	SSE2
	SSE3
	SSSE3
	SSE4
	SSE4A
	SSE42
	AVX
	AVX2
	FMA3
	FMA4
	XOP
	F16C
	BMI1
	BMI2
	TBM
	LZCNT
	POPCNT
	AESNI
	CLMUL
	HTT
	HLE
	RTM
	RDRAND
	RDSEED
	ADX
	SHA
	AVX512F
	AVX512DQ
	AVX512IFMA
	AVX512PF
	AVX512ER
	AVX512CD
	AVX512BW
	AVX512VL
	AVX512VBMI
	MPX
	ERMS
	RDTSCP
	CX16
	SGX

	SSE2SLOW
	SSE3SLOW
	ATOM
)

var flagNames = map[Flags]string{
	CMOV:        "CMOV",
	NX:          "NX",
	AMD3DNOW:    "AMD3DNOW",
	AMD3DNOWEXT: "AMD3DNOWEXT",
	MMX:         "MMX",
	MMXEXT:      "MMXEXT",
	SSE:         "SSE",
	SSE2:        "SSE2",
	SSE3:        "SSE3",
	SSSE3:       "SSSE3",
	SSE4:        "SSE4.1",
	SSE4A:       "SSE4A",
	SSE42:       "SSE4.2",
	AVX:         "AVX",
	AVX2:        "AVX2",
	FMA3:        "FMA3",
	FMA4:        "FMA4",
	XOP:         "XOP",
	F16C:        "F16C",
	BMI1:        "BMI1",
	BMI2:        "BMI2",
	TBM:         "TBM",
	LZCNT:       "LZCNT",
	POPCNT:      "POPCNT",
	AESNI:       "AESNI",
	CLMUL:       "CLMUL",
	HTT:         "HTT",
	HLE:         "HLE",
	RTM:         "RTM",
	RDRAND:      "RDRAND",
	RDSEED:      "RDSEED",
	ADX:         "ADX",
	SHA:         "SHA",
	AVX512F:     "AVX512F",
	AVX512DQ:    "AVX512DQ",
	AVX512IFMA:  "AVX512IFMA",
	AVX512PF:    "AVX512PF",
	AVX512ER:    "AVX512ER",
	AVX512CD:    "AVX512CD",
	AVX512BW:    "AVX512BW",
	AVX512VL:    "AVX512VL",
	AVX512VBMI:  "AVX512VBMI",
	MPX:         "MPX",
	ERMS:        "ERMS",
	RDTSCP:      "RDTSCP",
	CX16:        "CX16",
	SGX:         "SGX",

	SSE2SLOW: "SSE2SLOW",
	SSE3SLOW: "SSE3SLOW",
	ATOM:     "ATOM",
}

type CPUInfo struct {
	BrandName      string
	VendorID       Vendor
	Features       Flags
	PhysicalCores  int
	ThreadsPerCore int
	LogicalCores   int
	Family         int
	Model          int
	CacheLine      int
	Cache          struct {
		L1I int
		L1D int
		L2  int
		L3  int
	}
	SGX       SGXSupport
	maxFunc   uint32
	maxExFunc uint32
}

var cpuid func(op uint32) (eax, ebx, ecx, edx uint32)
var cpuidex func(op, op2 uint32) (eax, ebx, ecx, edx uint32)
var xgetbv func(index uint32) (eax, edx uint32)
var rdtscpAsm func() (eax, ebx, ecx, edx uint32)

var CPU CPUInfo

func init() {
	initCPU()
	Detect()
}

func Detect() {
	CPU.maxFunc = maxFunctionID()
	CPU.maxExFunc = maxExtendedFunction()
	CPU.BrandName = brandName()
	CPU.CacheLine = cacheLine()
	CPU.Family, CPU.Model = familyModel()
	CPU.Features = support()
	CPU.SGX = sgx(CPU.Features&SGX != 0)
	CPU.ThreadsPerCore = threadsPerCore()
	CPU.LogicalCores = logicalCores()
	CPU.PhysicalCores = physicalCores()
	CPU.VendorID = vendorID()
	CPU.cacheSize()
}

func (c CPUInfo) Cmov() bool {
	return c.Features&CMOV != 0
}

func (c CPUInfo) Amd3dnow() bool {
	return c.Features&AMD3DNOW != 0
}

func (c CPUInfo) Amd3dnowExt() bool {
	return c.Features&AMD3DNOWEXT != 0
}

func (c CPUInfo) MMX() bool {
	return c.Features&MMX != 0
}

func (c CPUInfo) MMXExt() bool {
	return c.Features&MMXEXT != 0
}

func (c CPUInfo) SSE() bool {
	return c.Features&SSE != 0
}

func (c CPUInfo) SSE2() bool {
	return c.Features&SSE2 != 0
}

func (c CPUInfo) SSE3() bool {
	return c.Features&SSE3 != 0
}

func (c CPUInfo) SSSE3() bool {
	return c.Features&SSSE3 != 0
}

func (c CPUInfo) SSE4() bool {
	return c.Features&SSE4 != 0
}

func (c CPUInfo) SSE42() bool {
	return c.Features&SSE42 != 0
}

func (c CPUInfo) AVX() bool {
	return c.Features&AVX != 0
}

func (c CPUInfo) AVX2() bool {
	return c.Features&AVX2 != 0
}

func (c CPUInfo) FMA3() bool {
	return c.Features&FMA3 != 0
}

func (c CPUInfo) FMA4() bool {
	return c.Features&FMA4 != 0
}

func (c CPUInfo) XOP() bool {
	return c.Features&XOP != 0
}

func (c CPUInfo) F16C() bool {
	return c.Features&F16C != 0
}

func (c CPUInfo) BMI1() bool {
	return c.Features&BMI1 != 0
}

func (c CPUInfo) BMI2() bool {
	return c.Features&BMI2 != 0
}

func (c CPUInfo) TBM() bool {
	return c.Features&TBM != 0
}

func (c CPUInfo) Lzcnt() bool {
	return c.Features&LZCNT != 0
}

func (c CPUInfo) Popcnt() bool {
	return c.Features&POPCNT != 0
}

func (c CPUInfo) HTT() bool {
	return c.Features&HTT != 0
}

func (c CPUInfo) SSE2Slow() bool {
	return c.Features&SSE2SLOW != 0
}

func (c CPUInfo) SSE3Slow() bool {
	return c.Features&SSE3SLOW != 0
}

func (c CPUInfo) AesNi() bool {
	return c.Features&AESNI != 0
}

func (c CPUInfo) Clmul() bool {
	return c.Features&CLMUL != 0
}

func (c CPUInfo) NX() bool {
	return c.Features&NX != 0
}

func (c CPUInfo) SSE4A() bool {
	return c.Features&SSE4A != 0
}

func (c CPUInfo) HLE() bool {
	return c.Features&HLE != 0
}

func (c CPUInfo) RTM() bool {
	return c.Features&RTM != 0
}

func (c CPUInfo) Rdrand() bool {
	return c.Features&RDRAND != 0
}

func (c CPUInfo) Rdseed() bool {
	return c.Features&RDSEED != 0
}

func (c CPUInfo) ADX() bool {
	return c.Features&ADX != 0
}

func (c CPUInfo) SHA() bool {
	return c.Features&SHA != 0
}

func (c CPUInfo) AVX512F() bool {
	return c.Features&AVX512F != 0
}

func (c CPUInfo) AVX512DQ() bool {
	return c.Features&AVX512DQ != 0
}

func (c CPUInfo) AVX512IFMA() bool {
	return c.Features&AVX512IFMA != 0
}

func (c CPUInfo) AVX512PF() bool {
	return c.Features&AVX512PF != 0
}

func (c CPUInfo) AVX512ER() bool {
	return c.Features&AVX512ER != 0
}

func (c CPUInfo) AVX512CD() bool {
	return c.Features&AVX512CD != 0
}

func (c CPUInfo) AVX512BW() bool {
	return c.Features&AVX512BW != 0
}

func (c CPUInfo) AVX512VL() bool {
	return c.Features&AVX512VL != 0
}

func (c CPUInfo) AVX512VBMI() bool {
	return c.Features&AVX512VBMI != 0
}

func (c CPUInfo) MPX() bool {
	return c.Features&MPX != 0
}

func (c CPUInfo) ERMS() bool {
	return c.Features&ERMS != 0
}

func (c CPUInfo) RDTSCP() bool {
	return c.Features&RDTSCP != 0
}

func (c CPUInfo) CX16() bool {
	return c.Features&CX16 != 0
}

func (c CPUInfo) Atom() bool {
	return c.Features&ATOM != 0
}

func (c CPUInfo) Intel() bool {
	return c.VendorID == Intel
}

func (c CPUInfo) AMD() bool {
	return c.VendorID == AMD
}

func (c CPUInfo) Transmeta() bool {
	return c.VendorID == Transmeta
}

func (c CPUInfo) NSC() bool {
	return c.VendorID == NSC
}

func (c CPUInfo) VIA() bool {
	return c.VendorID == VIA
}

func (c CPUInfo) RTCounter() uint64 {
	if !c.RDTSCP() {
		return 0
	}
	a, _, _, d := rdtscpAsm()
	return uint64(a) | (uint64(d) << 32)
}

func (c CPUInfo) Ia32TscAux() uint32 {
	if !c.RDTSCP() {
		return 0
	}
	_, _, ecx, _ := rdtscpAsm()
	return ecx
}

func (c CPUInfo) LogicalCPU() int {
	if c.maxFunc < 1 {
		return -1
	}
	_, ebx, _, _ := cpuid(1)
	return int(ebx >> 24)
}

func (c CPUInfo) VM() bool {
	switch c.VendorID {
	case MSVM, KVM, VMware, XenHVM:
		return true
	}
	return false
}

type Flags uint64

func (f Flags) String() string {
	return strings.Join(f.Strings(), ",")
}

func (f Flags) Strings() []string {
	s := support()
	r := make([]string, 0, 20)
	for i := uint(0); i < 64; i++ {
		key := Flags(1 << i)
		val := flagNames[key]
		if s&key != 0 {
			r = append(r, val)
		}
	}
	return r
}

func maxExtendedFunction() uint32 {
	eax, _, _, _ := cpuid(0x80000000)
	return eax
}

func maxFunctionID() uint32 {
	a, _, _, _ := cpuid(0)
	return a
}

func brandName() string {
	if maxExtendedFunction() >= 0x80000004 {
		v := make([]uint32, 0, 48)
		for i := uint32(0); i < 3; i++ {
			a, b, c, d := cpuid(0x80000002 + i)
			v = append(v, a, b, c, d)
		}
		return strings.Trim(string(valAsString(v...)), " ")
	}
	return "unknown"
}

func threadsPerCore() int {
	mfi := maxFunctionID()
	if mfi < 0x4 || vendorID() != Intel {
		return 1
	}

	if mfi < 0xb {
		_, b, _, d := cpuid(1)
		if (d & (1 << 28)) != 0 {

			v := (b >> 16) & 255
			if v > 1 {
				a4, _, _, _ := cpuid(4)

				v2 := (a4 >> 26) + 1
				if v2 > 0 {
					return int(v) / int(v2)
				}
			}
		}
		return 1
	}
	_, b, _, _ := cpuidex(0xb, 0)
	if b&0xffff == 0 {
		return 1
	}
	return int(b & 0xffff)
}

func logicalCores() int {
	mfi := maxFunctionID()
	switch vendorID() {
	case Intel:

		if mfi < 0xb {
			if mfi < 1 {
				return 0
			}

			_, ebx, _, _ := cpuid(1)
			logical := (ebx >> 16) & 0xff
			return int(logical)
		}
		_, b, _, _ := cpuidex(0xb, 1)
		return int(b & 0xffff)
	case AMD:
		_, b, _, _ := cpuid(1)
		return int((b >> 16) & 0xff)
	default:
		return 0
	}
}

func familyModel() (int, int) {
	if maxFunctionID() < 0x1 {
		return 0, 0
	}
	eax, _, _, _ := cpuid(1)
	family := ((eax >> 8) & 0xf) + ((eax >> 20) & 0xff)
	model := ((eax >> 4) & 0xf) + ((eax >> 12) & 0xf0)
	return int(family), int(model)
}

func physicalCores() int {
	switch vendorID() {
	case Intel:
		return logicalCores() / threadsPerCore()
	case AMD:
		if maxExtendedFunction() >= 0x80000008 {
			_, _, c, _ := cpuid(0x80000008)
			return int(c&0xff) + 1
		}
	}
	return 0
}

var vendorMapping = map[string]Vendor{
	"AMDisbetter!": AMD,
	"AuthenticAMD": AMD,
	"CentaurHauls": VIA,
	"GenuineIntel": Intel,
	"TransmetaCPU": Transmeta,
	"GenuineTMx86": Transmeta,
	"Geode by NSC": NSC,
	"VIA VIA VIA ": VIA,
	"KVMKVMKVMKVM": KVM,
	"Microsoft Hv": MSVM,
	"VMwareVMware": VMware,
	"XenVMMXenVMM": XenHVM,
}

func vendorID() Vendor {
	_, b, c, d := cpuid(0)
	v := valAsString(b, d, c)
	vend, ok := vendorMapping[string(v)]
	if !ok {
		return Other
	}
	return vend
}

func cacheLine() int {
	if maxFunctionID() < 0x1 {
		return 0
	}

	_, ebx, _, _ := cpuid(1)
	cache := (ebx & 0xff00) >> 5
	if cache == 0 && maxExtendedFunction() >= 0x80000006 {
		_, _, ecx, _ := cpuid(0x80000006)
		cache = ecx & 0xff
	}

	return int(cache)
}

func (c *CPUInfo) cacheSize() {
	c.Cache.L1D = -1
	c.Cache.L1I = -1
	c.Cache.L2 = -1
	c.Cache.L3 = -1
	vendor := vendorID()
	switch vendor {
	case Intel:
		if maxFunctionID() < 4 {
			return
		}
		for i := uint32(0); ; i++ {
			eax, ebx, ecx, _ := cpuidex(4, i)
			cacheType := eax & 15
			if cacheType == 0 {
				break
			}
			cacheLevel := (eax >> 5) & 7
			coherency := int(ebx&0xfff) + 1
			partitions := int((ebx>>12)&0x3ff) + 1
			associativity := int((ebx>>22)&0x3ff) + 1
			sets := int(ecx) + 1
			size := associativity * partitions * coherency * sets
			switch cacheLevel {
			case 1:
				if cacheType == 1 {

					c.Cache.L1D = size
				} else if cacheType == 2 {

					c.Cache.L1I = size
				} else {
					if c.Cache.L1D < 0 {
						c.Cache.L1I = size
					}
					if c.Cache.L1I < 0 {
						c.Cache.L1I = size
					}
				}
			case 2:
				c.Cache.L2 = size
			case 3:
				c.Cache.L3 = size
			}
		}
	case AMD:

		if maxExtendedFunction() < 0x80000005 {
			return
		}
		_, _, ecx, edx := cpuid(0x80000005)
		c.Cache.L1D = int(((ecx >> 24) & 0xFF) * 1024)
		c.Cache.L1I = int(((edx >> 24) & 0xFF) * 1024)

		if maxExtendedFunction() < 0x80000006 {
			return
		}
		_, _, ecx, _ = cpuid(0x80000006)
		c.Cache.L2 = int(((ecx >> 16) & 0xFFFF) * 1024)
	}

	return
}

type SGXSupport struct {
	Available           bool
	SGX1Supported       bool
	SGX2Supported       bool
	MaxEnclaveSizeNot64 int64
	MaxEnclaveSize64    int64
}

func sgx(available bool) (rval SGXSupport) {
	rval.Available = available

	if !available {
		return
	}

	a, _, _, d := cpuidex(0x12, 0)
	rval.SGX1Supported = a&0x01 != 0
	rval.SGX2Supported = a&0x02 != 0
	rval.MaxEnclaveSizeNot64 = 1 << (d & 0xFF)
	rval.MaxEnclaveSize64 = 1 << ((d >> 8) & 0xFF)

	return
}

func support() Flags {
	mfi := maxFunctionID()
	vend := vendorID()
	if mfi < 0x1 {
		return 0
	}
	rval := uint64(0)
	_, _, c, d := cpuid(1)
	if (d & (1 << 15)) != 0 {
		rval |= CMOV
	}
	if (d & (1 << 23)) != 0 {
		rval |= MMX
	}
	if (d & (1 << 25)) != 0 {
		rval |= MMXEXT
	}
	if (d & (1 << 25)) != 0 {
		rval |= SSE
	}
	if (d & (1 << 26)) != 0 {
		rval |= SSE2
	}
	if (c & 1) != 0 {
		rval |= SSE3
	}
	if (c & 0x00000200) != 0 {
		rval |= SSSE3
	}
	if (c & 0x00080000) != 0 {
		rval |= SSE4
	}
	if (c & 0x00100000) != 0 {
		rval |= SSE42
	}
	if (c & (1 << 25)) != 0 {
		rval |= AESNI
	}
	if (c & (1 << 1)) != 0 {
		rval |= CLMUL
	}
	if c&(1<<23) != 0 {
		rval |= POPCNT
	}
	if c&(1<<30) != 0 {
		rval |= RDRAND
	}
	if c&(1<<29) != 0 {
		rval |= F16C
	}
	if c&(1<<13) != 0 {
		rval |= CX16
	}
	if vend == Intel && (d&(1<<28)) != 0 && mfi >= 4 {
		if threadsPerCore() > 1 {
			rval |= HTT
		}
	}

	if c&(1<<26) != 0 && c&(1<<27) != 0 && c&(1<<28) != 0 {

		eax, _ := xgetbv(0)
		if (eax & 0x6) == 0x6 {
			rval |= AVX
			if (c & 0x00001000) != 0 {
				rval |= FMA3
			}
		}
	}

	if mfi >= 7 {
		_, ebx, ecx, _ := cpuidex(7, 0)
		if (rval&AVX) != 0 && (ebx&0x00000020) != 0 {
			rval |= AVX2
		}
		if (ebx & 0x00000008) != 0 {
			rval |= BMI1
			if (ebx & 0x00000100) != 0 {
				rval |= BMI2
			}
		}
		if ebx&(1<<2) != 0 {
			rval |= SGX
		}
		if ebx&(1<<4) != 0 {
			rval |= HLE
		}
		if ebx&(1<<9) != 0 {
			rval |= ERMS
		}
		if ebx&(1<<11) != 0 {
			rval |= RTM
		}
		if ebx&(1<<14) != 0 {
			rval |= MPX
		}
		if ebx&(1<<18) != 0 {
			rval |= RDSEED
		}
		if ebx&(1<<19) != 0 {
			rval |= ADX
		}
		if ebx&(1<<29) != 0 {
			rval |= SHA
		}

		if c&((1<<26)|(1<<27)) == (1<<26)|(1<<27) {

			eax, _ := xgetbv(0)

			if (eax>>5)&7 == 7 && (eax>>1)&3 == 3 {
				if ebx&(1<<16) != 0 {
					rval |= AVX512F
				}
				if ebx&(1<<17) != 0 {
					rval |= AVX512DQ
				}
				if ebx&(1<<21) != 0 {
					rval |= AVX512IFMA
				}
				if ebx&(1<<26) != 0 {
					rval |= AVX512PF
				}
				if ebx&(1<<27) != 0 {
					rval |= AVX512ER
				}
				if ebx&(1<<28) != 0 {
					rval |= AVX512CD
				}
				if ebx&(1<<30) != 0 {
					rval |= AVX512BW
				}
				if ebx&(1<<31) != 0 {
					rval |= AVX512VL
				}

				if ecx&(1<<1) != 0 {
					rval |= AVX512VBMI
				}
			}
		}
	}

	if maxExtendedFunction() >= 0x80000001 {
		_, _, c, d := cpuid(0x80000001)
		if (c & (1 << 5)) != 0 {
			rval |= LZCNT
			rval |= POPCNT
		}
		if (d & (1 << 31)) != 0 {
			rval |= AMD3DNOW
		}
		if (d & (1 << 30)) != 0 {
			rval |= AMD3DNOWEXT
		}
		if (d & (1 << 23)) != 0 {
			rval |= MMX
		}
		if (d & (1 << 22)) != 0 {
			rval |= MMXEXT
		}
		if (c & (1 << 6)) != 0 {
			rval |= SSE4A
		}
		if d&(1<<20) != 0 {
			rval |= NX
		}
		if d&(1<<27) != 0 {
			rval |= RDTSCP
		}

		if vendorID() != Intel &&
			rval&SSE2 != 0 && (c&0x00000040) == 0 {
			rval |= SSE2SLOW
		}

		if (rval & AVX) != 0 {
			if (c & 0x00000800) != 0 {
				rval |= XOP
			}
			if (c & 0x00010000) != 0 {
				rval |= FMA4
			}
		}

		if vendorID() == Intel {
			family, model := familyModel()
			if family == 6 && (model == 9 || model == 13 || model == 14) {

				if (rval & SSE2) != 0 {
					rval |= SSE2SLOW
				}
				if (rval & SSE3) != 0 {
					rval |= SSE3SLOW
				}
			}

			if family == 6 && model == 28 {
				rval |= ATOM
			}
		}
	}
	return Flags(rval)
}

func valAsString(values ...uint32) []byte {
	r := make([]byte, 4*len(values))
	for i, v := range values {
		dst := r[i*4:]
		dst[0] = byte(v & 0xff)
		dst[1] = byte((v >> 8) & 0xff)
		dst[2] = byte((v >> 16) & 0xff)
		dst[3] = byte((v >> 24) & 0xff)
		switch {
		case dst[0] == 0:
			return r[:i*4]
		case dst[1] == 0:
			return r[:i*4+1]
		case dst[2] == 0:
			return r[:i*4+2]
		case dst[3] == 0:
			return r[:i*4+3]
		}
	}
	return r
}
