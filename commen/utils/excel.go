package utils

import (
	"fmt"
	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
	"log"
)

var (
	FilePath = "./tmp/"
)

func ReadXls(filename string) (res [][]string, err error) {
	if xlFile, err := xls.Open(FilePath+filename, "utf-8"); err == nil {
		fmt.Println(xlFile.Author)
		//第一个sheet
		sheet := xlFile.GetSheet(0)
		if sheet.MaxRow != 0 {
			temp := make([][]string, sheet.MaxRow)
			for i := 0; i < int(sheet.MaxRow); i++ {
				row := sheet.Row(i)
				data := make([]string, 0)
				if row.LastCol() > 0 {
					for j := 0; j < row.LastCol(); j++ {
						col := row.Col(j)
						data = append(data, col)
					}
					temp[i] = data
				}
			}
			res = append(res, temp...)
		}
	} else {
		log.Printf("read %s failed, err:%v", filename, err)
		return nil, err
	}
	return res, nil
}

func ReadXlsx(filename string) (res [][]string, err error) {
	if xlFile, err := xlsx.OpenFile(FilePath + filename); err == nil {
		for index, sheet := range xlFile.Sheets {
			//第一个sheet
			if index == 0 {
				temp := make([][]string, len(sheet.Rows))
				for k, row := range sheet.Rows {
					var data []string
					for _, cell := range row.Cells {
						data = append(data, cell.Value)
					}
					temp[k] = data
				}
				res = append(res, temp...)
			}
		}
	} else {
		log.Printf("read %s failed, err:%v", filename, err)
		return nil, err
	}
	return res, nil
}
