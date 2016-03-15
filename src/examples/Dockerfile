FROM 123.59.151.144:5000/centos

MAINTAINER wangfei "wangfei@chunyu.me"


RUN mkdir -p /root/goprojects/iplocation/src/examples/ENV/
RUN mkdir -p /root/goprojects/iplocation/src/examples/ip_service/
RUN mkdir -p /root/goprojects/iplocation/src/examples/log/

ADD ENV /root/goprojects/iplocation/src/examples/ENV/
ADD ip_service /root/goprojects/iplocation/src/examples/ip_service/
ADD config.test.ini /root/goprojects/iplocation/src/examples/config.ini
ADD rpc_*.py /root/goprojects/iplocation/src/examples/

# VOLUME ["/usr/local/rpc_proxy"]

