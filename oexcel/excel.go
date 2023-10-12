package oexcel

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const stName = "Sheet1"

type coder struct {
	f              *excelize.File
	cells          [][]string
	rowIdx         int
	colIdx         int
	endFlag        bool
	sheetName      string
	isIgnoreHeader bool
}

// XExcelReader excel 读解码器
// XExcelWriter excel 读解码器
type XExcelReader coder
type XExcelWriter coder

// NewPicXExcelReader 获取一个能够获取图片的excel读解码器
// picColNames 图片所在的列名 for example ： ”AA“
// 单元格默认可以存放多张图片，解析的图片以 url1/url2 的 string 格式存储
func NewPicXExcelReader(b []byte, picColNames []string) (*XExcelReader, error) {
	r := bytes.NewReader(b)
	fl, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}

	// 暂且只取第一个 sheet
	st := fl.GetSheetName(0)
	if st == "" {
		return nil, errors.New("get excel file first sheet fail")
	}

	cl, err := fl.GetRows(st)
	if err != nil {
		return nil, err
	}

	PicIndexMap := make(map[string]int)
	for _, colName := range picColNames {
		PicIndexMap[colName] = excelNum2Digit(colName)
	}
	for i := 0; i < len(cl); i++ {
		for _, colName := range picColNames {
			pics, err := fl.GetPictures(st, colName+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}
			if len(pics) == 0 {
				continue
			}
			var picsURL strings.Builder
			for idx, pic := range pics {
				name := fmt.Sprintf("./temp/image%s%d-%d%s", colName, i+1, idx+1, pic.Extension)
				if err := os.WriteFile(name, pic.File, 0644); err != nil {
					return nil, err
				}
				_, err := picsURL.WriteString(name + "/")
				if err != nil {
					return nil, err
				}
			}
			colNum := PicIndexMap[colName]
			cl[i][colNum] = picsURL.String()
		}
	}

	return &XExcelReader{
		f:              fl,
		cells:          cl,
		rowIdx:         0,
		colIdx:         0,
		endFlag:        false,
		sheetName:      st,
		isIgnoreHeader: false,
	}, nil
}

func excelNum2Digit(excelNum string) int {
	digit := 0
	for _, ch := range excelNum {
		digit = digit*26 + (int(ch) - 'A' + 1)
	}
	return digit - 1
}

// UnMarshal data 必须是个指针
func (d *XExcelReader) UnMarshal(data interface{}) (e error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("data must be a pri and not nil")
	}
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("parse error: %+v", err)
		}
	}()

	v = v.Elem()
	if !v.IsValid() {
		panic("cannot serialize (deserialize) zero value")
	}
	if v.Kind() == reflect.Slice {
		rows := 0
		if d.isIgnoreHeader && len(d.cells) > 0 {
			rows = len(d.cells) - 1
		} else {
			rows = len(d.cells)
		}

		sc := reflect.MakeSlice(v.Type(), rows, rows)
		i := 0
		for d.Valid() {
			d.Value(sc.Index(i))
			if d.Valid() {
				i++
			}
		}
		v.Set(sc)
	} else {
		//d.Value()
	}
	return nil
}

func (d *XExcelReader) IgnoreHeader() {
	d.nextLine()
	d.isIgnoreHeader = true
}

func (d *XExcelReader) nextLine() {

}
func (d *XExcelReader) Valid() bool {
	return true
}

func (d *XExcelReader) Value(v reflect.Value) {
	//if !v.IsValid() {
	//	panic("cannot serialize (deserialize) zero value")
	//}
	//switch v.Kind() {
	//case reflect.Array:
	//	i := 0
	//	for d.Valid() {
	//		d.Value(v.Index(i))
	//		if d.Valid() {
	//			i++
	//		}
	//	}
	//}
	//
	//case reflect.Slice:

}
