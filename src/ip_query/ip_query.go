package ip_query

import (
	"encoding/binary"
	"fmt"
	"github.com/qiniu/iconv"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

type IpIndex struct {
	Ip     uint32
	Offset uint32
}

type IpRecord struct {
	Ip          uint32
	AreaCountry string
	AreaCity    string
}

type IpInfoService struct {
	IpIndexes []*IpIndex
	IpRecords []*IpRecord

	termBuff []byte
	cd       *iconv.Iconv
	fbuf     []byte
}

func (p *IpInfoService) Ip2Address(ip string) (country string, city string) {
	var mid int
	ipInt := inet_aton(ip)

	start := 0
	end := len(p.IpIndexes) - 1

	for end-start > 1 {
		mid = (start + end) / 2
		if ipInt <= p.IpIndexes[mid].Ip {
			start = mid
		} else {
			end = mid
		}
	}
	index := p.IpIndexes[start]
	offset := index.Offset + 4
	return p.getAddr(offset)
}

func (p *IpInfoService) LoadData(filename string) error {

	fid, err := os.Open(filename)
	if err != nil {
		return err
		//		panic(fmt.Sprintf("File Open Error: %s", filename))
	}
	defer fid.Close()

	p.fbuf, err = ioutil.ReadAll(fid)

	var indexStart uint32
	var indexEnd uint32

	indexStart = binary.LittleEndian.Uint32(p.fbuf[0:4])
	fmt.Printf("Index Start: %d\n", indexStart)

	indexEnd = binary.LittleEndian.Uint32(p.fbuf[4:8])

	fmt.Printf("Index End: %d\n", indexEnd)

	ipStr := "192.168.0.1"
	fmt.Printf("IP: %s --> %d\n", ipStr, inet_aton(ipStr))

	// 读取到了 indexStart, indexEnd

	itemNum := (indexEnd - indexStart) / 7
	p.IpIndexes = make([]*IpIndex, itemNum, itemNum)
	p.IpRecords = make([]*IpRecord, itemNum, itemNum)

	var i uint32
	offset := indexStart
	for i = 0; i < itemNum; i++ {
		offset += 7
		index := &IpIndex{}
		index.Ip = binary.LittleEndian.Uint32(p.fbuf[offset:(offset + 4)])
		index.Offset = byte3ToUint32(p.fbuf[(offset + 4):(offset + 7)])

		p.IpIndexes[i] = index
	}

	// https://github.com/qiniu/iconv
	cd, err := iconv.Open("utf-8", "gbk") // convert gbk to utf-8
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return err
	}
	p.cd = &cd

	var country string
	var city string

	p.termBuff = make([]byte, 2000)

	for i = 0; i < itemNum; i++ {
		index := p.IpIndexes[i]
		recordStart := index.Offset + 4
		country, city = p.getAddr(recordStart)
		record := &IpRecord{
			Ip:          index.Ip,
			AreaCountry: country,
			AreaCity:    city,
		}
		p.IpRecords[i] = record
	}

	p.termBuff = nil
	p.fbuf = nil
	p.cd.Close()
	p.cd = nil

	return nil

}

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
