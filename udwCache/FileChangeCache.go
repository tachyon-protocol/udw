package udwCache

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwCryptoMd5"
	"github.com/tachyon-protocol/udw/udwDebug"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwProjectPath"
	"os"
	"path/filepath"
)

const debugFileNotMatch = false
const debugCheckedFile = false

func MustMd5FileChangeCacheClean() {
	udwFile.MustDelete(filepath.Join(udwProjectPath.MustGetProjectPath(), "tmp/FileChangeCache"))
	udwFile.MustDelete(filepath.Join(udwProjectPath.MustGetProjectPath(), "tmp/fileMd5GetDb.gob"))
}

func getFileChangeCachePath(key string) string {

	return filepath.Join(udwProjectPath.MustGetProjectPath(), "tmp/FileChangeCache", udwCryptoMd5.Md5HexFromString(key))
}

func MustMd5FileChangeCache(key string, pathList []string, f func()) (hasChange bool) {
	return MustMd5FileChangeCacheV2(key, func() (pathListF []string, hasChange bool) {
		return pathList, false
	}, f)
}

func MustMd5FileChangeCacheV2(key string, pathListF func() (pathListF []string, hasChange bool), f func()) (hasChange bool) {
	return MustMd5FileChangeCacheV3(MustMd5FileChangeCacheV3Request{
		MetaCacheFilePath: getFileChangeCachePath(key),
		Md5CacheFilePath:  filepath.Join(udwProjectPath.MustGetProjectPath(), "tmp/fileMd5GetDb.gob"),
		PathListF:         pathListF,
		F:                 f,
		Key:               key,
	})
}

type MustMd5FileChangeCacheV3Request struct {
	Key               string
	MetaCacheFilePath string
	Md5CacheFilePath  string
	PathListF         func() (pathListF []string, hasChange bool)
	F                 func()
	Ignore            func(fullPath string) bool
}

func MustMd5FileChangeCacheV3(req MustMd5FileChangeCacheV3Request) (hasChange bool) {

	toChange := false
	cacheInfo := map[string]string{}
	cacheFilePath := req.MetaCacheFilePath
	fileMd5GetDbPath := req.Md5CacheFilePath

	func() {
		content, err := udwFile.ReadFile(cacheFilePath)
		if err != nil {
			return
		}
		m, ok := unmarshalMapStringString(content)
		if ok == false {
			return
		}
		cacheInfo = m
	}()

	hasReadFileMap := map[string]bool{}
	pathList, thisHasChange := req.PathListF()
	if thisHasChange {
		toChange = true
	} else {
		FileMd5Get(fileMd5GetDbPath, func(getter FileMd5Getter) {
			mark := map[string]bool{}
			for _, path := range pathList {
				path = udwFile.MustGetFullPath(path)
				if mark[path] {
					continue
				} else {
					mark[path] = true
				}
				statList, err := udwFile.GetAllFileAndDirectoryStat(path)
				if err != nil {
					if udwFile.ErrorIsFileNotFound(err) {
						toChange = true
						if debugFileNotMatch {
							fmt.Println("[MustMd5FileChangeCacheV2] path not exist", path)
						}
						return
					}
					panic(err)
				}

				for _, stat := range statList {
					if stat.Fi.IsDir() {
						continue
					}
					if req.Ignore != nil && req.Ignore(stat.FullPath) {
						if debugCheckedFile {
							fmt.Println("dddu7ku8q9", stat.FullPath)
						}

						continue
					}
					hasReadFileMap[stat.FullPath] = true
					if stat.Fi.Mode()&os.ModeSymlink == os.ModeSymlink {
						if "symlink_"+udwFile.MustReadSymlink(stat.FullPath) != cacheInfo[stat.FullPath] {
							toChange = true
							if debugFileNotMatch {
								fmt.Println("[MustMd5FileChangeCache] Symlink not match", stat.FullPath, cacheInfo[stat.FullPath])
							}
							return
						}
						continue
					}
					cacheInfoMd5 := cacheInfo[stat.FullPath]

					if cacheInfoMd5 == "" || getter.GetMd5ByStatAndFullPath__NOLOCK(stat) != cacheInfoMd5 {
						toChange = true
						if debugFileNotMatch {
							fmt.Println("uhhqw4awce [MustMd5FileChangeCache] path mod md5 not match save file", stat.FullPath, "["+udwCryptoMd5.MustMd5File(stat.FullPath)+"]", "["+cacheInfoMd5+"]", udwDebug.GetCurrentAllStackString(0))
						}

						return
					}
				}
				if toChange {
					break
				}
			}
		})
	}
	if toChange == false {

		for fullPath := range cacheInfo {
			if !hasReadFileMap[fullPath] {
				toChange = true
				if debugFileNotMatch {
					fmt.Println("[MustMd5FileChangeCache] path been delete", fullPath)
				}
				break
			}
		}
	}
	if toChange == false {
		return false
	}
	req.F()
	cacheInfo = map[string]string{}
	pathList, thisHasChange = req.PathListF()
	if thisHasChange {
		if debugFileNotMatch {
			fmt.Println("[MustMd5FileChangeCache] PathListF return invalid2")
		}

		udwFile.MustDelete(cacheFilePath)
		return true
	}
	FileMd5Get(fileMd5GetDbPath, func(getter FileMd5Getter) {
		mark := map[string]bool{}
		for _, path := range pathList {
			path = udwFile.MustGetFullPath(path)
			if mark[path] {
				continue
			} else {
				mark[path] = true
			}
			statList, err := udwFile.GetAllFileAndDirectoryStat(path)
			if err != nil {
				panic(err)
			}
			for _, stat := range statList {
				if stat.Fi.IsDir() {
					continue
				}
				if req.Ignore != nil && req.Ignore(stat.FullPath) {

					continue
				}
				if stat.Fi.Mode()&os.ModeSymlink == os.ModeSymlink {
					linkToPath := udwFile.MustReadSymlink(stat.FullPath)
					cacheInfo[stat.FullPath] = "symlink_" + linkToPath
					continue
				}
				cacheInfo[stat.FullPath] = getter.GetMd5ByStatAndFullPath__NOLOCK(stat)
			}
		}
	})
	udwFile.MustMkdirForFile(cacheFilePath)
	udwFile.MustWriteFile(cacheFilePath, marshalMapStringString(cacheInfo))

	return true
}

func marshalMapStringString(m map[string]string) []byte {
	if m == nil || len(m) == 0 {
		return []byte{0}
	}
	_buf := udwBytes.BufWriter{}
	_buf.WriteUvarint(uint64(len(m)))
	for k, v := range m {
		_buf.WriteStringLenUvarint(k)
		_buf.WriteStringLenUvarint(v)
	}
	return _buf.GetBytes()
}

func unmarshalMapStringString(b []byte) (m map[string]string, ok bool) {
	_buf := udwBytes.NewBufReaderWithOk(b)
	l := _buf.ReadUvarint()
	if _buf.IsOk() == false {
		return nil, false
	}
	if l == 0 {
		return map[string]string{}, true
	}
	if l > uint64(_buf.GetRemainSize()/2) {
		return nil, false
	}
	m = make(map[string]string, int(l))
	for i := 0; i < int(l); i++ {
		k := _buf.ReadStringLenUvarint()
		v := _buf.ReadStringLenUvarint()
		if _buf.IsOk() == false {
			return nil, false
		}
		m[k] = v
	}
	return m, true
}
