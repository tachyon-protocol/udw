package udwDnsPacket

import (
	"encoding/binary"
	"github.com/tachyon-protocol/udw/udwFmt"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoWriter/udwGoTypeMarshal"
	"net"
	"strconv"
	"strings"
)

func PacketReadFromByteSlice(buf []byte) (packet DnsPacket, errMsg string) {
	if len(buf) < 12 {
		return packet, "[NewDnsPacketFromBuf] len(buf)<12 " + strconv.Itoa(len(buf))
	}
	packet.Id = binary.BigEndian.Uint16(buf[:2])
	packet.Bits = binary.BigEndian.Uint16(buf[2:4])
	Qdcount := binary.BigEndian.Uint16(buf[4:6])
	packet.Ancount = binary.BigEndian.Uint16(buf[6:8])

	packet.Nscount = binary.BigEndian.Uint16(buf[8:10])

	packet.Arcount = binary.BigEndian.Uint16(buf[10:12])

	rrCount := int(packet.Ancount + packet.Nscount + packet.Arcount)
	pos := 12

	addCompressFn := func(domainResp readDomainFromDnsPacketResp) {
		if domainResp.hasCompress {
			packet.HasCompressDomain = true
		}
	}
	if Qdcount > 1 {
		return packet, "[NewDnsPacketFromBuf] Qdcount>1 " + strconv.Itoa(int(Qdcount))
	} else if Qdcount == 0 {
		return packet, udwFmt.Sprint("[NewDnsPacketFromBuf] Qdcount==0")
	} else if Qdcount == 1 {
		domainResp := readDomainFromDnsPacket(buf, pos)
		if domainResp.errMsg != "" {
			return packet, domainResp.errMsg
		}
		addCompressFn(domainResp)
		packet.QuestionName = domainResp.domain
		pos = domainResp.outPos
		if len(buf) < pos+4 {
			return packet, udwFmt.Sprint("[NewDnsPacketFromBuf]1 len(buf)<pos+4", len(buf), pos)
		}
		packet.QuestionType = DnsRRType(binary.BigEndian.Uint16(buf[pos : pos+2]))

		pos += 4
	}
	if rrCount > 0 {
		packet.AnswerList = make([]DnsRR, rrCount)
	}
	for i := 0; i < rrCount; i++ {

		domainResp := readDomainFromDnsPacket(buf, pos)
		if domainResp.errMsg != "" {
			return packet, domainResp.errMsg
		}
		addCompressFn(domainResp)
		packet.AnswerList[i].Name = domainResp.domain
		pos = domainResp.outPos
		if len(buf) < pos+10 {
			return packet, udwFmt.Sprint("[NewDnsPacketFromBuf]2 len(buf)<pos+10", len(buf), pos)
		}

		packet.AnswerList[i].Type = DnsRRType(binary.BigEndian.Uint16(buf[pos : pos+2]))
		pos += 2

		qclass := binary.BigEndian.Uint16(buf[pos : pos+2])
		pos += 2
		packet.AnswerList[i].Ttl = binary.BigEndian.Uint32(buf[pos : pos+4])
		pos += 4
		dataLen := int(binary.BigEndian.Uint16(buf[pos : pos+2]))
		pos += 2
		if len(buf) < pos+dataLen {
			return packet, udwFmt.Sprint("[NewDnsPacketFromBuf] llen(buf)<pos+dataLen", len(buf), pos, dataLen)
		}
		switch packet.AnswerList[i].Type {
		case DnsTypeA:
			if dataLen != 4 {
				return packet, udwFmt.Sprint("[NewDnsPacketFromBuf] DnsTypeA dataLen!=4", dataLen)
			}
			ip := net.IP(buf[pos : pos+4])
			pos += 4
			packet.AnswerList[i].Result = ip.String()
		case DnsTypeNS, DnsTypeCNAME, DnsTypePTR, DnsTypeDNAME:
			domainResp := readDomainFromDnsPacket(buf, pos)
			if domainResp.errMsg != "" {
				return packet, domainResp.errMsg
			}
			addCompressFn(domainResp)
			pos = domainResp.outPos
			packet.AnswerList[i].Result = domainResp.domain
		case DnsTypeSOA:
			if dataLen < 1+1+5*4 {
				return packet, "[NewDnsPacketFromBuf] 7f679f49nb dataLen<22"
			}
			startPos := pos
			soa := DnsRRSoa{}
			domainResp := readDomainFromDnsPacket(buf, pos)
			if domainResp.errMsg != "" {
				return packet, domainResp.errMsg
			}
			addCompressFn(domainResp)
			pos = domainResp.outPos
			soa.Ns = domainResp.domain
			domainResp = readDomainFromDnsPacket(buf, pos)
			if domainResp.errMsg != "" {
				return packet, domainResp.errMsg
			}
			addCompressFn(domainResp)
			pos = domainResp.outPos
			soa.Mbox = domainResp.domain
			if startPos+dataLen-pos < 5*4 {
				return packet, "[NewDnsPacketFromBuf] qz5ppu9wyt startPos + dataLen - pos<20"
			}
			soa.Serial = binary.BigEndian.Uint32(buf[pos : pos+4])
			pos += 4
			soa.Refresh = binary.BigEndian.Uint32(buf[pos : pos+4])
			pos += 4
			soa.Retry = binary.BigEndian.Uint32(buf[pos : pos+4])
			pos += 4
			soa.Expire = binary.BigEndian.Uint32(buf[pos : pos+4])
			pos += 4
			soa.Minttl = binary.BigEndian.Uint32(buf[pos : pos+4])
			pos += 4
			if startPos+dataLen != pos {
				return packet, udwFmt.Sprint("[NewDnsPacketFromBuf] DnsTypeMX startPos + dataLen != pos", pos, startPos, dataLen)
			}
			packet.AnswerList[i].SetSoa(soa)
		case DnsTypeMX:
			if dataLen < 1+2 {
				return packet, "[NewDnsPacketFromBuf] tvb895tx6a dataLen<3"
			}
			mx := DnsRRMX{}
			startPos := pos
			mx.Preference = binary.BigEndian.Uint16(buf[pos : pos+2])
			pos += 2
			domainResp := readDomainFromDnsPacket(buf, pos)
			if domainResp.errMsg != "" {
				return packet, domainResp.errMsg
			}
			addCompressFn(domainResp)
			mx.Mx = domainResp.domain
			pos = domainResp.outPos

			if startPos+dataLen != pos {
				return packet, udwFmt.Sprint("[NewDnsPacketFromBuf] DnsTypeMX startPos + dataLen != pos", pos, startPos, dataLen)
			}
			packet.AnswerList[i].SetMx(mx)
		case DnsTypeTXT:
			txtList := []string{}
			startPos := pos
			for {
				if pos-startPos >= dataLen {
					break
				}
				var thisTxt string
				thisTxt, pos, errMsg = readTxtStringFromDnsPacket(buf, pos)
				if errMsg != "" {
					return packet, errMsg
				}
				txtList = append(txtList, thisTxt)
			}
			packet.AnswerList[i].SetTxt(txtList)
		case DnsTypeSRV:
			if dataLen < 7 {
				return packet, "[NewDnsPacketFromBuf] dataLen<7"
			}
			srv := DnsRRSrv{}
			startPos := pos
			srv.Priority = binary.BigEndian.Uint16(buf[pos : pos+2])
			pos += 2
			srv.Weight = binary.BigEndian.Uint16(buf[pos : pos+2])
			pos += 2
			srv.Port = binary.BigEndian.Uint16(buf[pos : pos+2])
			pos += 2
			domainResp := readDomainFromDnsPacket(buf, pos)
			if domainResp.errMsg != "" {
				return packet, domainResp.errMsg
			}
			addCompressFn(domainResp)
			pos = domainResp.outPos
			srv.Target = domainResp.domain
			if startPos+dataLen != pos {
				return packet, udwFmt.Sprint("[NewDnsPacketFromBuf] DnsTypeSRV startPos + dataLen != pos", pos, startPos, dataLen)
			}
			packet.AnswerList[i].SetSrv(srv)
		case DnsTypeAAAA:
			if dataLen != 16 {
				return packet, udwFmt.Sprint("[NewDnsPacketFromBuf] DnsTypeAAAA dataLen!=16", dataLen)
			}
			ip := net.IP(buf[pos : pos+16])
			pos += 16
			packet.AnswerList[i].Result = ip.String()
		case DnsTypeOPT:
			pos += dataLen

			if qclass == 4096 && dataLen == 0 && packet.AnswerList[i].Ttl == 0 && packet.AnswerList[i].Name == "" {
				packet.HasEdns4096 = true
			}
			i--
			rrCount--
			packet.AnswerList = packet.AnswerList[:len(packet.AnswerList)-1]
		default:
			if debugLog {
				return packet, udwFmt.Sprint("[NewDnsPacketFromBuf] unknow type", packet.AnswerList[i].Type, pos, udwGoTypeMarshal.WriteByteSlice(buf))
			} else {
				return packet, udwFmt.Sprint("[NewDnsPacketFromBuf] unknow type", packet.AnswerList[i].Type, pos)
			}
		}
	}
	return packet, ""
}

