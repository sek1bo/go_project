package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Vertex структура для хранения координат вершины
type Vertex struct {
	X, Y, Z float64
}

// TexCoord структура для хранения текстурных координат
type TexCoord struct {
	U, V float64
}

// Normal структура для хранения нормалей
type Normal struct {
	X, Y, Z float64
}

// Face структура для хранения вершин, текстур и нормалей
type Face struct {
	Vertices  []int
	TexCoords []int
	Normals   []int
}

// Model структура для хранения всех данных модели
type Model struct {
	Vertices   []Vertex
	TexCoords  []TexCoord
	Normals    []Normal
	Faces      []Face
	Scale      float64
	TranslateX float64 // Сдвиг по оси X
	TranslateY float64 // Сдвиг по оси Y
	TranslateZ float64 // Сдвиг по оси Z
}

// Метод для загрузки модели из файла
func (m *Model) LoadFromFile(path string) error {
	inputFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "v "):
			v, err := m.parseVertex(line)
			if err == nil {
				m.Vertices = append(m.Vertices, v)
			}
		case strings.HasPrefix(line, "vt "):
			vt, err := m.parseTexCoord(line)
			if err == nil {
				m.TexCoords = append(m.TexCoords, vt)
			}
		case strings.HasPrefix(line, "vn "):
			vn, err := m.parseNormal(line)
			if err == nil {
				m.Normals = append(m.Normals, vn)
			}
		case strings.HasPrefix(line, "f "):
			face, err := m.parseFace(line)
			if err == nil {
				m.Faces = append(m.Faces, face)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}

// Метод для парсинга вершины
func (m *Model) parseVertex(line string) (Vertex, error) {
	var v Vertex
	parts := strings.Fields(line)[1:]
	if len(parts) != 3 {
		return v, fmt.Errorf("invalid vertex line: %s", line)
	}

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return v, fmt.Errorf("invalid X coordinate: %v", err)
	}
	v.X = x

	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return v, fmt.Errorf("invalid Y coordinate: %v", err)
	}
	v.Y = y

	z, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return v, fmt.Errorf("invalid Z coordinate: %v", err)
	}
	v.Z = z

	return v, nil
}

// Метод для парсинга текстурных координат
func (m *Model) parseTexCoord(line string) (TexCoord, error) {
	var vt TexCoord
	parts := strings.Fields(line)[1:]
	if len(parts) != 2 {
		return vt, fmt.Errorf("invalid texture coordinate line: %s", line)
	}

	u, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return vt, fmt.Errorf("invalid U coordinate: %v", err)
	}
	vt.U = u

	v, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return vt, fmt.Errorf("invalid V coordinate: %v", err)
	}
	vt.V = v

	return vt, nil
}

// Метод для парсинга нормалей
func (m *Model) parseNormal(line string) (Normal, error) {
	var vn Normal
	parts := strings.Fields(line)[1:]
	if len(parts) != 3 {
		return vn, fmt.Errorf("invalid normal line: %s", line)
	}

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return vn, fmt.Errorf("invalid X coordinate: %v", err)
	}
	vn.X = x

	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return vn, fmt.Errorf("invalid Y coordinate: %v", err)
	}
	vn.Y = y

	z, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return vn, fmt.Errorf("invalid Z coordinate: %v", err)
	}
	vn.Z = z

	return vn, nil
}

// Метод для парсинга грани
func (m *Model) parseFace(line string) (Face, error) {
	var f Face
	parts := strings.Fields(line)[1:]
	for _, part := range parts {
		vertexData := strings.Split(part, "/")

		vIdx, err := strconv.Atoi(vertexData[0])
		if err != nil {
			return f, fmt.Errorf("invalid vertex index: %v", err)
		}
		f.Vertices = append(f.Vertices, vIdx-1)

		if len(vertexData) > 1 && vertexData[1] != "" {
			tIdx, err := strconv.Atoi(vertexData[1])
			if err != nil {
				return f, fmt.Errorf("invalid texture coordinate index: %v", err)
			}
			f.TexCoords = append(f.TexCoords, tIdx-1)
		}

		if len(vertexData) > 2 && vertexData[2] != "" {
			nIdx, err := strconv.Atoi(vertexData[2])
			if err != nil {
				return f, fmt.Errorf("invalid normal index: %v", err)
			}
			f.Normals = append(f.Normals, nIdx-1)
		}
	}

	return f, nil
}

// Метод для записи модели в файл
func (m *Model) SaveToFile(path string) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	writer.WriteString("# Parsed .obj file\n\n")

	// Запись вершин
	writer.WriteString("# Vertices\n")
	for _, v := range m.Vertices {
		writer.WriteString(fmt.Sprintf("v %f %f %f\n", v.X*m.Scale+m.TranslateX, v.Y*m.Scale+m.TranslateY, v.Z*m.Scale+m.TranslateZ))
	}

	// Запись текстурных координат (если они есть)
	if len(m.TexCoords) > 0 {
		writer.WriteString("\n# Texture Coordinates\n")
		for _, vt := range m.TexCoords {
			writer.WriteString(fmt.Sprintf("vt %f %f\n", vt.U, vt.V))
		}
	}

	// Запись нормалей (если они есть)
	if len(m.Normals) > 0 {
		writer.WriteString("\n# Normals\n")
		for _, vn := range m.Normals {
			writer.WriteString(fmt.Sprintf("vn %f %f %f\n", vn.X, vn.Y, vn.Z))
		}
	}

	// Запись граней
	if len(m.Faces) > 0 {
		writer.WriteString("\n# Faces\n")
		for _, f := range m.Faces {
			writer.WriteString("f")
			for i := range f.Vertices {
				writer.WriteString(fmt.Sprintf(" %d", f.Vertices[i]+1))
				if len(f.TexCoords) > i {
					writer.WriteString(fmt.Sprintf("/%d", f.TexCoords[i]+1))
				}
				if len(f.Normals) > i {
					writer.WriteString(fmt.Sprintf("/%d", f.Normals[i]+1))
				}
			}
			writer.WriteString("\n")
		}
	}

	writer.Flush()
	return nil
}

// Метод для трансляции модели по заданным сдвигам по осям X, Y, Z
func (m *Model) Translate(tx, ty, tz float64) {
	m.TranslateX += tx
	m.TranslateY += ty
	m.TranslateZ += tz
}

// Метод для вращения модели вокруг осей X, Y и Z
func (m *Model) Rotate(rx, ry, rz float64) {
	// Вращение вокруг оси X
	for i, v := range m.Vertices {
		y := v.Y*math.Cos(rx) - v.Z*math.Sin(rx)
		z := v.Y*math.Sin(rx) + v.Z*math.Cos(rx)
		m.Vertices[i].Y = y
		m.Vertices[i].Z = z
	}

	// Вращение вокруг оси Y
	for i, v := range m.Vertices {
		x := v.X*math.Cos(ry) + v.Z*math.Sin(ry)
		z := -v.X*math.Sin(ry) + v.Z*math.Cos(ry)
		m.Vertices[i].X = x
		m.Vertices[i].Z = z
	}

	// Вращение вокруг оси Z
	for i, v := range m.Vertices {
		x := v.X*math.Cos(rz) - v.Y*math.Sin(rz)
		y := v.X*math.Sin(rz) + v.Y*math.Cos(rz)
		m.Vertices[i].X = x
		m.Vertices[i].Y = y
	}
}

func getUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
