package myutil

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"os"

	"github.com/seehuhn/mt19937"
)

// cacheデータを保存するディレクトリを作成
func MakeCacheDirectry(cache *bool) {
	if !*cache {
		return
	}

	if !FileExist("cache") {
		_ = os.Mkdir("cache", 0755)
	}

	if !FileExist("cache/characters") {
		_ = os.Mkdir("cache/characters", 0755)
	}

	if !FileExist("cache/languages") {
		_ = os.Mkdir("cache/languages", 0755)
	}
}

// ファイルの存在チェック
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ランダムなdim個の要素を持つ，[]uint8を返す
func Random(dim int) []uint8 {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64)) // セキュアな乱数生成器
	rng := rand.New(mt19937.New())                                // メルセンヌツイスタ
	rng.Seed(seed.Int64())                                        // seed設定

	tmp := make([]uint8, dim)
	for i := 0; i < dim; i++ {
		tmp[i] = uint8(rng.Intn(2))
	}

	return tmp
}
