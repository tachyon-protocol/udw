package udwDbI

type IDb interface {
	Close()
	MustEmptyK1(k1 string)
	IMustGetRangeCallback(req GetRangeReq, cb func(k string, v string))
	INewMsb(threadNum int) IMsb

	MustDeleteK1(k1 string)
	MustSet(k1 string, k2 string, v string)
	MustGet(k1 string, k2 string) string
	IMustDelete(k1 string, k2 string)
}

type IMsb interface {
	Set(k1 string, k2 string, v string)

	MustClose()
}

type GetRangeReq struct {
	K1           string `json:",omitempty"`
	IsDescOrder  bool   `json:",omitempty"`
	DisableOrder bool   `json:",omitempty"`

	MinValue           string `json:",omitempty"`
	MaxValue           string `json:",omitempty"`
	MinValueNotInclude string `json:",omitempty"`
	MaxValueNotInclude string `json:",omitempty"`
	Prefix             string `json:",omitempty"`
	Limit              int    `json:",omitempty"`
}
