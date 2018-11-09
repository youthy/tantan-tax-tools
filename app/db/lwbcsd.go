package db

import (
	"errors"
	"fmt"
	"math/big"
)

// Lwbcsd 劳务报酬所得
type Lwbcsd struct {
	XM       string `mysql:"XM"`
	ZZLX     string `mysql:"ZZLX"`
	ZZHM     string `mysql:"ZZHM"`
	SRE      string `mysql:"SRE"`
	NSRSBH   string `mysql:"NSRSBH"`
	QYNSRSBH string `mysql:"QYNSRSBH"`
}

// ParseLwbcsdFile parse text from file to userdomain PrivateQuestion
func ParseLwbcsdFile(fileName string, delimeter string) ([]*Lwbcsd, error) {

	var ret []*Lwbcsd
	values, err := ParseFile(fileName, delimeter)
	if err != nil {
		return nil, err
	}
	for _, v := range values {
		lwb, err := newLwbcsd(v)
		if err != nil {
			continue
		}
		ret = append(ret, lwb)
	}
	return ret, nil
}

func newLwbcsd(values []string) (*Lwbcsd, error) {
	if len(values) != 4 {
		fmt.Printf("bad format data, len:%d, data:%v\n", len(values), values)
		return nil, errors.New("bad format")
	}
	oriSre := new(big.Float)
	_, err := oriSre.SetString(values[3])
	if !err {
		return nil, errors.New("parse sre failed")
	}
	div := big.NewFloat(1000000)
	sre := new(big.Float)
	sre.Quo(oriSre, div)
	return &Lwbcsd{
		XM:       values[0],
		ZZHM:     values[1],
		ZZLX:     idCardMap[values[2]],
		SRE:      fmt.Sprintf("%.2f", sre),
		QYNSRSBH: CompanySbh,
	}, nil
}
