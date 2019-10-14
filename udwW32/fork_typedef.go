// +build windows

package udwW32

type SECTION_IMAGE_INFORMATION struct {
	TransferAddress             uintptr
	ZeroBits                    uint32
	MaximumStackSize            uint64
	CommittedStackSize          uint64
	SubSystemType               uint32
	SubSystemMinorVersion       uint16
	SubSystemMajorVersion       uint16
	SubSystemVersion            uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	OperatingSystemVersion      uint32
	ImageCharacteristics        uint16
	DllCharacteristics          uint16
	Machine                     uint16
	ImageContainsCode           uint8
	ImageFlags                  uint8
	ComPlusFlags                uint8
	LoaderFlags                 uint32
	ImageFileSize               uint32
	CheckSum                    uint32
}

func (si *SECTION_IMAGE_INFORMATION) ComPlusNativeReady() bool {
	return (si.ComPlusFlags & (1 << 0)) == 1
}

func (si *SECTION_IMAGE_INFORMATION) ComPlusILOnly() bool {
	return (si.ComPlusFlags & (1 << 1)) == 1
}

func (si *SECTION_IMAGE_INFORMATION) ImageDynamicallyRelocated() bool {
	return (si.ComPlusFlags & (1 << 2)) == 1
}

func (si *SECTION_IMAGE_INFORMATION) ImageMappedFlat() bool {
	return (si.ComPlusFlags & (1 << 3)) == 1
}

func (si *SECTION_IMAGE_INFORMATION) BaseBelow4gb() bool {
	return (si.ComPlusFlags & (1 << 4)) == 1
}

func (si *SECTION_IMAGE_INFORMATION) ComPlusPrefer32bit() bool {
	return (si.ComPlusFlags & (1 << 5)) == 1
}

type RTL_USER_PROCESS_INFORMATION struct {
	Length           uint32
	Process          HANDLE
	Thread           HANDLE
	ClientId         CLIENT_ID
	ImageInformation SECTION_IMAGE_INFORMATION
}
