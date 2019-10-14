package udwDnsPacket

import (
	"strconv"
	"strings"
)

type DnsRRType uint16

const (
	DnsPort = 53

	DnsTypeNone DnsRRType = 0
	DnsTypeA    DnsRRType = 1
	DnsTypeNS   DnsRRType = 2

	DnsTypeCNAME DnsRRType = 5
	DnsTypeSOA   DnsRRType = 6
	DnsTypePTR   DnsRRType = 12
	DnsTypeMX    DnsRRType = 15
	DnsTypeTXT   DnsRRType = 16

	DnsTypeAAAA DnsRRType = 28
	DnsTypeSRV  DnsRRType = 33

	DnsTypeMD         DnsRRType = 3
	DnsTypeMF         DnsRRType = 4
	DnsTypeMB         DnsRRType = 7
	DnsTypeMG         DnsRRType = 8
	DnsTypeMR         DnsRRType = 9
	DnsTypeNULL       DnsRRType = 10
	DnsTypeWKS        DnsRRType = 11
	DnsTypeHINFO      DnsRRType = 13
	DnsTypeMINFO      DnsRRType = 14
	DnsTypeRP         DnsRRType = 17
	DnsTypeAFSDB      DnsRRType = 18
	DnsTypeX25        DnsRRType = 19
	DnsTypeISDN       DnsRRType = 20
	DnsTypeRT         DnsRRType = 21
	DnsTypeNSAPPTR    DnsRRType = 23
	DnsTypeSIG        DnsRRType = 24
	DnsTypeKEY        DnsRRType = 25
	DnsTypePX         DnsRRType = 26
	DnsTypeGPOS       DnsRRType = 27
	DnsTypeLOC        DnsRRType = 29
	DnsTypeNXT        DnsRRType = 30
	DnsTypeEID        DnsRRType = 31
	DnsTypeNIMLOC     DnsRRType = 32
	DnsTypeATMA       DnsRRType = 34
	DnsTypeNAPTR      DnsRRType = 35
	DnsTypeKX         DnsRRType = 36
	DnsTypeCERT       DnsRRType = 37
	DnsTypeDNAME      DnsRRType = 39
	DnsTypeOPT        DnsRRType = 41
	DnsTypeDS         DnsRRType = 43
	DnsTypeSSHFP      DnsRRType = 44
	DnsTypeIPSECKEY   DnsRRType = 45
	DnsTypeRRSIG      DnsRRType = 46
	DnsTypeNSEC       DnsRRType = 47
	DnsTypeDNSKEY     DnsRRType = 48
	DnsTypeDHCID      DnsRRType = 49
	DnsTypeNSEC3      DnsRRType = 50
	DnsTypeNSEC3PARAM DnsRRType = 51
	DnsTypeTLSA       DnsRRType = 52
	DnsTypeHIP        DnsRRType = 55
	DnsTypeNINFO      DnsRRType = 56
	DnsTypeRKEY       DnsRRType = 57
	DnsTypeTALINK     DnsRRType = 58
	DnsTypeCDS        DnsRRType = 59
	DnsTypeCDNSKEY    DnsRRType = 60
	DnsTypeOPENPGPKEY DnsRRType = 61
	DnsTypeSPF        DnsRRType = 99
	DnsTypeUINFO      DnsRRType = 100
	DnsTypeUID        DnsRRType = 101
	DnsTypeGID        DnsRRType = 102
	DnsTypeUNSPEC     DnsRRType = 103
	DnsTypeNID        DnsRRType = 104
	DnsTypeL32        DnsRRType = 105
	DnsTypeL64        DnsRRType = 106
	DnsTypeLP         DnsRRType = 107
	DnsTypeEUI48      DnsRRType = 108
	DnsTypeEUI64      DnsRRType = 109
	DnsTypeTKEY       DnsRRType = 249
	DnsTypeTSIG       DnsRRType = 250
	DnsTypeIXFR       DnsRRType = 251
	DnsTypeAXFR       DnsRRType = 252
	DnsTypeMAILB      DnsRRType = 253
	DnsTypeMAILA      DnsRRType = 254
	DnsTypeANY        DnsRRType = 255
	DnsTypeURI        DnsRRType = 256
	DnsTypeCAA        DnsRRType = 257
	DnsTypeTA         DnsRRType = 32768
	DnsTypeDLV        DnsRRType = 32769
)

