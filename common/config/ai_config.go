package config

type AIConfig struct {
	Command      string
	Car_Image    *DetectConfig
	Car_Video    *DetectConfig
	Object_Image *DetectConfig
	Object_Video *DetectConfig
	Mask_Image   *DetectConfig
	FFMPEGPath   string
}

type DetectConfig struct {
	DetectPath string
	InputPath  string
	OutPath    string
}

func GetAiConfig() *AIConfig {
	command := ITHINGS_CONFIG.GetString("ai.command")
	ffmpegPath := ITHINGS_CONFIG.GetString("ai.ffmpegPath")
	car_image_detectPath := ITHINGS_CONFIG.GetString("ai.car_image.detectPath")
	car_image_inputPath := ITHINGS_CONFIG.GetString("ai.car_image.inputPath")
	car_image_outPath := ITHINGS_CONFIG.GetString("ai.car_image.outPath")
	car_video_detectPath := ITHINGS_CONFIG.GetString("ai.car_video.detectPath")
	car_video_inputPath := ITHINGS_CONFIG.GetString("ai.car_video.inputPath")
	car_video_outPath := ITHINGS_CONFIG.GetString("ai.car_video.outPath")
	object_image_detectPath := ITHINGS_CONFIG.GetString("ai.object_image.detectPath")
	object_image_inputPath := ITHINGS_CONFIG.GetString("ai.object_image.inputPath")
	object_image_outPath := ITHINGS_CONFIG.GetString("ai.object_image.outPath")
	object_video_detectPath := ITHINGS_CONFIG.GetString("ai.object_video.detectPath")
	object_video_inputPath := ITHINGS_CONFIG.GetString("ai.object_video.inputPath")
	object_video_outPath := ITHINGS_CONFIG.GetString("ai.object_video.outPath")
	mask_image_detectPath := ITHINGS_CONFIG.GetString("ai.mask_image.detectPath")
	mask_image_inputPath := ITHINGS_CONFIG.GetString("ai.mask_image.inputPath")
	mask_image_outPath := ITHINGS_CONFIG.GetString("ai.mask_image.outPath")
	if ffmpegPath == "" {
		ffmpegPath = "./ffmpeg/ffmpeg"
	}
	car_image := &DetectConfig{
		DetectPath: car_image_detectPath,
		InputPath:  car_image_inputPath,
		OutPath:    car_image_outPath,
	}
	object_image := &DetectConfig{
		DetectPath: object_image_detectPath,
		InputPath:  object_image_inputPath,
		OutPath:    object_image_outPath,
	}
	car_video := &DetectConfig{
		DetectPath: car_video_detectPath,
		InputPath:  car_video_inputPath,
		OutPath:    car_video_outPath,
	}
	object_video := &DetectConfig{
		DetectPath: object_video_detectPath,
		InputPath:  object_video_inputPath,
		OutPath:    object_video_outPath,
	}
	mask_image := &DetectConfig{
		DetectPath: mask_image_detectPath,
		InputPath:  mask_image_inputPath,
		OutPath:    mask_image_outPath,
	}
	return &AIConfig{
		Command:      command,
		Car_Image:    car_image,
		Car_Video:    car_video,
		Object_Image: object_image,
		Object_Video: object_video,
		Mask_Image:   mask_image,
		FFMPEGPath:   ffmpegPath,
	}
}
