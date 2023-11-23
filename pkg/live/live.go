package live

import (
	"crypto/md5"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"ruoyi/settings"
	"time"
)

type LiveInterface interface {
	//1. 获取推流地址鉴权
	GetPushUrl(streamName string) map[string]string
	//2.获取播放地址鉴权
	GetPullUrl(streamName string) map[string]string
	//3.查询直播中的流程
	GetLiveStreamOnlineList(pageNum, pageSize uint64, streamName ...string) (*live.DescribeLiveStreamOnlineListResponse, error)
	//4.获取禁推流列表
	GetLiveForbidStreamList(pageNum, pageSize uint64, streamName ...string) (*live.DescribeLiveForbidStreamListResponse, error)
	//5.禁推直播流
	ForbidLiveStream(streamName string, resumeTime time.Time, reason string) (*live.ForbidLiveStreamResponse, error)
	//6.断开直播流
	DropLiveStream(streamName string) (*live.DropLiveStreamResponse, error)
	//7.查询流状态
	GetLiveStreamState(streamName string) (*live.DescribeLiveStreamStateResponse, error)
	//8.恢复直播流
	ResumeLiveStream(streamName string) (*live.ResumeLiveStreamResponse, error)
	//9.获取延迟播放列表
	GetLiveDelayInfoList() (*live.DescribeLiveDelayInfoListResponse, error)
	//10.设置延迟直播
	AddDelayLiveStream(streamName string, delayTime uint64, expireTime time.Time) (*live.AddDelayLiveStreamResponse, error)
	//11.取消延迟直播
	ResumeDelayLiveStream(streamName string) (*live.ResumeDelayLiveStreamResponse, error)
}

type Live struct {
	AppName    string
	PushDomain string
	PushKey    string
	PullDomain string
	PullKey    string
	Ts         time.Duration //时长
}

func getCredential() *common.Credential {
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取

	credential := common.NewCredential(
		settings.Conf.TencentCloudSecret.SecretId,
		settings.Conf.TencentCloudSecret.SecretKey,
	)
	return credential
}

func getClient() *live.Client {
	credential := getCredential()
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = settings.Conf.LiveConfig.EndPoint
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := live.NewClient(credential, "", cpf)
	return client
}

var lv LiveInterface

func NewLive() LiveInterface {
	ts := settings.Conf.LiveConfig.Ts
	if lv != nil {
		return lv
	}
	lv = &Live{
		AppName:    settings.Conf.LiveConfig.AppName,
		PushDomain: settings.Conf.LiveConfig.PushDomain,
		PushKey:    settings.Conf.LiveConfig.PushKey,
		PullDomain: settings.Conf.LiveConfig.PullDomain,
		PullKey:    settings.Conf.LiveConfig.PullKey,
		Ts:         time.Second * time.Duration(ts),
	}
	return lv
}

func (l *Live) GetPushUrl(streamName string) map[string]string {
	mp := map[string]string{
		"rtmp":   fmt.Sprintf("rtmp://%s/%s/%s", l.PushDomain, l.AppName, streamName),
		"webrtc": fmt.Sprintf("webrtc://%s/%s/%s", l.PushDomain, l.AppName, streamName),
		"srt":    fmt.Sprintf("srt://%s:9000?streamid=#!::h=%s,r=%s/%s", l.PullDomain, l.PushDomain, l.AppName, streamName),
	}
	if l.PushKey != "" {
		txSecret, txTime := sign(l.PushKey, streamName, l.Ts)
		for k, v := range mp {
			if k == "srt" {
				mp[k] = fmt.Sprintf("%s?txSecret=%s&txTime=%s", v, txSecret, txTime)
			} else {
				mp[k] = fmt.Sprintf("%s?txSecret=%s&txTime=%s", v, txSecret, txTime)

			}
		}
	}

	return mp
}

func (l *Live) GetPullUrl(streamName string) map[string]string {
	mp := map[string]string{
		"rtmp":   fmt.Sprintf("rtmp://%s/%s/%s", l.PullDomain, l.AppName, streamName),
		"webrtc": fmt.Sprintf("webrtc://%s/%s/%s", l.PullDomain, l.AppName, streamName),
		"flv":    fmt.Sprintf("http://%s/%s/%s.flv", l.PullDomain, l.AppName, streamName),
		"hls":    fmt.Sprintf("http://%s/%s/%s.m3u8", l.PullDomain, l.AppName, streamName),
	}
	if l.PullKey != "" {
		txSecret, txTime := sign(l.PullKey, streamName, l.Ts)
		for k, v := range mp {
			if k == "srt" {
				mp[k] = fmt.Sprintf("%s?txSecret=%s&txTime=%s", v, txSecret, txTime)
			} else {
				mp[k] = fmt.Sprintf("%s?txSecret=%s&txTime=%s", v, txSecret, txTime)

			}
		}
	}

	return mp
}

