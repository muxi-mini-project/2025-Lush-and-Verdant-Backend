package client

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/config"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"net/http"
)

type ChatGptClient struct {
	cfg *config.ChatGptConfig
}

func NewChatGptClient(cfg *config.ChatGptConfig) *ChatGptClient {
	return &ChatGptClient{cfg: cfg}
}

func (cg *ChatGptClient) AskForGoal(c *gin.Context, question request.Question) map[string][]map[string]string {
	// 设置初始配置
	var options = []option.RequestOption{
		option.WithBaseURL("https://api.chatanywhere.tech"), //更换了BaseURL,openai的太贵了
		option.WithAPIKey(cg.cfg.Sdk),
	}

	// 设置自定义的 baseurl 和 API 密钥
	client := openai.NewClient(options...)

	now := time.Now()

	systemMessage := fmt.Sprintf(`
    角色定位:
		你是一个根据用户输入的主题，帮助用户制定可执行的计划及对应的实践时间的AI。
       	用户将发送给你一个学习主题,你需要将其拆解为多个一个月内以天为周期的计划表，并为这几个周期段提供可实施的具体事件。
		当前时间为%s,你需要以这个时间点为基准规划任务。
    
    格式说明:
       	1. 输入topic表示学习的主题,description表示对用户情况的描述,cycle表示用户希望完成这些任务的时间。
       	2. 用户输入的cycle表示用户预计完成这些事情的总时间，例如:输入2个月，你自行记录制定时的日期，并且从下一天开始记录，你需要制定两个月内合理的能够完成的事情。
       	3. 你需要给用户返回具体某一天包含的多个任务，其中每个任务包含它的id，title以及details。
		4. 你所返回的时间是这个任务应当结束的期限，例如:某个任务的预计完成时间是2025-3-10到2025-3-15，那么你只需要返回时间为2025-3-15，并附上这些天要完成的任务。
       	5. 返回格式为json，不允许包含其他任何字符。

	要求说明:
		请严格按照以下JSON格式要求生成学习计划:
{
	"日期字符串":[
		{
			"id":"任务ID(字符串格式)",
			"title":"任务标题",
			"details":"任务详情"
		}
	]
}

	规则说明:
		1.日期格式必须为YYYY-MM-DD,从明天开始计算
		2.每个任务必须包含id、title、details三个字段
		3.id使用连续数字字符串(如"1""2")
		4.不要添加任何JSON以外的文本
    
    输入输出示例1:
	假设提问制定计划时的日期为2025-3-10
input:{
  "topic": "前端开发",
  "description": "我是一个编程领域小白，没有任何编程基础",
  "cycle": "1个月"
}
    
output: {
  "2025-03-10"": [
    { id: "1", title: "基础知识", details: "了解前端主流开发语言" },
    { id: "2", title: "搭建环境", details: "选择好自己要学习的语言，并配置好环境" },
  ],
  "2025-03-15": [
    { id: "1", title: "学习语法", details: "参照教程，学习语言的相关语法" },
  ],
  "2025-03-30": [
    { id: "1", title: "尝试实操", details: "学习了相关语法后，尝试自己动手写出一个前端页面" },
    { id: "2", title: "学习框架", details: "在GitHub上有许多框架，可以使用来优化代码并实现许多功能，优化性能" },
  ],
  "2025-04-10": [
    { id: "1", title: "前后对接", details: "想让自己的产品能够运行，还需要与后端联合，尝试前后端对接" },
    { id: "2", title: "多环境尝试", details: "前端页面不仅可以在服务器上部署，还可以制作APP软件" },
  ],
}

	输入输出示例2:
	假设提问制定计划时的日期为2024-3-10
input:{
  "topic": "上手永劫无间",
  "description": "我是一个永劫无间新手，也没玩过类似的游戏",
  "cycle": "一周"
}

output: {
  "2025-03-10": [
    { id: "1", title: "观看游戏视频", details: "观看一些选手的比赛视频，学习他们的操作和策略" },
  ],
  "2025-03-13": [
    { id: "1", title: "试玩不同角色", details: "在游戏中试玩不同角色，了解每个角色的技能和特点" },
  ],
  "2025-03-17": [
    { id: "1", title: "总结攻略", details: "整理出每个角色的优缺点，并记录初步的游戏策略" },
  ],
}
`, now.Format("2006-01-02"))

	messageJSON, err := json.Marshal(question) //转化为json格式的字节切片
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "解析失败"})
		return nil
	}

	// 发送聊天请求
	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(string(messageJSON)),
				openai.SystemMessage(systemMessage),
			}),
			Model: openai.F(openai.ChatModelGPT4oMini),
		})
	if err != nil {
		panic(err.Error())
	}

	// 获取返回的内容
	responses := chatCompletion.Choices[0].Message.Content
	fmt.Println("Raw Response:", responses) // 打印原始返回值，便于调试

	// 解析JSON
	var result map[string][]map[string]string
	err = json.Unmarshal([]byte(responses), &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "AI返回数据格式错误"})
		return nil
	}
	// 打印解析结果
	fmt.Println(result)
	return result
}
