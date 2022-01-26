package hv

import (
	"encoding/json"
	"io/ioutil"
	"math"

	"github.com/yuya-isaka/langHD/myutil"
)

// ハイパーベクトル
type HyperVector struct {
	self   *HyperVector
	values []uint8
	length int
	norm   float32
}

// コンストラクタ
func NewHyperVector(dim int) *HyperVector {
	values := make([]uint8, dim)
	result := &HyperVector{
		values: values,
		length: dim,
		norm:   -1, // default value
	}
	result.self = result
	return result
}

// dim次元のランダムなハイパーベクトルを生成
func (v *HyperVector) Generate() *HyperVector {
	if v != v.self {
		panic("should not copy without Copy()")
	}

	v.values = myutil.Random(v.length)
	return v
}

// dim次元のランダムなハイパーベクトルをファイルから読み取り
func (v *HyperVector) GenerateFromFile(filePath string) *HyperVector {
	if v != v.self {
		panic("should not copy without Copy()")
	}

	tmp, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(tmp, &v.values)
	if err != nil {
		panic(err)
	}
	return v
}

// dim次元のランダムなハイパーベクトルを保存
func (v *HyperVector) WriteCache(filePath string) {
	if v != v.self {
		panic("should not copy without Copy()")
	}

	tmp, _ := json.Marshal(&v.values)
	_ = ioutil.WriteFile(filePath, tmp, 0666)
}

func (v *HyperVector) NotMuch() {
	if v != v.self {
		panic("should not copy without Copy()")
	}
	for i := 0; i < v.length; i++ {
		v.values[i] = 1
	}
}

func (v *HyperVector) Rotate(num int, dim int) *HyperVector {
	if v != v.self {
		panic("should not copy without Copy()")
	}
	vector := NewHyperVector(dim)
	tmp := v.length - num%v.length
	vector.values = append(v.values[tmp:], v.values[:tmp]...)
	return vector
}

func (v *HyperVector) Xor(v2 *HyperVector) {
	if v != v.self {
		panic("should not copy without Copy()")
	}
	for i := range v.values {
		v.values[i] = v.values[i] ^ v2.values[i]
	}
	v.norm = -1
}

func (v *HyperVector) Add(textHyperVectors ...*HyperVector) {
	if v != v.self {
		panic("should not copy without Copy()")
	}
	if len(textHyperVectors) == 0 {
		return
	}

	if len(textHyperVectors) == 1 {
		v.values = textHyperVectors[0].values
	}

	thr := len(textHyperVectors)/2 + 1

	for i := range v.values {
		sum := 0
		for _, vector := range textHyperVectors {
			sum += int(vector.values[i])
		}

		if sum >= thr {
			v.values[i] = 1
		} else {
			v.values[i] = 0
		}
	}
}

func (v *HyperVector) Cosine(v2 *HyperVector) float32 {
	return float32(v.dot(v2)) / (v.normCheck() * v2.normCheck())
}

func (v *HyperVector) dot(v2 *HyperVector) int {
	var result int
	for i := range v.values {
		result += int(v.values[i] & v2.values[i])
	}
	return result
}

func (v *HyperVector) normCheck() float32 {
	if v.norm != -1 {
		return v.norm
	}

	var sum int
	for _, value := range v.values {
		sum += int(value)
	}

	v.norm = float32(math.Sqrt(float64(sum)))
	return v.norm
}
