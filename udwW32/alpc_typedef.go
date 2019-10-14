// +build windows

package udwW32

import (
	"errors"
)

type ALPC_MESSAGE_ATTRIBUTES struct {
	AllocatedAttributes uint32
	ValidAttributes     uint32
}

type ALPC_CONTEXT_ATTR struct {
	PortContext    *AlpcPortContext
	MessageContext uintptr
	Sequence       uint32
	MessageId      uint32
	CallbackId     uint32
}

type ALPC_HANDLE_ATTR struct {
	Flags         uint32
	Handle        HANDLE
	ObjectType    uint32
	DesiredAccess uint32
}

type CLIENT_ID struct {
	UniqueProcess uintptr
	UniqueThread  uintptr
}

type UNICODE_STRING struct {
	Length        uint16
	MaximumLength uint16
	_             [4]byte
	Buffer        *uint16
}

type OBJECT_ATTRIBUTES struct {
	Length                   uint32
	_                        [4]byte
	RootDirectory            HANDLE
	ObjectName               *UNICODE_STRING
	Attributes               uint32
	_                        [4]byte
	SecurityDescriptor       *SECURITY_DESCRIPTOR
	SecurityQualityOfService *SECURITY_QUALITY_OF_SERVICE
}

type PORT_MESSAGE struct {
	DataLength     uint16
	TotalLength    uint16
	Type           uint16
	DataInfoOffset uint16
	ClientId       CLIENT_ID
	MessageId      uint32
	_              [4]byte
	ClientViewSize uint64
}

func (pm PORT_MESSAGE) CallbackId() uint32 {
	return uint32(pm.ClientViewSize >> 32)
}

func (pm PORT_MESSAGE) DoNotUseThisField() float64 {
	panic("WE TOLD YOU NOT TO USE THIS FIELD")
}

const PORT_MESSAGE_SIZE = 0x28

type SECURITY_QUALITY_OF_SERVICE struct {
	Length              uint32
	ImpersonationLevel  uint32
	ContextTrackingMode byte
	EffectiveOnly       byte
	_                   [2]byte
}

const SECURITY_QOS_SIZE = 12

type ALPC_PORT_ATTRIBUTES struct {
	Flags               uint32
	SecurityQos         SECURITY_QUALITY_OF_SERVICE
	MaxMessageLength    uint64
	MemoryBandwidth     uint64
	MaxPoolUsage        uint64
	MaxSectionSize      uint64
	MaxViewSize         uint64
	MaxTotalSectionSize uint64
	DupObjectTypes      uint32
	Reserved            uint32
}

const SHORT_MESSAGE_MAX_SIZE uint16 = 65535
const SHORT_MESSAGE_MAX_PAYLOAD uint16 = SHORT_MESSAGE_MAX_SIZE - PORT_MESSAGE_SIZE

type AlpcShortMessage struct {
	PORT_MESSAGE
	Data [SHORT_MESSAGE_MAX_PAYLOAD]byte
}

func NewAlpcShortMessage() AlpcShortMessage {
	sm := AlpcShortMessage{}
	sm.TotalLength = SHORT_MESSAGE_MAX_SIZE
	return sm
}

func (sm *AlpcShortMessage) SetData(d []byte) (e error) {

	copy(sm.Data[:], d)
	if len(d) > int(SHORT_MESSAGE_MAX_PAYLOAD) {
		e = errors.New("data too big - truncated")
		sm.DataLength = SHORT_MESSAGE_MAX_PAYLOAD
		sm.TotalLength = SHORT_MESSAGE_MAX_SIZE
		return
	}
	sm.TotalLength = uint16(PORT_MESSAGE_SIZE + len(d))
	sm.DataLength = uint16(len(d))
	return

}

func (sm *AlpcShortMessage) GetData() []byte {
	if int(sm.DataLength) > int(SHORT_MESSAGE_MAX_PAYLOAD) {
		return sm.Data[:]
	}
	return sm.Data[:sm.DataLength]
}

func (sm *AlpcShortMessage) Reset() {

	sm.PORT_MESSAGE = PORT_MESSAGE{}
	sm.TotalLength = SHORT_MESSAGE_MAX_SIZE
	sm.DataLength = 0
}

type AlpcPortContext struct {
	Handle HANDLE
}
