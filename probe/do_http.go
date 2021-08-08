package probe

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"simple-http-probe/config"
	"strings"
	"time"
)

func DoHttpProbe(url string) string {
	twSec := config.GlobalTwsec

	// 结果
	dnsStr := ""
	targetAddr := ""

	// 提前定义好这些计算时间的对象
	var dnsstartTime, dnsoripTime, t3, t4 time.Time
	// 全局或出错使用
	start := time.Now()

	// 初始化request对象
	req, _ := http.NewRequest("GET", url, nil)
	// 开始trace探测的地址
	trace := &httptrace.ClientTrace{
		DNSStart: func(httptrace.DNSStartInfo) {
			// 开始dns解析的时刻
			dnsstartTime = time.Now()
			fmt.Printf("%v+++++++++++++++++\n", dnsstartTime)
		},
		DNSDone: func(dnsDoneInfo httptrace.DNSDoneInfo) {
			// 完成dns解析时间
			dnsoripTime = time.Now()
			fmt.Printf("%v----------------\n", dnsoripTime)
			// 解析得到的ip地址
			ips := make([]string, 0)
			for _, d := range dnsDoneInfo.Addrs {
				ips = append(ips, d.IP.String())
			}
			dnsStr = strings.Join(ips, ",")
		},
		ConnectStart: func(network, addr string) {
			// 没有域名解析，直接ip地址的
			if dnsoripTime.IsZero() {
				dnsoripTime = time.Now()

			}
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				log.Printf("无法建立和探测目标的连接地址：%v  错误：%v", addr, err)
				return
			}
			targetAddr = addr
		},
		GotConn: func(httptrace.GotConnInfo) {
			t3 = time.Now()
		},
		GotFirstResponseByte: func() {
			t4 = time.Now()
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	// 超时的客户端
	client := http.Client{
		Timeout: time.Duration(twSec) * time.Second,
	}
	resp, err := client.Do(req)
	// 出错的返回
	if err != nil {
		msg := fmt.Sprintf("[http探测出错]\n"+
			"[http探测的目标：%s]\n"+
			"[错误详细：%v]\n"+
			"[总耗时：%s]\n",
			url,
			err,
			msDurationStr(time.Now().Sub(start)),
		)
		log.Printf(msg)
		return msg
	}
	// 延迟关闭链接
	defer resp.Body.Close()
	end := time.Now()

	// 没有dns解析
	if dnsstartTime.IsZero() {
		dnsstartTime = dnsoripTime
	}

	// dns连接的耗时
	dnsLookup := msDurationStr(dnsoripTime.Sub(dnsstartTime))
	// tcp连接的耗时
	tcpConnection := msDurationStr(t3.Sub(dnsoripTime))

	// 服务端处理的耗时
	serverProcessing := msDurationStr(t4.Sub(t3))
	// 总耗时
	totoal := msDurationStr(end.Sub(dnsstartTime))

	probeResStr := fmt.Sprintf(
		"[http探测的目标：%s]\n"+
			"[dns解析的结果：%s]\n"+
			"[连接的ip和端口：%s]\n"+
			"[状态码：%d]\n"+
			"[dns解析耗时：%s]\n"+
			"[tcp连接耗时：%s]\n"+
			"[服务端处理耗时：%s]\n"+
			"[总耗时：%s]\n",
		url,
		dnsStr,
		targetAddr,
		resp.StatusCode,
		dnsLookup,
		tcpConnection,
		serverProcessing,
		totoal,
	)
	return probeResStr
}

// 将秒转成毫秒
func msDurationStr(d time.Duration) string {
	return fmt.Sprintf("%dms", int(d/time.Millisecond))
}
