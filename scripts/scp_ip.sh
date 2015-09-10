if [ "$#" -ne 1 ]; then
    echo "Please input hostname"
    exit -1
fi

host_name=$1

ssh root@${host_name} "rm -f /usr/local/ip/iplocation"
scp iplocation root@${host_name}:/usr/local/ip/iplocation

