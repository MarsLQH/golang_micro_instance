#micro 实战示例



##目录结构
services --服务端 提供服务  
handler --具体服务实现  
proto-proto文件  

models --对应表结构及request/response结构体  
core --核心组件/中间件【DB,mysql...】  
dao --操作数据库层 存放操作数据库的逻辑 对数据库的操作仅可存放于此目录  
--base 对数据库操作公共部分的封装 如connect Db  
config --配置文件

client --服务请求部分

##请求规范

1，任何调用方调用本服务时必须指定timeout时间，具体时长可依据业务来定  
eg: client.NewClient(&client.Options{Timeout: 5 * time.Second}) //设置5秒超时  

###注意
1，定义proto字段类型时，int64的字段类型会出错，会报转化为String失败 ,所以uint32应该够用了

#####advise:  
1，任何调用方都加入链路追踪埋点

    //方式一： 记录请求的根数据  
    span := tracer.StartSpan("some_operation")  
    //working  
    span.Finish()  

	//方式二：如果需要记录请求的上一步和下一步操作，则需要传入上下文
	
	childSpan := tracer.StartSpan("some_operation2", zipkin.Parent(span.Context()))  
	//operations
	childSpan.Finish()  


	//====为了快速排查问题，可以为某个记录添加一些自定义标签，如记录是否发生错误/请求的返回值等
	//childSpan.Tag("http.status_code",statusCode)
    childSpan.Tag("error",err.Error())





     
    
