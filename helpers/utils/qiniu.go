package utils

import (
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/xifengzhu/eshop/helpers/setting"
)

type QiniuMeta struct {
	Token        string `json:"token"`
	BucketDomain string `json:"bucket_domain"`
	ExpiresIn    int    `json:"expires_in"`
	Bucket       string `json:"bucket"`
}

var (
	bucket    = setting.QiniuBucket
	domain    = setting.QiniuBucketDomain
	expiresIn = 600 // 10分钟有效期
	accessKey = setting.QiniuAccessKey
	secretKey = setting.QiniuSecretKey
)

func GetUploadMeta() (meta QiniuMeta) {
	putPolicy := storage.PutPolicy{
		Scope:   bucket,
		Expires: uint64(expiresIn),
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	meta = QiniuMeta{Token: upToken, BucketDomain: domain, ExpiresIn: expiresIn, Bucket: bucket}
	return
}
