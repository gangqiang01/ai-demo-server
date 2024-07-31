package v1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/edgehook/ithings/common/config"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"

	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
)

var detectStatus = make(map[string]string)

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

// /home/gangqiangsun/MYPROJECT/TOOLS/ffmpeg-6.0/ffmpeg -i ./demo.mp4 -vcodec libx264 -acodec copy ./demo-convert.mp4
func ffmpegConvert(ffmpegPath, filePath, outPath string) error {
	cmd := exec.Command(ffmpegPath, "-i", filePath, "-vcodec", "libx264", "-acodec", "copy", "-y", outPath)
	klog.Infof("ffmpeg cmd path: %s, args: %v", cmd.Path, cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		klog.Errorf("Failed to run ffmpeg command:%v ", err)
		return err
	}
	klog.Infof("FFmpeg output:%s, error: %v", out.String(), stderr.String())
	return nil
}

func GetDetectStatus(c *gin.Context) {
	filename := c.Param("filename")
	val, ok := detectStatus[filename]
	if !ok {
		responce.FailWithMessage("file does not exist", c)
		return

	}
	responce.OkWithData(val, c)
}

func GetDetectImage(c *gin.Context) {
	filename := c.Param("filename")
	tp := c.Query("type")
	aiCfg := config.GetAiConfig()
	filePath := path.Join(aiCfg.OutPath, filename)
	if tp == "data" {
		filePath = path.Join(aiCfg.InputPath, filename)
	}

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
	tp := c.Query("type")
	aiCfg := config.GetAiConfig()
	filePath := path.Join(aiCfg.OutPath, "display", filename)
	if tp == "data" {
		filePath = path.Join(aiCfg.InputPath, filename)
	}

	file, err := os.Open(filePath)
	if err != nil {
		responce.FailWithMessage("File not found", c)
		return
	}
	defer file.Close()

	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		responce.FailWithMessage("Read file error", c)
		return
	}
}

func GetDetectedFiles(c *gin.Context) {
	tp := c.Query("type")
	aiCfg := config.GetAiConfig()
	files, err := os.ReadDir(aiCfg.OutPath)
	if err != nil {
		responce.FailWithMessage(err.Error(), c)
		return
	}
	fileData := make([]map[string]string, 0)
	for _, file := range files {
		filename := file.Name()
		fileInfo := make(map[string]string)
		if tp == "video" {
			if !file.IsDir() && strings.HasSuffix(filename, ".mp4") {
				fileInfo["filename"] = filename
				fileData = append(fileData, fileInfo)
			}
		} else {
			if !file.IsDir() && (strings.HasSuffix(filename, ".png") || strings.HasSuffix(filename, ".jpg")) {
				fileInfo["filename"] = filename
				fileData = append(fileData, fileInfo)
			}
		}

	}
	responce.OkWithData(fileData, c)
}

func AddAiDetect(c *gin.Context) {

	file, err := c.FormFile("file")
	tp := c.Param("type")
	// 获取上传的文件
	if err != nil {
		responce.FailWithMessage(fmt.Sprintf("Upload file error: %v", err.Error()), c)
		return
	}
	filename := file.Filename
	aiCfg := config.GetAiConfig()
	dstPath := path.Join(aiCfg.InputPath, filename)
	klog.Infof("dstPath: %s", dstPath)
	if !utils.DirIsExist(aiCfg.InputPath) {
		os.MkdirAll(aiCfg.InputPath, os.ModePerm)
	}

	err = c.SaveUploadedFile(file, dstPath)
	if err != nil {
		responce.FailWithMessage(fmt.Sprintf("Save file error: %v", err.Error()), c)
		return
	}
	detectPath := aiCfg.InputPath
	if tp == "video" {
		detectPath = path.Join(detectPath, filename)
	}
	detectStatus[filename] = "0"
	if err := detect(aiCfg.Command, aiCfg.DetectPath, detectPath); err != nil {
		detectStatus[filename] = "1"
		responce.FailWithMessage(err.Error(), c)
		return
	}
	if tp == "video" {
		displayDir := path.Join(aiCfg.OutPath, "display")
		if !utils.DirIsExist(displayDir) {
			os.MkdirAll(displayDir, os.ModePerm)
		}
		outPath := path.Join(aiCfg.OutPath, filename)
		displayPath := path.Join(displayDir, filename)
		time.Sleep(3 * time.Second)
		if err := ffmpegConvert(aiCfg.FFMPEGPath, outPath, displayPath); err != nil {
			detectStatus[filename] = "1"
			responce.FailWithMessage(err.Error(), c)
			return
		}
	}

	detectStatus[filename] = "1"
	responce.Ok(c)
}
func ConvertFormat(c *gin.Context) {
	filename := c.Param("filename")
	aiCfg := config.GetAiConfig()
	displayDir := path.Join(aiCfg.OutPath, "display")
	if !utils.DirIsExist(displayDir) {
		os.MkdirAll(displayDir, os.ModePerm)
	}
	outPath := path.Join(aiCfg.OutPath, filename)
	displayPath := path.Join(displayDir, filename)
	if err := ffmpegConvert(aiCfg.FFMPEGPath, outPath, displayPath); err != nil {
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
	defer os.Remove(inputPath)
	defer os.Remove(outputPath)
	responce.Ok(c)
}
