# -*- coding: utf-8 -*-
from __future__ import absolute_import
import time

from zerothrift import (parse_config, RPC_PROXY_ADDRESS,
                        RPC_DEFAULT_CONFIG, RPC_SERVICE)
from zerothrift import (TimeoutException, get_transport, TZmqBinaryProtocol)
from ip_service.IpService import Client


class TZmqBinaryProtocolEx(TZmqBinaryProtocol):
    elapsed = 0
    name = ""
    def writeMessageBegin(self, name, type, seqid):
        self.elapsed = time.time()
        self.name = name

        TZmqBinaryProtocol.writeMessageBegin(self, name, type, seqid)

    def readMessageEnd(self):
        TZmqBinaryProtocol.readMessageEnd(self)
        elapsed = time.time() - self.elapsed


        print "Zerothrift %s, delay: %.2fms" % (self.name, elapsed * 1000)

def main():

    config_path = RPC_DEFAULT_CONFIG
    config = parse_config(config_path)

    # 从配置文件读取配置
    endpoint = config[RPC_PROXY_ADDRESS]
    service = config[RPC_SERVICE]



    # 获取protocol
    transport = get_transport(endpoint)
    # protocol = get_protocol(service)
    protocol = TZmqBinaryProtocolEx(transport, service=service)

    # 获取Client
    client = Client(protocol)


    total_times = 10
    t1 = time.time()
    result_set = set()
    for i  in range(0, total_times):
        print "index: ", i
        try:

            result = client.IpToLocation("60.29.255.197")
            print result.city
            print result.province
            print result.detail


        except TimeoutException as e:
            print "TimeoutException: ", e
        except Exception as e:
            print "Exception: ", e

        if i % 200 == 0:
            print "QPS: %.2f" % (i / (time.time() - t1), )

    print "Total Result: ", len(result_set)
    t = time.time() - t1
    print "Elapsed: ",  t / total_times
if __name__ == "__main__":
    main()
