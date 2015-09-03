#!/bin/bash

WORKSPACE=$(cd $(dirname $0)/; pwd)
cd $WORKSPACE

mkdir -p log

app=./iplocation
conf=config.ini
proxy_log=log/ip.log
pidfile=log/app.pid
logfile=log/app.log

function check_pid() {
    if [ -f $pidfile ];then
        pid=`cat $pidfile`
        if [ -n $pid ]; then
            running=`ps -p $pid|grep -v "PID TTY" |wc -l`
            return $running
        fi
    fi
    return 0
}

function start() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo -n "$app now is running already, pid="
        cat $pidfile
        return 1
    fi

    if ! [ -f $conf ];then
        echo "Config file $conf doesn't exist"
        exit -1
    fi
    nohup $app -c $conf -L $proxy_log &> $logfile &
    echo $! > $pidfile
    echo "$app started..., pid=$!"
}

function stop() {
	check_pid
    running=$?
	if [ $running -gt 0 ];then
	    pid=`cat $pidfile`
		kill -15 $pid
		status="0"
		while [ "$status" == "0" ];do
			echo "Waiting for process ${pid} ..."
			sleep 1
			ps -p$pid 2>&1 > /dev/null
			status=$?
		done
	    echo "$app stoped..."
	else
		echo "$app already stoped..."
	fi
}

function restart() {
    stop
    sleep 1
    start
}

function status() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo started
    else
        echo stoped
    fi
}

function tailf() {	
	date=`date +"%Y%m%d"`
	tail -Fn 200 "${proxy_log}-${date}"	
}


function help() {
    echo "$0 start|stop|restart|status|tail"
}

if [ "$1" == "" ]; then
    help
elif [ "$1" == "stop" ];then
    stop
elif [ "$1" == "start" ];then
    start
elif [ "$1" == "restart" ];then
    restart
elif [ "$1" == "status" ];then
    status
elif [ "$1" == "tail" ];then
    tailf
else
    help
fi
