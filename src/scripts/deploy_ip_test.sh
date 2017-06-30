if [ "$#" -ne 1 ]; then
    echo "Please input hostname"
    exit -1
fi


host_name=$1

ssh root@${host_name} "mkdir -p /usr/local/ip/"
ssh root@${host_name} "mkdir -p /usr/local/ip/log/"

# 拷贝: iplocation
ssh root@${host_name} "rm -f /usr/local/ip/iplocation"
scp iplocation root@${host_name}:/usr/local/ip/iplocation


# 拷贝脚本
scp control.sh  root@${host_name}:/usr/local/ip/
scp ip_query/qqwry.dat  root@${host_name}:/usr/local/ip/qqwry.dat
scp config.test.ini   root@${host_name}:/usr/local/ip/config.ini

scp scripts/iplocation.service root@${host_name}:/lib/systemd/system/iplocation.service