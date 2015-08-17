package main

import (
	"fmt"
	ip_query "git.chunyu.me/infra/ip_utils/src/ip_query"
	"time"
)

func main() {
	filename := "/Users/feiwang/gowork/src/git.chunyu.me/infra/ip_utils/src/ip_query/qqwry.dat"

	service := ip_query.IpInfoService{}
	t1 := time.Now()
	err := service.LoadData(filename)
	t2 := time.Now()

	fmt.Printf("Elapsed: %.4fs\n", float64(t2.UnixNano()-t1.UnixNano())*10e-9)

	if err != nil {
		panic(fmt.Sprintf("File Open Error: %s", filename))
	}

	country, city := service.Ip2Address("60.29.255.197")

	fmt.Println("Country: ", country, ", City: ", city)

	//city, _detail = ip_utils.get_default_ip_info().getIPAddr('60.29.255.197')
	//        self.assertEqual(u'天津市', ensure_unicode(city))

}
