package main

import (
	"fmt"
	"strconv"
)

func main() {
	model := &Model{}

	// Получение коэффициента масштабирования от пользователя
	scaleFactorStr := getUserInput("Введите коэффициент масштабирования (например, 1.5): ")
	scaleFactor, err := strconv.ParseFloat(scaleFactorStr, 64)
	if err != nil {
		fmt.Println("Неверный коэффициент масштабирования:", err)
		return
	}

	model.Scale = scaleFactor

	// Получение пути к входному файлу от пользователя
	inputFile := getUserInput("Введите путь к входному файлу (например, model.obj): ")

	// Загрузка модели из файла
	err = model.LoadFromFile(inputFile)
	if err != nil {
		fmt.Println("Ошибка загрузки модели:", err)
		return
	}

	// Получение координат сдвига от пользователя
	txStr := getUserInput("Введите сдвиг по оси X (например, 0.0): ")
	tyStr := getUserInput("Введите сдвиг по оси Y (например, 0.0): ")
	tzStr := getUserInput("Введите сдвиг по оси Z (например, 0.0): ")

	tx, err := strconv.ParseFloat(txStr, 64)
	if err != nil {
		fmt.Println("Неверный сдвиг по оси X:", err)
		return
	}

	ty, err := strconv.ParseFloat(tyStr, 64)
	if err != nil {
		fmt.Println("Неверный сдвиг по оси Y:", err)
		return
	}

	tz, err := strconv.ParseFloat(tzStr, 64)
	if err != nil {
		fmt.Println("Неверный сдвиг по оси Z:", err)
		return
	}

	// Применение трансляции к модели
	model.Translate(tx, ty, tz)

	// Получение углов вращения от пользователя
	rxStr := getUserInput("Введите угол вращения по оси X в радианах (например, 3.14): ")
	ryStr := getUserInput("Введите угол вращения по оси Y в радианах (например, 3.14): ")
	rzStr := getUserInput("Введите угол вращения по оси Z в радианах (например, 3.14): ")

	rx, err := strconv.ParseFloat(rxStr, 64)
	if err != nil {
		fmt.Println("Неверный угол вращения по оси X:", err)
		return
	}

	ry, err := strconv.ParseFloat(ryStr, 64)
	if err != nil {
		fmt.Println("Неверный угол вращения по оси Y:", err)
		return
	}

	rz, err := strconv.ParseFloat(rzStr, 64)
	if err != nil {
		fmt.Println("Неверный угол вращения по оси Z:", err)
		return
	}

	// Применение вращения к модели
	model.Rotate(rx, ry, rz)

	// Сохранение масштабированной, транслированной и вращенной модели в файл
	outputFile := "parsed_model.obj"
	err = model.SaveToFile(outputFile)
	if err != nil {
		fmt.Println("Ошибка сохранения модели:", err)
		return
	}

	fmt.Println("Завершено. Разобранные данные записаны в", outputFile)
}
