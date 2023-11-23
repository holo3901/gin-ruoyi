package pinyin

import (
	"github.com/mozillazg/go-pinyin"
)

func PinYin(a string) [][]string {

	// 默认
	x := pinyin.NewArgs()

	// 包含声调
	x.Style = pinyin.Tone
	return pinyin.Pinyin(a, x)
	// [[zhōng] [guó] [rén]]

}
