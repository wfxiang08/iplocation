FROM 123.59.151.144:5000/centos

MAINTAINER luohongwei "luohongwei@chunyu.me"

RUN mkdir -p /usr/local/ip/

ADD ip_query/qqwry.dat /usr/local/ip/qqwry.dat
ADD iplocation /usr/local/ip/

ADD control.sh /usr/local/ip/

CMD /usr/local/ip/control.sh start

