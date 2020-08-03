package export

import (
	"github.com/tealeg/xlsx"
	"github.com/xifengzhu/eshop/initializers/setting"
	"strconv"
	"time"
)

func Exec(title string, headers []string, data [][]string) (string, error) {

	file := xlsx.NewFile()
	sheet, err := file.AddSheet(title)
	if err != nil {
		return "", err
	}

	// write header
	row := sheet.AddRow()
	var cell *xlsx.Cell
	for _, title := range headers {
		cell = row.AddCell()
		cell.Value = title
	}

	// write data row
	for _, values := range data {
		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	// set filename
	time := strconv.Itoa(int(time.Now().Unix()))
	filename := title + "-" + time + ".xlsx"
	fullPath := GetExcelFullPath() + filename

	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func GetExcelFullUrl(name string) string {
	return setting.Domain + "/" + GetExcelPath() + name
}

func GetExcelPath() string {
	return setting.PUBLIC_SAVE_PATH
}

func GetExcelFullPath() string {
	return setting.RuntimeRootPath + GetExcelPath()
}
