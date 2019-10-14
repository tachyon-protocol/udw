package rsrc

import (
	"encoding/binary"
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwWindowsUI/rsrc/binutil"
	"github.com/tachyon-protocol/udw/udwWindowsUI/rsrc/coff"
	"github.com/tachyon-protocol/udw/udwWindowsUI/rsrc/ico"
	"io"
	"os"
	"reflect"
)

const (
	RT_ICON       = coff.RT_ICON
	RT_GROUP_ICON = coff.RT_GROUP_ICON
	RT_MANIFEST   = coff.RT_MANIFEST
)

type GRPICONDIR struct {
	ico.ICONDIR
	Entries []GRPICONDIRENTRY
}

func (group GRPICONDIR) Size() int64 {
	return int64(binary.Size(group.ICONDIR) + len(group.Entries)*binary.Size(group.Entries[0]))
}

type GRPICONDIRENTRY struct {
	ico.IconDirEntryCommon
	Id uint16
}

func run(req MustBuildRequest) error {
	nextId := uint16(0)
	newId := func() uint16 {
		nextId++
		return nextId
	}

	coffObj := coff.NewRSRC()
	err := coffObj.Arch(req.Arch)
	if err != nil {
		return err
	}

	if len(req.ManifestFileContent) > 0 {
		manifest := binutil.SizeFileFromBuffer([]byte(req.ManifestFileContent))

		coffObj.AddResource(RT_MANIFEST, newId(), manifest)
	}
	if len(req.IconContentList) > 0 {
		for _, iconContent := range req.IconContentList {
			err := addicon(coffObj, iconContent, newId)
			if err != nil {
				return err
			}
		}
	}

	coffObj.Freeze()

	return write(coffObj, req.OutputFilePath)
}

func addicon(coff *coff.Coff, iconContent []byte, newid func() uint16) error {
	f := udwBytes.NewBufReader(iconContent)

	icons, err := ico.DecodeHeaders(f)
	if err != nil {
		return err
	}

	if len(icons) > 0 {

		group := GRPICONDIR{ICONDIR: ico.ICONDIR{
			Reserved: 0,
			Type:     1,
			Count:    uint16(len(icons)),
		}}
		for _, icon := range icons {
			id := newid()
			r := io.NewSectionReader(f, int64(icon.ImageOffset), int64(icon.BytesInRes))
			coff.AddResource(RT_ICON, id, r)
			group.Entries = append(group.Entries, GRPICONDIRENTRY{icon.IconDirEntryCommon, id})
		}
		id := newid()
		coff.AddResource(RT_GROUP_ICON, id, group)

	}

	return nil
}

func write(coff *coff.Coff, fnameout string) error {
	udwFile.MustMkdirForFile(fnameout)
	out, err := os.Create(fnameout)
	if err != nil {
		return err
	}
	defer out.Close()
	w := binutil.Writer{W: out}

	binutil.Walk(coff, func(v reflect.Value, path string) error {
		if binutil.Plain(v.Kind()) {
			w.WriteLE(v.Interface())
			return nil
		}
		vv, ok := v.Interface().(binutil.SizedReader)
		if ok {
			w.WriteFromSized(vv)
			return binutil.WALK_SKIP
		}
		return nil
	})

	if w.Err != nil {
		return fmt.Errorf("Error writing output file: %s", w.Err)
	}

	return nil
}
