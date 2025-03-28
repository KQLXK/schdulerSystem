package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/extrame/xls"
	xlsReader "github.com/shakinm/xlsReader/xls"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	FilePath       = "./tmp/"
	ExcelFormatErr = errors.New("不支持的文件格式")
)

func ReadExcel(filename string) ([][]string, error) {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".xls":
		return readXls(filename)
	case ".xlsx":
		return readXlsx(filename)
	default:
		return nil, ExcelFormatErr
	}
}

// 使用shakinm/xlsReader库实现，支持 '.xls' 格式
func readXls(filename string) ([][]string, error) {
	outputFile := filepath.Join("./tmp", fmt.Sprintf("%d_%v.csv", time.Now().Unix(), filename))
	workbook, err := xlsReader.OpenFile(FilePath + filename)
	if err != nil {
		log.Println("open excel file failed, err:", err)
	}

	sheet, err := workbook.GetSheet(0)
	maxRows := sheet.GetNumberRows()
	result := make([][]string, 0, maxRows)

	for i := 0; i <= maxRows; i++ {
		if row, err := sheet.GetRow(i); err == nil {
			if cols := row.GetCols(); err == nil {
				rowData := make([]string, len(cols))
				for j, col := range cols {
					rowData[j] = col.GetString()
				}
				result = append(result, rowData)
			}
		}
	}

	// 写入 CSV 文件
	if err := WriteToCSV(result, outputFile); err != nil {
		log.Println("write to csv failed, err :", err)
		return nil, fmt.Errorf("CSV 写入失败: %v", err)
	}

	log.Printf("成功转换 .xls 文件到: %s", outputFile)
	return result, nil

}

// 使用 extrame/xls 库完成，但会出现乱码bug
func readxls(filename string) ([][]string, error) {
	// 生成带时间戳的输出文件名
	outputFile := filepath.Join("./tmp", fmt.Sprintf("%d_%s.csv", time.Now().Unix(), filename))

	// 打开 .xls 文件
	xlFile, err := xls.Open(filepath.Join(FilePath, filename), "utf-8")
	if err != nil {
		log.Println("open excel file failed, err:", err)
		return nil, fmt.Errorf("打开 .xls 文件失败: %v", err)
	}

	// 获取第一个工作表
	sheet := xlFile.GetSheet(0)
	if sheet == nil {
		return nil, fmt.Errorf("工作表不存在或文件为空")
	}

	// 预分配结果切片
	maxRow := int(sheet.MaxRow)
	result := make([][]string, 0, maxRow)

	// 遍历所有行
	for i := 0; i < maxRow; i++ {
		row := sheet.Row(i)
		if row == nil {
			continue // 跳过空行
		}

		// 读取单元格数据
		colCount := row.LastCol()
		rowData := make([]string, colCount)
		for j := 0; j < colCount; j++ {
			rowData[j] = row.Col(j)
		}
		result = append(result, rowData)
	}

	// 写入 CSV 文件
	if err := WriteToCSV(result, outputFile); err != nil {
		log.Printf("CSV 写入失败: %v", err)
		return nil, fmt.Errorf("CSV 写入失败: %v", err)
	}

	log.Printf("成功转换 .xls 文件到: %s", outputFile)
	return result, nil
}

// ReadXlsx 使用 excelize 读取 Excel 文件（支持 .xlsx 格式）
func readXlsx(filename string) (res [][]string, err error) {
	// 生成输出路径
	outputFile := filepath.Join("./tmp", fmt.Sprintf("%d_%s.csv", time.Now().Unix(), filename))

	// 打开 Excel 文件
	f, err := excelize.OpenFile(filepath.Join(FilePath, filename))
	if err != nil {
		log.Printf("打开文件失败: %v", err)
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("关闭文件失败: %v", err)
		}
	}()

	// 获取第一个工作表
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return nil, fmt.Errorf("文件中没有工作表")
	}
	firstSheet := sheetList[0]

	// 读取所有行数据
	rows, err := f.GetRows(firstSheet)
	if err != nil {
		log.Printf("读取工作表数据失败: %v", err)
		return nil, err
	}

	// 写入 CSV 文件
	if err := WriteToCSV(rows, outputFile); err != nil {
		log.Printf("CSV 写入失败: %v", err)
		return nil, err
	}

	log.Printf("文件已转换保存至: %s", outputFile)
	return rows, nil
}

// WriteToCSV 通用 CSV 写入函数
func WriteToCSV(data [][]string, outputPath string) error {
	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建并写入 CSV 文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("写入行失败: %v", err)
		}
	}
	return nil
}
