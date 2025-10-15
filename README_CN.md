# go-plugify 例子

比较于 rpc 有什么不一样。

- rpc 没法直接获取server端对象，go-plugify只要挂载了，可以通过interface或unsafe方式调用方法或属性
- rpc 要实现类似效果，需要server端做很多映射代码，go-plugify则不需要server有额外过多和复杂的侵入逻辑
- 但也有缺点，就是plugin没法卸载，go-plugify每次执行会导致server端进程rss内存增加，只能通过重启优化
- 编译.so以及上传稍微比较费时，不能用于实时服务

比较于 yaegi 有什么不一样。

- 可以直接调用server端对象，但需要手动注册，但是写代码的时候是没有补全和提示的
- 调用方法有限制，比如反射等是用不了的，而 go-plugify 支持原生go语法，有补全和提示，也能使用反射

不建议用在生产环境