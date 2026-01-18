package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type googleSpeechRequest struct {
	Config googleSpeechConfig `json:"config"`
	Audio  googleSpeechAudio  `json:"audio"`
}

type googleSpeechConfig struct {
	Encoding     string `json:"encoding,omitempty"`
	LanguageCode string `json:"languageCode"`
}

type googleSpeechAudio struct {
	Content string `json:"content"`
}

type googleSpeechResponse struct {
	Results []struct {
		Alternatives []struct {
			Transcript string `json:"transcript"`
		} `json:"alternatives"`
	} `json:"results"`
}

func TranscribeAudio(reader io.Reader, filename, contentType string) (string, error) {
	apiKey := os.Getenv("GCP_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GCP_API_KEY not set")
	}
	return transcribeWithGoogle(apiKey, reader, filename, contentType)
}

func transcribeWithGoogle(apiKey string, reader io.Reader, filename, contentType string) (string, error) {
	audioBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	encoding := googleEncodingFor(contentType, filename)
	payload := googleSpeechRequest{
		Config: googleSpeechConfig{
			Encoding:     encoding,
			LanguageCode: "en-US",
		},
		Audio: googleSpeechAudio{
			Content: base64.StdEncoding.EncodeToString(audioBytes),
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://speech.googleapis.com/v1/speech:recognize?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google speech error: %s", string(body))
	}

	var speechResp googleSpeechResponse
	if err := json.Unmarshal(body, &speechResp); err != nil {
		return "", err
	}

	for _, result := range speechResp.Results {
		for _, alt := range result.Alternatives {
			if strings.TrimSpace(alt.Transcript) != "" {
				return strings.TrimSpace(alt.Transcript), nil
			}
		}
	}

	return "", fmt.Errorf("google speech returned no transcript")
}

func googleEncodingFor(contentType, filename string) string {
	if strings.Contains(contentType, "webm") || strings.EqualFold(filepath.Ext(filename), ".webm") {
		return "WEBM_OPUS"
	}
	if strings.Contains(contentType, "ogg") || strings.EqualFold(filepath.Ext(filename), ".ogg") {
		return "OGG_OPUS"
	}
	if strings.Contains(contentType, "wav") || strings.EqualFold(filepath.Ext(filename), ".wav") {
		return "LINEAR16"
	}
	return ""
}

func transcribeWithWhisper(reader io.Reader, filename, contentType string) (string, error) {
	whisperPath := os.Getenv("WHISPER_CPP_PATH")
	modelPath := os.Getenv("WHISPER_MODEL_PATH")
	if whisperPath == "" || modelPath == "" {
		return "", fmt.Errorf("GOOGLE_SPEECH_API_KEY not set and WHISPER_CPP_PATH/WHISPER_MODEL_PATH not set")
	}

	inputExt := filepath.Ext(filename)
	if inputExt == "" {
		inputExt = ".webm"
	}

	inputFile, err := os.CreateTemp("", "speech-input-*"+inputExt)
	if err != nil {
		return "", err
	}
	defer os.Remove(inputFile.Name())

	if _, err := io.Copy(inputFile, reader); err != nil {
		inputFile.Close()
		return "", err
	}
	if err := inputFile.Close(); err != nil {
		return "", err
	}

	workingInput := inputFile.Name()
	var convertedFile string
	if needsConversion(contentType, inputExt) {
		ffmpegPath, err := exec.LookPath("ffmpeg")
		if err != nil {
			return "", fmt.Errorf("ffmpeg not found; install it to convert audio for Whisper")
		}
		convertedFile = strings.TrimSuffix(inputFile.Name(), inputExt) + ".wav"
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, ffmpegPath, "-y", "-i", inputFile.Name(), "-ar", "16000", "-ac", "1", convertedFile)
		if output, err := cmd.CombinedOutput(); err != nil {
			_ = os.Remove(convertedFile)
			return "", fmt.Errorf("ffmpeg conversion failed: %s", string(output))
		}
		defer os.Remove(convertedFile)
		workingInput = convertedFile
	}

	outputBase := filepath.Join(os.TempDir(), fmt.Sprintf("speech-%d", time.Now().UnixNano()))
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, whisperPath,
		"--model", modelPath,
		"--output-txt",
		"--output-file", outputBase,
		"--no-timestamps",
		"--language", "en",
		workingInput,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("whisper.cpp failed: %s", string(output))
	}

	textBytes, err := os.ReadFile(outputBase + ".txt")
	if err != nil {
		return "", err
	}
	_ = os.Remove(outputBase + ".txt")

	return strings.TrimSpace(string(textBytes)), nil
}

func needsConversion(contentType, ext string) bool {
	if strings.Contains(contentType, "wav") {
		return false
	}
	if strings.EqualFold(ext, ".wav") {
		return false
	}
	return true
}
