package encrypt

import (
	"go-admin/common/core/sdk/pkg"
	"golang.org/x/crypto/bcrypt"
)

//
//  HashEncrypt
//  @Description: 将字符串单向加密
//  @param value
//  @return string
//  @return error
//
func HashEncrypt(value string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//
//  HashCmpPassword
//  @Description: 比较密码
//  @param pwd1Hash 已经hash的密码
//  @param pwd2  未hash的密码
//  @return bool
//  @return error
//
func HashCmpPassword(pwd1Hash, pwd2 string) (bool, error) {
	_, err := pkg.CompareHashAndPassword(pwd1Hash, pwd2)
	if err != nil {
		return false, err
	}
	return true, nil
}
