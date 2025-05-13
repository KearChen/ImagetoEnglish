# ImagetoEnglish — 拍照生成英文的轻量程序

🧠 基于大语言模型，**ImagetoEnglish** 是一个「拍照即可记单词」的程序 Demo。

用户只需打开页面拍一张照片，程序会自动识别图中物品，并输出目标语言下的结构化描述（如中英文名称、英文释义等）。

> ✅ 轻量、易改、免登录、免费开源
> 📱 微信 / 浏览器环境均可用，前端适配微信授权逻辑
> 🤖 后端默认接入 GLM-4V 图文识别模型

---

## ✨ 功能亮点

* 📷 打开即拍，快速识别
* 🌐 中文 + 英文词汇 / 释义自动生成
* 🎯 输出结构化结果，适合词库构建
* ⚙️ 图片自动压缩，节省带宽
* 🌱 微信拍照兼容处理（避免黑屏 / 授权失败）
* 📦 易于扩展，如添加“单词本”“用户系统”等功能

---

## 🧠 AI Prompt（图像识别 → 英文信息）

```text
你是一个图像识别和语言学习助手。用户上传一张图片，你的任务是：
判断图中是什么物体（例如：苹果、书包、椅子、眼镜、狗、猫等）；
给出该物体的标准中文名称和对应英文名称；
提供一句最简单的英文描述。
请严格按照如下格式仅输出 JSON，不要输出代码或其他信息，JSON 字段使用顿号【、】区隔：

{
  "中文名称": "",
  "英文名称": "",
  "英文描述": ""
}
```

---

## 🚀 快速开始

> ⚠️ **注意：程序需要部署在公网服务器上运行，确保图片 URL 可被大模型访问。**
> 本地开发环境生成的链接（如 `localhost`、`127.0.0.1`）无法被 AI 模型读取，请使用内网穿透或部署到云服务器。


### 1. 克隆项目

```bash
git clone https://github.com/KearChen/ImagetoEnglish.git
```

---

### 2. 配置说明

运行前，请创建一个 `config.yaml` 文件，并填入以下内容：

```yaml
app:
  name: ImagetoEnglish           # 应用名称，可自定义
  port: 8090                     # 后端监听端口
  SetMode: release               # 运行模式：debug / release

Ai:
  models:
    glm-4v-flash:
      apiKey: <你的API密钥>
      base_url: https://open.bigmodel.cn/api/paas/v4/chat/completions
      model: glm-4v-flash
      temperature: 0.95
```

#### 📌 字段说明：

| 字段            | 说明                                |
| ------------- | --------------------------------- |
| `apiKey`      | 替换为你在智谱 AI（BigModel）平台申请的 API Key |
| `base_url`    | 模型接口地址（默认使用智谱 GLM-4V Flash）       |
| `model`       | 使用的模型名称，例如 `glm-4v-flash`         |
| `temperature` | 控制模型生成结果的多样性，建议在 0.7～1.0 之间       |

配置完成后，即可运行后端程序，配合前端进行图像上传与分析。

---

### 3. 启动服务

```bash
go mod tidy
go run main.go
```
---
### 4. 部署与调用注意事项

#### ✅ 你可以使用 `curl` 手动测试接口：

```bash
curl --location --request POST 'https://你的域名/v1/analyze' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Connection: keep-alive' \
--data-raw '{
  "image_url": "https://xxx.com/your-image.jpg"
}'
```

📌 `image_url` 必须是外网可访问的图片链接，如对象存储、CDN、图床等地址。

---
## 🧩 TODO

* [ ] 添加“我的单词本”本地存储功能
* [ ] 提供多语言切换（日语 / 法语等）
* [ ] 增加识别历史记录管理
* [ ] 适配更多模型

---

## 📄 License

本项目采用 **MIT License**，欢迎二次开发和商用，但请保留原始出处链接。