package main

import (
	"crypto/md5"
	"encoding/hex"
)

//MD5 Calculate md5 of a string
func MD5(input string) string {
    hash := md5.Sum([]byte(input))
    return hex.EncodeToString(hash[:])
}
