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

func getDetect(module string) (string, string, string) {
	aiCfg := config.GetAiConfig()
	detectPath := ""
	outPath := ""
	inputPath := ""
	switch module {
	case "car_image":
		outPath = aiCfg.Car_Image.OutPath
		inputPath = aiCfg.Car_Image.InputPath
		detectPath = aiCfg.Car_Image.DetectPath
	case "car_video":
		outPath = aiCfg.Car_Video.OutPath
		inputPath = aiCfg.Car_Video.InputPath
		detectPath = aiCfg.Car_Video.DetectPath
	case "object_image":
		outPath = aiCfg.Object_Image.OutPath
		inputPath = aiCfg.Object_Image.InputPath
		detectPath = aiCfg.Object_Image.DetectPath
	case "object_video":
		outPath = aiCfg.Object_Video.OutPath
		inputPath = aiCfg.Object_Video.InputPath
		detectPath = aiCfg.Object_Video.DetectPath
	case "mask_image":
		outPath = aiCfg.Mask_Image.OutPath
		inputPath = aiCfg.Mask_Image.InputPath
		detectPath = aiCfg.Mask_Image.DetectPath

	}

	return detectPath, inputPath, outPath
}

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
	module := c.Query("module")
	_, inputPath, outPath := getDetect(module)
	filePath := path.Join(outPath, filename)
	if tp == "data" {
		filePath = path.Join(inputPath, filename)
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
	module := c.Query("module")
	_, inputPath, outPath := getDetect(module)
	filePath := path.Join(outPath, "display", filename)
	if tp == "data" {
		filePath = path.Join(inputPath, filename)
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
	module := c.Query("module")
	_, _, outPath := getDetect(module)
	files, err := os.ReadDir(outPath)
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
	module := c.Param("module")
	// 获取上传的文件
	if err != nil {
		responce.FailWithMessage(fmt.Sprintf("Upload file error: %v", err.Error()), c)
		return
	}
	detectPath, inputPath, outPath := getDetect(module)
	filename := file.Filename
	aiCfg := config.GetAiConfig()
	dstPath := path.Join(inputPath, filename)
	klog.Infof("dstPath: %s", dstPath)
	if !utils.DirIsExist(inputPath) {
		os.MkdirAll(inputPath, os.ModePerm)
	}

	err = c.SaveUploadedFile(file, dstPath)
	if err != nil {
		responce.FailWithMessage(fmt.Sprintf("Save file error: %v", err.Error()), c)
		return
	}

	if tp == "video" {
		inputPath = path.Join(detectPath, filename)
	}
	detectStatus[filename] = "0"
	if err := detect(aiCfg.Command, detectPath, inputPath); err != nil {
		detectStatus[filename] = "1"
		responce.FailWithMessage(err.Error(), c)
		return
	}
	if tp == "video" {
		displayDir := path.Join(outPath, "display")
		if !utils.DirIsExist(displayDir) {
			os.MkdirAll(displayDir, os.ModePerm)
		}
		outPath := path.Join(outPath, filename)
		displayPath := path.Join(displayDir, filename)
		time.Sleep(1 * time.Second)
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
	module := c.Param("module")
	_, _, outPath := getDetect(module)
	aiCfg := config.GetAiConfig()
	displayDir := path.Join(outPath, "display")
	if !utils.DirIsExist(displayDir) {
		os.MkdirAll(displayDir, os.ModePerm)
	}
	outPath = path.Join(outPath, filename)
	displayPath := path.Join(displayDir, filename)
	if err := ffmpegConvert(aiCfg.FFMPEGPath, outPath, displayPath); err != nil {
		responce.FailWithMessage(err.Error(), c)
		return
	}
	responce.Ok(c)
}
func DeleteAiDetect(c *gin.Context) {
	filename := c.Param("filename")
	module := c.Param("module")
	_, inputPath, outPath := getDetect(module)
	inputPath = path.Join(inputPath, filename)
	outPath = path.Join(outPath, filename)
	defer os.Remove(inputPath)
	defer os.Remove(outPath)
	responce.Ok(c)
}
