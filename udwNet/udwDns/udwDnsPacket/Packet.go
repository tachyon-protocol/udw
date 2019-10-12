package udwDnsPacket

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwBitwise"
	"github.com/tachyon-protocol/udw/udwBytes"
	"github.com/tachyon-protocol/udw/udwDebug"
	"github.com/tachyon-protocol/udw/udwFmt"
	"github.com/tachyon-protocol/udw/udwJson"
	"github.com/tachyon-protocol/udw/udwRand"
	"strconv"
	"strings"
)

type DnsRR struct {
	Name   string
	Type   DnsRRType
	Ttl    uint32
	Result string
}

func (rr DnsRR) GetNormalizeName() string {
	return strings.ToLower(rr.Name)
}

type DnsPacket struct {
	Id                    uint16
	Bits                  uint16
	QuestionName          string
	QuestionType          DnsRRType
	AnswerList            []DnsRR
	HasEdns4096           bool
	IsWriteAnsToAuthority bool

	Ancount           uint16
	Nscount           uint16
	Arcount           uint16
	HasCompressDomain bool
}

func (packet DnsPacket) GetNormalizeQuestionName() string {
	return strings.ToLower(packet.QuestionName)
}

func (packet DnsPacket) IsResponse() bool {
	return (packet.Bits & (1 << 15)) != 0
}
func (packet *DnsPacket) SetIsResponse(is bool) {
	packet.Bits = udwBitwise.Uint16SetBitWithBool(packet.Bits, 15, is)
}
func (packet DnsPacket) GetOpcode() int {
	return int(packet.Bits>>11) & 0xF
}
func (packet DnsPacket) IsAuthoritative() bool {

	return (packet.Bits & (1 << 10)) != 0
}
func (packet DnsPacket) IsTruncated() bool {

	return (packet.Bits & (1 << 9)) != 0
}
func (packet DnsPacket) IsRecursionDesired() bool {

	return (packet.Bits & (1 << 8)) != 0
}
func (packet *DnsPacket) SetIsRecursionDesired(is bool) {
	packet.Bits = udwBitwise.Uint16SetBitWithBool(packet.Bits, 8, is)
}
func (packet DnsPacket) IsRecursionAvailable() bool {

	return udwBitwise.Uint16GetBitToBool(packet.Bits, 7)
}
func (packet *DnsPacket) SetIsRecursionAvailable(is bool) {
	packet.Bits = udwBitwise.Uint16SetBitWithBool(packet.Bits, 7, is)
}
func (packet DnsPacket) IsBitReserverd9() bool {

	return (packet.Bits & (1 << 6)) != 0
}
func (packet DnsPacket) IsAuthenticatedData() bool {

	return (packet.Bits & (1 << 5)) != 0
}
func (packet *DnsPacket) SetIsAuthenticatedData(is bool) {
	packet.Bits = udwBitwise.Uint16SetBitWithBool(packet.Bits, 5, is)
}
func (packet DnsPacket) IsCheckingDisabled() bool {

	return (packet.Bits & (1 << 4)) != 0
}
func (packet DnsPacket) GetRcode() Rcode {
	return Rcode(packet.Bits & 0xF)
}
func (packet *DnsPacket) SetRcode(code Rcode) {
	packet.Bits = udwBitwise.Uint16SetPartWithMask(packet.Bits, uint16(code), 0xf)
}
func (packet *DnsPacket) HasPlainResult(result string) bool {
	if len(packet.AnswerList) == 0 {
		return false
	}
	for _, anwser := range packet.AnswerList {
		switch anwser.Type {
		case DnsTypeA, DnsTypeAAAA, DnsTypeNS, DnsTypeCNAME:
			if anwser.Result == result {
				return true
			}
		case DnsTypeTXT:
			outList := anwser.GetTxt()
			if len(outList) == 1 && outList[0] == result {
				return true
			}
		case DnsTypeMX:
			mxObj := anwser.GetMx()
			if mxObj.Mx == result {
				return true
			}
		default:
			fmt.Println("ynzdyf2z76", anwser.Type.String(), result)
		}
	}
	return false
}

func (packet *DnsPacket) ShuffleAnswer() {
	if len(packet.AnswerList) <= 1 {
		return
	}
	for i := range packet.AnswerList {
		j := udwRand.Intn(len(packet.AnswerList))
		packet.AnswerList[i], packet.AnswerList[j] = packet.AnswerList[j], packet.AnswerList[i]
	}
	return
}

func (rr DnsRR) GetSoa() (out DnsRRSoa) {
	udwJson.MustUnmarshal([]byte(rr.Result), &out)
	return out
}

func (rr *DnsRR) SetSoa(in DnsRRSoa) {
	rr.Result = udwJson.MustMarshalToString(in)
}

func (rr DnsRR) GetMx() (out DnsRRMX) {
	udwJson.MustUnmarshal([]byte(rr.Result), &out)
	return out
}

