package ip_query

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/qiniu/iconv"
	"github.com/wfxiang08/rpc_proxy/utils/log"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

type ipIndex struct {
	Ip     uint32
	Offset uint32
}

type IpRecord struct {
	Ip     uint32
	City   string
	Detail string
}

type IpInfoService struct {
	IpIndexes []*ipIndex
	IpRecords []*IpRecord

	termBuff []byte
	cd       *iconv.Iconv
	fbuf     []byte
}

func (p *IpInfoService) Ip2Address(ip string) (city string, detail string) {

	intIP := inet_aton(ip)

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
			result = 0
		} else {
			result = end
		}
	}

	log.Printf("Binary Search End: %d %d\n", len(p.IpIndexes), start)
	// 最终的结果：
	// IP[end] <= mid
	return p.IpRecords[start].City, p.IpRecords[start].Detail
}

func (p *IpInfoService) LoadData(filename string) error {

	fid, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("File Not Found: %s", filename))
		return err
	}

	p.fbuf, err = ioutil.ReadAll(bufio.NewReaderSize(fid, 7*1024*1024))
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
	p.IpRecords = make([]*IpRecord, itemNum, itemNum)

	var i uint32
	offset := indexStart
	for i = 0; i < itemNum; i++ {
		offset += 7
		index := &ipIndex{}
		index.Ip = binary.LittleEndian.Uint32(p.fbuf[offset:(offset + 4)])
		index.Offset = byte3ToUint32(p.fbuf[(offset + 4):(offset + 7)])

		p.IpIndexes[i] = index
	}

	log.Println("Index Decoding Succeed")

	// https://github.com/qiniu/iconv
	cd, err := iconv.Open("utf-8", "gbk") // 从GBK转换成为utf8
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return err
	}
	p.cd = &cd

	p.termBuff = make([]byte, 2000)
	for i = 0; i < itemNum; i++ {
		index := p.IpIndexes[i]

		recordStart := index.Offset + 4
		country, city := p.getAddr(recordStart)

		//		fmt.Println("Country: ", country, "City: ", city)

		record := &IpRecord{
			Ip:     index.Ip,
			City:   country,
			Detail: city,
		}
		p.IpRecords[i] = record
	}

	// 清空缓存
	p.termBuff = nil
	p.fbuf = nil
	p.cd.Close()
	p.cd = nil

	log.Printf("Data Load Complete: %d items: ", itemNum)
	return nil

}

// offset: 为country, city信息所在位置的offset, 跳过了之前的IP
func (p *IpInfoService) getAddr(offset uint32) (country string, city string) {

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

func (p *IpInfoService) getAreaAddr(offset uint32) string {
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

func (p *IpInfoService) getString(slice []byte) (term string, idx int) {

	for idx < len(slice) {
		if slice[idx] == 0 {
			break
		} else {
			idx += 1
		}
	}

	str := slice[0:idx]
	out, _, _ := p.cd.Conv(str, p.termBuff)
	return string(out), idx + 1
}

func byte3ToUint32(offset []byte) uint32 {
	return uint32(offset[0]) + (uint32(offset[1]) << 8) + (uint32(offset[2]) << 16)
}

func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func inet_aton(ipnr string) uint32 {
	bits := strings.Split(ipnr, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}
