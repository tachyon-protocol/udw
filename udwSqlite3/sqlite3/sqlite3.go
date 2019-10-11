package sqlite3

/*
#cgo CFLAGS: -std=gnu99 -Os
#cgo CFLAGS: -DSQLITE_THREADSAFE=1 -DSQLITE_USE_URI=0 -DSQLITE_OMIT_DECLTYPE -DSQLITE_DEFAULT_MEMSTATUS=0 -DSQLITE_OMIT_DEPRECATED -DSQLITE_OMIT_UTF16 -DSQLITE_OMIT_COMPLETE -DSQLITE_OMIT_DECLTYPE -DSQLITE_OMIT_LOAD_EXTENSION -DSQLITE_OMIT_LOCALTIME -DSQLITE_OMIT_AUTHORIZATION -DSQLITE_OMIT_COMPLETE -DSQLITE_OMIT_GET_TABLE -DSQLITE_OMIT_TRACE -DSQLITE_OMIT_AUTOINIT
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
} _zxr_slice;

typedef struct {
	const _zxr_slice *p;
	uintptr_t len;
	uintptr_t cap;
} _zxr_slice_slice;

static int
_zxr_bind_text_mulit(zxr_stmt *stmt, uintptr_t inData2 ) {
	static char* emptyP = "";
	int rv = zxr_reset(stmt);
	if (rv!=SQLITE_OK && rv!=SQLITE_ROW && rv!=SQLITE_DONE){
		rv = SQLITE_OK;
		return rv;
	}
	if (inData2!=0){
		_zxr_slice_slice inData = *(_zxr_slice_slice*)inData2;
		for (int i =0;i<inData.len;i++){
			int n = i+1;
			char* p = inData.p[i].p;
			int np = (int)(inData.p[i].len);
			if (np==0){
				// never pass null into the query argument
				rv = zxr_bind_blob(stmt, n, emptyP, 0, SQLITE_TRANSIENT);
			}else{
				rv = zxr_bind_blob(stmt, n, p, np, SQLITE_TRANSIENT);
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
} _zxr_slice_cacher;
static void _zxr_cacher_write_slice(_zxr_slice_cacher *cacher,const void *inP,int len){
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
static void _zxr_cacher_reset(_zxr_slice_cacher *cacher){
	cacher->len = 0;
}
static void _zxr_cacher_free(_zxr_slice_cacher *cacher){
	if (cacher->p!=NULL){
		free((void*)cacher->p);
		cacher->p = NULL;
	}
	cacher->len = 0;
	cacher->cap = 0;
}

typedef struct{
	zxr_stmt *stmt;
	int ColumnCount;
	bool autoCloseStmt;
	_zxr_slice_cacher* cacher;
} _zxr_query_2_req;

typedef struct{
	bool hasMoreRow;
	const void *rowP;
	int rowLen;
	int rv; // SQLITE_OK normal others error,will close stmt after error
} _zxr_query_2_resp;
static _zxr_query_2_resp _zxr_query_2(_zxr_query_2_req req){
	_zxr_query_2_resp resp = {};
	resp.hasMoreRow = false;
	resp.rv = SQLITE_OK;
	_zxr_cacher_reset(req.cacher);
	for (int i = 0;i<req.ColumnCount;i++){
		const void *p = zxr_column_blob(req.stmt,i);
		int thisLen = 0;
		if (p!=NULL){
			thisLen = zxr_column_bytes(req.stmt,i);
		}
		_zxr_cacher_write_slice(req.cacher,p,thisLen);
	}
	resp.rowP = req.cacher->p;
	resp.rowLen = req.cacher->len;
	resp.rv = zxr_step(req.stmt);
	if (resp.rv==SQLITE_DONE){
		if (req.autoCloseStmt==1){
			zxr_finalize(req.stmt);
		}
		resp.rv = SQLITE_OK;
		return resp;
	}
	if (resp.rv != SQLITE_ROW){
		resp.rv = zxr_reset(req.stmt);
		if (resp.rv == SQLITE_OK){
			if (req.autoCloseStmt==1){
				zxr_finalize(req.stmt);
			}
		}else{
			zxr_finalize(req.stmt);
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
	zxr* db;
	char *query;
	int queryLen;
	zxr_stmt *stmt;
	uintptr_t bind_slice_slice;
	int ColumnCount;
	bool needResult;
	bool autoCloseStmt;
	_zxr_slice_cacher* cacher;
} _zxr_query_1_req;

typedef struct{
	int rv; // SQLITE_OK normal others error,will close stmt after error
	zxr_stmt *stmt;
	bool hasMoreRow;
	const void *rowP;
	int rowLen;
	bool hasTail;
	int ColumnCount;
} _zxr_query_1_resp;

static _zxr_query_1_resp _zxr_query_1(uintptr_t preq){
	_zxr_query_1_req req = *(_zxr_query_1_req*)preq;
	_zxr_query_1_resp resp = {};
	resp.hasMoreRow = false;
	resp.rv = SQLITE_OK;
	if (req.stmt==NULL){
		const char* tail;
		resp.rv = zxr_prepare_v2(req.db,req.query,req.queryLen,&req.stmt,&tail);
		if (resp.rv!=SQLITE_OK){
			return resp;
		}
		if (tail!=NULL && tail[0]!=0){
			resp.hasTail = true;
			zxr_finalize(req.stmt);
			return resp;
		}
	}
	resp.stmt = req.stmt;
	resp.rv = _zxr_bind_text_mulit(req.stmt,req.bind_slice_slice);
	if (resp.rv!=SQLITE_OK){
		zxr_finalize(resp.stmt);
		return resp;
	}
	resp.rv = zxr_step(req.stmt);
	if (resp.rv==SQLITE_DONE){
		if (req.autoCloseStmt==true){
			zxr_finalize(req.stmt);
		}
		resp.rv = SQLITE_OK;
		return resp;
	}else if (resp.rv != SQLITE_ROW){
		resp.rv = zxr_reset(req.stmt);
		if (resp.rv == SQLITE_OK){
			if (req.autoCloseStmt==true){
				zxr_finalize(req.stmt);
			}
		}else{
			zxr_finalize(resp.stmt);
		}
		return resp;
	}
	if (req.needResult==0){
		if (req.autoCloseStmt==true){
			zxr_finalize(req.stmt);
		}
		return resp;
	}
	if (req.ColumnCount==0){
		req.ColumnCount = zxr_column_count(req.stmt);
	}
	resp.ColumnCount = req.ColumnCount;
	_zxr_query_2_req req2 = {};
	req2.stmt = req.stmt;
	req2.ColumnCount = req.ColumnCount;
	req2.autoCloseStmt = req.autoCloseStmt;
	req2.cacher = req.cacher;
	_zxr_query_2_resp resp2 = _zxr_query_2(req2);
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
	"github.com/tachyon-protocol/udw/udwBytes"
	"sync"
	"unsafe"
)

func errorString(err Error) string {
	return C.GoString(C.zxr_errstr(C.int(err.Code)))
}

func lastError(db *C.zxr) error {
	rv := C.zxr_errcode(db)
	if rv == C.SQLITE_OK {
		return nil
	}
	return Error{
		Code:         ErrNo(rv),
		ExtendedCode: ErrNoExtended(C.zxr_extended_errcode(db)),
		err:          C.GoString(C.zxr_errmsg(db)),
	}
}

var gSqliteInit sync.Once
var gSqliteInitErrMsg string

type NewDbReq struct {
	UsingMemory bool
	FilePath    string
	BusyTimeout int
}

func NewDb(req NewDbReq) (myDb *Db, errMsg string) {
	if req.BusyTimeout == 0 {
		req.BusyTimeout = 5000
	}
	name := req.FilePath
	if req.UsingMemory {
		name = ":memory:"
	}
	gSqliteInit.Do(func() {
		rv := C.zxr_initialize()
		if rv != 0 {
			gSqliteInitErrMsg = Error{Code: ErrNo(rv)}.Error()
		}
	})
	if gSqliteInitErrMsg != "" {
		return nil, gSqliteInitErrMsg
	}
	myDb = &Db{}
	myDb.stmtCache = map[string]*C.zxr_stmt{}
	myDb.cacheBuf.WriteString(name)
	myDb.cacheBuf.WriteByte(0)
	myDb.cCacher = &C._zxr_slice_cacher{}
	nameB := myDb.cacheBuf.GetBytes()
	var db *C.zxr
	rv := C.zxr_open_v2((*C.char)(unsafe.Pointer(&nameB[0])), &db,
		C.SQLITE_OPEN_NOMUTEX|
			C.SQLITE_OPEN_READWRITE|
			C.SQLITE_OPEN_CREATE,
		nil)
	if rv != 0 {
		return nil, Error{Code: ErrNo(rv)}.Error()
	}
	if db == nil {
		return nil, "jg2jtzbarh"
	}
	rv = C.zxr_busy_timeout(db, C.int(req.BusyTimeout))
	if rv != C.SQLITE_OK {
		C.zxr_close_v2(db)
		return nil, Error{Code: ErrNo(rv)}.Error()
	}
	myDb.db = db

	return myDb, ""
}

type Db struct {
	db              *C.zxr
	locker          sync.Mutex
	cacheBuf        udwBytes.BufWriter
	resultListCache [][]byte
	stmtCache       map[string]*C.zxr_stmt
	cCacher         *C._zxr_slice_cacher
}

func (c *Db) Close() error {
	if c.db == nil {
		return nil
	}
	c.cleanStmtCache()
	c.stmtCache = nil
	if c.cCacher != nil {
		C._zxr_cacher_free(c.cCacher)
		c.cCacher = nil
	}
	rv := C.zxr_close_v2(c.db)
	if rv != C.SQLITE_OK {
		return c.lastError()
	}
	c.db = nil

	return nil
}

func (c *Db) lastError() error {
	return lastError(c.db)
}

func (c *Db) lastErrorMsg() string {
	err := lastError(c.db)
	if err == nil {
		return ""
	}
	return err.Error()
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
	c.locker.Lock()
	rpc1_req := C._zxr_query_1_req{}
	rpc1_req.db = c.db
	isInStmt := false
	if req.UseStmtCache {
		stmt, ok := c.stmtCache[req.Query]
		if ok {
			rpc1_req.stmt = stmt
			isInStmt = true
		}
	}
	if isInStmt == false {
		c.updateQueryString(&rpc1_req, req)
	}
	args := req.Args
	if len(args) > 0 {
		rpc1_req.bind_slice_slice = C.uintptr_t(uintptr(unsafe.Pointer(&args)))
	}
	rpc1_req.ColumnCount = C.int(req.ColumnCount)
	rpc1_req.needResult = C.bool(req.RespDataCb != nil)
	rpc1_req.autoCloseStmt = C.bool(req.UseStmtCache == false && req.ColumnsCb == nil)
	rpc1_req.cacher = c.cCacher
	rpc1_resp := C._zxr_query_1(C.uintptr_t(uintptr(unsafe.Pointer(&rpc1_req))))
	if rpc1_resp.rv != C.SQLITE_OK || rpc1_resp.hasTail {
		if req.UseStmtCache && isInStmt == true {
			delete(c.stmtCache, req.Query)
		}
		if rpc1_resp.hasTail {
			c.locker.Unlock()
			return "desfgdu59n"
		}
		errMsg = c.lastErrorMsg()
		c.locker.Unlock()
		return errMsg
	}
	req.ColumnCount = int(rpc1_resp.ColumnCount)
	if req.ColumnsCb != nil {
		resultList := c.getResultListCache(req.ColumnCount)
		for i := 0; i < req.ColumnCount; i++ {
			columnS := C.GoString(C.zxr_column_name(rpc1_resp.stmt, C.int(i)))
			resultList[i] = []byte(columnS)
		}
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
			rpc2_req := C._zxr_query_2_req{}
			rpc2_req.stmt = rpc1_resp.stmt
			rpc2_req.ColumnCount = C.int(req.ColumnCount)
			rpc2_req.autoCloseStmt = rpc1_req.autoCloseStmt
			rpc2_req.cacher = c.cCacher
			rpc2_resp := C._zxr_query_2(rpc2_req)
			if rpc2_resp.rv != C.SQLITE_OK {
				if req.UseStmtCache && isInStmt == true {
					delete(c.stmtCache, req.Query)
				}
				return c.lastError().Error()
			}
			resultList := c.unmarshalResultListFromC__NOLOCK(req.ColumnCount, rpc2_resp.rowP, int(rpc2_resp.rowLen))
			req.RespDataCb(resultList)
			hasMoreRow = rpc2_resp.hasMoreRow
		}
	}
	if req.RespStatusCb != nil {
		req.RespStatusCb(c.getRespStatusCb())
	}
	if req.UseStmtCache && isInStmt == false && rpc1_resp.stmt != nil {
		if len(c.stmtCache) > 100 {
			c.cleanStmtCache()
		}
		queryCopy := make([]byte, len(req.Query))
		copy(queryCopy, req.Query)
		c.stmtCache[string(queryCopy)] = rpc1_resp.stmt
	}
	c.locker.Unlock()
	return ""
}

func (c *Db) updateQueryString(rpc1_req *C._zxr_query_1_req, req QueryReq) {
	var l int
	var p unsafe.Pointer
	var queryBuf []byte

	if req.Query[len(req.Query)-1] == 0 {
		p = unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&req.Query)))
		l = len(req.Query)
	} else {
		c.cacheBuf.Reset()
		c.cacheBuf.WriteString(req.Query)
		c.cacheBuf.WriteByte(0)
		queryBuf = c.cacheBuf.GetBytes()
		p = unsafe.Pointer(&queryBuf[0])
		l = len(queryBuf)
	}
	rpc1_req.query = (*C.char)(p)
	rpc1_req.queryLen = C.int(l)
}

func (c *Db) cleanStmtCache() {
	for key, stmt := range c.stmtCache {
		c.closeStmt(stmt)
		delete(c.stmtCache, key)
	}
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
	AffectedRows := uint64(C.zxr_changes(c.db))
	LastInsertId := uint64(C.zxr_last_insert_rowid(c.db))
	return QueryRespStatus{
		AffectedRows: AffectedRows,
		LastInsertId: LastInsertId,
	}
}

func (c *Db) getResultListCache(needSize int) [][]byte {
	if cap(c.resultListCache) < needSize {
		c.resultListCache = make([][]byte, needSize)
	}
	return c.resultListCache[:needSize]
}

func (c *Db) closeStmt(stmt *C.zxr_stmt) (errMsg string) {
	rv := C.zxr_finalize(stmt)
	if rv != C.SQLITE_OK {
		return c.lastError().Error()
	}
	return ""
}
