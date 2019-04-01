package api

func init() {
	SetRouterRegister(func(router *RouterGroup) {
		flowRouteGroup := router.Group("/oneflow")
		flowRouteGroup.StdGET("getInst", DoGetInst)
	})
}

func DoGetInst(c *Context) (code int, message string, data interface{}) {
	c.Error("aaa")
	c.Info("ffff")

	return 0, "", "aaaa"
}
