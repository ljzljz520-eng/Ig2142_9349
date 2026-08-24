package rules

type direction struct {
	rowDelta    int
	columnDelta int
}

var directions = [...]direction{
	{rowDelta: -1, columnDelta: -1},
	{rowDelta: -1, columnDelta: 0},
	{rowDelta: -1, columnDelta: 1},
	{rowDelta: 0, columnDelta: -1},
	{rowDelta: 0, columnDelta: 1},
	{rowDelta: 1, columnDelta: -1},
	{rowDelta: 1, columnDelta: 0},
	{rowDelta: 1, columnDelta: 1},
}

func inside(row, column int) bool {
	return row >= 0 && row < 8 && column >= 0 && column < 8
}
