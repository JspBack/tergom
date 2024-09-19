package level

const (
	NOTHING = 0
	WALL    = 10
)

type Level struct {
	Width, Height int
	Data          [][]byte
}

func NewLevel(width, height int) *Level {
	data := make([][]byte, height)
	for h := 0; h < height; h++ {
		data[h] = make([]byte, width)
		for w := 0; w < width; w++ {
			if w == 0 || w == width-1 || h == 0 || h == height-1 {
				data[h][w] = WALL
			} else {
				data[h][w] = NOTHING
			}
		}
	}
	return &Level{
		Width:  width,
		Height: height,
		Data:   data,
	}
}

func (l *Level) IsWall(x, y int) bool {
	if x < 0 || x >= l.Width || y < 0 || y >= l.Height {
		return true
	}
	return l.Data[y][x] == WALL
}
