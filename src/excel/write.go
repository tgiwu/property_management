package excel

import (
	"fmt"
	"property/src/config"
	"property/src/types"
	"property/src/utils"

	"github.com/xuri/excelize/v2"
)

var (
	styleCellTitle,
	styleCellCompanyAndMonth,
	styleCellHead,
	styleCellHeadSmallDate,
	styleCellHeadSmallWeekDay,
	styleCellData,
	styleCellHeadCross,
	lastDayInMonth int
	lastColStr,
	dateStart,
	dateEnd string
)

func initVar(excel *excelize.File) {
	lastDayInMonth = utils.LastDayInMonth(config.GetConfig().TargetYear, config.GetConfig().TargetMonth)
	colBytes := []byte{65, byte(lastDayInMonth + 5 - 26 + 64)}
	lastColStr = string(colBytes)
	dateStart = "E"
	dateEnd = string([]byte{65, byte(lastDayInMonth + 5 - 26 + 64 - 1)})
	setUpCellStyle(excel)
}

func setUpCellStyle(excel *excelize.File) {
	styleCellTitle, _ = excel.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "left", Color: "#FFFFFF", Style: 1},
			{Type: "right", Color: "#FFFFFF", Style: 1},
			{Type: "top", Color: "#FFFFFF", Style: 1},
			{Type: "bottom", Color: "#FFFFFF", Style: 1}},
		Font: &excelize.Font{
			Bold:   true,
			Color:  "#000000",
			Family: "黑体",
			Size:   16,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFFFF"}, Pattern: 1},
	})

	styleCellCompanyAndMonth, _ = excel.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "left", Color: "#FFFFFF", Style: 1},
			{Type: "right", Color: "#FFFFFF", Style: 1},
			{Type: "top", Color: "#FFFFFF", Style: 1},
			{Type: "bottom", Color: "#FFFFFF", Style: 1}},
		Font: &excelize.Font{
			Bold:   true,
			Color:  "#000000",
			Family: "宋体",
			Size:   14,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFFFF"}, Pattern: 1},
	})

	styleCellHead, _ = excel.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1}},
		Font: &excelize.Font{
			Bold:   true,
			Color:  "#000000",
			Family: "宋体",
			Size:   11,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFFFF"}, Pattern: 1},
	})

	styleCellHeadSmallDate, _ = excel.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1}},
		Font: &excelize.Font{
			Bold:   true,
			Color:  "#000000",
			Family: "黑体",
			Size:   10,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFFFF"}, Pattern: 1},
	})

	styleCellHeadSmallWeekDay, _ = excel.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1}},
		Font: &excelize.Font{
			Bold:   false,
			Color:  "#000000",
			Family: "黑体",
			Size:   10,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFFFF"}, Pattern: 1},
	})

	styleCellData, _ = excel.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1}},
		Font: &excelize.Font{
			Bold:   true,
			Color:  "#000000",
			Family: "黑体",
			Size:   10,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFFFF"}, Pattern: 1},
	})

	styleCellHeadCross, _ = excel.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "diagonalDown", Color: "#000000", Style: 1}},
		Font: &excelize.Font{
			Bold:   true,
			Color:  "#000000",
			Family: "宋体",
			Size:   10,
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFFFF"}, Pattern: 1},
	})
}

func ConstructExcel(data *map[string][]types.Attendance) error {
	excel := excelize.NewFile()
	defer excel.Close()

	initVar(excel)

	for key, atts := range *data {
		constructSingleExcel(excel, key, &atts)
	}

	return nil
}

func constructSingleExcel(excel *excelize.File, sheetName string, data *[]types.Attendance) error {
	if len(*data) == 0 || len(sheetName) == 0 {
		return nil
	}

	_, err := excel.NewSheet(sheetName)
	if err != nil {
		return err
	}

	return nil
}

func title(excel *excelize.File, sheetName string) error {

	excel.MergeCell(sheetName, "A1", lastColStr+"1")
	excel.SetCellStr(sheetName, "A1", fmt.Sprintf("合同名称：《综合服务（%s）承包合同》", sheetName))
	excel.SetCellStyle(sheetName, "A1", lastColStr+"1", styleCellTitle)

	excel.MergeCell(sheetName, "A2", lastColStr+"2")
	excel.SetCellStr(sheetName, "A2", fmt.Sprintf("%d年度考勤记录表", config.GetConfig().TargetYear))
	excel.SetCellStyle(sheetName, "A2", lastColStr+"2", styleCellTitle)

	return nil
}

func companyAndMonth(excel *excelize.File, sheetName string) error {
	excel.MergeCell(sheetName, "A3", lastColStr+"3")
	excel.SetCellStr(sheetName, "A3", fmt.Sprintf("用工单位名称：%s", ""))
	excel.SetCellStyle(sheetName, "A3", lastColStr+"3", styleCellCompanyAndMonth)

	excel.MergeCell(sheetName, "A4", lastColStr+"4")
	excel.SetCellStr(sheetName, "A4", fmt.Sprintf("月份：%d 月", config.GetConfig().TargetMonth))
	excel.SetCellStyle(sheetName, "A4", lastColStr+"4", styleCellCompanyAndMonth)
	return nil
}

func head(excel *excelize.File, sheetName string) error {
	excel.MergeCell(sheetName, "A5", "A6")
	excel.SetCellStr(sheetName, "A5", "序号")
	excel.SetCellStyle(sheetName, "A5", "A5", styleCellHead)

	excel.MergeCell(sheetName, "B5", "B6")
	excel.SetCellStr(sheetName, "B5", "日期\n姓名")
	excel.SetCellStyle(sheetName, "B5", "B5", styleCellHeadCross)

	excel.MergeCell(sheetName, "C5", "C6")
	excel.SetCellStr(sheetName, "C5", "岗位")
	excel.SetCellStyle(sheetName, "C5", "C5", styleCellHead)

	excel.MergeCell(sheetName, "D5", "D6")
	excel.SetCellStr(sheetName, "D5", "出勤\n天数")
	excel.SetCellStyle(sheetName, "D5", "D5", styleCellHead)

	for i := 0; i < lastDayInMonth; i++ {
		b := 'A' + 4 + i
		var cellIndexStr string
		if b <= 'Z' {
			cellIndexStr = string(byte(b))
		} else {
			cellIndexStr = string([]byte{'A',byte(b)})
		}

		excel.SetCellStr(sheetName, cellIndexStr+"5", string(rune(i+1)))
		excel.SetCellStyle(sheetName, cellIndexStr+"5", cellIndexStr+"5", styleCellHeadSmallDate)

		excel.SetCellStr(sheetName, cellIndexStr+"6", utils.WeekDay(config.GetConfig().TargetYear, config.GetConfig().TargetMonth, i+1))
		excel.SetCellStyle(sheetName, cellIndexStr+"6", cellIndexStr+"6", styleCellHeadSmallWeekDay)
	}

	excel.MergeCell(sheetName, lastColStr+"5", lastColStr+"6")
	excel.SetCellStr(sheetName, lastColStr+"5", "本人签字")
	excel.SetCellStyle(sheetName, lastColStr+"5", lastColStr+"6", styleCellHead)

	return nil
}
