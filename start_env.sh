if [ ! -d bin ] || [ ! -d src ] || [ ! -d pkg ]; then
	echo "----------------"
	echo "ERROR: 根目录应该包含 bin, src, pkg 这三个子目录; 如果没有这个目录，可能在错误的位置调用脚本"
	echo "----------------"
	exit -1
fi
source gvp
echo "----------------"
echo "\$GOPATH = ${GOPATH}"
echo "----------------"