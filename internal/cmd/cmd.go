package cmd

import (
	"context"

	"knowledge-system-api/internal/controller/knowledge"
	"knowledge-system-api/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "启动知识库检索系统API服务",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// 全局中间件
			s.Use(ghttp.MiddlewareHandlerResponse)
			s.Use(ghttp.MiddlewareCORS) // 允许跨域请求

			// 注册API路由
			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				// 知识库API
				group.Group("/knowledge", func(group *ghttp.RouterGroup) {
					// 绑定控制器
					group.Bind(
						knowledge.NewV1(),
					)
				})

				// 这里可以添加其他模块的API路由
			})

			// 设置Swagger UI
			s.SetSwaggerUITemplate(ScalarUITemplate)

			// 将任务恢复放在服务启动前
			g.Log().Info(ctx, "服务即将启动，初始化任务恢复...")
			// 异步执行任务恢复，避免阻塞主线程
			go service.RecoverUnfinishedTasks(ctx)

			// 启动服务
			s.Run()
			return nil
		},
	}
)

const (
	SwaggerUITemplate = `
<!DOCTYPE HTML>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI"/>
    <title>SwaggerUI</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui.min.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui-bundle.js" crossorigin></script>
<script>
    window.onload = () => {
        window.ui = SwaggerUIBundle({
            url:    '{SwaggerUIDocUrl}',
            dom_id: '#swagger-ui',
        });
    };
</script>
</body>
</html>
`

	OpenapiUITemplate = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>openAPI UI</title>
  </head>
  <body>
    <div id="openapi-ui-container" spec-url="{SwaggerUIDocUrl}" theme="light"></div>
    <script src="https://cdn.jsdelivr.net/npm/openapi-ui-dist@latest/lib/openapi-ui.umd.js"></script>
  </body>
</html>
`

	ScalarUITemplate = `
<!doctype html>
<html>
  <head>
    <title>知识库检索系统 API 文档</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="{SwaggerUIDocUrl}"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>
`
)