func (l *Live) GetLiveStreamOnlineList(pageNum, pageSize uint64, streamName ...string) (*live.DescribeLiveStreamOnlineListResponse, error) {
	client := getClient()

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := live.NewDescribeLiveStreamOnlineListRequest()

	request.DomainName = common.StringPtr(l.PushDomain)
	request.AppName = common.StringPtr(l.AppName)
	request.PageNum = common.Uint64Ptr(pageNum)
	request.PageSize = common.Uint64Ptr(pageSize)
	if len(streamName) > 0 {
		tmp := streamName[0]
		request.StreamName = common.StringPtr(tmp)

	}
	// 返回的resp是一个DescribeRegionsResponse的实例，与请求对象对应
	response, err := client.DescribeLiveStreamOnlineList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}

	return response, nil
}

func (l *Live) GetLiveForbidStreamList(pageNum, pageSize uint64, streamName ...string) (*live.DescribeLiveForbidStreamListResponse, error) {
	client := getClient()
	request := live.NewDescribeLiveForbidStreamListRequest()

	request.PageNum = common.Int64Ptr(int64(pageNum))
	request.PageSize = common.Int64Ptr(int64(pageSize))
	if len(streamName) > 0 {
		tmp := streamName[0]
		request.StreamName = common.StringPtr(tmp)

	}
	response, err := client.DescribeLiveForbidStreamList(request)
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)

		return nil, err
	}
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	return response, nil
}

func (l *Live) ForbidLiveStream(streamName string, resumeTime time.Time, reason string) (*live.ForbidLiveStreamResponse, error) {
	client := getClient()
	// 实例化一个请求对象，每个接口都会对应一个request对象
	request := live.NewForbidLiveStreamRequest()

	request.AppName = common.StringPtr(l.AppName)
	request.DomainName = common.StringPtr(l.PushDomain)
	request.StreamName = common.StringPtr(streamName)
	if !resumeTime.IsZero() {
		request.ResumeTime = common.StringPtr(resumeTime.UTC().String())
	}
	request.Reason = common.StringPtr(reason)

	// 返回的resp是一个ForbidLiveStreamResponse的实例，与请求对象对应
	response, err := client.ForbidLiveStream(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	return response, nil
}

func (l *Live) DropLiveStream(streamName string) (*live.DropLiveStreamResponse, error) {
	client := getClient()
	requset := live.NewDropLiveStreamRequest()
	requset.StreamName = common.StringPtr(streamName)
	requset.DomainName = common.StringPtr(l.PushDomain)
	requset.AppName = common.StringPtr(l.AppName)

	response, err := client.DropLiveStream(requset)
	if err != nil {
		return nil, err
	}
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	return response, nil
}

func (l *Live) GetLiveStreamState(streamName string) (*live.DescribeLiveStreamStateResponse, error) {
	client := getClient()
	request := live.NewDescribeLiveStreamStateRequest()
	request.AppName = common.StringPtr(l.AppName)
	request.StreamName = common.StringPtr(streamName)
	request.DomainName = common.StringPtr(l.PushDomain)

	response, err := client.DescribeLiveStreamState(request)
	if err != nil {
		return nil, err
	}
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	return response, nil
}

func (l *Live) ResumeLiveStream(streamName string) (*live.ResumeLiveStreamResponse, error) {
	client := getClient()
	requset := live.NewResumeLiveStreamRequest()
	requset.AppName = common.StringPtr(l.AppName)
	requset.StreamName = common.StringPtr(streamName)
	requset.DomainName = common.StringPtr(l.PushDomain)

	response, err := client.ResumeLiveStream(requset)
	if err != nil {
		return nil, err
	}
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	return response, nil
}

func (l *Live) GetLiveDelayInfoList() (*live.DescribeLiveDelayInfoListResponse, error) {
	client := getClient()
	request := live.NewDescribeLiveDelayInfoListRequest()

	response, err := client.DescribeLiveDelayInfoList(request)
	if err != nil {
		return nil, err
	}
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	return response, nil
}

func (l *Live) AddDelayLiveStream(streamName string, delayTime uint64, expireTime time.Time) (*live.AddDelayLiveStreamResponse, error) {
	client := getClient()
	request := live.NewAddDelayLiveStreamRequest()
	request.AppName = common.StringPtr(l.AppName)
	request.DomainName = common.StringPtr(l.PushDomain)
	request.StreamName = common.StringPtr(streamName)
	request.DelayTime = common.Uint64Ptr(delayTime)
	request.ExpireTime = common.StringPtr(expireTime.Format("2006-01-02 15:04:05"))
	response, err := client.AddDelayLiveStream(request)
	if err != nil {
		return nil, err
	}
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	return response, nil
}

func (l *Live) ResumeDelayLiveStream(streamName string) (*live.ResumeDelayLiveStreamResponse, error) {
	client := getClient()
	request := live.NewResumeDelayLiveStreamRequest()
	request.AppName = common.StringPtr(l.AppName)
	request.DomainName = common.StringPtr(l.PushDomain)
	request.StreamName = common.StringPtr(streamName)
	response, err := client.ResumeDelayLiveStream(request)
	if err != nil {
		return nil, err
	}
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	return response, nil
}

func sign(key, streamName string, duration time.Duration) (txSecret, txTime string) {
	endTime := time.Now().Add(duration).Unix()
	txTime = fmt.Sprintf("%x", endTime)
	bytes := md5.Sum([]byte(fmt.Sprintf("%s%s%s", key, streamName, txTime)))
	txSecret = fmt.Sprintf("%x", bytes)
	return txSecret, txTime
}
