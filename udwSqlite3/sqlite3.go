package udwSqlite3

/*
#cgo CFLAGS: -std=gnu99 -Os
#cgo CFLAGS: -DSQLITE_THREADSAFE=0 -DSQLITE_USE_URI=0 -DSQLITE_OMIT_DECLTYPE -DSQLITE_DEFAULT_MEMSTATUS=0 -DSQLITE_OMIT_DEPRECATED -DSQLITE_OMIT_UTF16 -DSQLITE_OMIT_COMPLETE -DSQLITE_OMIT_DECLTYPE -DSQLITE_OMIT_LOAD_EXTENSION -DSQLITE_OMIT_LOCALTIME -DSQLITE_OMIT_AUTHORIZATION -DSQLITE_OMIT_COMPLETE -DSQLITE_OMIT_GET_TABLE -DSQLITE_OMIT_TRACE -DSQLITE_OMIT_AUTOINIT
//#cgo CFLAGS: -DSQLITE_ENABLE_RTREE -DSQLITE_THREADSAFE
//#cgo CFLAGS: -DSQLITE_ENABLE_FTS3 -DSQLITE_ENABLE_FTS3_PARENTHESIS -DSQLITE_ENABLE_FTS4_UNICODE61
#cgo CFLAGS: -DSQLITE_TRACE_SIZE_LIMIT=15
#cgo CFLAGS: -DSQLITE_DISABLE_INTRINSIC
#cgo CFLAGS: -Wno-deprecated-declarations
#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include <sqlite3.h>
#include <stdlib.h>
#include <string.h>

#ifdef __CYGWIN__
# include <errno.h>
#endif

typedef struct{
	char* p;
	uintptr_t len;
	uintptr_t cap;
} _wgf_slice;

typedef struct {
	const _wgf_slice *p;
	uintptr_t len;
	uintptr_t cap;
} _wgf_slice_slice;

static int
_wgf_bind_text_mulit(wgf_stmt *stmt, uintptr_t inData2 ) {
	static char* emptyP = "";
	int rv = wgf_reset(stmt);
	if (rv!=SQLITE_OK && rv!=SQLITE_ROW && rv!=SQLITE_DONE){
		rv = SQLITE_OK;
		return rv;
	}
	if (inData2!=0){
		_wgf_slice_slice inData = *(_wgf_slice_slice*)inData2;
		for (int i =0;i<inData.len;i++){
			int n = i+1;
			char* p = inData.p[i].p;
			int np = (int)(inData.p[i].len);
			if (np==0){
				// never pass null into the query argument
				rv = wgf_bind_blob(stmt, n, emptyP, 0, SQLITE_TRANSIENT);
			}else{
				rv = wgf_bind_blob(stmt, n, p, np, SQLITE_TRANSIENT);
			}
			if (rv!=SQLITE_OK){
				return rv;
			}
		}
	}
	return SQLITE_OK;
}

typedef struct{
	uint8_t *p;
	int len;
	int cap;
} _wgf_slice_cacher;
static void _wgf_cacher_write_slice(_wgf_slice_cacher *cacher,const void *inP,int len){
	int thisSize = len+4;
	int needSize = cacher->len+thisSize;
	if (cacher->cap < needSize){
		int toAllocSize = needSize*2;
		if (toAllocSize<64){
			toAllocSize = 64;
		}
		cacher->p = (uint8_t*)realloc((void*)cacher->p,toAllocSize);
		cacher->cap = toAllocSize;
	}
	int cLen = cacher->len;
	cacher->p[cLen+0] = (uint8_t)(len);
	cacher->p[cLen+1] = (uint8_t)(len >> 8);
	cacher->p[cLen+2] = (uint8_t)(len >> 16);
	cacher->p[cLen+3] = (uint8_t)(len >> 24);
	if (len>0){
		memcpy((void*)(&cacher->p[cLen+4]),inP,len);
	}
	cacher->len = cacher->len+len+4;
}
static void _wgf_cacher_reset(_wgf_slice_cacher *cacher){
	cacher->len = 0;
}
static void _wgf_cacher_free(_wgf_slice_cacher *cacher){
	if (cacher->p!=NULL){
		free((void*)cacher->p);
		cacher->p = NULL;
	}
	cacher->len = 0;
	cacher->cap = 0;
}

typedef struct{
	wgf_stmt *stmt;
	int ColumnCount;
	bool autoCloseStmt;
	_wgf_slice_cacher* cacher;
} _wgf_query_2_req;

typedef struct{
	bool hasMoreRow;
	const void *rowP;
	int rowLen;
	int rv; // SQLITE_OK normal others error,will close stmt after error
} _wgf_query_2_resp;
static _wgf_query_2_resp _wgf_query_2(_wgf_query_2_req req){
	_wgf_query_2_resp resp = {};
	resp.hasMoreRow = false;
	resp.rv = SQLITE_OK;
	_wgf_cacher_reset(req.cacher);
	for (int i = 0;i<req.ColumnCount;i++){
		const void *p = wgf_column_blob(req.stmt,i);
		int thisLen = 0;
		if (p!=NULL){
			thisLen = wgf_column_bytes(req.stmt,i);
		}
		_wgf_cacher_write_slice(req.cacher,p,thisLen);
	}
	resp.rowP = req.cacher->p;
	resp.rowLen = req.cacher->len;
	resp.rv = wgf_step(req.stmt);
	if (resp.rv==SQLITE_DONE){
		if (req.autoCloseStmt==1){
			wgf_finalize(req.stmt);
		}
		resp.rv = SQLITE_OK;
		return resp;
	}
	if (resp.rv != SQLITE_ROW){
		resp.rv = wgf_reset(req.stmt);
		if (resp.rv == SQLITE_OK){
			if (req.autoCloseStmt==1){
				wgf_finalize(req.stmt);
			}
		}else{
			wgf_finalize(req.stmt);
		}
		return resp;
	}
	resp.hasMoreRow = true;
	if (resp.rv == SQLITE_ROW){
		resp.rv = SQLITE_OK;
	}
	return resp;
}

typedef struct{
	wgf* db;
	char *query;
	int queryLen;
	wgf_stmt *stmt;
	uintptr_t bind_slice_slice;
	int ColumnCount;
	bool needResult;
	bool autoCloseStmt;
	_wgf_slice_cacher* cacher;
} _wgf_query_1_req;

typedef struct{
	int rv; // SQLITE_OK normal others error,will close stmt after error
	wgf_stmt *stmt;
	bool hasMoreRow;
	const void *rowP;
	int rowLen;
	bool hasTail;
	int ColumnCount;
} _wgf_query_1_resp;

static _wgf_query_1_resp _wgf_query_1(uintptr_t preq){
	_wgf_query_1_req req = *(_wgf_query_1_req*)preq;
	_wgf_query_1_resp resp = {};
	resp.hasMoreRow = false;
	resp.rv = SQLITE_OK;
	if (req.stmt==NULL){
		const char* tail;
		resp.rv = wgf_prepare_v2(req.db,req.query,req.queryLen,&req.stmt,&tail);
		if (resp.rv!=SQLITE_OK){
			return resp;
		}
		if (tail!=NULL && tail[0]!=0){
			resp.hasTail = true;
			wgf_finalize(req.stmt);
			return resp;
		}
	}
	resp.stmt = req.stmt;
	resp.rv = _wgf_bind_text_mulit(req.stmt,req.bind_slice_slice);
	if (resp.rv!=SQLITE_OK){
		wgf_finalize(resp.stmt);
		return resp;
	}
	resp.rv = wgf_step(req.stmt);
	if (resp.rv==SQLITE_DONE){
		if (req.autoCloseStmt==true){
			wgf_finalize(req.stmt);
		}
		resp.rv = SQLITE_OK;
		return resp;
	}else if (resp.rv != SQLITE_ROW){
		resp.rv = wgf_reset(req.stmt);
		if (resp.rv == SQLITE_OK){
			if (req.autoCloseStmt==true){
				wgf_finalize(req.stmt);
			}
		}else{
			wgf_finalize(resp.stmt);
		}
		return resp;
	}
	if (req.needResult==0){
		if (req.autoCloseStmt==true){
			wgf_finalize(req.stmt);
		}
		return resp;
	}
	if (req.ColumnCount==0){
		req.ColumnCount = wgf_column_count(req.stmt);
	}
	resp.ColumnCount = req.ColumnCount;
	_wgf_query_2_req req2 = {};
	req2.stmt = req.stmt;
	req2.ColumnCount = req.ColumnCount;
	req2.autoCloseStmt = req.autoCloseStmt;
	req2.cacher = req.cacher;
	_wgf_query_2_resp resp2 = _wgf_query_2(req2);
	resp.hasMoreRow = resp2.hasMoreRow;
	resp.rowP = resp2.rowP;
	resp.rowLen = resp2.rowLen;
	resp.rv = resp2.rv;
	return resp;
}
*/
import "C"
import (
	"encoding/binary"
	"errors"
	"github.com/tachyon-protocol/udw/udwBytes"
	"runtime"
	"strconv"
	"sync"
	"unsafe"
)

