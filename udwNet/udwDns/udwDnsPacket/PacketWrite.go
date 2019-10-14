package udwDnsPacket

import (
	"encoding/binary"
	"github.com/tachyon-protocol/udw/udwFmt"
	"github.com/tachyon-protocol/udw/udwRand"
	"net"
	"strings"
)

func (packet DnsPacket) GetBinaryLen() int {
	outSize := 12
	outSize += getDomainStringBianryLen(packet.QuestionName) + 4
	for i := range packet.AnswerList {
		outSize += getDomainStringBianryLen(packet.AnswerList[i].Name) + 10
		switch packet.AnswerList[i].Type {
		case DnsTypeA:
			outSize += 4
		case DnsTypeAAAA:
			outSize += 16
		case DnsTypeNS, DnsTypeCNAME, DnsTypePTR, DnsTypeMX:
			outSize += getDomainStringBianryLen(packet.AnswerList[i].Result)
		case DnsTypeSOA:
			soa := packet.AnswerList[i].GetSoa()
			outSize += getSoaBinaryLen(soa)
		case DnsTypeTXT:
			txt := packet.AnswerList[i].GetTxt()
			outSize += getTxtBinaryLen(txt)
		case DnsTypeSRV:
			txt := packet.AnswerList[i].GetSrv()
			outSize += getSrvBinaryLen(txt)
		default:

		}
	}
	if packet.HasEdns4096 {
		outSize += gEdns4096Len
	}
	return outSize
}

func getDomainStringBianryLen(domain string) int {
	return len(domain) + 2
}

func writeDomainToDnsPacket(domain string, buf []byte, pos int) (outPos int, errMsg string) {
	if len(domain) > 255 {
		return pos, udwFmt.Sprint(`[writeDomainToDnsPacket] len(domain)>255`, len(domain), domain)
	}
	sPos := 0
	for {

		dotPos := strings.IndexByte(domain[sPos:], '.')
		if dotPos == -1 {
			break
		}
		if dotPos > 63 {
			return pos, udwFmt.Sprint("[writeDomainToDnsPacket] dotPos>63", dotPos, domain)
		}
		buf[pos] = byte(dotPos)
		pos++
		copy(buf[pos:], domain[sPos:sPos+dotPos])
		pos += dotPos
		sPos += dotPos + 1
	}
	if sPos < len(domain) {
		dotPos := len(domain) - sPos

		buf[pos] = byte(dotPos)
		pos++
		copy(buf[pos:], domain[sPos:sPos+dotPos])
		pos += dotPos
	}

	buf[pos] = 0
	pos++
	return pos, ""
}

func writeTxtStringToDnsPacket(domain string, buf []byte, pos int) (outPos int, errMsg string) {
	if len(domain) > 255 {
		return pos, udwFmt.Sprint(`[writeTxtStringToDnsPacket] len(domain)>255`, len(domain), domain)
	}
	buf[pos] = byte(len(domain))
	pos++
	copy(buf[pos:pos+len(domain)], domain)
	pos += len(domain)
	return pos, ""
}

func getSoaBinaryLen(soa DnsRRSoa) int {
	return getDomainStringBianryLen(soa.Ns) + getDomainStringBianryLen(soa.Mbox) + 20
}

func getMxBinaryLen(soa DnsRRMX) int {
	return getDomainStringBianryLen(soa.Mx) + 2
}

func getSrvBinaryLen(soa DnsRRSrv) int {
	return getDomainStringBianryLen(soa.Target) + 6
}

func getTxtBinaryLen(in []string) int {
	outSize := 0
	for _, s := range in {
		outSize += 1 + len(s)
	}
	return outSize
}

