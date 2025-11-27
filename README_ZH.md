# 🎮 EasyGames

一个简洁的小游戏合集，用于休闲时打发时间。

## 📝 项目介绍

EasyGames是一个集合了各种经典小游戏的项目，旨在提供简单有趣的休闲游戏体验。每个游戏都是独立的模块，采用前后端分离架构实现。

## 📁 项目结构

```
EasyGames/
├── README.md           # 项目说明文档
└── snake/              # 贪吃蛇游戏
    ├── backend/        # 后端代码（Go语言）
    └── frontend/       # 前端代码（Vue 3）
```

## 🎯 游戏列表

### 🐍 1. 贪吃蛇游戏 (snake)

- **技术栈**：
  - 后端：Go + Beego框架 + Redis + MySQL
  - 前端：Vue 3 + Vite
- **游戏特性**：
  - 经典贪吃蛇玩法
  - 实时排行榜
  - 自动生成障碍物
  - 积分系统（时间*1 + 豆子数量*10）

## 💻 开发说明

### 📋 环境要求

- Go 1.16+
- Node.js 16+
- MySQL
- Redis

### 👨‍💻 参与开发

1. Fork 本仓库
2. 创建你的功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的修改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🤝 贡献

欢迎各种形式的贡献，包括但不限于 bug 修复、功能增强、文档改进等。