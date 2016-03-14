# -*- coding: utf-8 -*-
from __future__ import absolute_import
import time

from rpc_thrift.utils import get_service_protocol, get_fast_transport, get_base_protocol

from ip_service.IpService import Client


FAST = True
def main():

    endpoint = "127.0.0.1:5563"
    service = ""

    if FAST:
        get_fast_transport(endpoint)
        protocol =  get_service_protocol(service, fast=True)
    else:
        get_base_protocol(endpoint)
        protocol =  get_service_protocol(service)
    client = Client(protocol)


    total_times = 1000
    print "测试 IpToLocation 接口的时延:"
    t1 = time.time()
    for i in range(0, total_times):
        try:
            result = client.IpToLocation("60.29.255.197")
        except Exception as e:
            print "Exception: ", e
    t = time.time() - t1
    print "Elapsed: %.3fms" % (t / total_times * 1000)
    
    print "测试 ping 接口的时延:"
    t1 = time.time()
    for i  in range(0, total_times):
        try:
            result = client.ping()
        except Exception as e:
            print "Exception: ", e
    t = time.time() - t1
    print "Elapsed: %.3fms" % (t / total_times * 1000)
    
if __name__ == "__main__":
    main()
