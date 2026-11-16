package AossGoSdk

import (
	"encoding/hex"
	"testing"
)

func TestV3EnvelopeAgainstPHPVector(t *testing.T) {
	secret := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	// PHP 服务端参考值
	if got := V3Kid(secret); got != "a8ae6e6e" {
		t.Fatalf("kid = %s, want a8ae6e6e", got)
	}
	wantKey := "7e0f60893205dccc27dd9bced0766dc945fc54af80de401c88ffb4550695fd39"
	if got := hex.EncodeToString(v3DeriveKey(secret)); got != wantKey {
		t.Fatalf("key = %s, want %s", got, wantKey)
	}

	// PHP seal 产生的信封，Go 不可直接解密（本 SDK 只负责 seal），
	// 但可用 Go seal 后交由 PHP open 验证。此处断言 seal 的结构正确。
	auth := V3Seal(secret, 1770000000)
	if len(auth) < 4 || auth[:3] != "v1." {
		t.Fatalf("seal format invalid: %s", auth)
	}
}
