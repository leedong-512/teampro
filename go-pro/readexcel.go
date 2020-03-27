package main

import (
    "flag"
    "fmt"
    "github.com/360EntSecGroup-Skylar/excelize"
    "github.com/extrame/xls"
    "github.com/tealeg/xlsx"
    "log"
    "os"
    "strconv"
    "strings"
    "time"
)

func main() {
    var xlsxfileName    string
    var xlsfileName    string
    flag.StringVar(&xlsxfileName, "xlsx", "", "读取xlsx文件地址不能为空")
    flag.StringVar(&xlsfileName, "xls", "", "读取xls文件地址不能为空")
    flag.Parse()
    limit := 100
    if xlsxfileName != "" {
        getXlsx(xlsxfileName, limit)
    } else if xlsfileName != "" {
        getXls(xlsfileName, limit)
    } else {
        fmt.Printf("%s\n", "请填写正确的文件地址")
    }
}

/*
获取 xlsx格式的数据
*/
func getXlsx(filename string, limit int) {
   _, err := os.Stat(filename)
   if err != nil {
       log.Fatal(err)
   }
   xlFile, err := xlsx.OpenFileWithRowLimit(filename, limit)
   if err != nil {
       log.Fatal(err)
   }
   for k, row := range xlFile.Sheets[0].Rows {
       for _, cell := range row.Cells {
           text:= ""
           if k > 0 {
               if cell.IsTime() {
                   val := cell.Value
                   str := cell.String()
                   if strings.Index(str, "\\ AM") > 0 {
                       f_inde := strings.Index(str, "]")
                       l_inde := strings.Index(str, "\\")
                       text = string([]byte(str)[f_inde+1:l_inde])
                   } else {
                       ind := strings.Index(cell.String(), ",")
                       if ind > 0 {
                           text = convertToFormatDate(val)
                       } else {
                           text = convertToFormatDay(val)
                       }
                   }
               } else {
                   nfmt := cell.GetNumberFormat()
                   if nfmt == "h:mm:ss;@" {
                       cell.SetFormat(clearFmt(nfmt))
                       v, err := cell.FormattedValue()
                       if err != nil {
                           fmt.Println(err)
                       }
                       text = v
                   } else if nfmt == "yyyy/m/d;@" {
                       cell.SetFormat(clearFmt(nfmt))
                       v, err := cell.FormattedValue()
                       if err != nil {
                           fmt.Println(err)
                       }
                       text = v
                   } else if nfmt == "yyyy\"年\"m\"月\"d\"日\";@" {
                       cell.SetFormat(clearFmt(nfmt))
                       v, err := cell.FormattedValue()
                       if err != nil {
                           fmt.Println(err)
                       }
                       text = strings.ReplaceAll(v, "\"", "")
                   } else {
                       text = cell.String()
                   }
               }
           } else {
               text = cell.String()
           }
           fmt.Printf("%s\t", text)
       }
       fmt.Printf("\n")
   }
}

func setCellType() {

}

func clearFmt(nfmt string) string {
    nfmt = strings.ReplaceAll(nfmt, ";@", "")
    return nfmt
}

func getxlsx2(filename string, limit int) {
    xlsx, err := excelize.OpenFile(filename)
    if err != nil {
        fmt.Println(err)
        return
    }
    // Get all the rows in the Sheet1.
    rows, _ := xlsx.GetRows("Sheet1")
    for k, row := range rows {
        for _, colCell := range row {
            if k == 0 {
                fmt.Print(colCell, "\t")
            } else {
                fmt.Print(colCell, "\t")
            }
        }
        fmt.Println()
    }
}
func convertToFormatDay(excelDaysString string)string{
    baseDiffDay := 38719
    curDiffDay := excelDaysString
    b,_ := strconv.Atoi(curDiffDay)
    realDiffDay := b - baseDiffDay
    realDiffSecond := realDiffDay * 24 * 3600
    baseOriginSecond := 1136185445
    resultTime := time.Unix(int64(baseOriginSecond + realDiffSecond), 0).Format("2006/01/02")

    return resultTime
}
func convertToFormatDate(excelDaysString string)string{
    baseDiffDay := 38719
    curDiffDay := excelDaysString
    b,_ := strconv.Atoi(curDiffDay)
    realDiffDay := b - baseDiffDay
    realDiffSecond := realDiffDay * 24 * 3600
    baseOriginSecond := 1136185445
    resultTime := time.Unix(int64(baseOriginSecond + realDiffSecond), 0).Format("2006年01月02日")

    return resultTime
}
/*
获取 xls格式的数据
*/
func getXls(filename string, limit int) {
    _, err := os.Stat(filename)
    if err != nil {
        log.Fatal(err)
    }
    xlsfile, closer, err := xls.OpenWithCloser(filename, "utf-8")
    if err != nil {
        log.Fatal(err)
    }
    sheet := xlsfile.GetSheet(0)
    MaxRow := uint(sheet.MaxRow + 1)
    if int(MaxRow) < limit {
        limit = int(MaxRow)
    }
    for i := 0; i < limit; i++ {
        xlsRow := sheet.Row(i)
        for j := 0; j < int(xlsRow.LastCol()); j++ {
            fmt.Printf("%s\t", xlsRow.Col(j))
        }
        fmt.Printf("\n")
    }
    closer.Close()
}
