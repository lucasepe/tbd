package table

import (
	"testing"
)

//
// Public Method
//

func TestDraw(t *testing.T) {
	tbl := &TextTable{}

	expected := `+------+----------+
| 名前 | ふりがな |
+------+----------+
| foo  | ふう     |
| hoge | ほげ     |
+------+----------+`

	tbl.SetHeader("名前", "ふりがな")

	tbl.AddRow("foo", "ふう")
	tbl.AddRow("hoge", "ほげ")

	got := tbl.Draw()
	if got != expected {
		t.Errorf("[got]\n%s\n\n[expected]\n%s\n", got, expected)
	}
}

func TestSetHeader(t *testing.T) {
	tbl := &TextTable{}

	err := tbl.SetHeader()
	if err == nil {
		t.Errorf("SetHeader should take one argument at least")
	}
}

func TestAddRow(t *testing.T) {
	tbl := &TextTable{}
	tbl.SetHeader("name", "age")

	err := tbl.AddRow("bob", "30", "182")
	if err == nil {
		t.Errorf("row length should be smaller than equal header length")
	}

	err = tbl.AddRow()
	if err == nil {
		t.Errorf("AddRow should take one argument at least")
	}
}

//
// Private Function/Methods
//

func Test_calcMaxHeight(t *testing.T) {
	input := []string{
		"hello", "apple\nmelon\norange", "1\n2",
	}

	got := calcMaxHeight(input)
	if got != 3 {
		t.Errorf("calcMaxHeight(%s) != 3(got=%d)", input, got)
	}
}

func Test_decideAlignment(t *testing.T) {
	got := decideAlignment("102948")
	if got != ALIGN_RIGHT {
		t.Errorf("decimal string of integer alighment is 'right'")
	}

	got = decideAlignment("01234")
	if got != ALIGN_RIGHT {
		t.Errorf("octal string of integer alighment is 'right'")
	}

	got = decideAlignment("ff")
	if got != ALIGN_RIGHT {
		t.Errorf("hex string without '0x' of integer alighment is 'right'")
	}

	got = decideAlignment("0xaabbccdd")
	if got != ALIGN_RIGHT {
		t.Errorf("hex string of integer alighment is 'right'")
	}

	got = decideAlignment("1.245")
	if got != ALIGN_RIGHT {
		t.Errorf("string of float alighment is 'right'")
	}

	got = decideAlignment("foo")
	if got != ALIGN_LEFT {
		t.Errorf("string  alighment is 'left'")
	}
}

func Test_stringsToTableRow(t *testing.T) {
	input := []string{
		"apple", "orange\nmelon\ngrape\nnuts", "peach\nbanana",
	}

	tableRows := stringsToTableRow(input)
	if len(tableRows) != 4 {
		t.Errorf("returned table height=%d(Expected 4)", len(tableRows))
	}

	for i, row := range tableRows {
		if len(row.cellUnits) != len(input) {
			t.Errorf("width of tableRows[%d]=%d(Expected %d)",
				i, len(row.cellUnits), len(input))
		}
	}
}

func Test_borderString(t *testing.T) {
	tbl := new(TextTable)
	tbl.maxWidths = []int{4, 5, 3, 2}

	expected := "+------+-------+-----+----+"

	border := tbl.borderString()
	if border != expected {
		t.Errorf("got %s(Expected %s)", border, expected)
	}

	tbl.maxWidths = []int{0}
	expected = "+--+"
	border = tbl.borderString()
	if border != expected {
		t.Errorf("got %s(Expected %s)", border, expected)
	}
}

func Test_formatCellUnit(t *testing.T) {
	cell := cellUnit{content: "apple", alignment: ALIGN_RIGHT}

	expected := " apple "
	got := formatCellUnit(&cell, 5)
	if got != expected {
		t.Errorf("got '%s'(Expected '%s')", got, expected)
	}

	expected = "      apple "
	got = formatCellUnit(&cell, 10)
	if got != expected {
		t.Errorf("got '%s'(Expected '%s')", got, expected)
	}

	cellLeft := cellUnit{content: "orange", alignment: ALIGN_LEFT}
	expected = " orange "
	got = formatCellUnit(&cellLeft, 6)
	if got != expected {
		t.Errorf("got '%s'(Expected '%s')", got, expected)
	}

	expected = " orange     "
	got = formatCellUnit(&cellLeft, 10)
	if got != expected {
		t.Errorf("got '%s'(Expected '%s')", got, expected)
	}
}

func Test_generateRowString(t *testing.T) {
	tbl := TextTable{}
	tbl.maxWidths = []int{8, 5}
	cells := []*cellUnit{
		{content: "apple", alignment: ALIGN_RIGHT},
		{content: "melon", alignment: ALIGN_RIGHT},
	}

	row := tableRow{cellUnits: cells, kind: ROW_CELLS}
	got := tbl.generateRowString(&row)

	expected := "|    apple | melon |"
	if got != expected {
		t.Errorf("got '%s'(Expected '%s')", got, expected)
	}
}
