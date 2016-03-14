# thrift -r --gen py geolocation_service.thrift
thrift --gen go:package_prefix=git.chunyu.me/infra/rpc_proxy/gen-go/  ip_service.thrift
# mv gen-py/cy_user_service/* cy_user_service