var gSqliteInit sync.Once
var gSqliteInitErrMsg string

type newDbReq struct {
	UsingMemory bool
	FilePath    string
	BusyTimeout int
}

func (db *Db) newInnerDb(req newDbReq) (errMsg string) {

	if req.BusyTimeout == 0 {
		req.BusyTimeout = 5000
	}
	name := req.FilePath
	if req.UsingMemory {
		name = ":memory:"
	}
	gSqliteInit.Do(func() {
		rv := C.wgf_initialize()
		if rv != 0 {
			gSqliteInitErrMsg = "janvcb9ey5 " + errorCodeToString__NOLOCK(nil, rv)
		}
	})
	if gSqliteInitErrMsg != "" {
		return gSqliteInitErrMsg
	}
	cacheBuf := &udwBytes.BufWriter{}
	cacheBuf.WriteString(name)
	cacheBuf.WriteByte(0)
	nameB := cacheBuf.GetBytes()
	var cDb *C.wgf
	rv := C.wgf_open_v2((*C.char)(unsafe.Pointer(&nameB[0])), &cDb,
		C.SQLITE_OPEN_NOMUTEX|
			C.SQLITE_OPEN_READWRITE|
			C.SQLITE_OPEN_CREATE,
		nil)
	runtime.KeepAlive(nameB)
	if rv != 0 {
		return "yrs42vdqts " + errorCodeToString__NOLOCK(nil, rv)
	}
	if cDb == nil {
		return "jg2jtzbarh"
	}
	db.db = cDb
	rv = C.wgf_extended_result_codes(cDb, C.int(1))
	if rv != C.SQLITE_OK {
		C.wgf_close_v2(cDb)
		return "nfxqrhzhqc " + errorCodeToString__NOLOCK(db, rv)
	}
	rv = C.wgf_busy_timeout(cDb, C.int(req.BusyTimeout))
	if rv != C.SQLITE_OK {
		C.wgf_close_v2(cDb)
		return "hsqcwp8h5p " + errorCodeToString__NOLOCK(db, rv)
	}

	return ""
}

