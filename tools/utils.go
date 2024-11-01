package tools

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func GetSubnetFromMask(mask int) [4]byte {
	a := mask / 8
	b := float64(mask % 8) // math.Pow only takes float64 has arguments
	idx := 0

	var ans [4]byte
	for i := 0; i < a; i++ {
		ans[idx] = uint8(math.Pow(2, 8))
		idx++
	}
	if idx < 4 && b > 0 {
		ans[idx] = uint8(math.Pow(2, b))
	}

	return ans
}

func RandomMacAddr() ([]byte, error) {
	bytes := make([]byte, 6)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("Coudn't generate Mac Addr")
	}
	return bytes, nil
}

func ConvertStrToIp(addr string) [4]byte {
	bytes := strings.Split(addr, ".")
	var ans [4]byte
	for i, val := range bytes {
		num, _ := strconv.Atoi(val)
		ans[i] = uint8(num)
	}

	return ans
}

func ConvertAddrToStr(addr []byte) string {
	var ans string
	if len(addr) == 4 {
		for _, val := range addr {
			ans += strconv.Itoa(int(val)) + "."
		}
	} else {
		for _, val := range addr {
			ans += strconv.FormatInt(int64(val), 16) + ":"
		}
	}
	return ans[:len(ans)-1]
}

func StructToByte[T any](data T) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("Can't convert struct to byte array")
	}

	return buf.Bytes(), nil
}

func ByteToStruct[T any](data []byte, tmp T) (*T, error) {
	var ans T
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&ans); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Can't decode strcut from byte array")
	}

	return &ans, nil
}
