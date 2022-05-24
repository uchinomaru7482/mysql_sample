package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
)

func GetLowerCaseRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	return makeRandomStr(digit, letters)
}

func makeRandomStr(digit uint32, letters string) (string, error) {
	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("Get random string error")
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

// ランダムな数値を生成する
func GetRandomInt(max int64) (int64, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, errors.New("Get random int error")
	}
	convNum := num.Int64() + 1
	return convNum, nil
}
