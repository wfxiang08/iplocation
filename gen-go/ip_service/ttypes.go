// Autogenerated by Thrift Compiler (1.0.0-dev)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package ip_service

import (
	"bytes"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var GoUnusedProtection__ int

// Attributes:
//  - Code
//  - Msg
type RpcException struct {
	Code int32  `thrift:"code,1" json:"code"`
	Msg  string `thrift:"msg,2" json:"msg"`
}

func NewRpcException() *RpcException {
	return &RpcException{}
}

func (p *RpcException) GetCode() int32 {
	return p.Code
}

func (p *RpcException) GetMsg() string {
	return p.Msg
}
func (p *RpcException) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *RpcException) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Code = v
	}
	return nil
}

func (p *RpcException) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Msg = v
	}
	return nil
}

func (p *RpcException) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("RpcException"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *RpcException) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("code", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:code: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Code)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.code (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:code: ", p), err)
	}
	return err
}

func (p *RpcException) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("msg", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:msg: ", p), err)
	}
	if err := oprot.WriteString(string(p.Msg)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.msg (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:msg: ", p), err)
	}
	return err
}

func (p *RpcException) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RpcException(%+v)", *p)
}

func (p *RpcException) Error() string {
	return p.String()
}

// 输入和输出的结果
//
// Attributes:
//  - City
//  - Province
//  - Detail
type Location struct {
	City     string `thrift:"city,1" json:"city"`
	Province string `thrift:"province,2" json:"province"`
	Detail   string `thrift:"detail,3" json:"detail"`
}

func NewLocation() *Location {
	return &Location{}
}

func (p *Location) GetCity() string {
	return p.City
}

func (p *Location) GetProvince() string {
	return p.Province
}

func (p *Location) GetDetail() string {
	return p.Detail
}
func (p *Location) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *Location) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.City = v
	}
	return nil
}

func (p *Location) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Province = v
	}
	return nil
}

func (p *Location) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.Detail = v
	}
	return nil
}

func (p *Location) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Location"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Location) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("city", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:city: ", p), err)
	}
	if err := oprot.WriteString(string(p.City)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.city (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:city: ", p), err)
	}
	return err
}

func (p *Location) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("province", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:province: ", p), err)
	}
	if err := oprot.WriteString(string(p.Province)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.province (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:province: ", p), err)
	}
	return err
}

func (p *Location) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("detail", thrift.STRING, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:detail: ", p), err)
	}
	if err := oprot.WriteString(string(p.Detail)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.detail (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:detail: ", p), err)
	}
	return err
}

func (p *Location) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Location(%+v)", *p)
}
