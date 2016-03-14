FROM 123.59.151.144:5000/centos

MAINTAINER wangfei "wangfei@chunyu.me"


RUN mkdir -p /root/goprojects/iplocation/src/examples/ENV/

ADD ENV /root/goprojects/iplocation/src/examples/ENV/
ADD ip_service /root/goprojects/iplocation/src/examples/
ADD config.test.ini /root/goprojects/iplocation/src/examples/config.ini
ADD rpc_*.py /root/goprojects/iplocation/src/examples/


VOLUME ["/usr/local/rpc_proxy"]

