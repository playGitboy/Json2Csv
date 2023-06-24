# Json2Csv
Golang实现的一款通用型JSON数据提取工具，支持自动识别JSON数据节点并有序提取为CSV文件。  
Convenient JSON data extraction tool.

# 使用简介
```
> Json2Csv：请指定JSON格式文件路径（支持批量）...
> Json2Csv [-k root.data.items] data.json data2.txt ...
Usage of main.exe:
  -h    显示帮助
  -k string
        设置Json中数据所处路径，如'-k root.topics.data'
```  

支持以下常见JSON数据格式：  
### *数据位于根/一级节点下*
```json
[{"ID":0,"Name":"Lucy","Age":17,"Granted":true},{"ID":1,"Name":"Lily","Age":20,"Granted":false}]

{"part":1,"items":[{"title":"one","price":23},{"title":"two","price":92},{"title":"three","price":5623}]}
```
数据提取命令：`Json2Csv test1.json test2.json`
> 拖放json文件到主程序或命令行均可运行，如JSON数据位于根/一级节点下程序可自动检测并提取
### *数据位于任意多级节点下*
```json
{"data":{"items":[{"title":"one","price":23},{"title":"two","price":92},{"title":"three","price":5623}]}}
```
数据位于"data.items"多级节点下，-k参数简单指定数据路径即可，如  
数据提取命令：`Json2Csv -k data.items test.json`  
测试文件：[-k参数 JSON示例](https://danjuanfunds.com/djapi/v3/filter/fund?type=1&order_by=2y&size=200&page=1)  
# 编译依赖
Golang 1.18+
