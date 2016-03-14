# thrift -r --gen py geolocation_service.thrift
thrift --gen go:package_prefix="git.chunyu.me/golang/rpc_proxy_base/src/",thrift_import="git.chunyu.me/infra/go_thrift/thrift" ip_service.thrift
find gen-go -name '*-remote.go' | xargs perl -pi -e 's|git.chunyu.me/golang/rpc_proxy_base/src/||g'
