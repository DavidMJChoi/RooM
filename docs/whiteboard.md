{{template "footer.html" .}}
{{template "nav.html" .}}

constructing workspace
* chrome (1. web search and ai chat; 2. page rendering)
* vs code (explorer + terminal + editor + ai agent)
* whiteboard
* dev journal and documentations (UML diagrams et al...)

```mermaid
stateDiagram-v2
    [*] --> Server_Starting
    
    Server_Starting --> Server_Running: 启动成功
    Server_Starting --> Server_Error: 启动失败
    
    Server_Running --> Request_Received: 收到HTTP请求
    Request_Received --> Route_Matching: 路由匹配
    
    Route_Matching --> Home_Page: 路径 "/"
    Route_Matching --> Posts_Page: 路径 "/posts/"
    Route_Matching --> Post_Detail: 路径 "/posts/{slug}"
    Route_Matching --> Static_File: 路径 "/static/*"
    Route_Matching --> Not_Found: 无效路径
    
    Home_Page --> Template_Rendering: 加载模板
    Posts_Page --> Load_Posts: 加载文章列表
    Post_Detail --> Load_Post: 加载单篇文章
    Static_File --> File_Server: 静态文件服务
    
    Load_Posts --> Template_Rendering: 文章列表渲染
    Load_Post --> Markdown_Parse: 解析Markdown
    Markdown_Parse --> Template_Rendering: 文章详情渲染
    File_Server --> Response_Sent: 返回文件内容
    Not_Found --> Error_Response: 返回404
    
    Template_Rendering --> Response_Sent: HTML响应
    Error_Response --> Response_Sent: 错误响应
    
    Response_Sent --> Request_Received: 等待下一个请求
    Server_Running --> Server_Shutdown: 服务器关闭
    
    Server_Shutdown --> [*]
    Server_Error --> [*]
```