func (c *Db) cClose() error {
	c.cLocker.Lock()
	if c.db == nil {
		c.cLocker.Unlock()
		return nil
	}
	rv := C.wgf_close_v2(c.db)
	c.db = nil
	if rv != C.SQLITE_OK {
		err := errors.New("pveze8cy8g " + errorCodeToString__NOLOCK(nil, rv))
		c.cLocker.Unlock()
		return err
	}
	c.cLocker.Unlock()
	return nil
}

type QueryReq struct {
	Query        string
	Args         [][]byte
	ColumnsCb    func(columnNameList [][]byte)
	RespStatusCb func(status QueryRespStatus)
	RespDataCb   func(result [][]byte)
	ColumnCount  int
	UseStmtCache bool
}
type QueryRespStatus struct {
	AffectedRows uint64
	LastInsertId uint64
}

func (c *Db) Query(req QueryReq) (errMsg string) {
	req.UseStmtCache = false
	rpc1_req := C._wgf_query_1_req{}
	var cacheBuf udwBytes.BufWriter
	c.updateQueryString(&rpc1_req, req, &cacheBuf)
	args := req.Args
	if len(args) > 0 {
		rpc1_req.bind_slice_slice = C.uintptr_t(uintptr(unsafe.Pointer(&args)))
	}
	cCacher := &C._wgf_slice_cacher{}
	rpc1_req.ColumnCount = C.int(req.ColumnCount)
	rpc1_req.needResult = C.bool(req.RespDataCb != nil)
	rpc1_req.autoCloseStmt = C.bool(req.UseStmtCache == false && req.ColumnsCb == nil)
	rpc1_req.cacher = cCacher
	c.cLocker.Lock()
	rpc1_req.db = c.db
	if c.db == nil {
		c.cLocker.Unlock()
		return "tfyjvsyw2g"
	}
	rpc1_resp := C._wgf_query_1(C.uintptr_t(uintptr(unsafe.Pointer(&rpc1_req))))
	runtime.KeepAlive(cacheBuf)
	if rpc1_resp.rv != C.SQLITE_OK || rpc1_resp.hasTail {
		if rpc1_resp.hasTail {
			c.cLocker.Unlock()
			return "desfgdu59n"
		}
		if cCacher != nil {
			C._wgf_cacher_free(cCacher)
			cCacher = nil
		}
		if rpc1_resp.rv == C.SQLITE_ROW {
			c.cLocker.Unlock()
			return ""
		}
		errMsg = "znnsuywp72 " + errorCodeToString__NOLOCK(c, rpc1_resp.rv)
		c.cLocker.Unlock()
		return errMsg
	}
	c.cLocker.Unlock()
	req.ColumnCount = int(rpc1_resp.ColumnCount)
	if req.ColumnsCb != nil {
		resultList := c.getResultListCache(req.ColumnCount)
		c.cLocker.Lock()
		for i := 0; i < req.ColumnCount; i++ {
			columnS := C.GoString(C.wgf_column_name(rpc1_resp.stmt, C.int(i)))
			resultList[i] = []byte(columnS)
		}
		c.cLocker.Unlock()
		req.ColumnsCb(resultList)
		if rpc1_resp.hasMoreRow == false && req.UseStmtCache == false {
			c.closeStmt(rpc1_resp.stmt)
		}
	}
	if req.RespDataCb != nil {
		hasMoreRow := rpc1_resp.hasMoreRow
		if rpc1_resp.rowLen > 0 {
			resultList := c.unmarshalResultListFromC__NOLOCK(req.ColumnCount, rpc1_resp.rowP, int(rpc1_resp.rowLen))
			req.RespDataCb(resultList)
		}
		for {
			if hasMoreRow == false {
				break
			}
			rpc2_req := C._wgf_query_2_req{}
			rpc2_req.stmt = rpc1_resp.stmt
			rpc2_req.ColumnCount = C.int(req.ColumnCount)
			rpc2_req.autoCloseStmt = rpc1_req.autoCloseStmt
			rpc2_req.cacher = cCacher
			c.cLocker.Lock()
			rpc2_resp := C._wgf_query_2(rpc2_req)
			if rpc2_resp.rv != C.SQLITE_OK {
				if cCacher != nil {
					C._wgf_cacher_free(cCacher)
					cCacher = nil
				}
				errMsg := "hue5scz2dw " + errorCodeToString__NOLOCK(c, rpc2_resp.rv)
				c.cLocker.Unlock()
				return errMsg
			}
			c.cLocker.Unlock()
			resultList := c.unmarshalResultListFromC__NOLOCK(req.ColumnCount, rpc2_resp.rowP, int(rpc2_resp.rowLen))
			req.RespDataCb(resultList)
			hasMoreRow = rpc2_resp.hasMoreRow
		}
	}
	if req.RespStatusCb != nil {
		vQueryRespStatus := c.getRespStatusCb()
		req.RespStatusCb(vQueryRespStatus)
	}
	if cCacher != nil {
		C._wgf_cacher_free(cCacher)
		cCacher = nil
	}
	return ""
}

