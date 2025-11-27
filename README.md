# 🎮 EasyGames

A concise collection of small games for passing time during leisure moments.

## 📝 Project Introduction

EasyGames is a project that brings together various classic small games, aiming to provide simple and enjoyable casual gaming experiences. Each game is an independent module implemented with a frontend-backend separation architecture.

## 📁 Project Structure

```
EasyGames/
├── README.md           # Project documentation
└── snake/              # Snake game
    ├── backend/        # Backend code (Go language)
    └── frontend/       # Frontend code (Vue 3)
```

## 🎯 Game List

### 🐍 1. Snake Game (snake)

- **Tech Stack**:
  - Backend: Go + Beego framework + Redis + MySQL
  - Frontend: Vue 3 + Vite
- **Game Features**:
  - Classic snake gameplay
  - Real-time leaderboard
  - Auto-generated obstacles
  - Scoring system (time*1 + bean count*10)

## 💻 Development Guide

### 📋 Environment Requirements

- Go 1.16+
- Node.js 16+
- MySQL
- Redis

### 👨‍💻 Contributing

1. Fork this repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributions

All forms of contributions are welcome, including but not limited to bug fixes, feature enhancements, and documentation improvements.