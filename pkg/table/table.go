package table

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
)

type cellAlignment int

const (
	ALIGN_LEFT cellAlignment = iota
	ALIGN_RIGHT
)

type rowType int

const (
	ROW_LINE rowType = iota
	ROW_CELLS
)

type cellUnit struct {
	content   string
	alignment cellAlignment
}

type tableRow struct {
	cellUnits []*cellUnit
	kind      rowType
}

type tableLine struct{}

type TextTable struct {
	header    []*tableRow
	rows      []*tableRow
	width     int
	maxWidths []int
}

func (t *TextTable) updateColumnWidth(rows []*tableRow) {
	for _, row := range rows {
		for i, unit := range row.cellUnits {
			width := stringWidth(unit.content)
			if t.maxWidths[i] < width {
				t.maxWidths[i] = width
			}
		}
	}
}

/*

SetHeader adds header row from strings given

*/
func (t *TextTable) SetHeader(headers ...string) error {
	if len(headers) == 0 {
		return errors.New("no headers")
	}

	columnSize := len(headers)

	t.width = columnSize
	t.maxWidths = make([]int, columnSize)

	rows := stringsToTableRow(headers)
	t.updateColumnWidth(rows)

	t.header = rows

	return nil
}

/*

AddRow adds column from strings given

*/
func (t *TextTable) AddRow(strs ...string) error {
	if len(strs) == 0 {
		return errors.New("no rows")
	}

	if len(strs) > t.width {
		return errors.New("row width should be less than header width")
	}

	padded := make([]string, t.width)
	copy(padded, strs)
	rows := stringsToTableRow(padded)
	t.rows = append(t.rows, rows...)

	t.updateColumnWidth(rows)

	return nil
}

/*

AddRowLine adds row border

*/
func (t *TextTable) AddRowLine() error {
	rowLine := &tableRow{kind: ROW_LINE}
	t.rows = append(t.rows, rowLine)

	return nil
}

func (t *TextTable) borderString() string {
	borderString := "+"
	margin := 2

	for _, width := range t.maxWidths {
		for i := 0; i < width+margin; i++ {
			borderString += "-"
		}
		borderString += "+"
	}

	return borderString
}

func stringsToTableRow(strs []string) []*tableRow {
	maxHeight := calcMaxHeight(strs)
	strLines := make([][]string, maxHeight)

	for i := 0; i < maxHeight; i++ {
		strLines[i] = make([]string, len(strs))
	}

	alignments := make([]cellAlignment, len(strs))
	for i := range strs {
		alignments[i] = ALIGN_LEFT // decideAlignment(str)
	}

	for i, str := range strs {
		divideds := strings.Split(str, "\n")
		for j, line := range divideds {
			strLines[j][i] = line
		}
	}

	rows := make([]*tableRow, maxHeight)
	for j := 0; j < maxHeight; j++ {
		row := new(tableRow)
		row.kind = ROW_CELLS
		for i := 0; i < len(strs); i++ {
			content := strLines[j][i]
			unit := &cellUnit{content: content}
			unit.alignment = alignments[i]
			row.cellUnits = append(row.cellUnits, unit)
		}

		rows[j] = row
	}

	return rows
}

var hexRegexp = regexp.MustCompile("^0x")

func decideAlignment(str string) cellAlignment {
	// decimal/octal number
	_, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return ALIGN_RIGHT
	}

	// hex number
	_, err = strconv.ParseInt(str, 16, 64)
	if err == nil {
		return ALIGN_RIGHT
	}

	if hexRegexp.MatchString(str) {
		tmp := str[2:]
		_, err := strconv.ParseInt(tmp, 16, 64)
		if err == nil {
			return ALIGN_RIGHT
		}
	}

	_, err = strconv.ParseFloat(str, 64)
	if err == nil {
		return ALIGN_RIGHT
	}

	return ALIGN_LEFT
}

func calcMaxHeight(strs []string) int {
	max := -1

	for _, str := range strs {
		lines := strings.Split(str, "\n")
		height := len(lines)
		if height > max {
			max = height
		}
	}

	return max
}

func stringWidth(str string) int {
	return runewidth.StringWidth(str)
}

/*

Draw constructs text table from receiver and returns it as string

*/
func (t *TextTable) Draw() string {
	drawedRows := make([]string, len(t.header)+len(t.rows)+3)
	index := 0

	border := t.borderString()

	// top line
	drawedRows[index] = border
	index++

	for _, row := range t.header {
		drawedRows[index] = t.generateRowString(row)
		index++
	}

	drawedRows[index] = border
	index++

	for _, row := range t.rows {
		var rowStr string
		if row.kind == ROW_CELLS {
			rowStr = t.generateRowString(row)
		} else {
			rowStr = border
		}
		drawedRows[index] = rowStr
		index++
	}

	// bottom line
	if len(t.rows) != 0 {
		drawedRows[index] = border
		index++
	}

	return strings.Join(drawedRows[:index], "\n")
}

func formatCellUnit(unit *cellUnit, maxWidth int) string {
	str := unit.content
	width := stringWidth(unit.content)

	padding := strings.Repeat(" ", maxWidth-width)

	var ret string
	if unit.alignment == ALIGN_RIGHT {
		ret = padding + str
	} else {
		ret = str + padding
	}

	return " " + ret + " "
}

func (t *TextTable) generateRowString(row *tableRow) string {
	separator := "|"

	str := separator
	for i, unit := range row.cellUnits {
		str += formatCellUnit(unit, t.maxWidths[i])
		str += separator
	}

	return str
}
