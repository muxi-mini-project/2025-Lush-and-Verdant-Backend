package client

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/config"
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"net/http"
	"regexp"
)

type ChatGptClient struct {
	cfg *config.ChatGptConfig
}

func NewChatGptClient(cfg *config.ChatGptConfig) *ChatGptClient {
	return &ChatGptClient{cfg: cfg}
}

func (cg *ChatGptClient) AskForGoal(c *gin.Context) map[string]string {
	// 设置初始配置
	var options = []option.RequestOption{
		option.WithBaseURL("https://api.chatanywhere.tech"), //更换了BaseURL,openai的太贵了
		option.WithAPIKey(cg.cfg.Sdk),
	}

	// 设置自定义的 baseurl 和 API 密钥
	client := openai.NewClient(options...)

	systemMessage := `
    角色定位:
		你是一个根据用户输入的主题，帮助用户制定可执行的计划及对应的实践时间的AI
       	用户将发送给你一个学习主题,你需要将其拆解为多个一个月内以天为周期的计划表，并为这几个周期段提供可实施的具体事件。
    
    格式说明:
       	1. 输入topic表示学习的主题,description表示对用户情况的描述,cycle表示用户希望分解成的几个阶段。
       	2. 用户输入的cycle表示用户预计完成这些事情的总时间，例如输入2个月，你自行记录制定时的日期，并且从下一天开始记录，你需要制定两个月内合理的能够完成的事情。
       	3. 你需要给用户返回总的topic，以及制定的计划，计划包含event表示需要做的事件，并附上它的开始结束时间。
		4. 返回的event又包含多个具体实施的事情，这个内容及时间由你来定义，我相信你的能力。
       	5. 返回格式为json，不允许包含其他任何字符。
    
    输入输出示例1:
	假设提问制定计划时的日期为2024-10-17
input:{
  "topic": "后端开发",
  "description": "我是一个编程领域小白，没有任何编程基础",
  "cycle": "2个月"
}
    
output:{
  "tasks": [
    {
      "name": "了解程序员基础",
      "description": "学习程序员必备基础",
      "start_time": "2024-10-18 00:00:00",
      "end_time": "2024-10-26 23:59:59",
      "events": [
        {
          "name": "版本控制",
          "description": "学习Git指令，并创建一个GitHub账号",
          "start_time": "2024-10-18 00:00:00",
          "end_time": "2024-10-19 23:59:59"
        },
        {
          "name": "Markdown",
          "description": "学习轻量级标记语言Markdown，并多加练习",
          "start_time": "2024-10-20 00:00:00",
          "end_time": "2024-10-22 23:59:59"
        },
        {
          "name": "Linux",
          "description": "安装Linux操作系统，了解相关使用方法",
          "start_time": "2024-10-23 00:00:00",
          "end_time": "2024-10-26 23:59:59"
        }
      ]
    },
    {
      "name": "本地和云服务之间程序的传输",
      "description": "学习如何使用docker在云服务器部署自己的程序",
      "start_time": "2024-10-27 00:00:00",
      "end_time": "2024-11-2 23:59:59",
      "events": [
        {
          "name": "云服务申请",
          "description": "申请一台自己的云服务器",
          "start_time": "2024-10-27 00:00:00",
          "end_time": "2024-10-27 23:59:59"
        },
        {
          "name": "Docker学习",
          "description": "学习Docker相关指令，尝试用它在云服务上部署镜像容器",
          "start_time": "2024-10-28 00:00:00",
          "end_time": "2024-10-29 23:59:59"
        },
        {
          "name": "Go语言基础",
          "description": "学习Go语言基础语法",
          "start_time": "2024-10-30 00:00:00",
          "end_time": "2024-11-2 23:59:59"
        }
      ]
    },
    {
      "name": "深入学习语言",
      "description": "了解基础语法后，深入学习相关知识",
      "start_time": "2024-11-3 00:00:00",
      "end_time": "2024-11-9 23:59:59",
      "events": [
        {
          "name": "Go模块化管理",
          "description": "Go Model可以对文件进行模块化管理，学习相关知识",
          "start_time": "2024-11-3 00:00:00",
          "end_time": "2024-11-4 23:59:59"
        },
        {
          "name": "Go并发",
          "description": "Go语言的并发很强，学习并发的语法和用法",
          "start_time": "2024-11-5 00:00:00",
          "end_time": "2024-11-9 23:59:59"
        }
      ]
    },
    {
      "name": "部署网络环境",
      "description": "对语言编辑深入了解后，可以尝试在网络上部署自己的URL链接",
      "start_time": "2024-11-10 00:00:00",
      "end_time": "2024-11-16 23:59:59",
      "events": [
        {
          "name": "http基础",
          "description": "在部署网络前，要先了解网络的运行机制",
          "start_time": "2024-11-10 00:00:00",
          "end_time": "2024-11-11 23:59:59"
        },
        {
          "name": "JSON",
          "description": "学习Go语言结构体对应的JSON格式",
          "start_time": "2024-11-12 00:00:00",
          "end_time": "2024-11-14 23:59:59"
        },
        {
          "name": "鉴权",
          "description": "学习网络保存会话的机制，Cookie、Session、Token、JWT的作用和区别",
          "start_time": "2024-11-15 00:00:00",
          "end_time": "2024-11-16 23:59:59"
        }
      ]
    },
    {
      "name": "简单web服务",
      "description": "Go是一个为web而生的语言，尝试用Go语言完成一个完整的web服务",
      "start_time": "2024-11-17 00:00:00",
      "end_time": "2024-11-24 23:59:59",
      "events": [
        {
          "name": "Go Web",
          "description": "学习Go Web并尝试搭建一个用户管理系统",
          "start_time": "2024-11-17 00:00:00",
          "end_time": "2024-11-24 23:59:59"
        }
      ]
    },
    {
      "name": "爬虫模拟请求",
      "description": "Go语言还能实现模拟登录和爬取数据进行存储",
      "start_time": "2024-11-25 00:00:00",
      "end_time": "2024-12-1 23:59:59",
      "events": [
        {
          "name": "模拟登录",
          "description": "学习使用Go语言模拟登录",
          "start_time": "2024-11-25 00:00:00",
          "end_time": "2024-11-27 23:59:59"
        },
        {
          "name": "爬虫",
          "description": "学习如何使用Go语言爬取数据",
          "start_time": "2024-11-27 00:00:00",
          "end_time": "2024-12-1 23:59:59"
        }
      ]
    },
    {
      "name": "数据库基础",
      "description": "作为一名后端工程师，学习数据库并存储数据是必不可少的",
      "start_time": "2024-12-2 00:00:00",
      "end_time": "2024-12-8 23:59:59",
      "events": [
        {
          "name": "数据库基础",
          "description": "学习数据的分类，了解关系型数据库和非关系型数据库",
          "start_time": "2024-12-2 00:00:00",
          "end_time": "2024-12-3 23:59:59"
        },
        {
          "name": "数据库指令",
          "description": "学习对数据库的操作",
          "start_time": "2024-12-4 00:00:00",
          "end_time": "2024-12-6 23:59:59"
        },
        {
          "name": "数据库使用",
          "description": "尝试使用数据库保存之前用户管理系统传入的数据",
          "start_time": "2024-12-7 00:00:00",
          "end_time": "2024-12-8 23:59:59"
        }
      ]
    },
    {
      "name": "综合项目",
      "description": "学习必备的知识，可以尝试完成一个综合项目",
      "start_time": "2024-12-9 00:00:00",
      "end_time": "2024-12-16 23:59:59",
      "events": [
        {
          "name": "学习swagger文档",
          "description": "学习使用swagger文档是和前端对接的关键",
          "start_time": "2024-12-9 00:00:00",
          "end_time": "2024-12-10 23:59:59"
        },
        {
          "name": "学习对接",
          "description": "学习使用Apifox完成对接",
          "start_time": "2024-12-11 00:00:00",
          "end_time": "2024-12-11 23:59:59"
        },
        {
          "name": "部署系统",
          "description": "与一个前端伙伴完成对接并部署",
          "start_time": "2024-12-11 00:00:00",
          "end_time": "2024-12-16 23:59:59"
        }
      ]
    }
  ]
}

	输入输出示例2:
	假设提问制定计划时的日期为2024-10-17
input:{
  "topic": "上手永劫无间",
  "description": "我是一个永劫无间新手，也没玩过类似的游戏",
  "cycle": "一周"
}

output:{
  "tasks": [
    {
      "name": "熟悉永劫无间",
      "description": "了解永劫无间的基本玩法和角色特点",
      "start_time": "2024-10-18 00:00:00",
      "end_time": "2024-10-24 23:59:59",
      "events": [
        {
          "name": "观看游戏视频",
          "description": "观看一些顶级选手的比赛视频，学习他们的操作和策略",
          "start_time": "2024-10-18 00:00:00",
          "end_time": "2024-10-18 23:59:59"
        },
        {
          "name": "试玩不同角色",
          "description": "在游戏中试玩不同角色，了解每个角色的技能和特点",
          "start_time": "2024-10-19 00:00:00",
          "end_time": "2024-10-21 23:59:59"
        },
        {
          "name": "总结攻略",
          "description": "整理出每个角色的优缺点，并记录初步的游戏策略",
          "start_time": "2024-10-22 00:00:00",
          "end_time": "2024-10-24 23:59:59"
        }
      ]
    }
  ]
}
`
	var message request.Question
	if err := c.ShouldBind(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return nil
	}

	messageJSON, err := json.Marshal(message) //转化为json格式的字节切片
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
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
	response := chatCompletion.Choices[0].Message.Content
	//fmt.Println("Raw Response:", response)

	// 匹配 JSON 格式的 key-value
	re := regexp.MustCompile(`"([\w-]+)":\s*"([^"]+)"`)

	// 提取所有匹配结果
	matches := re.FindAllStringSubmatch(response, -1)

	// 构造 map
	result := make(map[string]string)
	for _, match := range matches {
		key := match[1]   // JSON 的 key
		value := match[2] // JSON 的 value
		result[key] = value
	}

	// 打印解析结果
	//fmt.Println(result)
	return result
}
