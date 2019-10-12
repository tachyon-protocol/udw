package udwCache

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwFile"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

type FileMd5Getter struct {
	imp *fileMd5GetterImp
}

func (getter FileMd5Getter) GetMd5ByFullPath__NOLOCK(path string) string {
	return getter.imp.getMd5ByFullPath__NOLOCK(path)
}
func (getter FileMd5Getter) GetMd5ByStatAndFullPath__NOLOCK(statAndFullPath udwFile.StatAndFullPath) string {
	return getter.imp.getMd5ByStatAndFullPath__NOLOCK(statAndFullPath)
}
func (getter FileMd5Getter) Close() {
	getter.imp.close__NOLOCK()
}

type FileMd5GetByFullPath func(path string) string

func FileMd5Get(dbFilePath string, cb func(getter FileMd5Getter)) {
	getterImp := getFileMd5GetterByDbFilePath(dbFilePath)
	getterImp.locker.Lock()
	getterImp.cacheNow = time.Now()
	getterImp.hasChange = false
	getterImp.dbFilePath = dbFilePath
	outGetter := FileMd5Getter{imp: getterImp}
	cb(outGetter)
	if getterImp.hasChange {
		writeMd5fileChangeCacheToFile(dbFilePath, getterImp)

	}
	getterImp.locker.Unlock()
}

type fileMd5GetterImp struct {
	CacheInfo  map[string]*cacheInfoEntry
	locker     sync.Mutex
	hasChange  bool
	cacheNow   time.Time
	innerBuf   []byte
	dbFilePath string
}
type cacheInfoEntry struct {
	MTime                time.Time
	Md5                  string
	md5B                 [16]byte
	HasCheckInNextSecond bool
	hasChange            bool
}

func (cie *cacheInfoEntry) SetMd5String(md5 string) {
	hex.Decode(cie.md5B[:], []byte(md5))
	cie.Md5 = ""
}
func (cie *cacheInfoEntry) GetMd5String() string {
	if cie.hasChange {
		cie.hasChange = false
		cie.Md5 = hex.EncodeToString(cie.md5B[:])
		return cie.Md5
	}
	if cie.Md5 != "" {
		return cie.Md5
	}
	cie.Md5 = hex.EncodeToString(cie.md5B[:])
	return cie.Md5
}

func FileMd5TryFreeMemory() {
	fileMd5GetterMapLock.Lock()
	fileMd5GetterMap = map[string]*fileMd5GetterImp{}
	fileMd5GetterMapLock.Unlock()
}

var fileMd5GetterMap map[string]*fileMd5GetterImp
var fileMd5GetterMapLock sync.Mutex

func (getter *fileMd5GetterImp) getMd5ByFullPath__NOLOCK(path string) string {
	fi, err := os.Lstat(path)
	if err != nil {
		return ""
	}
	return getter.getMd5ByStatAndFullPath__NOLOCK(udwFile.StatAndFullPath{
		Fi:       fi,
		FullPath: path,
	})
}
func (getter *fileMd5GetterImp) getMd5ByStatAndFullPath__NOLOCK(statAndFullPath udwFile.StatAndFullPath) string {
	path := statAndFullPath.FullPath
	fi := statAndFullPath.Fi
	var err error
	thisEntry := getter.CacheInfo[path]
	if thisEntry == nil {
		thisEntry = &cacheInfoEntry{}
		thisEntry.MTime = fi.ModTime()
		thisEntry.md5B, err = gMd5FileWithBuf(path, getter.innerBuf)
		if err != nil {
			return ""
		}
		thisEntry.HasCheckInNextSecond = (getter.cacheNow.Sub(thisEntry.MTime) > time.Second)
		getter.CacheInfo[path] = thisEntry
		getter.hasChange = true
		return thisEntry.GetMd5String()
	}

	fileMtime := fi.ModTime()
	if !fileMtime.Equal(thisEntry.MTime) {

		thisEntry.MTime = fileMtime
		thisEntry.md5B, err = gMd5FileWithBuf(path, getter.innerBuf)
		if err != nil {
			return ""
		}
		thisEntry.HasCheckInNextSecond = (getter.cacheNow.Sub(thisEntry.MTime) > time.Second)
		thisEntry.hasChange = true
		getter.hasChange = true
		return thisEntry.GetMd5String()
	}
	now := getter.cacheNow
	if now.Sub(thisEntry.MTime) < time.Second {

		thisMd5, err := gMd5FileWithBuf(path, getter.innerBuf)
		if err != nil {
			return ""
		}
		if bytes.Equal(thisMd5[:], thisEntry.md5B[:]) == false {
			getter.hasChange = true
			thisEntry.hasChange = true
		}
		thisEntry.md5B = thisMd5
		thisEntry.HasCheckInNextSecond = false
		return thisEntry.GetMd5String()
	}
	if thisEntry.HasCheckInNextSecond == false {

		thisEntry.MTime = fileMtime
		thisMd5, err := gMd5FileWithBuf(path, getter.innerBuf)
		if err != nil {
			return ""
		}
		thisEntry.HasCheckInNextSecond = (bytes.Equal(thisMd5[:], thisEntry.md5B[:]))
		getter.hasChange = true
		thisEntry.hasChange = true
		thisEntry.md5B = thisMd5
		return thisEntry.GetMd5String()
	}

	return thisEntry.GetMd5String()
}
func (getter *fileMd5GetterImp) close__NOLOCK() {
	if getter.hasChange {
		writeMd5fileChangeCacheToFile(getter.dbFilePath, getter)

	}
}

