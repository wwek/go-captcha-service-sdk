package grpcapi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tidwall/gjson"
	"github.com/wenlng/go-captcha-service-sdk/golang/consts"
	"github.com/wenlng/go-captcha-service-sdk/golang/proto"
	"github.com/wenlng/go-captcha-service-sdk/golang/sdlb"
	"github.com/wenlng/go-captcha-service-sdk/golang/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Client .
type Client interface {
	GetData(ctx context.Context, id string) (*types.CaptData, error)
	CheckData(ctx context.Context, id, captchaKey, value string) (bool, error)
	CheckStatus(ctx context.Context, captchaKey string) (bool, error)
	GetStatusInfo(ctx context.Context, captchaKey string) (*types.CaptStatusInfo, error)
	DelStatusInfo(ctx context.Context, captchaKey string) (bool, error)
}

// client ..
type client struct {
	sdlb   *sdlb.SDLB
	config ClientConfig

	retryCount       int
	retryWaitTime    time.Duration
	retryMaxWaitTime time.Duration
	timeout          time.Duration
}

// ClientConfig ..
type ClientConfig struct {
	BaseAddress string // 127.0.0.1:8080
	APIKey      string // APIKeys

	RetryCount       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
	Timeout          time.Duration
}

var _ Client = (*client)(nil)

// NewGRPCClient ..
func NewGRPCClient(cnf ClientConfig, sdlb *sdlb.SDLB) (Client, error) {
	c := &client{
		sdlb:             sdlb,
		config:           cnf,
		retryCount:       3,
		retryWaitTime:    500 * time.Millisecond,
		retryMaxWaitTime: 3 * time.Second,
		timeout:          10 * time.Second,
	}

	if cnf.RetryCount > 0 {
		c.retryCount = cnf.RetryCount
	}
	if cnf.RetryWaitTime > 0 {
		c.retryWaitTime = cnf.RetryWaitTime
	}
	if cnf.RetryMaxWaitTime > 0 {
		c.retryMaxWaitTime = cnf.RetryMaxWaitTime
	}
	if cnf.Timeout > 0 {
		c.timeout = cnf.Timeout
	}

	return c, nil
}

// SelectAddress ..
func (c *client) SelectAddress() (string, error) {
	hostname, _ := os.Hostname()
	return c.SelectAddressWithKey(hostname)
}

// SelectAddressWithKey ..
func (c *client) SelectAddressWithKey(key string) (string, error) {
	if c.sdlb == nil {
		return c.config.BaseAddress, nil
	}

	inst, err := c.sdlb.Select(key)
	if err != nil {
		if c.config.BaseAddress != "" {
			return c.config.BaseAddress, nil
		}
		return "", fmt.Errorf("failed to select instance: %v", err)
	}
	grpcPort, ok := inst.Metadata["grpc_port"]
	if !ok {
		if c.config.BaseAddress != "" {
			return c.config.BaseAddress, nil
		}
		return "", fmt.Errorf("grpc_port not found in instance metadata")
	}

	addr := fmt.Sprintf("%s:%s", inst.Host, grpcPort)
	return addr, nil
}

// getConnection ..
func (c *client) getConnection() (*grpc.ClientConn, error) {
	hostname, _ := os.Hostname()
	return c.getConnectionWithKey(hostname)
}

// getConnection ..
func (c *client) getConnectionWithKey(key string) (*grpc.ClientConn, error) {
	addr, err := c.SelectAddressWithKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to select instance: %v", err)
	}
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(c.unaryInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %v", addr, err)
	}
	return conn, nil
}

