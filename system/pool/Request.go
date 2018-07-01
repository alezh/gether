package pool

import (
	"compress/zlib"
	"compress/gzip"
	"io"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"golang.org/x/net/context"
)

type Request struct {
	Url           string
	Method        string
	PostData      string
	DialTimeout   time.Duration //创建连接超时 dial tcp: i/o timeout
	ConnTimeout   time.Duration //连接状态超时 WSARecv tcp: i/o timeout
	EnableCookie  bool
	RedirectTimes int //重定向的最大次数，为0时不限，小于0时禁止重定向
	CookieJar     *cookiejar.Jar
}

func (r *Request) DownLoad() (*goquery.Document, bool) {
	httpClient := http.Client{}
	ctx, cancel := context.WithCancel(context.TODO())
	timer := time.AfterFunc(8*time.Second, func() {
		fmt.Println("this url timeout " + r.Url)
		cancel()
	})
	req, err := http.NewRequest("GET", r.Url, nil)
	req.Header.Set("User-Agent", r.getAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err == nil {
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			var gzipReader *gzip.Reader
			gzipReader, err = gzip.NewReader(resp.Body)
			if err == nil {
				resp.Body = gzipReader
			}

		case "deflate":
			// resp.Body = flate.NewReader(resp.Body)

		case "zlib":
			var readCloser io.ReadCloser
			readCloser, err = zlib.NewReader(resp.Body)
			if err == nil {
				resp.Body = readCloser
			}
		}
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	defer resp.Body.Close()
	timer.Stop()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp StatusCode:", resp.StatusCode, r.Url)
		if resp.StatusCode == http.StatusNotFound {
			return nil, true
		}
		return nil, false
	}

	contentType := resp.Header.Get("Content-Type")

	if strings.LastIndex(strings.ToLower(contentType), "gbk") > 0 {
		utfBody := mahonia.NewDecoder("gbk").NewReader(resp.Body)
		doc, err := goquery.NewDocumentFromReader(utfBody)			
		doc.Url = resp.Request.URL
		if err != nil {		
			return nil, false
		}
		return doc, true
	}else{
		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {		
			return nil, false
		}
		return doc, true
	}	
}

func (r *Request) getAgent() string {
	agent := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
	}

	n := rand.New(rand.NewSource(time.Now().UnixNano()))
	lens := len(agent)
	return agent[n.Intn(lens)]
}

// buildClient creates, configures, and returns a *http.Client type.
func (r *Request) buildClient() *http.Client {
	client := &http.Client{
		CheckRedirect: r.checkRedirect,
	}

	if r.EnableCookie {
		client.Jar = r.CookieJar
	}

	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			var (
				c      net.Conn
				err    error
				ipPort = addr
			)
			c, err = net.DialTimeout(network, ipPort, r.DialTimeout)
			if err != nil {
				return nil, err
			}
			if r.ConnTimeout > 0 {
				c.SetDeadline(time.Now().Add(r.ConnTimeout))
			}
			return c, nil
		},
	}

	// if param.proxy != nil {
	// transport.Proxy = http.ProxyURL(param.proxy)
	// }
	if i := strings.Index(r.Url, "https"); i >= 0 {
		transport.TLSClientConfig = &tls.Config{RootCAs: nil, InsecureSkipVerify: true}
		transport.DisableCompression = true
	}
	client.Transport = transport
	return client
}

func (r *Request) checkRedirect(req *http.Request, via []*http.Request) error {
	if r.RedirectTimes == 0 {
		return nil
	}
	if len(via) >= r.RedirectTimes {
		if r.RedirectTimes < 0 {
			return fmt.Errorf("not allow redirects.")
		}
		return fmt.Errorf("stopped after %v redirects.", r.RedirectTimes)
	}
	return nil
}
