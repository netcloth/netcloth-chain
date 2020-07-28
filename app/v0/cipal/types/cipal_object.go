package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/netcloth/netcloth-chain/codec"
)

// CIPALObject defines the struct of cipal
type CIPALObject struct {
	UserAddress  string        `json:"user_address" yaml:"user_address"`
	ServiceInfos []ServiceInfo `json:"service_infos" yaml:"service_infos"`
}

type CIPALObjects []CIPALObject

// NewCIPALObject creates a new cipal object
func NewCIPALObject(userAddress string, serviceAddress string, serviceType uint64) CIPALObject {
	si := ServiceInfo{serviceType, serviceAddress}
	sis := make([]ServiceInfo, 0)
	sis = append(sis, si)
	return CIPALObject{
		UserAddress:  userAddress,
		ServiceInfos: sis,
	}
}

// MarshalYAML returns the YAML representation of an account.
func (obj CIPALObject) MarshalYAML() (interface{}, error) {
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

func MustMarshalCIPALObject(cdc *codec.Codec, obj CIPALObject) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(obj)
}

func MustUnmarshalCIPALObject(cdc *codec.Codec, value []byte) CIPALObject {
	validator, err := UnmarshalIPALObject(cdc, value)
	if err != nil {
		panic(err)
	}
	return validator
}

func UnmarshalIPALObject(cdc *codec.Codec, value []byte) (obj CIPALObject, err error) {
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

func (obj CIPALObject) String() string {
	return fmt.Sprintf(`CIPALObject
User Address:			%s
Service Infos:		    %s`, obj.UserAddress, getServiceInfosString(obj.ServiceInfos))
}