// unaryInterceptor ..
func (c *client) unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md := metadata.New(map[string]string{
		"X-API-Key": c.config.APIKey,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}

// GetData ..
func (c *client) GetData(ctx context.Context, id string) (*types.CaptData, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := proto.NewGoCaptchaServiceClient(conn)
	resp, err := cli.GetData(ctx, &proto.GetDataRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetData: %v", err)
	}

	return &types.CaptData{
		Id:                resp.Id,
		CaptchaKey:        resp.CaptchaKey,
		MasterImageBase64: resp.MasterImageBase64,
		ThumbImageBase64:  resp.ThumbImageBase64,
		MasterImageWidth:  int64(resp.MasterWidth),
		MasterImageHeight: int64(resp.MasterHeight),
		ThumbImageWidth:   int64(resp.ThumbWidth),
		ThumbImageHeight:  int64(resp.ThumbHeight),
		ThumbImageSize:    int64(resp.ThumbSize),
		DisplayX:          int64(resp.DisplayX),
		DisplayY:          int64(resp.DisplayY),
	}, nil
}

// CheckData ..
func (c *client) CheckData(ctx context.Context, id, captchaKey, value string) (bool, error) {
	conn, err := c.getConnection()
	if err != nil {
		return false, err
	}
	defer conn.Close()

	cli := proto.NewGoCaptchaServiceClient(conn)
	resp, err := cli.CheckData(ctx, &proto.CheckDataRequest{
		Id:         id,
		CaptchaKey: captchaKey,
		Value:      value,
	})
	if err != nil {
		return false, fmt.Errorf("failed to call CheckData: %v", err)
	}

	return resp.Data == "ok", nil
}

// CheckStatus ..
func (c *client) CheckStatus(ctx context.Context, captchaKey string) (bool, error) {
	conn, err := c.getConnection()
	if err != nil {
		return false, err
	}
	defer conn.Close()

	cli := proto.NewGoCaptchaServiceClient(conn)
	resp, err := cli.CheckStatus(ctx, &proto.StatusInfoRequest{
		CaptchaKey: captchaKey,
	})
	if err != nil {
		return false, fmt.Errorf("failed to call CheckStatus: %v", err)
	}

	return resp.Data == "ok", nil
}

// GetStatusInfo ..
func (c *client) GetStatusInfo(ctx context.Context, captchaKey string) (*types.CaptStatusInfo, error) {
	conn, err := c.getConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := proto.NewGoCaptchaServiceClient(conn)
	resp, err := cli.GetStatusInfo(ctx, &proto.StatusInfoRequest{
		CaptchaKey: captchaKey,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to call CheckStatus: %v", err)
	}

	ttype := gjson.Get(resp.Data, "type").Int()
	info := gjson.Get(resp.Data, "data").String()
	res := &types.CaptStatusInfo{
		Status: gjson.Get(resp.Data, "status").Int(),
		Type:   ttype,
		Info:   info,
	}

	switch ttype {
	case consts.GoCaptchaTypeClick, consts.GoCaptchaTypeClickShape:
		var data map[int]*types.ClickData
		err = json.Unmarshal([]byte(info), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to json unmarshal: %v", err)
		}

		res.ClickDataMaps = data
		break
	case consts.GoCaptchaTypeSlide, consts.GoCaptchaTypeDrag:
		var data map[int]*types.SlideData
		err = json.Unmarshal([]byte(info), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to json unmarshal: %v", err)
		}

		res.SlideDataMaps = data
		break
	case consts.GoCaptchaTypeRotate:
		var data map[int]*types.RotateData
		err = json.Unmarshal([]byte(info), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to json unmarshal: %v", err)
		}
		res.RotateDataMaps = data
		break
	}
	return res, nil
}

// DelStatusInfo ..
func (c *client) DelStatusInfo(ctx context.Context, captchaKey string) (bool, error) {
	conn, err := c.getConnection()
	if err != nil {
		return false, err
	}
	defer conn.Close()

	cli := proto.NewGoCaptchaServiceClient(conn)
	resp, err := cli.DelStatusInfo(ctx, &proto.StatusInfoRequest{
		CaptchaKey: captchaKey,
	})
	if err != nil {
		return false, fmt.Errorf("failed to call CheckStatus: %v", err)
	}

	return resp.Data == "ok", nil
}