func (c *Db) updateQueryString(rpc1_req *C._wgf_query_1_req, req QueryReq, cacheBuf *udwBytes.BufWriter) {
	var l int
	var p unsafe.Pointer
	var queryBuf []byte

	if req.Query[len(req.Query)-1] == 0 {
		p = unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&req.Query)))
		l = len(req.Query)
	} else {
		cacheBuf.Reset()
		cacheBuf.WriteString(req.Query)
		cacheBuf.WriteByte(0)
		queryBuf = cacheBuf.GetBytes()
		p = unsafe.Pointer(&queryBuf[0])
		l = len(queryBuf)
	}
	rpc1_req.query = (*C.char)(p)
	rpc1_req.queryLen = C.int(l)
}

func (c *Db) unmarshalResultListFromC__NOLOCK(ColumnCount int, rowP unsafe.Pointer, rowLen int) [][]byte {
	resultList := c.getResultListCache(ColumnCount)
	rowBinaryBuf := (*[1 << 30]byte)(unsafe.Pointer(rowP))[0:rowLen]
	pos := 0
	for i := 0; i < ColumnCount; i++ {
		thisLen := int(binary.LittleEndian.Uint32(rowBinaryBuf[pos:]))
		pos += 4
		resultList[i] = rowBinaryBuf[pos : pos+thisLen]
		pos += int(thisLen)
	}
	return resultList
}

func (c *Db) getRespStatusCb() QueryRespStatus {
	c.cLocker.Lock()
	AffectedRows := uint64(C.wgf_changes(c.db))
	LastInsertId := uint64(C.wgf_last_insert_rowid(c.db))
	c.cLocker.Unlock()
	return QueryRespStatus{
		AffectedRows: AffectedRows,
		LastInsertId: LastInsertId,
	}
}

func (c *Db) getResultListCache(needSize int) [][]byte {
	return make([][]byte, needSize)
}

func (c *Db) closeStmt(stmt *C.wgf_stmt) (errMsg string) {
	c.cLocker.Lock()
	rv := C.wgf_finalize(stmt)
	if rv != C.SQLITE_OK {
		errMsg = "rb8w2cw79g " + errorCodeToString__NOLOCK(c, rv)
	}
	c.cLocker.Unlock()
	return errMsg
}

func errorCodeToString__NOLOCK(c *Db, rv C.int) string {

	db := c.db
	errMsg := ""
	if db != nil {
		errMsg = " " + C.GoString(C.wgf_errmsg(db))
	}
	return strconv.Itoa(int(rv)) + errMsg
}

type Db struct {
	db                      *C.wgf
	cLocker                 sync.Mutex
	req                     NewDbRequest
	initLevel               int
	initLevelLocker         sync.Mutex
	tableNameCacheMapLocker sync.Mutex
	tableNameCacheMap       map[string]string
}

func MustGetSqlite3Version() string {
	return C.GoString(C.wgf_libversion())
}
