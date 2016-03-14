# -*- coding: utf-8 -*-
from __future__ import absolute_import
import time

from ip_service.IpService import Client
from rpc_thrift.config import RPC_DEFAULT_CONFIG, RPC_PROXY_ADDRESS, RPC_SERVICE, print_exception
from rpc_thrift.config import parse_config
from rpc_thrift.utils import get_base_protocol, get_service_protocol, get_fast_transport


FAST = True
def main():

    config_path = RPC_DEFAULT_CONFIG
    config = parse_config(config_path)

    # 从配置文件读取配置
    endpoint = config[RPC_PROXY_ADDRESS]
    service = config[RPC_SERVICE]

    if FAST:
        get_fast_transport(endpoint)
        protocol =  get_service_protocol(service, fast=True)
    else:
        get_base_protocol(endpoint)
        protocol =  get_service_protocol(service)
    client = Client(protocol)


    print "测试 IpToLocation 接口的时延:"
    total_times = 10000
  

    print "测试 ping1 接口的时延:"
    t1 = time.time()
    for i  in range(0, total_times):
        try:
            client.ping1()
        except Exception as e:
            print "Exception: ", e
            print_exception()
            break
    t = time.time() - t1
    print "Elapsed: %.3fms" % (t / total_times * 1000)

if __name__ == "__main__":
    main()
