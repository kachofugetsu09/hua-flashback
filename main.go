package main

import (
	"bytes"
	"fmt"
	"github.com/go-mysql-org/go-mysql/replication"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func printEventHeader(header *replication.EventHeader) {
	if header == nil {
		return
	}

	fmt.Printf(
		"时间戳=%d 事件=%s ServerID=%d NextLogPos=%d Size=%d\n",
		header.Timestamp,
		header.EventType,
		header.ServerID,
		header.LogPos,
		header.EventSize,
	)
}

func main() {
	file, err := os.Open("source/binlog.000002")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) < len(replication.BinLogFileHeader) || !bytes.Equal(data[:len(replication.BinLogFileHeader)], replication.BinLogFileHeader) {
		log.Fatalf("文件不是有效的 binlog，开头应为 %#v", replication.BinLogFileHeader)
	}

	reader := bytes.NewReader(data[len(replication.BinLogFileHeader):])

	// 初始化解析器
	parser := replication.NewBinlogParser()
	parser.SetFlavor("mysql")
	parser.SetVerifyChecksum(true)

	if err := parser.ParseReader(reader, func(event *replication.BinlogEvent) error {
		switch evt := event.Event.(type) {
		case *replication.FormatDescriptionEvent:
			header := event.Header
			printEventHeader(header)

			version := evt.Version                // binlog 版本号
			serverVersion := evt.ServerVersion    // 产生日志的 MySQL 版本
			headerLength := evt.EventHeaderLength // 事件头长度
			checkSum := evt.ChecksumAlgorithm     // 校验算法
			fmt.Printf(
				"  [Format] Version=%d ServerVersion=%s HeaderLength=%d ChecksumAlgo=%d\n",
				version,
				strings.TrimRight(string(serverVersion), "\x00"),
				headerLength,
				checkSum,
			)

		case *replication.QueryEvent:
			sqlText := strings.TrimSpace(string(evt.Query)) // 具体 SQL 文本
			if strings.EqualFold(sqlText, "BEGIN") {
				return nil
			}

			header := event.Header
			printEventHeader(header)

			schema := evt.Schema       // 目标数据库
			errorCode := evt.ErrorCode // 执行过程中返回的错误码

			fmt.Printf("  [Query] DB=%s Error=%d SQL=%s\n",
				string(schema),
				errorCode,
				sqlText,
			)

		default:
		}
		return nil
	}); err != nil {
		log.Fatalf("解析 binlog 失败: %v", err)
	}

}
