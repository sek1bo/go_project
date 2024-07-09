package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	// Инициализация флагов командной строки
	scaleFlag := flag.Float64("scale", 1.0, "коэффициент масштабирования")
	translateFlag := flag.String("translate", "", "сдвиг по осям X, Y, Z (разделенные пробелами)")
	rotateFlag := flag.String("rotate", "", "углы вращения по осям X, Y, Z в радианах (разделенные пробелами)")
	inputFileFlag := flag.String("input", "", "путь к входному файлу (например, model.obj)")
	outputFileFlag := flag.String("output", "parsed_model.obj", "путь к выходному файлу")

	// Парсинг флагов командной строки
	flag.Parse()

	model := &Model{}

	// Загрузка модели из файла, если указан путь к входному файлу
	if *inputFileFlag != "" {
		err := model.LoadFromFile(*inputFileFlag)
		if err != nil {
			fmt.Println("Ошибка загрузки модели:", err)
			return
		}
	} else {
		fmt.Println("Не указан путь к входному файлу.")
		flag.PrintDefaults()
		return
	}

	// Преобразование линий в грани
	model.ConvertLinesToFaces()

	// Применение масштабирования, если указан флаг scale
	if *scaleFlag != 1.0 {
		model.Scale(*scaleFlag, *scaleFlag, *scaleFlag)
		fmt.Printf("Масштабирование модели на коэффициент: %.2f\n", *scaleFlag)
	}

	// Применение трансляции, если указан флаг translate
	if *translateFlag != "" {
		translateValues := strings.Fields(*translateFlag)
		if len(translateValues) != 3 {
			fmt.Println("Неверный формат для сдвига по осям X, Y, Z.")
			flag.PrintDefaults()
			return
		}

		tx, err := strconv.ParseFloat(translateValues[0], 64)
		if err != nil {
			fmt.Println("Неверный сдвиг по оси X:", err)
			return
		}

		ty, err := strconv.ParseFloat(translateValues[1], 64)
		if err != nil {
			fmt.Println("Неверный сдвиг по оси Y:", err)
			return
		}

		tz, err := strconv.ParseFloat(translateValues[2], 64)
		if err != nil {
			fmt.Println("Неверный сдвиг по оси Z:", err)
			return
		}

		model.Translate(tx, ty, tz)
		fmt.Printf("Трансляция модели на сдвиг: %.2f, %.2f, %.2f\n", tx, ty, tz)
	}

	// Применение вращения, если указан флаг rotate
	if *rotateFlag != "" {
		rotateValues := strings.Fields(*rotateFlag)
		if len(rotateValues) != 3 {
			fmt.Println("Неверный формат для углов вращения по осям X, Y, Z.")
			flag.PrintDefaults()
			return
		}

		rx, err := strconv.ParseFloat(rotateValues[0], 64)
		if err != nil {
			fmt.Println("Неверный угол вращения по оси X:", err)
			return
		}

		ry, err := strconv.ParseFloat(rotateValues[1], 64)
		if err != nil {
			fmt.Println("Неверный угол вращения по оси Y:", err)
			return
		}

		rz, err := strconv.ParseFloat(rotateValues[2], 64)
		if err != nil {
			fmt.Println("Неверный угол вращения по оси Z:", err)
			return
		}

		model.Rotate(rx, ry, rz)
		fmt.Printf("Вращение модели на углы: %.2f, %.2f, %.2f (радианы)\n", rx, ry, rz)
	}

	// Сохранение модели в файл
	err := model.SaveToFile(*outputFileFlag)
	if err != nil {
		fmt.Println("Ошибка сохранения модели:", err)
		return
	}

	fmt.Println("Завершено. Разобранные данные записаны в", *outputFileFlag)
}
