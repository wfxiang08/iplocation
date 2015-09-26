# -*- coding:utf-8 -*-
#
# Autogenerated by Thrift Compiler (1.0.0-dev)
#
# DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
#
#  options string: py
#

from __future__ import absolute_import
from thrift.Thrift import TType, TMessageType, TException, TApplicationException
import rpc_thrift.services.RpcServiceBase
from ip_service.ttypes import *
from thrift.Thrift import TProcessor
from thrift.transport import TTransport
from thrift.protocol import TBinaryProtocol, TProtocol
try:
  from rpc_thrift.cython.cybinary_protocol import TCyBinaryProtocol
except:
  TCyBinaryProtocol = None


class Iface(rpc_thrift.services.RpcServiceBase.Iface):
  def IpToLocation(self, ip):
    """
    根据IP获取相关的Location

    Parameters:
     - ip
    """
    pass


class Client(rpc_thrift.services.RpcServiceBase.Client, Iface):
  def __init__(self, iprot, oprot=None):
    rpc_thrift.services.RpcServiceBase.Client.__init__(self, iprot, oprot)

  def IpToLocation(self, ip):
    """
    根据IP获取相关的Location

    Parameters:
     - ip
    """
    self.send_IpToLocation(ip)
    return self.recv_IpToLocation()

  def send_IpToLocation(self, ip):
    self._oprot.writeMessageBegin('IpToLocation', TMessageType.CALL, self._seqid)
    args = IpToLocation_args()
    args.ip = ip
    args.write(self._oprot)
    self._oprot.writeMessageEnd()
    self._oprot.trans.flush()

  def recv_IpToLocation(self):
    iprot = self._iprot
    (fname, mtype, rseqid) = iprot.readMessageBegin()
    if mtype == TMessageType.EXCEPTION:
      x = TApplicationException()
      x.read(iprot)
      iprot.readMessageEnd()
      raise x
    result = IpToLocation_result()
    result.read(iprot)
    iprot.readMessageEnd()
    if result.success is not None:
      return result.success
    if result.re is not None:
      raise result.re
    raise TApplicationException(TApplicationException.MISSING_RESULT, "IpToLocation failed: unknown result")


class Processor(rpc_thrift.services.RpcServiceBase.Processor, Iface, TProcessor):
  def __init__(self, handler):
    rpc_thrift.services.RpcServiceBase.Processor.__init__(self, handler)
    self._processMap["IpToLocation"] = Processor.process_IpToLocation

  def process(self, iprot, oprot):
    (name, type, seqid) = iprot.readMessageBegin()
    if name not in self._processMap:
      iprot.skip(TType.STRUCT)
      iprot.readMessageEnd()
      x = TApplicationException(TApplicationException.UNKNOWN_METHOD, 'Unknown function %s' % (name))
      oprot.writeMessageBegin(name, TMessageType.EXCEPTION, seqid)
      x.write(oprot)
      oprot.writeMessageEnd()
      oprot.trans.flush()
      return
    else:
      self._processMap[name](self, seqid, iprot, oprot)
    return True

  def process_IpToLocation(self, seqid, iprot, oprot):
    args = IpToLocation_args()
    args.read(iprot)
    iprot.readMessageEnd()
    result = IpToLocation_result()
    try:
      result.success = self._handler.IpToLocation(args.ip)
    except rpc_thrift.services.ttypes.RpcException as re:
      result.re = re
    oprot.writeMessageBegin("IpToLocation", TMessageType.REPLY, seqid)
    result.write(oprot)
    oprot.writeMessageEnd()
    oprot.trans.flush()


# HELPER FUNCTIONS AND STRUCTURES

