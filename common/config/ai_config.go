package config

type AIConfig struct {
	Command    string
	DetectPath string
	InputPath  string
	OutPath    string
}

func GetAiConfig() *AIConfig {
	command := ITHINGS_CONFIG.GetString("ai.command")
	detectPath := ITHINGS_CONFIG.GetString("ai.detectPath")
	inputPath := ITHINGS_CONFIG.GetString("ai.inputPath")
	outPath := ITHINGS_CONFIG.GetString("ai.outPath")
	return &AIConfig{
		Command:    command,
		DetectPath: detectPath,
		InputPath:  inputPath,
		OutPath:    outPath,
	}
}
