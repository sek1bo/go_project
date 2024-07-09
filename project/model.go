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

// Line структура для хранения линий
type Line struct {
	Start, End int
}

// Model структура для хранения всех данных модели
type Model struct {
	Vertices   []Vertex
	TexCoords  []TexCoord
	Normals    []Normal
	Faces      []Face
	Lines      []Line
	Scaler     float64
	TranslateX float64
	TranslateY float64
	TranslateZ float64
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
		case strings.HasPrefix(line, "l "):
			l, err := m.parseLine(line)
			if err == nil {
				m.Lines = append(m.Lines, l)
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

// Метод для парсинга линии
func (m *Model) parseLine(line string) (Line, error) {
	var l Line
	parts := strings.Fields(line)[1:]
	if len(parts) != 2 {
		return l, fmt.Errorf("invalid line: %s", line)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return l, fmt.Errorf("invalid start vertex index: %v", err)
	}
	l.Start = start - 1

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return l, fmt.Errorf("invalid end vertex index: %v", err)
	}
	l.End = end - 1

	return l, nil
}

// Метод для преобразования линий в грани
func (m *Model) ConvertLinesToFaces() {
	lineMap := make(map[int][]int)

	// Формирование карты связей между вершинами на основе линий
	for _, line := range m.Lines {
		lineMap[line.Start] = append(lineMap[line.Start], line.End)
		lineMap[line.End] = append(lineMap[line.End], line.Start)
	}

	visitedLines := make(map[int]bool)
	for start := range lineMap {
		if !visitedLines[start] {
			visited := make(map[int]bool)
			face := Face{}
			current := start
			for {
				// Добавляем текущую линию в грань
				face.Vertices = append(face.Vertices, current)
				visitedLines[current] = true
				visited[current] = true
				next := -1
				for _, neighbor := range lineMap[current] {
					if !visitedLines[neighbor] {
						next = neighbor
						break
					}
				}
				if next == -1 {
					break
				}
				current = next
				// Проверяем, достигли ли мы начальной точки
				if current == start {
					break
				}
			}
			// Добавляем грань в список, если у неё больше двух вершин
			if len(face.Vertices) > 2 {
				m.Faces = append(m.Faces, face)
			}
		}
	}
}

// Метод для масштабирования модели по заданным коэффициентам масштабирования по осям X, Y, Z
func (m *Model) Scale(sx, sy, sz float64) {
	m.Scaler *= sx
	m.TranslateX *= sx
	m.TranslateY *= sy
	m.TranslateZ *= sz

	// Коррекция координат вершин
	for i := range m.Vertices {
		m.Vertices[i].X *= sx
		m.Vertices[i].Y *= sy
		m.Vertices[i].Z *= sz
	}
}

// Метод для трансляции модели по заданным сдвигам по осям X, Y, Z
func (m *Model) Translate(tx, ty, tz float64) {
	m.TranslateX += tx
	m.TranslateY += ty
	m.TranslateZ += tz

	// Коррекция координат вершин
	for i := range m.Vertices {
		m.Vertices[i].X += tx
		m.Vertices[i].Y += ty
		m.Vertices[i].Z += tz
	}
}

// Метод для вращения модели вокруг осей X, Y и Z
func (m *Model) Rotate(rx, ry, rz float64) {
	cosRx, sinRx := math.Cos(rx), math.Sin(rx)
	cosRy, sinRy := math.Cos(ry), math.Sin(ry)
	cosRz, sinRz := math.Cos(rz), math.Sin(rz)

	for i := range m.Vertices {
		// Вращение вокруг оси X
		y := m.Vertices[i].Y*cosRx - m.Vertices[i].Z*sinRx
		z := m.Vertices[i].Y*sinRx + m.Vertices[i].Z*cosRx
		m.Vertices[i].Y, m.Vertices[i].Z = y, z

		// Вращение вокруг оси Y
		x := m.Vertices[i].X*cosRy + m.Vertices[i].Z*sinRy
		z = -m.Vertices[i].X*sinRy + m.Vertices[i].Z*cosRy
		m.Vertices[i].X, m.Vertices[i].Z = x, z

		// Вращение вокруг оси Z
		x = m.Vertices[i].X*cosRz - m.Vertices[i].Y*sinRz
		y = m.Vertices[i].X*sinRz + m.Vertices[i].Y*cosRz
		m.Vertices[i].X, m.Vertices[i].Y = x, y
	}
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
		writer.WriteString(fmt.Sprintf("v %f %f %f\n", v.X, v.Y, v.Z))
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

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
