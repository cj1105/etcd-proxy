// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: qos.proto

package qospb

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Subject struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Prefix               string   `protobuf:"bytes,2,opt,name=prefix,proto3" json:"prefix,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Subject) Reset()         { *m = Subject{} }
func (m *Subject) String() string { return proto.CompactTextString(m) }
func (*Subject) ProtoMessage()    {}
func (*Subject) Descriptor() ([]byte, []int) {
	return fileDescriptor_acb97775b1858998, []int{0}
}
func (m *Subject) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Subject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Subject.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Subject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Subject.Merge(m, src)
}
func (m *Subject) XXX_Size() int {
	return m.Size()
}
func (m *Subject) XXX_DiscardUnknown() {
	xxx_messageInfo_Subject.DiscardUnknown(m)
}

var xxx_messageInfo_Subject proto.InternalMessageInfo

type QoSRule struct {
	RuleName             string   `protobuf:"bytes,1,opt,name=rule_name,json=ruleName,proto3" json:"rule_name,omitempty"`
	RuleType             string   `protobuf:"bytes,2,opt,name=rule_type,json=ruleType,proto3" json:"rule_type,omitempty"`
	Subject              *Subject `protobuf:"bytes,3,opt,name=subject,proto3" json:"subject,omitempty"`
	Qps                  uint64   `protobuf:"varint,4,opt,name=qps,proto3" json:"qps,omitempty"`
	Threshold            uint64   `protobuf:"varint,5,opt,name=threshold,proto3" json:"threshold,omitempty"`
	Condition            string   `protobuf:"bytes,6,opt,name=condition,proto3" json:"condition,omitempty"`
	Priority             uint64   `protobuf:"varint,7,opt,name=priority,proto3" json:"priority,omitempty"`
	Ratelimiter          string   `protobuf:"bytes,8,opt,name=ratelimiter,proto3" json:"ratelimiter,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QoSRule) Reset()         { *m = QoSRule{} }
func (m *QoSRule) String() string { return proto.CompactTextString(m) }
func (*QoSRule) ProtoMessage()    {}
func (*QoSRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_acb97775b1858998, []int{1}
}
func (m *QoSRule) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QoSRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QoSRule.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QoSRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QoSRule.Merge(m, src)
}
func (m *QoSRule) XXX_Size() int {
	return m.Size()
}
func (m *QoSRule) XXX_DiscardUnknown() {
	xxx_messageInfo_QoSRule.DiscardUnknown(m)
}

var xxx_messageInfo_QoSRule proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Subject)(nil), "qospb.Subject")
	proto.RegisterType((*QoSRule)(nil), "qospb.QoSRule")
}

func init() { proto.RegisterFile("qos.proto", fileDescriptor_acb97775b1858998) }

var fileDescriptor_acb97775b1858998 = []byte{
	// 280 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xdd, 0x4a, 0xc3, 0x30,
	0x14, 0xc7, 0x1b, 0xd7, 0xf5, 0xe3, 0x0c, 0x44, 0x82, 0x48, 0x98, 0x12, 0xca, 0xae, 0x7a, 0x55,
	0x41, 0xf1, 0x05, 0x7c, 0x00, 0xc1, 0xce, 0x7b, 0x59, 0xb7, 0xb8, 0x45, 0xba, 0x9e, 0x34, 0x4d,
	0xc1, 0xbe, 0x89, 0x8f, 0xb4, 0xcb, 0x3d, 0x82, 0xab, 0x4f, 0xe1, 0x9d, 0x34, 0xeb, 0x3e, 0xee,
	0xfe, 0xf9, 0xfd, 0xf8, 0x9f, 0x73, 0x08, 0x84, 0x25, 0x56, 0x89, 0xd2, 0x68, 0x90, 0x0e, 0x4b,
	0xac, 0x54, 0x36, 0xbe, 0x5e, 0xe2, 0x12, 0x2d, 0xb9, 0xef, 0xd2, 0x5e, 0x4e, 0x9e, 0xc0, 0x9f,
	0xd6, 0xd9, 0xa7, 0x98, 0x1b, 0x4a, 0xc1, 0x2d, 0x66, 0x6b, 0xc1, 0x48, 0x44, 0xe2, 0x30, 0xb5,
	0x99, 0xde, 0x80, 0xa7, 0xb4, 0xf8, 0x90, 0x5f, 0xec, 0xc2, 0xd2, 0xfe, 0x35, 0xf9, 0x23, 0xe0,
	0xbf, 0xe2, 0x34, 0xad, 0x73, 0x41, 0x6f, 0x21, 0xd4, 0x75, 0x2e, 0xde, 0xcf, 0xca, 0x41, 0x07,
	0x5e, 0xba, 0x01, 0x07, 0x69, 0x1a, 0x25, 0xfa, 0x19, 0x56, 0xbe, 0x35, 0x4a, 0xd0, 0x18, 0xfc,
	0x6a, 0xbf, 0x9c, 0x0d, 0x22, 0x12, 0x8f, 0x1e, 0x2e, 0x13, 0x7b, 0x6b, 0xd2, 0x9f, 0x94, 0x1e,
	0x34, 0xbd, 0x82, 0x41, 0xa9, 0x2a, 0xe6, 0x46, 0x24, 0x76, 0xd3, 0x2e, 0xd2, 0x3b, 0x08, 0xcd,
	0x4a, 0x8b, 0x6a, 0x85, 0xf9, 0x82, 0x0d, 0x2d, 0x3f, 0x81, 0xce, 0xce, 0xb1, 0x58, 0x48, 0x23,
	0xb1, 0x60, 0x9e, 0x5d, 0x7b, 0x02, 0x74, 0x0c, 0x81, 0xd2, 0x12, 0xb5, 0x34, 0x0d, 0xf3, 0x6d,
	0xf5, 0xf8, 0xa6, 0x11, 0x8c, 0xf4, 0xcc, 0x88, 0x5c, 0xae, 0xa5, 0x11, 0x9a, 0x05, 0xb6, 0x7b,
	0x8e, 0x9e, 0xd9, 0x66, 0xc7, 0x9d, 0xed, 0x8e, 0x3b, 0x9b, 0x96, 0x93, 0x6d, 0xcb, 0xc9, 0x4f,
	0xcb, 0xc9, 0xf7, 0x2f, 0x77, 0x32, 0xcf, 0xfe, 0xe9, 0xe3, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x57, 0xba, 0xb8, 0xc5, 0x7d, 0x01, 0x00, 0x00,
}