func (typ DnsRRType) String() string {
	switch typ {
	case DnsTypeNone:
		return `None`
	case DnsTypeA:
		return `A`
	case DnsTypeNS:
		return `NS`
	case DnsTypeMD:
		return `MD`
	case DnsTypeMF:
		return `MF`
	case DnsTypeCNAME:
		return `CNAME`
	case DnsTypeSOA:
		return `SOA`
	case DnsTypeMB:
		return `MB`
	case DnsTypeMG:
		return `MG`
	case DnsTypeMR:
		return `MR`
	case DnsTypeNULL:
		return `NULL`
	case DnsTypeWKS:
		return `WKS`
	case DnsTypePTR:
		return `PTR`
	case DnsTypeHINFO:
		return `HINFO`
	case DnsTypeMINFO:
		return `MINFO`
	case DnsTypeMX:
		return `MX`
	case DnsTypeTXT:
		return `TXT`
	case DnsTypeRP:
		return `RP`
	case DnsTypeAFSDB:
		return `AFSDB`
	case DnsTypeX25:
		return `X25`
	case DnsTypeISDN:
		return `ISDN`
	case DnsTypeRT:
		return `RT`
	case DnsTypeNSAPPTR:
		return `NSAPPTR`
	case DnsTypeSIG:
		return `SIG`
	case DnsTypeKEY:
		return `KEY`
	case DnsTypePX:
		return `PX`
	case DnsTypeGPOS:
		return `GPOS`
	case DnsTypeAAAA:
		return `AAAA`
	case DnsTypeLOC:
		return `LOC`
	case DnsTypeNXT:
		return `NXT`
	case DnsTypeEID:
		return `EID`
	case DnsTypeNIMLOC:
		return `NIMLOC`
	case DnsTypeSRV:
		return `SRV`
	case DnsTypeATMA:
		return `ATMA`
	case DnsTypeNAPTR:
		return `NAPTR`
	case DnsTypeKX:
		return `KX`
	case DnsTypeCERT:
		return `CERT`
	case DnsTypeDNAME:
		return `DNAME`
	case DnsTypeOPT:
		return `OPT`
	case DnsTypeDS:
		return `DS`
	case DnsTypeSSHFP:
		return `SSHFP`
	case DnsTypeIPSECKEY:
		return `IPSECKEY`
	case DnsTypeRRSIG:
		return `RRSIG`
	case DnsTypeNSEC:
		return `NSEC`
	case DnsTypeDNSKEY:
		return `DNSKEY`
	case DnsTypeDHCID:
		return `DHCID`
	case DnsTypeNSEC3:
		return `NSEC3`
	case DnsTypeNSEC3PARAM:
		return `NSEC3PARAM`
	case DnsTypeTLSA:
		return `TLSA`
	case DnsTypeHIP:
		return `HIP`
	case DnsTypeNINFO:
		return `NINFO`
	case DnsTypeRKEY:
		return `RKEY`
	case DnsTypeTALINK:
		return `TALINK`
	case DnsTypeCDS:
		return `CDS`
	case DnsTypeCDNSKEY:
		return `CDNSKEY`
	case DnsTypeOPENPGPKEY:
		return `OPENPGPKEY`
	case DnsTypeSPF:
		return `SPF`
	case DnsTypeUINFO:
		return `UINFO`
	case DnsTypeUID:
		return `UID`
	case DnsTypeGID:
		return `GID`
	case DnsTypeUNSPEC:
		return `UNSPEC`
	case DnsTypeNID:
		return `NID`
	case DnsTypeL32:
		return `L32`
	case DnsTypeL64:
		return `L64`
	case DnsTypeLP:
		return `LP`
	case DnsTypeEUI48:
		return `EUI48`
	case DnsTypeEUI64:
		return `EUI64`
	case DnsTypeTKEY:
		return `TKEY`
	case DnsTypeTSIG:
		return `TSIG`
	case DnsTypeIXFR:
		return `IXFR`
	case DnsTypeAXFR:
		return `AXFR`
	case DnsTypeMAILB:
		return `MAILB`
	case DnsTypeMAILA:
		return `MAILA`
	case DnsTypeANY:
		return `ANY`
	case DnsTypeURI:
		return `URI`
	case DnsTypeCAA:
		return `CAA`
	case DnsTypeTA:
		return `TA`
	case DnsTypeDLV:
		return `DLV`
	default:
		return "DnsRRType: " + strconv.Itoa(int(typ))
	}
}

func CanProcessByType(Type DnsRRType) bool {
	switch Type {
	case DnsTypeA, DnsTypeNS, DnsTypeCNAME, DnsTypeSOA, DnsTypePTR, DnsTypeMX, DnsTypeTXT, DnsTypeAAAA, DnsTypeSRV:
		return true
	default:
		return false
	}
}

