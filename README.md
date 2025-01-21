# TikTokMall

依赖 / 版本：

- Go, 1.23.4
- Protocol Buffers, v29.3

仓库地址：[TikTok Mall](https://github.com/arthur-stat/TikTokMall)




注释：

- auth：认证服务
- cart：购物车服务
- checkout：结算服务
- order：订单服务
- payment：支付服务
- product：商品服务
- user：用户服务



### 创建路由注解文件api.proto
https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/usage-protobuf/
```protobuf
// idl/api.proto; 注解拓展
syntax = "proto2";

package api;

import "google/protobuf/descriptor.proto";

option go_package = "/api";

extend google.protobuf.FieldOptions {
    optional string raw_body = 50101;
    optional string query = 50102;
    optional string header = 50103;
    optional string cookie = 50104;
    optional string body = 50105;
    optional string path = 50106;
    optional string vd = 50107;
    optional string form = 50108;
    optional string js_conv = 50109;
    optional string file_name = 50110;
    optional string none = 50111;

    // 50131~50160 used to extend field option by hz
    optional string form_compatible = 50131;
    optional string js_conv_compatible = 50132;
    optional string file_name_compatible = 50133;
    optional string none_compatible = 50134;
    // 50135 is reserved to vt_compatible
    // optional FieldRules vt_compatible = 50135;

    optional string go_tag = 51001;
}

extend google.protobuf.MethodOptions {
    optional string get = 50201;
    optional string post = 50202;
    optional string put = 50203;
    optional string delete = 50204;
    optional string patch = 50205;
    optional string options = 50206;
    optional string head = 50207;
    optional string any = 50208;
    optional string gen_path = 50301; // The path specified by the user when the client code is generated, with a higher priority than api_version
    optional string api_version = 50302; // Specify the value of the :version variable in path when the client code is generated
    optional string tag = 50303; // rpc tag, can be multiple, separated by commas
    optional string name = 50304; // Name of rpc
    optional string api_level = 50305; // Interface Level
    optional string serializer = 50306; // Serialization method
    optional string param = 50307; // Whether client requests take public parameters
    optional string baseurl = 50308; // Baseurl used in ttnet routing
    optional string handler_path = 50309; // handler_path specifies the path to generate the method

    // 50331~50360 used to extend method option by hz
    optional string handler_path_compatible = 50331; // handler_path specifies the path to generate the method
}

extend google.protobuf.EnumValueOptions {
    optional int32 http_code = 50401;

// 50431~50460 used to extend enum option by hz
}

extend google.protobuf.ServiceOptions {
    optional string base_domain = 50402;

    // 50731~50760 used to extend service option by hz
    optional string base_domain_compatible = 50731;
}

extend google.protobuf.MessageOptions {
    // optional FieldRules msg_vt = 50111;

    optional string reserve = 50830;
    // 550831 is reserved to msg_vt_compatible
    // optional FieldRules msg_vt_compatible = 50831;
}

```
为官方提供的7个proto文件添加http注解以生成对应的router和handler代码, 具体看api/**/*.proto文件的注释

```bash
hz new -module TikTokMall -I api 
hz update -I api -idl api/auth/auto.proto
hz update -I api -idl api/cart/cart.proto
hz update -I api -idl api/checkout/checkout.proto
hz update -I api -idl api/order/order.proto
hz update -I api -idl api/payment/payment.proto
hz update -I api -idl api/product/product.proto
hz update -I api -idl api/user/user.proto
```
