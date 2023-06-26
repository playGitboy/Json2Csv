# Json2Csv
Golang实现的一款通用型JSON数据提取工具，支持自动识别JSON数据节点并有序提取为CSV文件。  
Convenient JSON data extraction tool.

# 使用简介
```
> Json2Csv：请指定JSON格式文件路径（支持批量）...
> Json2Csv [-d data.items] data.json data2.txt ...
> Json2Csv [-d data.items] [-k title/url] data.json ...
> Json2Csv [-d data.items] [-i 3] data.json ...
Usage of Json2Csv.win.amd64.exe:
 -d string 设置Json中数据区域所处路径，如'-d root.topics.data'
 -i int    指定从第N个对象中提取字段名 (default 1)
 -k string 设置Json数据字段名称(优先级高于-i参数)，如'-k title/media_url/desc'
 -h        显示帮助
 -v        显示版本信息
```  

支持以下常见JSON数据格式：  
### 1.*自动提取位于根/一级节点下的数组数据*
```json
[{"ID":0,"Name":"Lucy","Age":17,"Granted":true},{"ID":1,"Name":"Lily","Age":20,"Granted":false}]

{"part":1,"items":[{"title":"one","price":23},{"title":"two","price":92},{"title":"three","price":5623}]}
```
数据提取命令：`Json2Csv test1.json test2.json`
> 拖放json文件到主程序或命令行均可运行，如JSON数组数据位于根/一级节点下程序可自动检测并提取成同名csv文件
### 2.*手动提取位于任意多级节点下的数组数据*
```json
{"data":{"items":[{"title":"one","price":23},{"title":"two","price":92},{"title":"three","price":5623}]}}
```
数据位于"data.items"多级节点下，-k参数简单指定数据路径即可，如  
数据提取命令：`Json2Csv -d data.items test.json`  
测试文件：[-k参数 JSON示例](https://danjuanfunds.com/djapi/v3/filter/fund?type=1&order_by=2y&size=200&page=1)   
### 3.*手动提取多级节点下的对象数据*  
```json
{"part":1,"data":{"items":{"1":{"title":"one","name":"test1"},"2":{"title":"two","name":"test2"},{"3":{"title":"three","name":"test3"}}}}
```  
注意数据区域非数组结构而是一个对象，大多数在线网站和json工具都无法解析数据提取，为方便使用这里一并兼容处理了  
数据提取命令：`Json2Csv -d data.items test.json`  
### 4.*支持手动指定字段名/自动从第N个数据块读取字段名*
```json
{"status":"ok","data":{"list":[{"uuid":"0DC0002B","title":"前言","is_chapter":1},{"uuid":"8743CB8D","title":"前言讲义","type":"document","length":90,"weight":1,"media_uri":"a6283c64\/document\/BrDM.doc","course_title":"2016年司考","is_chapter":0}}}
```
有些json不太标准，类似上面这种直接用`-k data.list`解析出的数据默认会从第1个数据块读取字段即uuid/title/is_chapter，导致生成的csv文件缺失大量数据，此时可以设置"-i 2"参数从第2个数据块读取解析全部字段  
数据提取命令：`Json2Csv -d data.items -i 2 test.json`  

也可以使用"-k title/media_uri/course_title"参数手动设置待读取的字段名称，这样导出的数据将只包含title/media_uri/course_title几列  
数据提取命令：`Json2Csv -d data.items -k title/media_uri/course_title test.json`  

# 编译依赖
Golang 1.18+
