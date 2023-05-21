在shared目录下执行以下命令生成Kitex的服务依赖代码
kitex -module tiktok-demo -I  ./../idl ./../idl/UserServer.proto
kitex -module tiktok-demo -I  ./../idl ./../idl/VideoServer.proto
kitex -module tiktok-demo -I  ./../idl ./../idl/CommentServer.proto
kitex -module tiktok-demo -I  ./../idl ./../idl/RelationServer.proto
kitex -module tiktok-demo -I  ./../idl ./../idl/FavoriteServer.proto

在cmd目录中对应的服务目录下执行以下命令生成对应RPC的服务端代码
kitex -service UserService -module tiktok-demo -use tiktok-demo/shared/kitex_gen -I ./../../idl ./../../idl/UserServer.proto
kitex -service VideoService -module tiktok-demo -use tiktok-demo/shared/kitex_gen -I ./../../idl ./../../idl/VideoServer.proto
kitex -service CommentService -module tiktok-demo -use tiktok-demo/shared/kitex_gen -I ./../../idl ./../../idl/CommentServer.proto
kitex -service RelationService -module tiktok-demo -use tiktok-demo/shared/kitex_gen -I ./../../idl ./../../idl/RelationServer.proto
kitex -service FavoriteService -module tiktok-demo -use tiktok-demo/shared/kitex_gen -I ./../../idl ./../../idl/FavoriteServer.proto

在cmd的api目录生成Hertz的网关代码
hz new -idl ./../../idl/ApiServer.proto -mod tiktok-demo/cmd/api
rm .gitignore go.mod // 删除生成的.gitignore和go.mod文件

在项目根目录拉取并验证依赖
go mod tidy && go mod verify