func (rr *DnsRR) SetMx(in DnsRRMX) {
	rr.Result = udwJson.MustMarshalToString(in)
}

func (rr DnsRR) GetTxt() (out []string) {
	udwJson.MustUnmarshal([]byte(rr.Result), &out)
	return out
}

func (rr *DnsRR) SetTxt(in []string) {
	rr.Result = udwJson.MustMarshalToString(in)
}

func (rr DnsRR) GetSrv() (out DnsRRSrv) {
	udwJson.MustUnmarshal([]byte(rr.Result), &out)
	return out
}

func (rr *DnsRR) SetSrv(in DnsRRSrv) {
	rr.Result = udwJson.MustMarshalToString(in)
}

type Rcode uint8

const (
	RcodeNOERROR  Rcode = 0
	RcodeFORMERR  Rcode = 1
	RcodeSERVFAIL Rcode = 2
	RcodeNXDOMAIN Rcode = 3
	RcodeNOTIMPL  Rcode = 4
	RcodeREFUSED  Rcode = 5
)

func (rcode Rcode) String() string {
	switch rcode {
	case RcodeNOERROR:
		return "NOERROR"
	case RcodeFORMERR:
		return "FORMERR"
	case RcodeSERVFAIL:
		return "SERVFAIL"
	case RcodeNXDOMAIN:
		return "NXDOMAIN"
	case RcodeNOTIMPL:
		return "NOTIMPL"
	case RcodeREFUSED:
		return "REFUSED"
	default:
		return udwFmt.Sprint("Rcode:", int(rcode))
	}
}

func (rcode Rcode) IsOkStatus() bool {
	return rcode == RcodeNXDOMAIN || rcode == RcodeNOERROR
}

type DnsRRSoa struct {
	Ns      string
	Mbox    string
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	Minttl  uint32
}

type DnsRRMX struct {
	Preference uint16
	Mx         string
}

type DnsRRSrv struct {
	Priority uint16
	Weight   uint16
	Port     uint16
	Target   string
}

func (packet DnsPacket) String() (out string) {
	_buf := udwBytes.BufWriter{}
	_buf.WriteString(udwDebug.SepLineAndOneLBString)
	_buf.WriteString("question: ")
	_buf.WriteString(packet.QuestionName)
	_buf.WriteString(" type: ")
	_buf.WriteString(packet.QuestionType.String())
	_buf.WriteString(" rcode: ")
	_buf.WriteString(packet.GetRcode().String())
	_buf.WriteString(" opcode: ")
	_buf.WriteString(strconv.Itoa(packet.GetOpcode()))

	if packet.IsTruncated() {
		_buf.WriteString(" IsTruncated")

	}
	if packet.IsResponse() {
		_buf.WriteString(" IsResponse")

	} else {
		_buf.WriteString(" IsRequest")

	}
	if packet.IsAuthenticatedData() {
		_buf.WriteString(" IsAuthenticatedData")

	}
	if packet.IsAuthoritative() {
		_buf.WriteString(" IsAuthoritative")

	}
	if packet.IsCheckingDisabled() {
		_buf.WriteString(" IsCheckingDisabled")

	}
	if packet.IsRecursionAvailable() {
		_buf.WriteString(" IsRecursionAvailable")

	}
	if packet.IsRecursionDesired() {
		_buf.WriteString(" IsRecursionDesired")

	}
	if packet.IsBitReserverd9() {
		_buf.WriteString(" IsBitReserverd9")
	}
	if packet.IsAuthenticatedData() {
		_buf.WriteString(" IsAuthenticatedData")
	}
	if packet.IsCheckingDisabled() {
		_buf.WriteString(" IsCheckingDisabled")
	}
	if packet.HasEdns4096 {
		_buf.WriteString(" HasEdns4096")
	}
	if packet.HasCompressDomain {
		_buf.WriteString(" HasCompressDomain")
	}
	if packet.Ancount > 0 {
		_buf.WriteString(" An:")
		_buf.WriteString(strconv.Itoa(int(packet.Ancount)))
	}
	if packet.Nscount > 0 {
		_buf.WriteString(" Ns:")
		_buf.WriteString(strconv.Itoa(int(packet.Nscount)))
	}
	if packet.Arcount > 0 {
		_buf.WriteString(" Ar:")
		_buf.WriteString(strconv.Itoa(int(packet.Arcount)))
	}
	_buf.WriteString(" Id:")
	_buf.WriteString(strconv.Itoa(int(packet.Id)))

	_buf.WriteString("\n")

	for _, rr := range packet.AnswerList {
		_buf.WriteString("rr: " + rr.Name + " " + rr.Type.String() + " " + strconv.Itoa(int(rr.Ttl)) + " " + rr.Result + "\n")

	}
	_buf.WriteString(udwDebug.SepLineAndOneLBString)
	return _buf.GetString()
}
