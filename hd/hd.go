package hd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuya-isaka/langHD/hv"
	"github.com/yuya-isaka/langHD/myutil"
)

// HDC本体
type langHD struct {
	self   *langHD
	asciis map[int]*hv.HyperVector
	dim    int
	ngram  int
	langs  map[string]*hv.HyperVector
	tests  map[string]*hv.HyperVector
}

// コンストラクタ
func NewLangHD(dim *int, ngram *int) *langHD {
	asciis := make(map[int]*hv.HyperVector, 127)
	result := &langHD{
		asciis: asciis,
		dim:    *dim,
		ngram:  *ngram,
	}
	result.self = result
	return result
}

// ASCII127文字をハイパーベクトルとして本体に保存
func (h *langHD) EncodeAsciis(cache *bool) {
	// interfaceで共通化できる？
	if h != h.self {
		panic("should not copy without Copy()")
	}

	for i := 1; i <= 127; i++ {
		hvector := hv.NewHyperVector(h.dim)
		filePath := fmt.Sprintf("cache/characters/%s", strconv.Itoa(i))
		if *cache {
			if !myutil.FileExist(filePath) {
				h.asciis[i-1] = hvector.Generate()
				hvector.WriteCache(filePath)
				continue
			}
			h.asciis[i-1] = hvector.GenerateFromFile(filePath)
			continue
		}
		h.asciis[i-1] = hvector.Generate()
	}
}

func (h *langHD) EncodeTrainingData(cache *bool, train *string) {
	if h != h.self {
		panic("should not copy without Copy()")
	}

	h.langs = make(map[string]*hv.HyperVector)

	err := filepath.Walk(*train, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}

		filePath := fmt.Sprintf("cache/languages/%s", info.Name())
		if *cache {
			if myutil.FileExist(filePath) {
				hvector := hv.NewHyperVector(h.dim)
				h.langs[info.Name()] = hvector.GenerateFromFile(filePath)
				return nil
			}
		}

		tmp, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		text := string(tmp)

		h.makeTextHypervector(info.Name(), text, h.ngram, false)

		if *cache {
			h.langs[info.Name()].WriteCache(filePath)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (h *langHD) EncodeTestingData(test *string) {
	if h != h.self {
		panic("should not copy without Copy()")
	}

	h.tests = make(map[string]*hv.HyperVector)

	err := filepath.Walk(*test, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}

		tmp, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		text := string(tmp)

		h.makeTextHypervector(info.Name(), text, h.ngram, true)

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (h *langHD) makeTextHypervector(name string, text string, ngram int, flag bool) {
	if h != h.self {
		panic("should not copy without Copy()")
	}

	if ngram < 1 {
		panic("cannot encode")
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+") // a-zA-Z0-9以外の数は取り込まない
	if err != nil {
		panic(err)
	}
	text = reg.ReplaceAllString(text, "") // ひとつづきのtextを入手"MynameisYuyaIsaka...""
	text = strings.ToLower(text)          // "mynameisyuyaisaka..."

	if len(text) == 0 {
		log.Printf("nothing text")
		tmp := hv.NewHyperVector(h.dim)
		tmp.NotMuch()
		if flag {
			h.tests[name] = tmp
		} else {
			h.langs[name] = tmp
		}
	}

	numNgram := len(text) - ngram + 1
	if numNgram < 1 {
		log.Printf("nothing text")
		tmp := hv.NewHyperVector(h.dim)
		tmp.NotMuch()
		if flag {
			h.tests[name] = tmp
		} else {
			h.langs[name] = tmp
		}
	}

	textHyperVectors := make([]*hv.HyperVector, numNgram)
	for i := range text {
		if len(text)-ngram < i {
			break
		}

		asciiTexts := make([]int, ngram)
		for j := range asciiTexts {
			asciiTexts[j] = int(text[i+j])
		}

		var result *hv.HyperVector
		for index, asciiNum := range asciiTexts {
			if index == 0 {
				tmp := h.asciis[asciiNum]
				result = tmp.Rotate(len(asciiTexts)-1, h.dim)
				continue
			}

			next := h.asciis[asciiNum]
			if len(asciiTexts)-index-1 != 0 {
				next = next.Rotate(len(asciiTexts)-index-1, h.dim) // 新しくnextに割り当てた場合は，元nextの内容(h.asciis[asciiNum]は変更されない)
			}

			result.Xor(next)
		}

		textHyperVectors[i] = result
	}

	if len(textHyperVectors)%2 == 0 {
		tmp := hv.NewHyperVector(h.dim)
		textHyperVectors = append(textHyperVectors, tmp.Generate())
	}

	if flag {
		h.tests[name] = hv.NewHyperVector(h.dim)
		h.tests[name].Add(textHyperVectors...)
	} else {
		h.langs[name] = hv.NewHyperVector(h.dim)
		h.langs[name].Add(textHyperVectors...)
	}
}

func (h *langHD) Testing() {
	for testName, testVec := range h.tests {
		var match string
		var maxCosine float32 = -2
		for trainName, trainVec := range h.langs {
			cosine := testVec.Cosine(trainVec)
			if cosine > maxCosine {
				maxCosine = cosine
				match = trainName
			}
		}
		if match == "" {
			fmt.Println("could not find match language")
		} else {
			fmt.Println(testName + ": language is " + match)
		}
	}
}