func getFileMd5GetterByDbFilePath(dbFilePath string) *fileMd5GetterImp {
	fileMd5GetterMapLock.Lock()
	defer fileMd5GetterMapLock.Unlock()
	if fileMd5GetterMap == nil {
		fileMd5GetterMap = map[string]*fileMd5GetterImp{}
	}
	thisGetter := fileMd5GetterMap[dbFilePath]
	if thisGetter != nil {
		return thisGetter
	}
	thisGetter = &fileMd5GetterImp{}

	isValid := readMd5fileChangeCacheFromFile(dbFilePath, thisGetter)

	if !isValid {

		thisGetter = &fileMd5GetterImp{}
	}
	if thisGetter.CacheInfo == nil {
		thisGetter.CacheInfo = map[string]*cacheInfoEntry{}
	}
	if thisGetter.innerBuf == nil {
		thisGetter.innerBuf = make([]byte, 32*1024)
	}
	fileMd5GetterMap[dbFilePath] = thisGetter
	return thisGetter
}

var gMd5FileWithBuf func(path string, innerBuf []byte) ([16]byte, error) = md5FileWithBuf

func SetMd5FileWithBuf(cb func(path string, innerBuf []byte) ([16]byte, error)) {
	gMd5FileWithBuf = cb
}

func md5FileWithBuf(path string, innerBuf []byte) (output [16]byte, err error) {
	hash := md5.New()
	f, err := os.Open(path)
	if err != nil {
		return output, err
	}
	_, err = io.CopyBuffer(hash, f, innerBuf)
	f.Close()
	if err != nil {
		return output, err
	}
	hashB := hash.Sum(innerBuf[0:0])
	copy(output[:], hashB)
	return output, nil
}

func readMd5fileChangeCacheFromFile(path string, imp *fileMd5GetterImp) (isValid bool) {
	content, err := udwFile.ReadFile(path)
	if err != nil {

		return false
	}
	pos := 0
	if len(content) < 4 {

		return false
	}
	length := len(content)
	size := int(binary.LittleEndian.Uint32(content[:4]))
	if size*18 > length {

		fmt.Println("[readMd5fileChangeCacheFromFile] size[" + strconv.Itoa(size) + "] too big")
		return false
	}
	imp.CacheInfo = make(map[string]*cacheInfoEntry, size)
	pos = 4
	for i := 0; i < size; i++ {
		entry := &cacheInfoEntry{}
		name, advance, isValid := readString(content[pos:])
		if !isValid {
			fmt.Printf("[readMd5fileChangeCacheFromFile] read string name fail pos[%d]\n", pos)
			return false
		}
		pos += advance
		if length < pos+8 {
			fmt.Printf("[readMd5fileChangeCacheFromFile] read time fail pos[%d]\n", pos)
			return false
		}
		unixNano := int64(binary.LittleEndian.Uint64(content[pos : pos+8]))
		pos += 8
		entry.MTime = time.Unix(unixNano/1e9, unixNano%1e9)

		if length < pos+16 {
			fmt.Println("[readMd5fileChangeCacheFromFile] 4fxc6g79dc ", pos)
			return false
		}
		copy(entry.md5B[:], content[pos:pos+16])

		pos += 16
		if length < pos+1 {
			fmt.Printf("[readMd5fileChangeCacheFromFile] read bool HasCheckInNextSecond fail pos[%d]\n", pos)
			return false
		}
		HasCheckInNextSecond := false
		if content[pos] != 0 {
			HasCheckInNextSecond = true
		}
		pos += 1

		entry.HasCheckInNextSecond = HasCheckInNextSecond
		imp.CacheInfo[name] = entry
	}
	if pos != length {
		fmt.Printf("[readMd5fileChangeCacheFromFile] pos!=length after read pos[%d] length[%d]]\n", pos, length)
		return false
	}
	return true
}

func writeMd5fileChangeCacheToFile(path string, imp *fileMd5GetterImp) {
	_buf := &udwBytes.BufWriter{}
	size := len(imp.CacheInfo)
	_buf.WriteByte(byte(size))
	_buf.WriteByte(byte(size >> 8))
	_buf.WriteByte(byte(size >> 16))
	_buf.WriteByte(byte(size >> 24))
	for name, entry := range imp.CacheInfo {
		writeString(_buf, name)
		unixNano := entry.MTime.UnixNano()
		_buf.WriteLittleEndUint64(uint64(unixNano))
		_buf.Write(entry.md5B[:])
		if entry.HasCheckInNextSecond {
			_buf.WriteByte(1)
		} else {
			_buf.WriteByte(0)
		}
	}
	udwFile.MustWriteFileWithMkdir(path, _buf.GetBytes())
}

func readString(content []byte) (out string, advance int, isValid bool) {
	if len(content) == 0 {
		return "", 0, false
	}
	size := int(content[0])*256 + int(content[1])
	if len(content) < size+2 {
		return "", 0, false
	}
	return string(content[2 : size+2]), size + 2, true
}

func writeString(_buf *udwBytes.BufWriter, s string) {
	if len(s) > 1<<16-1 {
		panic(fmt.Errorf("[writeString] can not write string length more than 2^16,s[%s]", s))
	}
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(len(s)))
	_buf.WriteByte(buf[0])
	_buf.WriteByte(buf[1])
	_buf.WriteString(s)
}
