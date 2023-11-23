package qiniu

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"path/filepath"
	"ruoyi/settings"
	"time"
)

func UploadToQiNiu(file multipart.File, fileSize int64, name string) (string, error) {
	a := settings.Conf.QiNiuConfig
	var AccessKey = a.AccessKey
	var SecretKey = a.SerectKey
	var Bucket = a.Bucket
	var ImgUrl = a.QiniuServe
	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPlicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	uniqueFileName := generateUniqueFileName(name) //生成唯一名字
	err := formUploader.Put(context.Background(), &ret, upToken, uniqueFileName, file, fileSize, &putExtra)
	if err != nil {
		return "", err
	}
	url := ImgUrl + "/" + ret.Key
	return url, nil
}

func generateUniqueFileName(originalName string) string {
	// 使用当前时间戳和原始文件名生成唯一文件名
	timestamp := time.Now().Unix()
	ext := filepath.Ext(originalName)
	nameWithoutExt := originalName[:len(originalName)-len(ext)]
	uniqueName := fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)
	return uniqueName
}
