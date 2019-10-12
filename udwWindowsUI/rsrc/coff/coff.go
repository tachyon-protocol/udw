package coff

import (
	"debug/pe"
	"encoding/binary"
	"errors"
	"github.com/tachyon-protocol/udw/udwWindowsUI/rsrc/binutil"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Dir struct {
	Characteristics      uint32
	TimeDateStamp        uint32
	MajorVersion         uint16
	MinorVersion         uint16
	NumberOfNamedEntries uint16
	NumberOfIdEntries    uint16
	DirEntries
	Dirs
}

type DirEntries []DirEntry
type Dirs []Dir

type DirEntry struct {
	NameOrId     uint32
	OffsetToData uint32
}

type DataEntry struct {
	OffsetToData uint32
	Size1        uint32
	CodePage     uint32
	Reserved     uint32
}

type RelocationEntry struct {
	RVA         uint32
	SymbolIndex uint32
	Type        uint16
}

const (
	_IMAGE_REL_AMD64_ADDR32NB = 0x03
	_IMAGE_REL_I386_DIR32NB   = 0x07
)

type Auxiliary [18]byte

type Symbol struct {
	Name           [8]byte
	Value          uint32
	SectionNumber  uint16
	Type           uint16
	StorageClass   uint8
	AuxiliaryCount uint8
	Auxiliaries    []Auxiliary
}

type StringsHeader struct {
	Length uint32
}

const (
	MASK_SUBDIRECTORY = 1 << 31

	RT_ICON       = 3
	RT_GROUP_ICON = 3 + 11
	RT_MANIFEST   = 24
)

const (
	DT_PTR  = 1
	T_UCHAR = 12
)

var (
	STRING_RSRC  = [8]byte{'.', 'r', 's', 'r', 'c', 0, 0, 0}
	STRING_RDATA = [8]byte{'.', 'r', 'd', 'a', 't', 'a', 0, 0}

	LANG_ENTRY = DirEntry{NameOrId: 0x0409}
)

type Sizer interface {
	Size() int64
}

type Coff struct {
	pe.FileHeader
	pe.SectionHeader32

	*Dir
	DataEntries []DataEntry
	Data        []Sizer

	Relocations []RelocationEntry
	Symbols     []Symbol
	StringsHeader
	Strings []Sizer
}

func NewRDATA() *Coff {
	return &Coff{
		pe.FileHeader{
			Machine:              pe.IMAGE_FILE_MACHINE_I386,
			NumberOfSections:     1,
			TimeDateStamp:        0,
			NumberOfSymbols:      2,
			SizeOfOptionalHeader: 0,
			Characteristics:      0x0105,
		},
		pe.SectionHeader32{
			Name:            STRING_RDATA,
			Characteristics: 0x40000040,
		},

		nil,
		[]DataEntry{},

		[]Sizer{},

		[]RelocationEntry{},

		[]Symbol{Symbol{
			Name:           STRING_RDATA,
			Value:          0,
			SectionNumber:  1,
			Type:           0,
			StorageClass:   3,
			AuxiliaryCount: 1,
			Auxiliaries:    []Auxiliary{{}},
		}},

		StringsHeader{
			Length: uint32(binary.Size(StringsHeader{})),
		},
		[]Sizer{},
	}
}

func (coff *Coff) Arch(arch string) error {
	switch arch {
	case "386":
		coff.Machine = pe.IMAGE_FILE_MACHINE_I386
	case "amd64":

		coff.Machine = pe.IMAGE_FILE_MACHINE_AMD64
	default:
		return errors.New("coff: unknown architecture: " + arch)
	}
	return nil
}

func (coff *Coff) AddData(symbol string, data Sizer) {
	coff.addSymbol(symbol)
	coff.Data = append(coff.Data, data)
	coff.SectionHeader32.SizeOfRawData += uint32(data.Size())
}

func (coff *Coff) addSymbol(s string) {
	coff.FileHeader.NumberOfSymbols++

	buf := strings.NewReader(s + "\000")
	r := io.NewSectionReader(buf, 0, int64(len(s)+1))
	coff.Strings = append(coff.Strings, r)

	coff.StringsHeader.Length += uint32(r.Size())

	coff.Symbols = append(coff.Symbols, Symbol{

		SectionNumber:  1,
		Type:           0,
		StorageClass:   2,
		AuxiliaryCount: 0,
	})
}

func NewRSRC() *Coff {
	return &Coff{
		pe.FileHeader{
			Machine:              pe.IMAGE_FILE_MACHINE_I386,
			NumberOfSections:     1,
			TimeDateStamp:        0,
			NumberOfSymbols:      1,
			SizeOfOptionalHeader: 0,
			Characteristics:      0x0104,
		},
		pe.SectionHeader32{
			Name:            STRING_RSRC,
			Characteristics: 0x40000040,
		},

		&Dir{},

		[]DataEntry{},
		[]Sizer{},

		[]RelocationEntry{},

		[]Symbol{Symbol{
			Name:           STRING_RSRC,
			Value:          0,
			SectionNumber:  1,
			Type:           0,
			StorageClass:   3,
			AuxiliaryCount: 0,
		}},

		StringsHeader{
			Length: uint32(binary.Size(StringsHeader{})),
		},
		[]Sizer{},
	}
}

func (coff *Coff) AddResource(kind uint32, id uint16, data Sizer) {
	re := RelocationEntry{

		SymbolIndex: 0,
	}
	switch coff.Machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		re.Type = _IMAGE_REL_I386_DIR32NB
	case pe.IMAGE_FILE_MACHINE_AMD64:
		re.Type = _IMAGE_REL_AMD64_ADDR32NB
	}
	coff.Relocations = append(coff.Relocations, re)
	coff.SectionHeader32.NumberOfRelocations++

	entries0 := coff.Dir.DirEntries
	dirs0 := coff.Dir.Dirs
	i0 := sort.Search(len(entries0), func(i int) bool {
		return entries0[i].NameOrId >= kind
	})
	if i0 >= len(entries0) || entries0[i0].NameOrId != kind {

		entries0 = append(entries0[:i0], append([]DirEntry{{NameOrId: kind}}, entries0[i0:]...)...)
		dirs0 = append(dirs0[:i0], append([]Dir{{}}, dirs0[i0:]...)...)
		coff.Dir.NumberOfIdEntries++
	}
	coff.Dir.DirEntries = entries0
	coff.Dir.Dirs = dirs0

	dirs0[i0].DirEntries = append(dirs0[i0].DirEntries, DirEntry{NameOrId: uint32(id)})
	dirs0[i0].Dirs = append(dirs0[i0].Dirs, Dir{
		NumberOfIdEntries: 1,
		DirEntries:        DirEntries{LANG_ENTRY},
	})
	dirs0[i0].NumberOfIdEntries++

	n := 0
	for _, dir0 := range dirs0[:i0+1] {
		n += len(dir0.DirEntries)
	}
	n--

	coff.DataEntries = append(coff.DataEntries[:n], append([]DataEntry{{Size1: uint32(data.Size())}}, coff.DataEntries[n:]...)...)
	coff.Data = append(coff.Data[:n], append([]Sizer{data}, coff.Data[n:]...)...)
}

