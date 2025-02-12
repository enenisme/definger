package utils

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

type FingerData struct {
	Protocol string
	Url      string
	Result   []string
	Title    string
}

// SaveExecl 保存指纹数据到Excel文件
// 参数:
//   - fingers: 指纹数据
//   - filename: 文件名
//
// 返回:
//   - error: 错误信息
func SaveExecl(fingers map[string]FingerData, filename string) error {
	file := excelize.NewFile()

	sheet := "LJ_Definger"
	file.NewSheet(sheet)
	file.SetCellValue(sheet, "A1", "Protocol")
	file.SetCellValue(sheet, "B1", "Url")
	file.SetCellValue(sheet, "C1", "Result")
	file.SetCellValue(sheet, "D1", "Title") // 修复了标题行的错误,将C1改为D1

	row := 2
	for _, finger := range fingers {
		// 将数据写入对应的单元格
		file.SetCellValue(sheet, fmt.Sprintf("A%d", row), finger.Protocol)
		file.SetCellValue(sheet, fmt.Sprintf("B%d", row), finger.Url)
		// 将Result数组转换为字符串后写入
		file.SetCellValue(sheet, fmt.Sprintf("C%d", row), strings.Join(finger.Result, ","))
		file.SetCellValue(sheet, fmt.Sprintf("D%d", row), finger.Title)
		row++
	}
	return file.SaveAs(filename)

}
