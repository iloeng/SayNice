SayNice 匿名情感倾诉社区。

> 线上地址：<https://saynice.knowlgraph.com/>

我们需要倾诉，需要发泄，需要有人关注，需要有人关怀，需要有人认同，需要有人支持，需要有人理解，需要有人陪伴。    
我们需要真实，我们也需要距离； 我们需要安全，我们也需要放肆。    
我们需要我们存在于这个世间，我们也需要我们不被人打扰。    
我们需要彼此，我们需要say nice。    

守约：
1. 不得发布政治、宗教、军事、性和毒品交易、赌博、暴力和血腥、未经证实的新闻、地区和种族以及性别歧视的言论，不得发布令人反感的广告，不文明用语请用 * 号替代；
2. 内容中不得出现任何联系方式，包含但不限于人名、地名、时间、地址、邮箱、手机号、微信号、Fackbook号、Youtube号等。
3. 发布新闻类主题时，请添加消息来源。

具体可见[SayNice守约](https://github.com/ThreeTenth/SayNice/blob/master/docs/articles.md)

寄语：
我们希望有一天，我们可以忘记这里，我们在这里，衷心的祝愿我们终会迎来那一天，我们不管早晚，我们满怀希望；
我们希望那一天，你不要悲伤，也不要留恋，你会如同清晨的阳光，拥抱一个全新的世界，我们不管早晚，我们满怀希望；

社区内容发布系统，发布分三步

1. 首先由社区智能过滤系统过滤
2. 然后由随机抽选的匿名用户组成的[随机匿名空间审查](https://github.com/ThreeTenth/SayNice/blob/master/docs/rnspace.md)
3. 最后发布到社区由所有用户监督；

用户只能使用指定的 emoji 语言回复；社区不校验用户，且社区没有用户系统，不保留任何个人信息，包括 IP。

## 使用

### 使用源码编译

在 `source` 文件夹下，分别有两个目录：`server` 和 `web`，使用 `golang` 编写，`web` 服务使用了 [packr](https://github.com/gobuffalo/packr) 作为静态资源打包工具，你可以将源码编译为目标平台下的可执行文件来搭建 SayNice 社区。

### GitHub Actions

本项目使用了 GitHub Actions 实现项目部署自动化。

需要在 GitHub 项目设置页 Settings 的 Secrets 选项页中设置 4 个变量：

变量名 | 含义
------ | -------
DEPLOY_KEY | 服务器私钥
SSH_HOST | ssh Host
SSH_PORT | ssh port
SSH_USERNAME | ssh username

------------------

项目代码及所有资源，采用 MIT 开源协议。