func (coff *Coff) Freeze() {
	switch coff.SectionHeader32.Name {
	case STRING_RSRC:
		coff.freezeRSRC()
	case STRING_RDATA:
		coff.freezeRDATA()
	}
}

func (coff *Coff) freezeCommon1(path string, offset, diroff uint32) (newdiroff uint32) {
	switch path {
	case "/Dir":
		coff.SectionHeader32.PointerToRawData = offset
		diroff = offset
	case "/Relocations":
		coff.SectionHeader32.PointerToRelocations = offset
		coff.SectionHeader32.SizeOfRawData = offset - diroff
	case "/Symbols":
		coff.FileHeader.PointerToSymbolTable = offset
	}
	return diroff
}

func freezeCommon2(v reflect.Value, offset *uint32) error {
	if binutil.Plain(v.Kind()) {
		*offset += uint32(binary.Size(v.Interface()))
		return nil
	}
	vv, ok := v.Interface().(Sizer)
	if ok {
		*offset += uint32(vv.Size())
		return binutil.WALK_SKIP
	}
	return nil
}

func (coff *Coff) freezeRDATA() {
	var offset, diroff, stringsoff uint32
	binutil.Walk(coff, func(v reflect.Value, path string) error {
		diroff = coff.freezeCommon1(path, offset, diroff)

		RE := regexp.MustCompile
		const N = `\[(\d+)\]`
		m := matcher{}

		switch {
		case m.Find(path, RE("^/Data"+N+"$")):
			n := m[0]
			coff.Symbols[1+n].Value = offset - diroff
			sz := uint64(coff.Data[n].Size())
			binary.LittleEndian.PutUint64(coff.Symbols[0].Auxiliaries[0][0:8], binary.LittleEndian.Uint64(coff.Symbols[0].Auxiliaries[0][0:8])+sz)
		case path == "/StringsHeader":
			stringsoff = offset
		case m.Find(path, RE("^/Strings"+N+"$")):
			binary.LittleEndian.PutUint32(coff.Symbols[m[0]+1].Name[4:8], offset-stringsoff)
		}

		return freezeCommon2(v, &offset)
	})
	coff.SectionHeader32.PointerToRelocations = 0
}

func (coff *Coff) freezeRSRC() {
	leafwalker := make(chan *DirEntry)
	go func() {
		for _, dir1 := range coff.Dir.Dirs {
			for _, dir2 := range dir1.Dirs {
				for i := range dir2.DirEntries {
					leafwalker <- &dir2.DirEntries[i]
				}
			}
		}
	}()

	var offset, diroff uint32
	binutil.Walk(coff, func(v reflect.Value, path string) error {
		diroff = coff.freezeCommon1(path, offset, diroff)

		RE := regexp.MustCompile
		const N = `\[(\d+)\]`
		m := matcher{}
		switch {
		case m.Find(path, RE("^/Dir/Dirs"+N+"$")):
			coff.Dir.DirEntries[m[0]].OffsetToData = MASK_SUBDIRECTORY | (offset - diroff)
		case m.Find(path, RE("^/Dir/Dirs"+N+"/Dirs"+N+"$")):
			coff.Dir.Dirs[m[0]].DirEntries[m[1]].OffsetToData = MASK_SUBDIRECTORY | (offset - diroff)
		case m.Find(path, RE("^/DataEntries"+N+"$")):
			direntry := <-leafwalker
			direntry.OffsetToData = offset - diroff
		case m.Find(path, RE("^/DataEntries"+N+"/OffsetToData$")):
			coff.Relocations[m[0]].RVA = offset - diroff
		case m.Find(path, RE("^/Data"+N+"$")):
			coff.DataEntries[m[0]].OffsetToData = offset - diroff
		}

		return freezeCommon2(v, &offset)
	})
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

type matcher []int

func (m *matcher) Find(s string, re *regexp.Regexp) bool {
	subs := re.FindStringSubmatch(s)
	if subs == nil {
		return false
	}

	*m = (*m)[:0]
	for i := 1; i < len(subs); i++ {
		*m = append(*m, mustAtoi(subs[i]))
	}
	return true
}