func (m *Subject) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Subject) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Subject) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Prefix) > 0 {
		i -= len(m.Prefix)
		copy(dAtA[i:], m.Prefix)
		i = encodeVarintQos(dAtA, i, uint64(len(m.Prefix)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintQos(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QoSRule) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QoSRule) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QoSRule) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Ratelimiter) > 0 {
		i -= len(m.Ratelimiter)
		copy(dAtA[i:], m.Ratelimiter)
		i = encodeVarintQos(dAtA, i, uint64(len(m.Ratelimiter)))
		i--
		dAtA[i] = 0x42
	}
	if m.Priority != 0 {
		i = encodeVarintQos(dAtA, i, uint64(m.Priority))
		i--
		dAtA[i] = 0x38
	}
	if len(m.Condition) > 0 {
		i -= len(m.Condition)
		copy(dAtA[i:], m.Condition)
		i = encodeVarintQos(dAtA, i, uint64(len(m.Condition)))
		i--
		dAtA[i] = 0x32
	}
	if m.Threshold != 0 {
		i = encodeVarintQos(dAtA, i, uint64(m.Threshold))
		i--
		dAtA[i] = 0x28
	}
	if m.Qps != 0 {
		i = encodeVarintQos(dAtA, i, uint64(m.Qps))
		i--
		dAtA[i] = 0x20
	}
	if m.Subject != nil {
		{
			size, err := m.Subject.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQos(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.RuleType) > 0 {
		i -= len(m.RuleType)
		copy(dAtA[i:], m.RuleType)
		i = encodeVarintQos(dAtA, i, uint64(len(m.RuleType)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.RuleName) > 0 {
		i -= len(m.RuleName)
		copy(dAtA[i:], m.RuleName)
		i = encodeVarintQos(dAtA, i, uint64(len(m.RuleName)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQos(dAtA []byte, offset int, v uint64) int {
	offset -= sovQos(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Subject) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovQos(uint64(l))
	}
	l = len(m.Prefix)
	if l > 0 {
		n += 1 + l + sovQos(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *QoSRule) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RuleName)
	if l > 0 {
		n += 1 + l + sovQos(uint64(l))
	}
	l = len(m.RuleType)
	if l > 0 {
		n += 1 + l + sovQos(uint64(l))
	}
	if m.Subject != nil {
		l = m.Subject.Size()
		n += 1 + l + sovQos(uint64(l))
	}
	if m.Qps != 0 {
		n += 1 + sovQos(uint64(m.Qps))
	}
	if m.Threshold != 0 {
		n += 1 + sovQos(uint64(m.Threshold))
	}
	l = len(m.Condition)
	if l > 0 {
		n += 1 + l + sovQos(uint64(l))
	}
	if m.Priority != 0 {
		n += 1 + sovQos(uint64(m.Priority))
	}
	l = len(m.Ratelimiter)
	if l > 0 {
		n += 1 + l + sovQos(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovQos(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQos(x uint64) (n int) {
	return sovQos(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Subject) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQos
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Subject: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Subject: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Prefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQos(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQos
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QoSRule) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQos
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QoSRule: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QoSRule: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RuleName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RuleName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RuleType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RuleType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subject", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQos
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Subject == nil {
				m.Subject = &Subject{}
			}
			if err := m.Subject.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Qps", wireType)
			}
			m.Qps = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Qps |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Threshold", wireType)
			}
			m.Threshold = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Threshold |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Condition", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Condition = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Priority", wireType)
			}
			m.Priority = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Priority |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ratelimiter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQos
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQos
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQos
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ratelimiter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQos(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQos
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQos(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQos
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQos
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQos
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQos
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQos
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQos
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQos        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQos          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQos = fmt.Errorf("proto: unexpected end of group")
)
