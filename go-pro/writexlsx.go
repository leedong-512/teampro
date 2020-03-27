package main


import (
	"bytes"
	"encoding/csv"
	"bufio"
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"io"
	"strings"
)

func main() {
	var xlsxfileName    string
	//var xlsfileName    string
	var csvfileName    string
	var txtfileName    string
	var value string
	var key int
	flag.StringVar(&xlsxfileName, "xlsx", "", "xlsx文件地址不能为空")
	flag.StringVar(&csvfileName, "csv", "", "csv文件地址不能为空")
	flag.StringVar(&txtfileName, "txt", "", "txt文件地址不能为空")
	flag.StringVar(&value, "value", "", "value不能为空")
	flag.StringVar(&key, "key", "", "key必须为整数")
	flag.Parse()
	if xlsxfileName != "" && value != "" && key >= 0 {
		writeXlsx(xlsxfileName)
	} else if csvfileName != "" && value != "" && key >= 0 {
		writeCsv(csvfileName)
	} else if txtfileName != "" && value != "" && key >= 0 {
		writeTxt(txtfileName)
	} else {
		fmt.Printf("%s\n", "请填写正确的文件地址")
	}

}

//  编辑xlsx单元格内容
func writeXlsx(filename string)  {
	_, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	f, _ := excelize.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := f.SetCellValue("Sheet1", "A1", "lidsdong"); err != nil {
		fmt.Println(err)
	}
	if err := f.SaveAs(filename); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", "ok")
}

//编辑csv
func writeCsv(filename string)  {
	fs, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer fs.Close()

	r := csv.NewReader(fs)
	i := 0
	//针对大文件，一行一行的读取文件
	var newContent [][]string
    for {
        row, err := r.Read()
        if err != nil && err != io.EOF {
            log.Fatal(err)
        }
        if err == io.EOF {
            break
		}
		if i == 0 {
			row[1] = "lidong"
		}
		i++
		newContent = append(newContent, row)
		doCsv(filename, newContent)
	}
	// fmt.Println(newContent)
}

func doCsv(filename string, row [][]string) {
	buf := new(bytes.Buffer)
	r2 := csv.NewWriter(buf)
	for i:=0;i<len(row);i++ {
		r2.Write(row[i])
	}
	
	r2.Flush()
	fout,err := os.Create(filename)
    if err != nil {
        log.Fatal(err)
    }
	defer fout.Close()
	if err != nil {
		fmt.Println(filename,err)
		return
	}
	fout.WriteString(buf.String())
}

// 编辑txt
func writeTxt(filename string) {
	f, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    buf := bufio.NewReader(f)
	var result  []string
	line := ""
	i := 0
    for {
        a, _, c := buf.ReadLine()
        if c == io.EOF {
            break
		}
		if i == 0 {
			result = strings.Split(string(a), "\t")
			result[1] = "21323"
			line += strings.Join(result, "\t") + "\n"
		} else {
			line += string(a) + "\n"
		}
		i++
	}
	doTxt(filename, strings.Trim(line, "\n"))
	fmt.Println(strings.Trim(line, "\n"))
}

func doTxt(filename string, line string) {
	f, err := os.Create(filename)
    if err != nil {
        fmt.Println(err)
        return
    }
    _, err = f.WriteString(line)
    if err != nil {
        fmt.Println(err)
        f.Close()
        return
    }
    // fmt.Println(l, "bytes written successfully")
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
}
