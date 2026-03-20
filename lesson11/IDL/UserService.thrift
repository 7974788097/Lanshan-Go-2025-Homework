namespace go UserService

include "Common.thrift"

service UserService{
   RegisterResp Register(1:RegisterReq req)
   LoginResp Login(1:LoginReq req)
}

struct RegisterReq{
    1:string user_name
    2:string password
}
struct RegisterResp{
    1:Common.Resp resp
}

struct LoginReq{
    1:string user_name
    2:string password
}
struct LoginResp{
    1:Common.Resp resp,
    2:string token
}
