package db

import (
	"errors"
	"fmt"
)

// Nsrjcxx  纳税人基础信息
type Nsrjcxx struct {
	Xm       string `mysql:"xm"`
	Zzlx     string `mysql:"zzlx"`
	Zzhm     string `mysql:"zzhm"`
	Gj       string `mysql:"gj"`
	Xb       string `mysql:"xb"`
	Nsrzt    string `mysql:"nsrzt"`
	Sfgd     string `mysql:"sfgd"`
	Sfgy     string `mysql:"sfgy"`
	Lxdh     string `mysql:"lxdh"`
	Sftdhy   string `mysql:"sftdhy"`
	Sftstzr  string `mysql:"sftstzr"`
	Bmbh     string `mysql:"BMBH"`
	Qynsrsbh string `mysql:"QYNSRSBH"`
}

// ParseNsrjcxxFile parse text from file to userdomain PrivateQuestion
func ParseNsrjcxxFile(fileName string, delimeter string) ([]*Nsrjcxx, error) {

	var ret []*Nsrjcxx
	values, err := ParseFile(fileName, delimeter)
	if err != nil {
		return nil, err
	}
	for _, v := range values {
		nsr, err := newNsrjcxx(v)
		if err != nil {
			fmt.Printf("wrong file data :%v\n", err)
			continue
		}
		ret = append(ret, nsr)
	}
	return ret, nil
}

func newNsrjcxx(values []string) (*Nsrjcxx, error) {
	var mobile string
	if len(values) != 6 {
		fmt.Printf("bad format data, len:%d, data:%v\n", len(values), values)
		return nil, errors.New("bad format")
	}
	if values[4] == "86" {
		mobile = values[5]
	}

	return &Nsrjcxx{
		Xm:       values[0],
		Zzhm:     values[1],
		Zzlx:     idCardMap[values[2]],
		Xb:       genderMap[values[3]],
		Gj:       "中国",
		Nsrzt:    "正常",
		Sfgy:     "否",
		Sfgd:     "否",
		Sftdhy:   "否",
		Sftstzr:  "否",
		Lxdh:     mobile,
		Qynsrsbh: CompanySbh,
	}, nil

}
