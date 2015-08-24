package ip_query

import (
	"fmt"
	"testing"
)

func TestLdapAccount(t *testing.T) {

	filename := "qqwry.dat"
	service := &IpInfoService{}
	_ = service.LoadData(filename)

	testcases := []struct {
		Ip     string
		Result string
	}{
		{Ip: "60.29.255.197",
			Result: "天津市",
		},
		{
			Ip:     "74.125.235.211",
			Result: "美国",
		},
	}
	for _, testcase := range testcases {
		city, _ := service.Ip2Address(testcase.Ip)
		if city != testcase.Result {
			fmt.Println("Invalid IP")
			t.Fail()
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
