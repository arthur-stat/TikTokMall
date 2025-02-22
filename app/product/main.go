package main
import (

	"net"
	"os"
	"time"

	"TikTokMall/app/product/conf"
	"TikTokMall/app/product/kitex_gen/product/productcatalogservice"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	_ = godotenv.Load()

	opts := kitexInit()

    productCatalogServiceImpl := NewProductCatalogServiceImpl()
	svr := productcatalogservice.NewServer(productCatalogServiceImpl, opts...)

	err := svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	// klog
	klog.SetLevel(conf.LogLevel())
	if conf.GetEnv() == "test" || conf.GetEnv() == "dev" {
		klog.SetOutput(os.Stdout)
	} else {
		// "online"
		asyncWriter := &zapcore.BufferedWriteSyncer{
			WS: zapcore.AddSync(&lumberjack.Logger{
				Filename:   conf.GetConf().Kitex.LogFileName,
				MaxSize:    conf.GetConf().Kitex.LogMaxSize,
				MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
				MaxAge:     conf.GetConf().Kitex.LogMaxAge,
			}),
			FlushInterval: time.Minute,
		}
		klog.SetOutput(asyncWriter)
		server.RegisterShutdownHook(func() {
			asyncWriter.Sync()
		})
	}
	return
}
