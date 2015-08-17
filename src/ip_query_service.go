package main

import (
	"fmt"
	ip_query "git.chunyu.me/infra/ip_utils/src/ip_query"
	"time"
)

//  		 	    IpInfoService   vs. IpInfoServiceV1
//    内存: 		     67M					28M
//    启动时间:     5.5206s			   0.3227s
// 单次处理时间:     0.0367s			   0.1508s * 10-4
//  对应的Python实现: 22M内存             0.30236s
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

	t1 = time.Now()
	for i := 0; i < 10000; i++ {
		_, _ = service.Ip2Address("60.29.255.197")
	}
	t2 = time.Now()
	fmt.Printf("Elapsed: %.4fms\n", float64(t2.UnixNano()-t1.UnixNano())*10e-10)

	country, city := service.Ip2Address("60.29.255.197")
	fmt.Println("Country: ", country, ", City: ", city)

	//city, _detail = ip_utils.get_default_ip_info().getIPAddr('60.29.255.197')
	//        self.assertEqual(u'天津市', ensure_unicode(city))

	serviceV1 := ip_query.IpInfoServiceV1{}
	t1 = time.Now()
	err = serviceV1.LoadData(filename)
	t2 = time.Now()

	// 10e-9 * 10e-4 * 1e3
	fmt.Printf("Elapsed: %.4fs\n", float64(t2.UnixNano()-t1.UnixNano())*10e-9)

	if err != nil {
		panic(fmt.Sprintf("File Open Error: %s", filename))
	}

	t1 = time.Now()
	for i := 0; i < 10000; i++ {
		_, _ = serviceV1.Ip2Address("60.29.255.197")
	}
	t2 = time.Now()
	fmt.Printf("Elapsed: %.4fms\n", float64(t2.UnixNano()-t1.UnixNano())*10e-10)

	country, city = serviceV1.Ip2Address("60.29.255.197")
	fmt.Println("Country: ", country, ", City: ", city)

}
