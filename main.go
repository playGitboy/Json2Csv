package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ecoshub/jin"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
)

// 将任意文本编码转为utf8编码(提升兼容性)
func ToUtf8(s []byte) []byte {
	d := chardet.NewTextDetector()
	var rs *chardet.Result
	var err1, err2 error
	if len(s) > 1024 {
		if utf8.Valid(s[:1024]) {
			return s
		}
		rs, err1 = d.DetectBest(s[:1024])
	} else {
		if utf8.Valid(s) {
			return s
		}
		rs, err1 = d.DetectBest(s)
	}

	var maps map[string]string = make(map[string]string)
	maps = map[string]string{
		"Shift_JIS":    "shift_jis",
		"EUC-JP":       "euc-jp",
		"EUC-KR":       "euc-kr",
		"Big5":         "big5",
		"GB18030":      "gb18030",
		"ISO-8859-2 ":  "iso-8859-2",
		"ISO-8859-5":   "iso-8859-5",
		"ISO-8859-6":   "iso-8859-6",
		"ISO-8859-7":   "iso-8859-7",
		"ISO-8859-8":   "iso-8859-8",
		"ISO-8859-8-I": "iso-8859-8-i",
		"ISO-8859-9":   "iso-8859-10",
		"windows-1256": "windows-1256",
		"windows-1251": "windows-1251",
		"KOI8-R":       "koi8-r",
		"ISO-2022-JP":  "iso-2022-jp",
		"UTF-16BE ":    "utf-16be",
		"UTF-16LE ":    "utf-16le",
	}

	ct := maps[rs.Charset]
	if ct == "" || err1 != nil {
		_, name, b := charset.DetermineEncoding([]byte(s), "utf-8")
		if b {
			return s
		}
		ct = name
	}

	byteReader := bytes.NewReader(s)
	reader, err1 := charset.NewReaderLabel(ct, byteReader)
	r, err2 := io.ReadAll(reader)

	if err1 != nil || err2 != nil {
		return s
	}
	return r
}

func readJSONFile(bJsonFile []byte) (*gmap.ListMap, error) {
	jsonListMap := gmap.NewListMap(true)
	if err := json.Unmarshal(bJsonFile, &jsonListMap); err != nil {
		return nil, err
	}
	return jsonListMap, nil
}

func getMaxNode(listMap *gmap.ListMap) string {
	maxSize := 0
	maxNode := ""

	listMap.Iterator(func(key interface{}, value interface{}) bool {
		if subList, ok := value.([]interface{}); ok {
			if len(subList) > maxSize {
				maxSize = len(subList)
				maxNode = key.(string)
			}
		}
		return true
	})

	return maxNode
}

