package types

type Attendance struct {
	No   int
	Name string
	Dept string

	Att       [31]int
	LookUpKey string
	FileMd5   string
}

const (
	ATTENDENT_COL_NO = iota
	ATTENDENT_COL_NAME
	ATTENDENT_COL_DEPT
)

const (
	ATTENDENT_ROW_TITLE = iota
	ATTENDENT_ROW_TIME
	ATTENDENT_ROW_TABLE_TITLE_1
	ATTENDENT_ROW_TABLE_TITLE_2
)

const (
	BYTE_MORNING = 0b1
	BYTE_EVENING = 0b10
)
