# -*- coding: utf-8 -*-
from __future__ import absolute_import
import time

from ip_service.IpService import Client
from rpc_thrift.config import RPC_DEFAULT_CONFIG, RPC_PROXY_ADDRESS, RPC_SERVICE
from rpc_thrift.config import parse_config
from rpc_thrift.utils import get_service_protocol
from rpc_thrift.utils import get_base_protocol


def main():

    config_path = RPC_DEFAULT_CONFIG
    config = parse_config(config_path)

    # 从配置文件读取配置
    endpoint = config[RPC_PROXY_ADDRESS]
    endpoint = "127.0.0.1:5550"
    service = config[RPC_SERVICE]

    get_base_protocol(endpoint)
    protocol =  get_service_protocol(service)
    client = Client(protocol)


    total_times = 100
    t1 = time.time()
    result_set = set()
    for i  in range(0, total_times):
        print "index: ", i
        try:

            result = client.IpToLocation("60.29.255.197")
            print result.city
            print result.province
            print result.detail

        except Exception as e:
            print "Exception: ", e

        if i % 200 == 0:
            print "QPS: %.2f" % (i / (time.time() - t1), )

    print "Total Result: ", len(result_set)
    t = time.time() - t1
    print "Elapsed: ",  t / total_times
if __name__ == "__main__":
    main()