// 兼容非数组对象数据提取{"items":{"a1":{"title":"one","name":"test"},"b2":{"title":"two","name":"test2"}}}
func writeObjToCSVFile(obj []string, csvHeader []string, outputFileName string) error {
	file, err := os.Create(outputFileName)
	checkErr(err)
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(csvHeader)
	for _, objItem := range obj {
		var record map[string]interface{}
		json.Unmarshal([]byte(objItem), &record)
		var row []string

		for _, k := range csvHeader {
			switch record[k].(type) {
			case []interface{}:
				data := record[k].([]interface{})
				if len(data) > 0 {
					str := ""
					for i := 0; i < len(data); i++ {
						str1 := fmt.Sprintf("%v", data[i])
						str += str1 + ","
					}
					row = append(row, str[:len(str)-1])
				}
			case string:
				row = append(row, record[k].(string))
			case nil:
				row = append(row, "")
			default:
				dataType, _ := json.Marshal(record[k])
				row = append(row, string(dataType))
			}
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func writeArrayToCSVFile(listMap *gmap.ListMap, node string, csvHeader []string, outputFileName string) error {
	data := listMap.Get(node)
	arr, ok := data.([]interface{})
	if !ok {
		return fmt.Errorf(" > 数据区域不是一个数组！")
	}
	file, err := os.Create(outputFileName)
	checkErr(err)
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(csvHeader)
	for _, item := range arr {
		if record, ok := item.(map[string]interface{}); ok {
			var row []string

			// 注意Map和json默认无序，会导致csv输出时列不对齐
			// 上面保留key顺序，然后从Map中按key顺序取值以保证value有序
			for _, k := range csvHeader {
				switch record[k].(type) {
				case []interface{}:
					data := record[k].([]interface{})
					if len(data) > 0 {
						str := ""
						for i := 0; i < len(data); i++ {
							str1 := fmt.Sprintf("%v", data[i])
							str += str1 + ","
						}
						row = append(row, str[:len(str)-1])
					}
				case string:
					row = append(row, record[k].(string))
				case nil:
					row = append(row, "")
				default:
					dataType, _ := json.Marshal(record[k])
					row = append(row, string(dataType))
				}
			}

			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	return nil
}

func isFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func splitString(s string, myStrings []rune) []string {
	Split := func(r rune) bool {
		for _, v := range myStrings {
			if v == r {
				return true
			}
		}
		return false
	}
	return strings.FieldsFunc(s, Split)
}

func process(jPath string) {
	if !isFileExist(jPath) {
		fmt.Printf(" > 指定文件不存在：%s\n", jPath)
		return
	}
	csvFilePath := strings.Split(jPath, filepath.Ext(jPath))[0] + ".csv"
	bJsonFile, err := os.ReadFile(jPath)
	checkErr(err)
	bJsonFile = bytes.ReplaceAll(ToUtf8(bJsonFile), []byte("\\'"), []byte("'"))

	// 如果设置-d参数，则按szkey指定路径提取json
	if len(szData) > 0 {
		path := []string{}
		for _, v := range strings.Split(szData, ".") {
			if len(v) > 0 {
				path = append(path, strings.TrimSpace(v))
			}
		}

		bJsonFile, err = jin.Get(bJsonFile, path...)
		if err != nil {
			fmt.Printf(" > %s 读取错误： %v\n", jPath, err)
			fmt.Printf(" > 请对照 %s 文件检查-d参数路径！\n", jPath)
			return
		}
	}

	// 兼容处理[{"ID":0,"Name":"Lucy"},{"ID":1,"Name":"Lily"}]类型
	if bJsonFile[0] == '[' {
		bJsonFile = bytes.Join([][]byte{[]byte("{\"兼容\":"), bJsonFile, []byte("}")}, []byte(""))
	}

	listMap, err := readJSONFile(bJsonFile)
	checkErr(err)

	var csvHeader, objValues []string
	maxNode := getMaxNode(listMap)
	// csv表头仅首行写入一次；用于后续保留key顺序

	if szKeys != "" {
		csvHeader = splitString(szKeys, []rune{'/', ','})
		if len(csvHeader) == 0 {
			fmt.Printf(" > %s 读取错误，-k指定字段名不存在！\n", jPath)
			flag.Usage()
			return
		}
	}

	if maxNode == "" {
		// 兼容处理非数组对象数据提取{"items":{"a1":{"title":"one","name":"test"},"b2":{"title":"two","name":"test2"}}}
		objValues, err = jin.GetValues(bJsonFile)
		checkErr(err)
		if len(csvHeader) == 0 {
			csvHeader, err = jin.GetKeys([]byte(objValues[iIndex-1]))
		}
	} else {
		fmt.Printf(" > 数据节点： %s\n", maxNode)
		if len(csvHeader) == 0 {
			// 因为map和json都无序，想保留json键顺序采用第三方库jin
			csvHeader, err = jin.GetKeys(bJsonFile, maxNode, strconv.Itoa(iIndex-1))
		}
	}
	checkErr(err)
	if len(csvHeader) > 0 {
		fmt.Printf(" > %s 字段列表： %v\n", jPath, csvHeader)
	} else {
		return
	}

	if maxNode == "" {
		err = writeObjToCSVFile(objValues, csvHeader, csvFilePath)
	} else {
		err = writeArrayToCSVFile(listMap, maxNode, csvHeader, csvFilePath)
	}
	if err != nil {
		fmt.Printf(" > CSV文件写入错误： %v\n", err)
	} else {
		fmt.Printf(" > CSV文件成功写入： %s\n", csvFilePath)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf(" > 读取错误： %v\n", err)
		flag.Usage()
		// os.Exit(0)
	}
}

func isRunFromCommandLine() (bool, error) {
	parentProcessID := os.Getppid()
	parentProcess, err := os.FindProcess(parentProcessID)
	if err != nil {
		return false, err
	}

	parentProcessName, err := getProcessName(parentProcess.Pid)
	if err != nil {
		return false, err
	}

	if isShell(parentProcessName) {
		return true, nil
	}
	return false, nil
}

func getProcessName(pid int) (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("wmic", "process", "where", "processid="+strconv.Itoa(pid), "get", "name")
	case "darwin", "linux":
		cmd = exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return "cmd", nil
	}

	processName := strings.TrimSpace(string(output))
	return processName, nil
}

func isShell(processName string) bool {
	shells := []string{
		"bash",
		"sh",
		"zsh",
		"powershell",
		"cmd",
	}

	for _, shell := range shells {
		if strings.Contains(strings.ToLower(processName), shell) {
			return true
		}
	}

	return false
}

var (
	bhelp    bool
	bVersion bool
	szData   string
	iIndex   int
	szKeys   string
)

func init() {
	flag.BoolVar(&bhelp, "h", false, "显示帮助")
	flag.BoolVar(&bVersion, "v", false, "显示版本信息")
	flag.StringVar(&szData, "d", "", "设置Json中数据区域所处路径，如'-d root.topics.data'")
	flag.IntVar(&iIndex, "i", 1, "指定从第N个对象中提取字段名")
	flag.StringVar(&szKeys, "k", "", "设置Json数据字段名称(分隔符'/'或','，优先级高于-i参数)，如'-k title/url/type'")
}

func main() {
	isCmdLine, err := isRunFromCommandLine()
	checkErr(err)

	flag.Parse()
	if bhelp {
		flag.Usage()
		return
	}
	if bVersion {
		fmt.Println(" > 版本：v0.7\n > 主页：https://github.com/playGitboy/Json2Csv")
		return
	}

	//process(`C:\Users\Administrator\Desktop\编程\Json2Csv-main\ss.txt`)
	if flag.NArg() > 0 {
		for _, jsonFilePath := range flag.Args() {
			process(jsonFilePath)
		}
		return
	}
	fmt.Println(" > Json2Csv：请指定JSON格式文件路径（支持批量）...")
	fmt.Println(" > Json2Csv [-d data.items] data.json data2.txt ...")
	fmt.Println(" > Json2Csv [-d data.items] [-k title/url] data.json ...")
	fmt.Println(" > Json2Csv [-d data.items] [-i 3] data.json ...")
	flag.Usage()
	if !isCmdLine {
		fmt.Println("\n > 请在终端/命令行使用本程序！")
		fmt.Scanln()
	}
}
