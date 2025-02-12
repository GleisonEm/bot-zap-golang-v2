package services

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

func ConvertMp4ToWebp(videoPath string, outputFile string, dimension string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		videoPath,
		"-vcodec", "libwebp",
		"-filter:v", "fps=fps=20",
		"-lossless", "0",
		"-compression_level", "3",
		"-q:v", "70",
		"-loop", "1",
		"-preset", "picture",
		"-an",
		"-vsync", "0",
		"-s",
		dimension,
		outputFile,
	)

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ConvertMp4ToGif(videoPath string, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath,
		"-vf", "fps=20,scale=320:-1:flags=lanczos", // Define a taxa de quadros e redimensiona o vídeo
		"-gifflags", "-transdiff", // Otimiza o GIF para melhor qualidade
		"-y", // Sobrescreve o arquivo de saída sem perguntar
		outputFile,
	)

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ConvertGifToWebp(videoPath string, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath,
		"-vcodec", "libwebp",
		"-lossless", "0", // Compressão com perdas (ajuste para 1 para sem perdas)
		"-q:v", "90", // Qualidade de compressão (0-100, sendo 100 a melhor qualidade)
		"-loop", "0", // O WebP gerado vai loopar indefinidamente (0)
		outputFile,
	)

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func GetVideoDimensions(videoPath string) (int, int, error) {
	// Comando FFprobe para obter as dimensões do vídeo
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", videoPath)

	// Buffer para armazenar a saída do comando
	var out bytes.Buffer
	cmd.Stdout = &out

	// Executa o comando
	err := cmd.Run()
	if err != nil {
		return 0, 0, err
	}

	// Usa uma expressão regular para capturar as dimensões
	re := regexp.MustCompile(`(\d+)x(\d+)`)
	match := re.FindStringSubmatch(out.String())
	if len(match) < 3 {
		return 0, 0, fmt.Errorf("não foi possível obter as dimensões")
	}

	// Converte os valores encontrados para inteiros
	var width, height int
	fmt.Sscanf(match[0], "%dx%d", &width, &height)

	return width, height, nil
}
