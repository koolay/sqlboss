package utils

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/oklog/ulid"
)

var (
	snowflakeNode *snowflake.Node
	src           = rand.NewSource(time.Now().UnixNano())
)

func init() {
	var err error
	nodeID, err := getIntPrivateIP()
	if err != nil {
		log.Fatal(err)
	}

	snowflakeNode, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	digitBytes  = "0123456789"
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func getIntPrivateIP() (int64, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}

	return int64(ip[2]) + int64(ip[3]), nil
}

// RandString 产生定长的随机字符串
// reference https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandDigit 产生定长的随机数字字符串
func RandDigit(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(digitBytes) {
			b[i] = digitBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func CombineOSSURL(endpoint, bucket string) (string, error) {
	URL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	if URL.Host == "" || URL.Scheme == "" {
		return "", fmt.Errorf("invalid endpoint: %s", endpoint)
	}
	return fmt.Sprintf("%s://%s.%s", URL.Scheme, bucket, URL.Host), nil
}

// ShortHash 生成 int32 hash
// 使用FNV算法
func ShortHash(source string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(source))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d", h.Sum32())
}

func SnowflakeID() string {
	return snowflakeNode.Generate().String()
}

func UUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func NewMysqlID() string {
	return uuid.Must(uuid.NewV1()).String()
}

// Ulid sortable, more shorter
func Ulid() string {

	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	entropy := rand.New(source)

	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

// Sha1 hash with sha1
func Sha1(source []byte) ([]byte, error) {

	// The pattern for generating a hash is `sha1.New()`,
	// `sha1.Write(bytes)`, then `sha1.Sum([]byte{})`.
	// Here we start with a new hash.
	h := sha1.New()

	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	_, err := h.Write(source)
	if err != nil {
		return nil, err
	}

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	return h.Sum(nil), nil
}

// JSONBytesEqual 比较两个json bytes是否相等
// 对于数组会忽略顺序, 数字和字符串不相等: 1!="1"
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}

	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}

	return cmp.Equal(j, j2, cmpopts.SortSlices(func(x, y interface{}) bool {
		return strings.Compare(fmt.Sprintf("%v", x), fmt.Sprintf("%v", y)) > 0
	})), nil
}
