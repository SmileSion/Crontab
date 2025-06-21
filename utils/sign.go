package utils

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/tjfoc/gmsm/sm3"
)

func GenerateSM3Sign(params map[string]string) string {
	var keys []string
	for k, v := range params {
		if k == "sign" || v == "" {
			continue
		}
		keys = append(keys, k)
	}
	fmt.Println("[Step1] 有效参数键名:", keys)
	sort.Strings(keys)
	fmt.Println("[Step2] 排序后的键名:", keys)

	var sb strings.Builder
	for i, k := range keys {
		sb.WriteString(k + "=" + params[k])
		if i < len(keys)-1 {
			sb.WriteString("&")
		}
	}
	stringA := sb.String()
	fmt.Println("[Step3] 待签名字符串 stringA:", stringA)

	digest := sm3.New()
	digest.Write([]byte(stringA))
	sign := digest.Sum(nil)

	result := strings.ToUpper(hex.EncodeToString(sign))
	fmt.Println("[Step4] 最终签名结果:", result)

	return result
}
