package ip_query

import (
	"encoding/binary"
	"fmt"
	"github.com/qiniu/iconv"
	"io/ioutil"
	"os"
)

type IpInfoServiceV1 struct {
	IpIndexes []*ipIndex

	cd   *iconv.Iconv
	fbuf []byte
}

const (
	EMPTY_STR = ""
)

func (p *IpInfoServiceV1) Ip2Address(ip string) (country string, city string) {

	intIP := inet_aton(ip)
	if intIP == 0 {
		return EMPTY_STR, EMPTY_STR
	}

	start := 0
	end := len(p.IpIndexes) - 1
	var mid int
	var result int = -1

	for start <= end {
		mid = (start + end) / 2
		if intIP < p.IpIndexes[mid].Ip {
			end = mid - 1
		} else if intIP == p.IpIndexes[mid].Ip {
			result = mid
			break
		} else {
			// intIP > p.IpIndexes[mid].Ip
			start = mid + 1
		}
	}
	if result == -1 {
		// start > end
		if end < 0 {
			//			fmt.Println("Result = 0")
			result = 0
		} else {
			//			fmt.Printf("Result = %d\n", end)
			//			fmt.Printf("Ip: %d, L: %d, U: %d\n", intIP, p.IpIndexes[end].Ip, p.IpIndexes[start].Ip)
			result = end
		}
	}

	// 最终的结果：
	// IP[end] <= intIP
	if result >= 0 && result < len(p.IpIndexes) {
		index := p.IpIndexes[result]
		return p.getAddr(index.Offset + 4)
	} else {
		return EMPTY_STR, EMPTY_STR
	}
}

func (p *IpInfoServiceV1) LoadData(filename string) error {

	fid, err := os.Open(filename)
	if err != nil {
		return err
	}

	p.fbuf, err = ioutil.ReadAll(fid)
	fid.Close()
	if err != nil {
		return err
	}

	var indexStart uint32
	var indexEnd uint32

	// 读取索引的起止位置
	indexStart = binary.LittleEndian.Uint32(p.fbuf[0:4])
	indexEnd = binary.LittleEndian.Uint32(p.fbuf[4:8])

	fmt.Printf("Index Start: %d, End: %d\n", indexStart, indexEnd)

	itemNum := (indexEnd - indexStart) / 7
	p.IpIndexes = make([]*ipIndex, itemNum, itemNum)
	//	p.IpRecords = make([]*IpRecord, itemNum, itemNum)

	var i uint32
	offset := indexStart
	for i = 0; i < itemNum; i++ {
		offset += 7
		index := &ipIndex{}
		index.Ip = binary.LittleEndian.Uint32(p.fbuf[offset:(offset + 4)])
		index.Offset = byte3ToUint32(p.fbuf[(offset + 4):(offset + 7)])

		p.IpIndexes[i] = index
	}

	fmt.Println("Index Decoding Succeed")

	// https://github.com/qiniu/iconv
	cd, err := iconv.Open("utf-8", "gbk") // 从GBK转换成为utf8
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return err
	}
	p.cd = &cd

	fmt.Printf("Data Load Complete: %d items: ", itemNum)
	return nil

}

func (p *IpInfoServiceV1) Close() {
	if p.cd != nil {
		p.cd.Close()
		p.cd = nil
	}
}

// offset: 为country, city信息所在位置的offset, 跳过了之前的IP
func (p *IpInfoServiceV1) getAddr(offset uint32) (country string, city string) {

	var idx int
	order := p.fbuf[offset]

	if order == 1 {
		offset += 1
		newOffset := byte3ToUint32(p.fbuf[offset:len(p.fbuf)])
		return p.getAddr(newOffset)

	} else if order == 2 {
		offset += 1
		country = p.getAreaAddr(byte3ToUint32(p.fbuf[offset : offset+3]))
		offset += 3
		city = p.getAreaAddr(offset)

	} else {
		country, idx = p.getString(p.fbuf[offset:len(p.fbuf)])
		city, _ = p.getString(p.fbuf[offset+uint32(idx) : len(p.fbuf)])
	}
	return
}

func (p *IpInfoServiceV1) getAreaAddr(offset uint32) string {
	order := p.fbuf[offset]
	if order == 1 || order == 2 {
		offset += 1
		offset = byte3ToUint32(p.fbuf[offset : offset+3])
		return p.getAreaAddr(offset)
	} else {
		result, _ := p.getString(p.fbuf[offset:len(p.fbuf)])
		return result
	}
}

func (p *IpInfoServiceV1) getString(slice []byte) (term string, idx int) {

	for idx < len(slice) {
		if slice[idx] == 0 {
			break
		} else {
			idx += 1
		}
	}

	str := slice[0:idx]
	termBuff := make([]byte, len(str)*3)
	out, _, _ := p.cd.Conv(str, termBuff)
	return string(out), idx + 1
}

//func byte3ToUint32(offset []byte) uint32 {
//	return uint32(offset[0]) + (uint32(offset[1]) << 8) + (uint32(offset[2]) << 16)
//}

//func inet_ntoa(ipnr int64) net.IP {
//	var bytes [4]byte
//	bytes[0] = byte(ipnr & 0xFF)
//	bytes[1] = byte((ipnr >> 8) & 0xFF)
//	bytes[2] = byte((ipnr >> 16) & 0xFF)
//	bytes[3] = byte((ipnr >> 24) & 0xFF)

//	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
//}

//func inet_aton(ipnr string) uint32 {
//	bits := strings.Split(ipnr, ".")

//	b0, _ := strconv.Atoi(bits[0])
//	b1, _ := strconv.Atoi(bits[1])
//	b2, _ := strconv.Atoi(bits[2])
//	b3, _ := strconv.Atoi(bits[3])

//	var sum uint32

//	sum += uint32(b0) << 24
//	sum += uint32(b1) << 16
//	sum += uint32(b2) << 8
//	sum += uint32(b3)

//	return sum
//}
