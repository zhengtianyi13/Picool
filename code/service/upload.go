package service

import (
	"context"
	"io/ioutil"
	"mall/conf"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// UploadToQiNiu 封装上传图片到七牛云然后返回状态和图片的url，单张
func UploadToQiNiu(file multipart.File, fileSize int64) (path string, err error) {
	var AccessKey = conf.AccessKey
	var SerectKey = conf.SerectKey
	var Bucket = conf.Bucket
	var ImgUrl = conf.QiniuServer
	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SerectKey)
	upToken := putPlicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err = formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", err
	}
	url := ImgUrl + ret.Key
	return url, nil
}

func UploadAvatarToLocalStatic(file multipart.File, uId uint, userName string) (filePath string, err error) {
	bId := strconv.Itoa(int(uId)) //将uint转换为string
	basePath := "." + conf.AvatarPath + "user" + bId + "/"
	if !DirExistOrNot(basePath) {
		err := CreateDir(basePath)
		if err != nil {
			return "", err
		}
	}

	avatraPath := basePath + userName + ".jpg"

	content, err := ioutil.ReadAll(file) //将file转换为[]byte
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(avatraPath, content, 0666) //将[]byte写入文件
	if err != nil {
		return "", err
	}
	return "user" + bId + "/" + userName + ".jpg", nil

}

func UploadProductToLocalStatic(file multipart.File, uId uint, productName string) (filePath string, err error) {
	bId := strconv.Itoa(int(uId)) //将uint转换为string
	basePath := "." + conf.ProductPath + "boss" + bId + "/"
	if !DirExistOrNot(basePath) {
		err := CreateDir(basePath)
		if err != nil {
			return "", err
		}
	}

	productPath := basePath + productName + ".jpg"

	content, err := ioutil.ReadAll(file) //将file转换为[]byte
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(productPath, content, 0666) //将[]byte写入文件
	if err != nil {
		return "", err
	}
	return "boss" + bId + "/" + productName + ".jpg", nil

}

func DirExistOrNot(dir string) bool {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//创建文件夹
func CreateDir(dir string) error {
	err := os.MkdirAll(dir, 755)
	if err != nil {
		return err
	}
	return nil
}
