package xray

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type XrayLink struct {
	Protocol   string
	User       string
	Host       string
	Port       int
	Parameters xrayLinkParameters
	LinkName   string
}

type xrayLinkParameters struct {
	Security    string `url:"security,omitempty"`
	SNI         string `url:"sni,omitempty"`
	Fingerprint string `url:"fp,omitempty"`
	PublicKey   string `url:"pbk,omitempty"`
	ShortID     string `url:"sid,omitempty"`
	Type        string `url:"type,omitempty"`
}

func NewXrayLink(protocol, user, host string, port int, linkname string) *XrayLink {
	newXL := XrayLink{
		Protocol:   protocol,
		User:       user,
		Host:       host,
		Port:       port,
		LinkName:   linkname,
		Parameters: xrayLinkParameters{},
	}
	return &newXL
}

func (xl *XrayLink) MarshalLink() string {
	if xl.Host == "" || xl.Port == 0 {
		panic(errors.New("tried marshaling an xray link without a host|port"))
	}
	u := &url.URL{
		Scheme: xl.Protocol,
		Host:   fmt.Sprintf("%s:%d", xl.Host, xl.Port),
	}
	u.User = url.User(xl.User)

	v := url.Values{}
	val := reflect.ValueOf(xl.Parameters)
	typ := reflect.TypeOf(xl.Parameters)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("url")
		if tag == "" {
			tag = fieldType.Name
		}

		parts := strings.Split(tag, ",")
		key := ""
		omitEmpty := false
		for _, part := range parts {
			if part == "omitempty" {
				omitEmpty = true
			} else if key == "" {
				key = part
			}
		}

		if key == "" {
			key = fieldType.Name
		}

		if omitEmpty && field.IsZero() {
			continue
		}

		v.Set(key, fmt.Sprintf("%v", field.Interface()))
	}
	u.RawQuery = v.Encode()

	return u.String() + "#" + xl.LinkName
}
