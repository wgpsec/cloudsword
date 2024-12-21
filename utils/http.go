package utils

import (
	"context"
	"fmt"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/hashicorp/go-version"
)

//type UserAgentTransport struct {
//	Base      http.RoundTripper
//	UserAgent string
//}
//
//func HttpTransport() func(*http.Transport) {
//	proxyStr := global.GetBasicOptionValue(global.Proxy)
//	if proxyStr == "" {
//		return nil
//	}
//	parsedURL, err := url.Parse(proxyStr)
//	logger.Println.Error(err.Error())
//
//	return func(transport *http.Transport) {
//		switch parsedURL.Scheme {
//		case "http", "https":
//			transport = &http.Transport{
//				Proxy: http.ProxyURL(parsedURL),
//			}
//		case "socks5":
//			auth := &proxy.Auth{}
//			if parsedURL.User != nil {
//				auth.User = parsedURL.User.Username()
//				auth.Password, _ = parsedURL.User.Password()
//			} else {
//				auth = nil
//			}
//
//			dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, auth, proxy.Direct)
//			if err == nil {
//				transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
//					return dialer.(proxy.ContextDialer).DialContext(ctx, network, addr)
//				}
//			} else {
//				logger.Println.Error(err.Error())
//			}
//		default:
//			logger.Println.Error("代理类型错误，请输入正确的代理类型。")
//		}
//	}
//}
//
//func (t *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
//	req.Header.Set("User-Agent", t.UserAgent)
//	return t.Base.RoundTrip(req)
//}

// 检查是否有新版本

func CheckVersion() {
	currentVer := global.Version
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), global.Team, global.Name)
	if err == nil {
		current, err := version.NewVersion(currentVer)
		if err == nil {
			if release != nil && release.TagName != nil {
				releaseVer := strings.TrimPrefix(*release.TagName, "v")
				latest, err := version.NewVersion(releaseVer)
				if err == nil && latest.GreaterThan(current) {
					logger.Println.Warnf(fmt.Sprintf("发现新版本! 当前版本: %s, 最新版本: %s\n\n", currentVer, releaseVer))
					return
				}
			}
		}
	}
}
