package AossGoSdk

import (
	"encoding/base64"
)

func decode(data string) (b []byte, err error) {
	ret, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return ret, err
}
