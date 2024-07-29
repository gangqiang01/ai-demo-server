package v1

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/edgehook/ithings/common/config"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"

	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
)

func detect(command, detectPath, inputPath string) error {
	cmd := exec.Command(command, detectPath, inputPath)
	klog.Infof("cmd path: %s, args: %v", cmd.Path, cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		klog.Errorf("Failed to run command:%v ", err)
		return err
	}
	klog.Infof("Detect output:%s", out.String())
	return nil
}
func GetDetectImage(c *gin.Context) {
	filename := c.Param("filename")
	aiCfg := config.GetAiConfig()
	filePath := path.Join(aiCfg.OutPath, filename)
	image, err := os.ReadFile(filePath)
	if err != nil {
		responce.FailWithMessage("Failed to read image file", c)
		return
	}
	// 设置HTTP响应头
	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Disposition", "inline")
	c.Data(http.StatusOK, "image/jpeg", image)
}

func GetDetectVideo(c *gin.Context) {
	filename := c.Param("filename")
	aiCfg := config.GetAiConfig()
	filePath := path.Join(aiCfg.OutPath, filename)
	videoData, err := os.ReadFile(filePath)
	if err != nil {
		responce.FailWithMessage("Failed to read video file", c)
		return
	}
	// 设置HTTP响应头
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Disposition", "inline")
	c.Data(http.StatusOK, "video/mp4", videoData)
}

func AddAiDetect(c *gin.Context) {

	file, err := c.FormFile("file") // 获取上传的文件
	if err != nil {
		responce.FailWithMessage(fmt.Sprintf("Upload file error: %v", err.Error()), c)
		return
	}
	filename := file.Filename
	aiCfg := config.GetAiConfig()
	dstPath := path.Join(aiCfg.InputPath, filename)
	klog.Infof("dstPath: %s", dstPath)
	if utils.DirIsExist(aiCfg.InputPath) {
		os.MkdirAll(aiCfg.InputPath, os.ModePerm)
	}
	err = c.SaveUploadedFile(file, dstPath)
	if err != nil {
		responce.FailWithMessage(fmt.Sprintf("Save file error: %v", err.Error()), c)
		return
	}
	if err := detect(aiCfg.Command, aiCfg.DetectPath, aiCfg.InputPath); err != nil {
		responce.FailWithMessage(err.Error(), c)
		return
	}

	responce.Ok(c)
}

func DeleteAiDetect(c *gin.Context) {
	filename := c.Param("filename")

	aiCfg := config.GetAiConfig()
	inputPath := path.Join(aiCfg.InputPath, filename)
	outputPath := path.Join(aiCfg.OutPath, filename)
	os.Remove(inputPath)
	os.Remove(outputPath)
	responce.Ok(c)
}