func readDomainFromDnsPacketL1(buf []byte, pos int, remainJmpNum int) (resp readDomainFromDnsPacketResp) {
	if remainJmpNum == 0 {
		resp.errMsg = "[readDomainFromDnsPacket] too many compression pointers"
		return resp
	}
	domainPartList := make([]string, 0, 4)
	for {
		if len(buf) < pos+1 {
			resp.errMsg = udwFmt.Sprint("[readDomainFromDnsPacket]1 len(buf)<pos+1", len(buf), pos)
			return resp
		}
		nextPartLen := buf[pos]
		pos++
		if nextPartLen == 0 {
			break
		}
		if nextPartLen&0xc0 == 0xc0 {
			if len(buf) < pos+1 {
				resp.errMsg = udwFmt.Sprint("[readDomainFromDnsPacket]2 len(buf)<pos+1", len(buf), pos)
				return resp
			}
			jumpOffset := int((nextPartLen^0xc0))<<8 | int(buf[pos])
			pos++
			thisResp := readDomainFromDnsPacketL1(buf, jumpOffset, remainJmpNum-1)
			if thisResp.errMsg != "" {
				return thisResp
			}
			domainPartList = append(domainPartList, string(thisResp.domain))
			resp.hasCompress = true

			break
		}
		if nextPartLen >= 64 {

			if debugLog {

				resp.errMsg = udwFmt.Sprint("[readDomainFromDnsPacket] nextPartLen>=64", nextPartLen, pos, udwGoTypeMarshal.WriteByteSlice(buf))
			} else {
				resp.errMsg = udwFmt.Sprint("[readDomainFromDnsPacket] nextPartLen>=64", nextPartLen, pos)
			}
			return resp
		}
		if len(buf) < (pos + int(nextPartLen)) {
			resp.errMsg = udwFmt.Sprint("[readDomainFromDnsPacket] len(buf)<(pos+nextPartLen)", len(buf), pos, nextPartLen)
			return resp
		}
		domainPartList = append(domainPartList, string(buf[pos:pos+int(nextPartLen)]))
		pos += int(nextPartLen)
	}
	domain := strings.Join(domainPartList, ".")
	resp.domain = domain
	resp.outPos = pos
	return resp
}

type readDomainFromDnsPacketResp struct {
	domain      string
	outPos      int
	errMsg      string
	hasCompress bool
}

func readDomainFromDnsPacket(buf []byte, pos int) (resp readDomainFromDnsPacketResp) {
	return readDomainFromDnsPacketL1(buf, pos, 5)
}

func readTxtStringFromDnsPacket(buf []byte, pos int) (domain string, outPos int, errMsg string) {
	if len(buf) < pos+1 {
		return "", pos, udwFmt.Sprint("[readTxtStringFromDnsPacket] len(buf)<pos+1", len(buf), pos)
	}
	nextPartLen := buf[pos]
	pos++
	if len(buf) < (pos + int(nextPartLen)) {
		return "", pos, udwFmt.Sprint("[readTxtStringFromDnsPacket] len(buf)<(pos+nextPartLen)", len(buf), pos, nextPartLen)
	}
	domain = string(buf[pos : pos+int(nextPartLen)])
	pos += int(nextPartLen)
	return domain, pos, ""
}
