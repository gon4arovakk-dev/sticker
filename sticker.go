// sticker.go
// Sticker Maker на Go (упрощённый алгоритм с использованием imaging)

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/disintegration/imaging"
)

func main() {
	var (
		input     string
		output    string
		colorStr  string
		threshold int
		verbose   bool
	)
	flag.StringVar(&input, "i", "", "Входной файл")
	flag.StringVar(&output, "o", "", "Выходной файл (PNG)")
	flag.StringVar(&colorStr, "c", "", "Цвет фона (HEX, например #FFFFFF)")
	flag.IntVar(&threshold, "t", 30, "Порог схожести (0-100)")
	flag.BoolVar(&verbose, "v", false, "Подробный вывод")
	flag.Parse()

	if input == "" && flag.NArg() > 0 {
		input = flag.Arg(0)
	}
	if input == "" {
		fmt.Println("Ошибка: не указан входной файл.")
		flag.Usage()
		os.Exit(1)
	}
	if output == "" {
		output = input[:len(input)-len(ext(input))] + "_sticker.png"
	}

	// Читаем изображение
	src, err := imaging.Open(input)
	if err != nil {
		fmt.Printf("Ошибка открытия: %v\n", err)
		os.Exit(1)
	}

	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Определяем цвет фона
	bgColor := color.RGBA{0, 0, 0, 255}
	if colorStr != "" {
		// Парсим HEX
		c, err := parseHexColor(colorStr)
		if err == nil {
			bgColor = c
		} else if verbose {
			fmt.Printf("Предупреждение: не удалось распарсить цвет, используем автоопределение\n")
		}
	} else {
		// Автоопределение: средний цвет по краям
		bgColor = detectBackgroundColor(src)
		if verbose {
			fmt.Printf("Определён цвет фона: RGB(%d,%d,%d)\n", bgColor.R, bgColor.G, bgColor.B)
		}
	}

	// Создаём новое изображение с альфа-каналом
	dst := imaging.New(width, height, color.Alpha16{}) // пустое

	// Порог в цветовом пространстве (макс расстояние 441.67)
	maxDist := math.Sqrt(3 * 255 * 255)
	thresholdNorm := float64(threshold) / 100.0 * maxDist

	// Проходим по пикселям и устанавливаем прозрачность
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			px := src.At(x, y)
			r, g, b, _ := px.RGBA()
			// Приводим к 8-бит
			pr := uint8(r >> 8)
			pg := uint8(g >> 8)
			pb := uint8(b >> 8)

			dist := colorDistance(pr, pg, pb, bgColor.R, bgColor.G, bgColor.B)
			if dist <= thresholdNorm {
				// Прозрачный
				dst.Set(x, y, color.RGBA{pr, pg, pb, 0})
			} else {
				dst.Set(x, y, color.RGBA{pr, pg, pb, 255})
			}
		}
	}

	// Сохраняем
	err = imaging.Save(dst, output, imaging.PNG)
	if err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
		os.Exit(1)
	}
	if verbose {
		fmt.Printf("Фон удалён. Результат сохранён в %s\n", output)
	}
}

func ext(filename string) string {
	for i := len(filename) - 1; i >= 0 && filename[i] != '.'; i-- {
		if filename[i] == '/' || filename[i] == '\\' {
			break
		}
		if i == 0 {
			return ""
		}
	}
	return filename[strings.LastIndex(filename, "."):]
}

func parseHexColor(s string) (color.RGBA, error) {
	// Упрощённый парсинг
	var r, g, b uint8
	_, err := fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	}
	return color.RGBA{r, g, b, 255}, err
}

func detectBackgroundColor(img image.Image) color.RGBA {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	// Собираем пиксели с краёв (верхняя, нижняя, левая, правая строки)
	var pixels []color.RGBA
	// Верхняя строка
	for x := 0; x < w; x++ {
		px := img.At(x, 0)
		r, g, b, _ := px.RGBA()
		pixels = append(pixels, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
	}
	// Нижняя строка
	for x := 0; x < w; x++ {
		px := img.At(x, h-1)
		r, g, b, _ := px.RGBA()
		pixels = append(pixels, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
	}
	// Левая строка
	for y := 0; y < h; y++ {
		px := img.At(0, y)
		r, g, b, _ := px.RGBA()
		pixels = append(pixels, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
	}
	// Правая строка
	for y := 0; y < h; y++ {
		px := img.At(w-1, y)
		r, g, b, _ := px.RGBA()
		pixels = append(pixels, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
	}
	// Находим средний цвет (упрощённо)
	var sumR, sumG, sumB uint64
	for _, p := range pixels {
		sumR += uint64(p.R)
		sumG += uint64(p.G)
		sumB += uint64(p.B)
	}
	n := uint64(len(pixels))
	if n == 0 {
		return color.RGBA{255, 255, 255, 255}
	}
	return color.RGBA{uint8(sumR / n), uint8(sumG / n), uint8(sumB / n), 255}
}

func colorDistance(r1, g1, b1, r2, g2, b2 uint8) float64 {
	dr := float64(r1) - float64(r2)
	dg := float64(g1) - float64(g2)
	db := float64(b1) - float64(b2)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}
