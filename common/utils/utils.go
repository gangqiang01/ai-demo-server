package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func NewUUID() string {
	uuidWithHyphen := uuid.New()

	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}

// return ms
func GetNowTimeStamp() int64 {
	return int64(time.Now().UnixNano() / 1e6)
}

func Md5V(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func FileIsExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

func DirIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			return false
		}
	}
	return true
}

func UpdateFileName(oldPath, newPath string) error {
	_, err := os.Stat(oldPath)
	if os.IsNotExist(err) {
		fmt.Printf("file %s not exist\n", oldPath)
		return err
	}

	// 重命名文件
	err = os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Printf("Rename filename error: %v", err.Error())
		return err
	}
	return nil
}

func ListDir(dirname string) ([]string, error) {
	infos, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(infos))
	for i, info := range infos {
		names[i] = info.Name()
	}
	return names, nil
}
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
