package AossGoSdk

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strconv"
	"time"
)

// V3 开放接口凭证信封（AES-256-GCM）
//
// 与 AOSSTP8 服务端 extend/Aoss/Envelope.php 严格对齐：
//
//	kid  = hex(sha256(secret))[:8]                     仅用于服务端定位记录
//	key  = HKDF-SHA256(ikm=secret, salt=zero[32], info="aoss.v3.envelope", 32)
//	auth = "v1." + base64url(iv[12] | ciphertext | tag[16])
//	plaintext = secret + "|" + ts, AAD = "v1"
//
// secret 本身即共享密钥，仅参与本地派生/加密，绝不明文传输。
const (
	v3EnvelopeVersion = "v1"
	v3KdfInfo         = "aoss.v3.envelope"
	v3IVLen           = 12
	v3TagLen          = 16
)

// V3Kid 由 secret 派生 8 位小写十六进制定位标识。
func V3Kid(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])[:8]
}

// v3DeriveKey 以 secret 为 IKM 派生 32 字节 AES-256 密钥。
//
// PHP hash_hkdf 在 salt 省略/空串时等价于全零 salt（RFC 5869），
// 这里显式用 32 个零字节，保证与服务端派生结果一致。
func v3DeriveKey(secret string) []byte {
	salt := make([]byte, sha256.Size) // 32 字节零值
	return hkdfSha256([]byte(secret), salt, []byte(v3KdfInfo), 32)
}

// hkdfSha256 实现 RFC 5869 的 HKDF（Extract + Expand），哈希为 SHA-256。
func hkdfSha256(ikm, salt, info []byte, length int) []byte {
	// Extract: PRK = HMAC(salt, IKM)
	mac := hmac.New(sha256.New, salt)
	mac.Write(ikm)
	prk := mac.Sum(nil)

	// Expand: OKM = T(1) | T(2) | ...，T(i) = HMAC(PRK, T(i-1) | info | i)
	out := make([]byte, 0, length)
	var prev []byte
	for i := byte(1); len(out) < length; i++ {
		block := hmac.New(sha256.New, prk)
		block.Write(prev)
		block.Write(info)
		block.Write([]byte{i})
		prev = block.Sum(nil)
		out = append(out, prev...)
	}
	return out[:length]
}

// V3Seal 生成凭证信封，形如 v1.xxxxx；失败返回空串。
//
// ts 为秒级 Unix 时间戳。每次调用使用随机 IV，相同 secret 的密文逐次不同。
func V3Seal(secret string, ts int64) string {
	if secret == "" {
		return ""
	}
	key := v3DeriveKey(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	gcm, err := cipher.NewGCMWithTagSize(block, v3TagLen)
	if err != nil {
		return ""
	}

	iv := make([]byte, v3IVLen)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}
	plaintext := []byte(secret + "|" + decimal(ts))

	// gcm.Seal 输出为 ciphertext | tag，AAD 用版本号绑定，防跨版本重放。
	sealed := gcm.Seal(nil, iv, plaintext, []byte(v3EnvelopeVersion))

	raw := make([]byte, 0, len(iv)+len(sealed))
	raw = append(raw, iv...)
	raw = append(raw, sealed...)
	return v3EnvelopeVersion + "." + base64.RawURLEncoding.EncodeToString(raw)
}

// V3AuthHeaders 生成一次请求的 V3 鉴权请求头（每次动态生成时间戳与随机 IV）。
//
// 服务端约定：X-Aoss-Ts / X-Aoss-Kid / X-Aoss-Auth，
// 服务端虽兼容表单/查询串传参，本 SDK 一律用请求头方式。
func V3AuthHeaders(secret string) map[string]string {
	ts := time.Now().Unix()
	return map[string]string{
		"X-Aoss-Ts":   strconv.FormatInt(ts, 10),
		"X-Aoss-Kid":  V3Kid(secret),
		"X-Aoss-Auth": V3Seal(secret, ts),
	}
}

// decimal 将非负 int64 转为十进制字符串（避免额外依赖）。
func decimal(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
