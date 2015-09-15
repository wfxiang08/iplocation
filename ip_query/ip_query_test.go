package ip_query

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"testing"
)

//
// go test git.chunyu.me/infra/iplocation/ip_query -v -run "TestIpLocation"
//

func TestIpLocation(t *testing.T) {

	filename := "qqwry.dat"
	service := &IpInfoService{}
	_ = service.LoadData(filename)

	testcases := []struct {
		Ip     string
		Result string
	}{
		{
			Ip:     "0.0.0.0",
			Result: "",
		},
		{
			Ip:     "119.189.65.58",
			Result: "山东省聊城市",
		},
		{
			Ip:     "60.29.255.197",
			Result: "天津市",
		},
		{
			Ip:     "74.125.235.211",
			Result: "美国",
		},

		{
			Ip:     "127.0.0.1",
			Result: "IANA",
		},
	}
	for _, testcase := range testcases {
		city, _ := service.Ip2Address(testcase.Ip)
		if city != testcase.Result {
			fmt.Printf("Invalid IP: %s --> %s, Exp: %s\n", testcase.Ip, city, testcase.Result)
			assert.Fail(t, "Error Ip Mapping")
		}
	}

	//city, _detail = ip_utils.get_default_ip_info().getIPAddr('60.29.255.197')
	//        self.assertEqual(u'天津市', ensure_unicode(city))

	//        city, _detail = ip_utils.get_default_ip_info().getIPAddr('74.125.235.211')
	//        self.assertEqual(u'美国', ensure_unicode(city))

	//        city, _detail = ip_utils.get_default_ip_info().getIPAddr('61.135.169.125')
	//        self.assertEqual(u'北京市', ensure_unicode(city))

	//        city, _detail = ip_utils.get_default_ip_info().getIPAddr('36.110.16.242')
	//        self.assertEqual(u'北京市', ensure_unicode(city))
}
