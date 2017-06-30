#!/usr/bin/env bash
# thrift -r --gen py geolocation_service.thrift
thrift --gen go:package_prefix="github.com/wfxiang08/thrift_rpc_base/",thrift_import="github.com/wfxiang08/go_thrift/thrift" ip_service.thrift
# find gen-go -name '*-remote.go' | xargs perl -pi -e 's|git.chunyu.me/golang/rpc_proxy_base/src/||g'
