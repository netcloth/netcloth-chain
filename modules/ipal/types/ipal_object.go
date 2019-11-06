package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/NetCloth/netcloth-chain/codec"
)

type IPALObject struct {
	UserAddress  string        `json:"user_address" yaml:"user_address"`
	ServiceInfos []ServiceInfo `json:"service_infos" yaml:"service_infos`
}

func NewIPALObject(userAddress string, serviceAddress string, serviceType uint64) IPALObject {
	si := ServiceInfo{serviceType, serviceAddress}
	sis := make([]ServiceInfo, 0)
	sis = append(sis, si)
	return IPALObject{
		UserAddress:  userAddress,
		ServiceInfos: sis,
	}
}

func (obj IPALObject) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(struct {
		UserAddress  string
		ServiceInfos []ServiceInfo
	}{
		UserAddress:  obj.UserAddress,
		ServiceInfos: obj.ServiceInfos,
	})

	if err != nil {
		return nil, err
	}
	return string(bs), nil
}

func MustMarshalIPALObject(cdc *codec.Codec, obj IPALObject) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(obj)
}

func MustUnmarshalIPALObject(cdc *codec.Codec, value []byte) IPALObject {
	validator, err := UnmarshalIPALObject(cdc, value)
	if err != nil {
		panic(err)
	}
	return validator
}

func UnmarshalIPALObject(cdc *codec.Codec, value []byte) (obj IPALObject, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &obj)
	return obj, err
}

func getServiceInfosString(infos []ServiceInfo) string {
	var s string
	for _, v := range infos {
		s = s + v.String() + "\n"
	}

	return s
}

func (obj IPALObject) String() string {
	return fmt.Sprintf(`IPALObject
User Address:			%s
Service Infos:		    %s`, getServiceInfosString(obj.ServiceInfos))
}
