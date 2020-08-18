# Golang 爬虫
使用goquery 和 chromedp 爬取知乎大v 想法以及下面的热评,共爬取50个页面,每个页面十条想法,共计500条,存入本地mongodb数据库中

####1
以爬取5个页面为一组goroutine,共开启10组,每组耗时大约1分钟  
```go
queue := make(chan bool)
	for j := 0; j <= 9; j++ {
		for i := 1; i <= 5; i++ {
			url := UrlZhihu + strconv.Itoa(j*5+i)
			go getHostComment(url,queue)
		}
		for i := 0; i < 5; i++ {
			<-queue
		}
	}
```

####2
在chromedp中添加了multiclick()以满足自身的功能需求
```go
func MultiClick(sel interface{}, opts ...QueryOption) QueryAction {
	return QueryAfter(sel, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector %q did not return any nodes", sel)
		}
		//fmt.Printf("nodes:%d\n",len(nodes))
		for i := 0; i < len(nodes); i+=2 {
			MouseClickNode(nodes[i]).Do(ctx)
			Query("div[class=\"Modal Modal--default signFlowModal\"]",append(opts,NodeVisible)...)
			Click("button[class=\"Button Modal-closeButton Button--plain\"]")
		}
		return nil
	}, append(opts, NodeVisible)...)
}
```

####3
使用后记得杀死所有的chrome进程
```shell script
ps -ef|grep Chrome |grep -v grep|cut -c 9-15|xargs kill -9
```







