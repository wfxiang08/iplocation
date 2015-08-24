const string VERSION = "0.0.1"


exception RpcException {
  1: i32  code,
  2: string msg
}

/**
 * 输入和输出的结果
 */
struct Location {
	1:string city,
	2:string province,
	3:string detail
}

service IpService {
	/**
	 * 根据IP获取相关的Location
	 */
    Location IpToLocation(1: string ip) throws (1: RpcException re),
}