func PacketWriteToByteSlice(packet DnsPacket, buf []byte) (outBuf []byte, errMsg string) {

	if packet.QuestionName == "" {

		return nil, `[WriteDnsPacketToBuf] packet.QuestionName==""`
	}
	packet.QuestionName = NormalizeDomain(packet.QuestionName)
	for i := range packet.AnswerList {
		packet.AnswerList[i].Name = NormalizeDomain(packet.AnswerList[i].Name)
	}
	needSize := packet.GetBinaryLen()
	if len(buf) < needSize {
		buf = make([]byte, needSize)
	}
	binary.BigEndian.PutUint16(buf[:2], packet.Id)
	binary.BigEndian.PutUint16(buf[2:4], packet.Bits)
	binary.BigEndian.PutUint16(buf[4:6], 1)
	if packet.IsWriteAnsToAuthority == false {

		binary.BigEndian.PutUint16(buf[6:8], uint16(len(packet.AnswerList)))
		buf[8] = 0
		buf[9] = 0
	} else {

		buf[6] = 0
		buf[7] = 0
		binary.BigEndian.PutUint16(buf[8:10], uint16(len(packet.AnswerList)))
	}

	additionalRRNumber := uint16(0)
	if packet.HasEdns4096 {
		additionalRRNumber = 1
	}
	binary.BigEndian.PutUint16(buf[10:12], additionalRRNumber)

	pos := 12
	pos, errMsg = writeDomainToDnsPacket(packet.QuestionName, buf, pos)
	if errMsg != "" {
		return nil, errMsg
	}
	binary.BigEndian.PutUint16(buf[pos:pos+2], uint16(packet.QuestionType))
	pos += 2
	buf[pos] = 0
	buf[pos+1] = 1
	pos += 2
	for i := range packet.AnswerList {
		pos, errMsg = writeDomainToDnsPacket(packet.AnswerList[i].Name, buf, pos)
		if errMsg != "" {
			return nil, errMsg
		}
		binary.BigEndian.PutUint16(buf[pos:pos+2], uint16(packet.AnswerList[i].Type))
		pos += 2
		buf[pos] = 0
		buf[pos+1] = 1
		pos += 2
		binary.BigEndian.PutUint32(buf[pos:pos+4], packet.AnswerList[i].Ttl)
		pos += 4
		switch packet.AnswerList[i].Type {
		case DnsTypeA:
			buf[pos] = 0
			buf[pos+1] = 4
			pos += 2
			ip := net.ParseIP(packet.AnswerList[i].Result)
			if ip == nil {
				return nil, `[WriteDnsPacketToBuf] DnsTypeA ip==nil`
			}
			ipv4 := ip.To4()
			if ipv4 == nil {
				return nil, `[WriteDnsPacketToBuf] DnsTypeA ipv4==nil`
			}
			copy(buf[pos:pos+4], []byte(ipv4))
			pos += 4
		case DnsTypeNS, DnsTypeCNAME, DnsTypePTR:
			domain := NormalizeDomain(packet.AnswerList[i].Result)
			binary.BigEndian.PutUint16(buf[pos:pos+2], uint16(getDomainStringBianryLen(domain)))
			pos += 2
			pos, errMsg = writeDomainToDnsPacket(domain, buf, pos)
			if errMsg != "" {
				return nil, errMsg
			}
		case DnsTypeSOA:
			soa := packet.AnswerList[i].GetSoa()
			binary.BigEndian.PutUint16(buf[pos:pos+2], uint16(getSoaBinaryLen(soa)))
			pos += 2
			pos, errMsg = writeDomainToDnsPacket(soa.Ns, buf, pos)
			if errMsg != "" {
				return nil, errMsg
			}
			pos, errMsg = writeDomainToDnsPacket(soa.Mbox, buf, pos)
			if errMsg != "" {
				return nil, errMsg
			}
			binary.BigEndian.PutUint32(buf[pos:pos+4], soa.Serial)
			pos += 4
			binary.BigEndian.PutUint32(buf[pos:pos+4], soa.Refresh)
			pos += 4
			binary.BigEndian.PutUint32(buf[pos:pos+4], soa.Retry)
			pos += 4
			binary.BigEndian.PutUint32(buf[pos:pos+4], soa.Expire)
			pos += 4
			binary.BigEndian.PutUint32(buf[pos:pos+4], soa.Minttl)
			pos += 4
		case DnsTypeMX:
			mx := packet.AnswerList[i].GetMx()
			binary.BigEndian.PutUint16(buf[pos:pos+2], uint16(getMxBinaryLen(mx)))
			pos += 2
			binary.BigEndian.PutUint16(buf[pos:pos+2], mx.Preference)
			pos += 2
			pos, errMsg = writeDomainToDnsPacket(mx.Mx, buf, pos)
			if errMsg != "" {
				return nil, errMsg
			}
		case DnsTypeTXT:
			txt := packet.AnswerList[i].GetTxt()
			binary.BigEndian.PutUint16(buf[pos:pos+2], uint16(getTxtBinaryLen(txt)))
			pos += 2
			for _, s := range txt {
				pos, errMsg = writeTxtStringToDnsPacket(s, buf, pos)
				if errMsg != "" {
					return nil, errMsg
				}
			}
		case DnsTypeSRV:
			srv := packet.AnswerList[i].GetSrv()
			binary.BigEndian.PutUint16(buf[pos:pos+2], uint16(getSrvBinaryLen(srv)))
			pos += 2
			binary.BigEndian.PutUint16(buf[pos:pos+2], srv.Priority)
			pos += 2
			binary.BigEndian.PutUint16(buf[pos:pos+2], srv.Weight)
			pos += 2
			binary.BigEndian.PutUint16(buf[pos:pos+2], srv.Port)
			pos += 2
			pos, errMsg = writeDomainToDnsPacket(srv.Target, buf, pos)
			if errMsg != "" {
				return nil, errMsg
			}
		case DnsTypeAAAA:
			buf[pos] = 0
			buf[pos+1] = 16
			pos += 2
			ip := net.ParseIP(packet.AnswerList[i].Result)
			if ip == nil {
				return nil, `[WriteDnsPacketToBuf] DnsTypeA ip==nil`
			}
			ipv6 := ip.To16()
			if ipv6 == nil {
				return nil, `[WriteDnsPacketToBuf] DnsTypeA ipv6==nil`
			}
			copy(buf[pos:pos+16], []byte(ipv6))
			pos += 16
		default:
			return nil, udwFmt.Sprint("[PacketWriteToByteSlice] unknow type", packet.AnswerList[i].Type)
		}
	}
	if packet.HasEdns4096 {
		copy(buf[pos:pos+gEdns4096Len], gEdns4096Content)
		pos += gEdns4096Len
	}
	return buf[:pos], ""
}

var gEdns4096Content = []byte{0x00, 0x00,
	0x29, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00}

const gEdns4096Len = 11

func NewDnsPacketForRequest(domain string, typ DnsRRType) DnsPacket {
	inMsg := DnsPacket{}
	inMsg.Id = udwRand.MustCryptoRandUint16()
	inMsg.SetIsRecursionDesired(true)
	inMsg.HasEdns4096 = true
	inMsg.QuestionName = domain
	inMsg.QuestionType = typ
	return inMsg
}
