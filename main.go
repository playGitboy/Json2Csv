package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ecoshub/jin"
	"github.com/gogf/gf/container/gmap"
)

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

func writeCSVFile(listMap *gmap.ListMap, node string, csvHeader []string, outputFileName string) error {
	data := listMap.Get(node)
	arr, ok := data.([]interface{})
	if !ok {
		return fmt.Errorf(" > Data is not an array")
	}

	file, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(csvHeader)
	for _, item := range arr {
		if record, ok := item.(map[string]interface{}); ok {
			var row []string
			// 注意Map和json默认无序导致csv输出各行的列值不对齐
			// 上面保留key顺序，然后从Map中按key顺序取值以保证value有序
			for _, k := range csvHeader {
				switch record[k].(type) {
				case []interface{}:
					data := record[k].([]interface{})
					str := ""
					for i := 0; i < len(data); i++ {
						str1 := fmt.Sprintf("%v", data[i])
						str += str1 + ","
					}
					row = append(row, str[:len(str)-1])
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

func process(jPath string, szkey string) {
	if isFileExist(jPath) == false {
		fmt.Printf(" > 指定文件不存在：%s\n", jPath)
		return
	}
	// 默认csv覆盖输出与原json文件同目录同名
	csvFilePath := strings.Split(jPath, filepath.Ext(jPath))[0] + ".csv"
	bJsonFile, err := ioutil.ReadFile(jPath)
	if err != nil {
		return
	}

	// 如果设置-k参数，则按szkey指定路径提取json数据区域
	if len(szkey) > 0 {
		path := []string{}
		for _, v := range strings.Split(szkey, ".") {
			if len(v) > 0 {
				path = append(path, strings.TrimSpace(v))
			}
		}

		bJsonFile, err = jin.Get(bJsonFile, path...)
		if err != nil {
			fmt.Printf(" > %s 读取错误： %v\n", jPath, err)
			fmt.Printf(" > 请对照 %s 文件检查-k参数路径！\n", jPath)
			return
		}
	}

	// 兼容处理[{"ID":0,"Name":"Lucy"},{"ID":1,"Name":"Lily"}]类型
	if bJsonFile[0] == '[' {
		bJsonFile = bytes.Join([][]byte{[]byte("{\"兼容\":"), bJsonFile, []byte("}")}, []byte(""))
	}

	listMap, err := readJSONFile(bJsonFile)
	if err != nil {
		fmt.Printf(" > %s读取错误： %v\n", jPath, err)
		return
	}

	maxNode := getMaxNode(listMap)
	fmt.Printf(" > 数据节点： %s\n", maxNode)

	// 因为map和json都无序，想保留json键顺序读取尝试多种方法被迫采用第三方库jin
	csvHeader, err := jin.GetKeys(bJsonFile, maxNode, "0")
	if err != nil {
		fmt.Printf(" > %s读取错误： %v\n", jPath, err)
		flag.Usage()
		return
	}
	fmt.Printf(" > %s字段列表： %v\n", jPath, csvHeader)

	err = writeCSVFile(listMap, maxNode, csvHeader, csvFilePath)
	if err != nil {
		fmt.Printf(" > CSV文件写入错误： %v\n\n", err)
	} else {
		fmt.Printf(" > CSV文件成功写入： %s\n\n", csvFilePath)
	}
}

var (
	bhelp bool
	szkey string
)

func init() {
	flag.BoolVar(&bhelp, "h", false, "显示帮助")
	flag.StringVar(&szkey, "k", "", "设置Json中数据所处路径，如'-k root.topics.data'")
}

func main() {
	flag.Parse()
	if bhelp {
		flag.Usage()
		return
	}

	if flag.NArg() > 0 {
		for _, jsonFilePath := range flag.Args() {
			process(jsonFilePath, szkey)
		}
	} else {
		fmt.Println(" > Json2Csv：请指定JSON格式文件路径（支持批量）...")
		fmt.Println(" > Json2Csv [-k root.data.items] data.json data2.txt ...")
		flag.Usage()
	}
}