class IpToLocation_args:
  """
  Attributes:
   - ip
  """

  thrift_spec = (
    None, # 0
    (1, TType.STRING, 'ip', None, None, ), # 1
  )

  def __init__(self, ip=None,):
    self.ip = ip

  def read(self, iprot):
    if iprot.__class__ == TCyBinaryProtocol and self.thrift_spec is not None:
      iprot.read_struct(self)
      return
    iprot.readStructBegin()
    while True:
      (fname, ftype, fid) = iprot.readFieldBegin()
      if ftype == TType.STOP:
        break
      if fid == 1:
        if ftype == TType.STRING:
          self.ip = iprot.readString()
        else:
          iprot.skip(ftype)
      else:
        iprot.skip(ftype)
      iprot.readFieldEnd()
    iprot.readStructEnd()

  def write(self, oprot):
    if oprot.__class__ == TCyBinaryProtocol and self.thrift_spec is not None:
      oprot.write_struct(self)
      return
    oprot.writeStructBegin('IpToLocation_args')
    if self.ip is not None:
      oprot.writeFieldBegin('ip', TType.STRING, 1)
      oprot.writeString(self.ip)
      oprot.writeFieldEnd()
    oprot.writeFieldStop()
    oprot.writeStructEnd()

  def validate(self):
    return


  def __hash__(self):
    value = 17
    value = (value * 31) ^ hash(self.ip)
    return value

  def __repr__(self):
    L = ['%s=%r' % (key, value)
      for key, value in self.__dict__.iteritems()]
    return '%s(%s)' % (self.__class__.__name__, ', '.join(L))

  def __eq__(self, other):
    return isinstance(other, self.__class__) and self.__dict__ == other.__dict__

  def __ne__(self, other):
    return not (self == other)

class IpToLocation_result:
  """
  Attributes:
   - success
   - re
  """

  thrift_spec = (
    (0, TType.STRUCT, 'success', (Location, Location.thrift_spec), None, ), # 0
    (1, TType.STRUCT, 're', (rpc_thrift.services.ttypes.RpcException, rpc_thrift.services.ttypes.RpcException.thrift_spec), None, ), # 1
  )

  def __init__(self, success=None, re=None,):
    self.success = success
    self.re = re

  def read(self, iprot):
    if iprot.__class__ == TCyBinaryProtocol and self.thrift_spec is not None:
      iprot.read_struct(self)
      return
    iprot.readStructBegin()
    while True:
      (fname, ftype, fid) = iprot.readFieldBegin()
      if ftype == TType.STOP:
        break
      if fid == 0:
        if ftype == TType.STRUCT:
          self.success = Location()
          self.success.read(iprot)
        else:
          iprot.skip(ftype)
      elif fid == 1:
        if ftype == TType.STRUCT:
          self.re = rpc_thrift.services.ttypes.RpcException()
          self.re.read(iprot)
        else:
          iprot.skip(ftype)
      else:
        iprot.skip(ftype)
      iprot.readFieldEnd()
    iprot.readStructEnd()

  def write(self, oprot):
    if oprot.__class__ == TCyBinaryProtocol and self.thrift_spec is not None:
      oprot.write_struct(self)
      return
    oprot.writeStructBegin('IpToLocation_result')
    if self.success is not None:
      oprot.writeFieldBegin('success', TType.STRUCT, 0)
      self.success.write(oprot)
      oprot.writeFieldEnd()
    if self.re is not None:
      oprot.writeFieldBegin('re', TType.STRUCT, 1)
      self.re.write(oprot)
      oprot.writeFieldEnd()
    oprot.writeFieldStop()
    oprot.writeStructEnd()

  def validate(self):
    return


  def __hash__(self):
    value = 17
    value = (value * 31) ^ hash(self.success)
    value = (value * 31) ^ hash(self.re)
    return value

  def __repr__(self):
    L = ['%s=%r' % (key, value)
      for key, value in self.__dict__.iteritems()]
    return '%s(%s)' % (self.__class__.__name__, ', '.join(L))

  def __eq__(self, other):
    return isinstance(other, self.__class__) and self.__dict__ == other.__dict__

  def __ne__(self, other):
    return not (self == other)
