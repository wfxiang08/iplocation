package ip_query

import (
	ips "git.chunyu.me/infra/iplocation/gen-go/ip_service"
	// services "git.chunyu.me/infra/rpc_proxy/gen-go/rpc_thrift/services"
	"git.chunyu.me/infra/rpc_proxy/utils/log"
	"strings"
)

type Handler struct {
	Service *IpInfoService
}

func NewHandler(ipdb string) *Handler {
	p := &Handler{
		Service: &IpInfoService{},
	}
	p.Service.LoadData(ipdb)
	return p
}
func (p *Handler) Ping() (err error) {
	return nil
}
func (p *Handler) IpToLocation(ip string) (r *ips.Location, err error) {
	log.Printf("IpToLocation: %s\n", ip)

	city, detail := p.Service.Ip2Address(ip)

	log.Printf("IpToLocation: City: %s,Detail: %s\n", city, detail)

	r = ips.NewLocation()
	r.City = city
	r.Detail = detail
	r.Province = getProvince(city)
	return r, nil
}

const (
	PROVINCE   = "省"
	ZIZHI_AREA = "自治区"
	SAR        = "特别行政区"
	CITY       = "市"
)

//
// 获取省份直辖市信息，传入的city为从省到市的字符，如浙江省杭州市，未找到匹配项则返回传入参数
//
func getProvince(city string) string {

	// 返回省
	idx := strings.Index(city, PROVINCE)
	if idx != -1 {
		return city[0 : idx+len(PROVINCE)]
	}

	// 返回自治区
	idx = strings.Index(city, ZIZHI_AREA)
	if idx != -1 {
		return city[0 : idx+len(ZIZHI_AREA)]
	}

	// 返回特别行政区
	idx = strings.Index(city, ZIZHI_AREA)
	if idx != -1 {
		return city[0 : idx+len(ZIZHI_AREA)]
	}

	// 返回特别行政区
	idx = strings.Index(city, SAR)
	if idx != -1 {
		return city[0 : idx+len(SAR)]
	}

	idx = strings.Index(city, CITY)
	if idx != -1 {
		return city[0 : idx+len(CITY)]
	}

	return city
}
