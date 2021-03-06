include "rpc_thrift.services.thrift"
const string VERSION = "0.0.1"

/**
 * 输入和输出的结果
 */
struct Location {
	1:string city,
	2:string province,
	3:string detail
}

service IpService extends rpc_thrift.services.RpcServiceBase {
	/**
	 * 根据IP获取相关的Location
	 */
    Location IpToLocation(1: string ip) throws (1: rpc_thrift.services.RpcException re),
}