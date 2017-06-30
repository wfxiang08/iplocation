# -*- coding: utf-8 -*-
from __future__ import absolute_import
import time

from ip_service.IpService import Client
from rpc_thrift.config import RPC_DEFAULT_CONFIG, RPC_PROXY_ADDRESS, RPC_SERVICE
from rpc_thrift.config import parse_config
from rpc_thrift.utils import get_service_protocol, get_fast_transport, get_base_protocol


FAST = True

def main():

    config_path = RPC_DEFAULT_CONFIG
    config = parse_config(config_path)

    # 从配置文件读取配置
    endpoint = config[RPC_PROXY_ADDRESS]
    endpoint = "60.29.255.199:5550"
    service = config[RPC_SERVICE]

    if FAST:
        get_fast_transport(endpoint)
        protocol =  get_service_protocol(service, fast=True)
    else:
        get_base_protocol(endpoint)
        protocol =  get_service_protocol(service)

    client = Client(protocol)


    total_times = 1000
    t1 = time.time()
    for i  in range(0, total_times):
        try:

            result = client.IpToLocation("60.29.255.197")
        except Exception as e:
            print "Exception: ", e

        if i % 200 == 0:
            print "QPS: %.2f" % (i / (time.time() - t1), )

    t = time.time() - t1
    print "Elapsed: %.3fms",  (t / total_times * 1000.0)

if __name__ == "__main__":
    main()
