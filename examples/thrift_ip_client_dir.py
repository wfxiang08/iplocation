# -*- coding: utf-8 -*-
from __future__ import absolute_import
import time

from ip_service.IpService import Client
from rpc_thrift.config import RPC_DEFAULT_CONFIG, RPC_PROXY_ADDRESS, RPC_SERVICE
from rpc_thrift.config import parse_config
from rpc_thrift.utils import get_service_protocol
from rpc_thrift.utils import get_base_protocol

"""
1. iplocation service Unix Domain Socket:
    测试 IpToLocation 接口的时延:
    Elapsed: 0.197ms
    测试 ping 接口的时延:
    Elapsed: 0.096ms

2. Local Loop:
    测试 IpToLocation 接口的时延:
    Elapsed: 0.211ms
    测试 ping 接口的时延:
    Elapsed: 0.128ms
"""

def main():
	# 直接连接 RPC服务器(不经过proxy这些环节)	
    config_path = RPC_DEFAULT_CONFIG
    config = parse_config(config_path)

    # 从配置文件读取配置
    endpoint = config[RPC_PROXY_ADDRESS]
    endpoint = "127.0.0.1:5563"
    # endpoint="/Users/feiwang/gowork/src/git.chunyu.me/infra/iplocation/aa.sock"
    service = config[RPC_SERVICE]

    get_base_protocol(endpoint)
    protocol =  get_service_protocol("")
    client = Client(protocol)


    total_times = 100
    print "测试 IpToLocation 接口的时延:"
    t1 = time.time()
    for i  in range(0, total_times):
        try:
            result = client.IpToLocation("60.29.255.197")
            # print result.city
            # print result.province
            # print result.detail
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
