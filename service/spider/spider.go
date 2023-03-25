package spider

import (
	"bytes"
	"fmt"
	musicspider "github.com/BeanWei/MusicSpider"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

func checkDicKey(dict map[string]int, thisKey ...string) bool {
	if dict == nil {
		return false
	}
	for _, k := range thisKey {
		_, ok := dict[k]
		if !ok {
			return ok
		}
	}
	return true
}

func Search(site, keyword string, option map[string]int) map[string]string {
	var (
		_type  = 1
		limit  = 30
		page   = 1
		offset = 0
	)
	switch site {
	case "wanyi":
		reqMethod := "POST"
		url := "http://music.163.com/weapi/search/get"
		total, n := "true", "1000"
		if checkDicKey(option, "type") {
			_type = option["type"]
		}
		if checkDicKey(option, "limit") {
			limit = option["limit"]
		}
		if checkDicKey(option, "page", "limit") {
			offset = (option["page"] - 1) * option["limit"]
		}
		data := fmt.Sprintf(`{"s": "%s", "offset": "%d", "limit": "%d", "type": "%d", "total": "%s", "n": "%s"}`,
			keyword, offset, limit, _type, total, n)
		return reqHandler(site, reqMethod, url, data, true)
	case "qq":
		reqMethod := "POST"
		url := "https://u.y.qq.com/cgi-bin/musicu.fcg"
		if checkDicKey(option, "limit") {
			limit = option["limit"]
		}
		if checkDicKey(option, "page", "limit") {
			offset = (option["page"] - 1) * option["limit"]
		}
		data := fmt.Sprintf(`{"comm":{"ct":"19","cv":"1859","uin":"0"},"req":{"method":"DoSearchForQQMusicDesktop","module":"music.search.SearchCgiService","param":{"grp":1,"num_per_page":%d,"page_num":%d,"query":"%s","search_type":0}}}`,
			limit, page, keyword)
		return reqHandler(site, reqMethod, url, data)
	case "xiami":
		reqMethod := "GET"
		//url := "http://h5api.m.xiami.com/h5/mtop.alimusic.search.searchservice.searchsongs/1.0/"
		url := "http://api.xiami.com/web"
		if checkDicKey(option, "page") {
			page = option["page"]
		}
		if checkDicKey(option, "limit") {
			limit = option["limit"]
		}
		v, r, appKey := "2.0", "search/songs", "1"
		data := fmt.Sprintf(`{"key": "%s", "page": "%d", "limit": "%d", "v": "%s", "r": "%s", "app_key": "%s"}`,
			keyword, page, limit, v, r, appKey)
		return reqHandler(site, reqMethod, url, data)
	case "kugou":
		reqMethod := "GET"
		url := "http://mobilecdn.kugou.com/api/v3/search/song"
		if checkDicKey(option, "limit") {
			limit = option["limit"]
		}
		if checkDicKey(option, "page") {
			page = option["page"]
		}
		api_ver, area_code, correct, plat, tag, sver, showtype, version := 1, 1, 1, 2, 1, 5, 10, 8990
		data := fmt.Sprintf(`{"keyword": "%s", "pagesize": "%d", "page": "%d", "api_ver": "%d", "area_code": "%d", "correct": "%d", "plat": "%d", "tag": "%d", "sver": "%d", "showtype": "%d", "version": "%d"}`,
			keyword, limit, page, api_ver, area_code, correct, plat, tag, sver, showtype, version)
		return reqHandler(site, reqMethod, url, data)
	//case "baidu":
	//	reqMethod := "GET"
	//	url :=
	//		fmt.Sprintf("http://musicmini.qianqian.com/v1/restserver/ting?method=baidu.ting.ugcdiy.getChannels&time=%d&timestamp=${data.timestamp}&param=${data.param}&sign=${data.sign}",
	//		,time.Now().UnixMilli(),)
	//
	//	if checkDicKey(option, "page") {
	//		page = option["page"]
	//	}
	//	if checkDicKey(option, "limit") {
	//		limit = option["limit"]
	//	}
	//	data := fmt.Sprintf(`{"query": "%s", "page_no": "%d", "page_size": "%d", "from": "%s", "method": "%s", "isNew": "%d", "platform": "%s", "version": "%s"}`,
	//		keyword, page, limit, from, method, isNew, platform, version)
	//	return reqHandler(site, reqMethod, url, data)
	default:
		return map[string]string{"status": "404", "result": "暂不支持此站点"}
	}
}

// 请求处理函数，所有的请求均通过此入口
// TODO: 更优雅的处理可变参数 args
// args 为一个bool值列表,
// 列表0为为是否需要加密参数处理的Bool值, 默认不处理加密
func reqHandler(site, reqmethod, url, data string, args ...bool) map[string]string {
	client := http.DefaultClient
	var dataByte io.Reader
	if data == "" {
		dataByte = nil
	} else {
		if len(args) > 0 && args[0] == true {
			if site == "wanyi" {
				dataByte = musicspider.NeteaseAESCBC(data)
			}
		} else {
			dataByte = bytes.NewBuffer([]byte(data))
		}
	}
	fmt.Println("the URL is: ", url)
	fmt.Println("the Params is: ", data)
	req, err := musicspider.RequestHandler(reqmethod, url, dataByte)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(site)
	switch site {
	case "wanyi":
		req.Header.Set("Host", "music.163.com")
		req.Header.Set("Referer", "http://music.163.com")
		//req.Header.Set("Cookie", "appver=1.5.9; os=osx; __remember_me=true; osver=%E7%89%88%E6%9C%AC%2010.13.5%EF%BC%88%E7%89%88%E5%8F%B7%2017F77%EF%BC%89")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5) AppleWebKit/605.1.15 (KHTML, like Gecko)")
		req.Header.Set("X-Real-IP", long2ip(randRange(1884815360, 1884890111)))
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,gl;q=0.6,zh-TW;q=0.4")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case "qq":
		req.Header.Set("Referer", "http://y.qq.com")
		//req.Header.Set("Cookie", "pgv_pvi=22038528; pgv_si=s3156287488; pgv_pvid=5535248600; yplayer_open=1; ts_last=y.qq.com/portal/player.html; ts_uid=4847550686; yq_index=0; qqmusic_fromtag=66; player_exist=1")
		req.Header.Set("User-Agent", "QQ%E9%9F%B3%E4%B9%90/54409 CFNetwork/901.1 Darwin/17.6.0 (x86_64)")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,gl;q=0.6,zh-TW;q=0.4")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case "xiami":
		req.Header.Set("Referer", "http://h.xiami.com/")
		req.Header.Set("Cookie", "_m_h5_tk=15d3402511a022796d88b249f83fb968_1511163656929; _m_h5_tk_enc=b6b3e64d81dae577fc314b5c5692df3c")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5) AppleWebKit/537.36 (KHTML, like Gecko) XIAMI-MUSIC/3.1.1 Chrome/56.0.2924.87 Electron/1.6.11 Safari/537.36")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Language", "zh-CN")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case "kugou":
		req.Header.Set("User-Agent", "IPhone-8990-searchSong")
		req.Header.Set("UNI-UserAgent", "iOS11.4-Phone8990-1009-0-WiFi")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case "baidu":
		req.Header.Set("Cookie", "BAIDUID='.$this->getRandomHex(32).':FG=1")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) baidu-music/1.2.1 Chrome/66.0.3359.181 Electron/3.0.5 Safari/537.36")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "zh-CN")
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	default:
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	fmt.Println("the ReqURL is: ", req.URL.String())
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	status := strconv.Itoa(resp.StatusCode)
	fmt.Println(req.Header)
	fmt.Println("status:" + status)
	if status == "200" {
		body, _ := ioutil.ReadAll(resp.Body)
		if site == "wanyi" {
			body, _ = musicspider.GzipDecode(body)
		}
		feedback := map[string]string{"status": status, "result": string(body)}
		return feedback
	}
	feedback := map[string]string{"status": status, "result": ""}
	return feedback
}

func randRange(min, max int64) int64 {
	diff := max - min
	move := rand.Int63n(diff)
	randNum := min + move
	return randNum
}

// Go版 IP2LONG, LONG2IP => https://blog.csdn.net/zengming00/article/details/80354248
func long2ip(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(ip>>24)&0xFF,
		(ip>>16)&0xFF,
		(ip>>8)&0xFF,
		ip&0xFF)
}