func GetDnsRRTypeFromString(s string) DnsRRType {
	switch strings.ToUpper(s) {
	case `A`:
		return DnsTypeA
	case `AAAA`:
		return DnsTypeAAAA
	case `AFSDB`:
		return DnsTypeAFSDB
	case `ANY`:
		return DnsTypeANY
	case `ATMA`:
		return DnsTypeATMA
	case `AXFR`:
		return DnsTypeAXFR
	case `CAA`:
		return DnsTypeCAA
	case `CDNSKEY`:
		return DnsTypeCDNSKEY
	case `CDS`:
		return DnsTypeCDS
	case `CERT`:
		return DnsTypeCERT
	case `CNAME`:
		return DnsTypeCNAME
	case `DHCID`:
		return DnsTypeDHCID
	case `DLV`:
		return DnsTypeDLV
	case `DNAME`:
		return DnsTypeDNAME
	case `DNSKEY`:
		return DnsTypeDNSKEY
	case `DS`:
		return DnsTypeDS
	case `EID`:
		return DnsTypeEID
	case `EUI48`:
		return DnsTypeEUI48
	case `EUI64`:
		return DnsTypeEUI64
	case `GID`:
		return DnsTypeGID
	case `GPOS`:
		return DnsTypeGPOS
	case `HINFO`:
		return DnsTypeHINFO
	case `HIP`:
		return DnsTypeHIP
	case `IPSECKEY`:
		return DnsTypeIPSECKEY
	case `ISDN`:
		return DnsTypeISDN
	case `IXFR`:
		return DnsTypeIXFR
	case `KEY`:
		return DnsTypeKEY
	case `KX`:
		return DnsTypeKX
	case `L32`:
		return DnsTypeL32
	case `L64`:
		return DnsTypeL64
	case `LOC`:
		return DnsTypeLOC
	case `LP`:
		return DnsTypeLP
	case `MAILA`:
		return DnsTypeMAILA
	case `MAILB`:
		return DnsTypeMAILB
	case `MB`:
		return DnsTypeMB
	case `MD`:
		return DnsTypeMD
	case `MF`:
		return DnsTypeMF
	case `MG`:
		return DnsTypeMG
	case `MINFO`:
		return DnsTypeMINFO
	case `MR`:
		return DnsTypeMR
	case `MX`:
		return DnsTypeMX
	case `NAPTR`:
		return DnsTypeNAPTR
	case `NID`:
		return DnsTypeNID
	case `NIMLOC`:
		return DnsTypeNIMLOC
	case `NINFO`:
		return DnsTypeNINFO
	case `NS`:
		return DnsTypeNS
	case `NSAPPTR`:
		return DnsTypeNSAPPTR
	case `NSEC`:
		return DnsTypeNSEC
	case `NSEC3`:
		return DnsTypeNSEC3
	case `NSEC3PARAM`:
		return DnsTypeNSEC3PARAM
	case `NULL`:
		return DnsTypeNULL
	case `NXT`:
		return DnsTypeNXT
	case `None`:
		return DnsTypeNone
	case `OPENPGPKEY`:
		return DnsTypeOPENPGPKEY
	case `OPT`:
		return DnsTypeOPT
	case `PTR`:
		return DnsTypePTR
	case `PX`:
		return DnsTypePX
	case `RKEY`:
		return DnsTypeRKEY
	case `RP`:
		return DnsTypeRP
	case `RRSIG`:
		return DnsTypeRRSIG
	case `RT`:
		return DnsTypeRT
	case `SIG`:
		return DnsTypeSIG
	case `SOA`:
		return DnsTypeSOA
	case `SPF`:
		return DnsTypeSPF
	case `SRV`:
		return DnsTypeSRV
	case `SSHFP`:
		return DnsTypeSSHFP
	case `TA`:
		return DnsTypeTA
	case `TALINK`:
		return DnsTypeTALINK
	case `TKEY`:
		return DnsTypeTKEY
	case `TLSA`:
		return DnsTypeTLSA
	case `TSIG`:
		return DnsTypeTSIG
	case `TXT`:
		return DnsTypeTXT
	case `UID`:
		return DnsTypeUID
	case `UINFO`:
		return DnsTypeUINFO
	case `UNSPEC`:
		return DnsTypeUNSPEC
	case `URI`:
		return DnsTypeURI
	case `WKS`:
		return DnsTypeWKS
	case `X25`:
		return DnsTypeX25
	default:
		return DnsTypeNone
	}
}

func GetAllDnsRRTypeStringList() []string {
	return []string{
		`A`,
		`AAAA`,
		`AFSDB`,
		`ANY`,
		`ATMA`,
		`AXFR`,
		`CAA`,
		`CDNSKEY`,
		`CDS`,
		`CERT`,
		`CNAME`,
		`DHCID`,
		`DLV`,
		`DNAME`,
		`DNSKEY`,
		`DS`,
		`EID`,
		`EUI48`,
		`EUI64`,
		`GID`,
		`GPOS`,
		`HINFO`,
		`HIP`,
		`IPSECKEY`,
		`ISDN`,
		`IXFR`,
		`KEY`,
		`KX`,
		`L32`,
		`L64`,
		`LOC`,
		`LP`,
		`MAILA`,
		`MAILB`,
		`MB`,
		`MD`,
		`MF`,
		`MG`,
		`MINFO`,
		`MR`,
		`MX`,
		`NAPTR`,
		`NID`,
		`NIMLOC`,
		`NINFO`,
		`NS`,
		`NSAPPTR`,
		`NSEC`,
		`NSEC3`,
		`NSEC3PARAM`,
		`NULL`,
		`NXT`,
		`None`,
		`OPENPGPKEY`,
		`OPT`,
		`PTR`,
		`PX`,
		`RKEY`,
		`RP`,
		`RRSIG`,
		`RT`,
		`SIG`,
		`SOA`,
		`SPF`,
		`SRV`,
		`SSHFP`,
		`TA`,
		`TALINK`,
		`TKEY`,
		`TLSA`,
		`TSIG`,
		`TXT`,
		`UID`,
		`UINFO`,
		`UNSPEC`,
		`URI`,
		`WKS`,
		`X25`,
	}
}
