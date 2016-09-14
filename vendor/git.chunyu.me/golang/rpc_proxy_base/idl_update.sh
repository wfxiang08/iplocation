thrift -r --gen go:thrift_import="git.chunyu.me/infra/go_thrift/thrift" rpc_thrift.services.thrift
rm -rf src/rpc_thrift
mv gen-go/rpc_thrift src/rpc_thrift
rm -rf gen-go