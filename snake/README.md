# 🐍 Snake Game

A classic snake game implementation based on a frontend-backend separation architecture.

## 🎮 Game Introduction

This is a modern implementation of the classic snake game, featuring the traditional gameplay while adding enhanced features such as leaderboards and obstacles. The game uses Go language for the backend and Vue 3 for the frontend, providing a smooth gaming experience.

## 💻 Tech Stack

- **Backend**:
  - Go Language
  - Beego Framework
  - MySQL Database
  - Redis Cache
- **Frontend**:
  - Vue 3
  - Vite

## ✨ Game Features

- Classic snake gameplay
- Real-time leaderboard functionality
- Auto-generated obstacles to increase challenge
- Scoring system (time*1 + bean count*10)
- Responsive design, supporting different devices

## 🚀 Installation and Launch

### 📋 Environment Requirements

- Go 1.16+
- Node.js 16+
- MySQL
- Redis

### 🔧 Backend Startup Steps

1. Enter the backend directory
```bash
cd ./snake/backend
```

2. Install dependencies
```bash
go mod tidy
```

3. Configure database connections
Edit the `conf/app.conf` file and set MySQL and Redis connection information:
```
mysql.user = root
mysql.password = your_password
mysql.db = snake_game

redis.host = 127.0.0.1
redis.port = 6379
```

4. Start the backend service
```bash
go run main.go
```
The backend service runs on http://localhost:8080 by default

### 🌐 Frontend Startup Steps

1. Enter the frontend directory
```bash
cd ./snake/frontend
```

2. Install dependencies
```bash
npm install
```

3. Start the development server
```bash
npm run dev
```
The frontend service runs on http://localhost:3001 by default

4. Build production version (optional)
```bash
npm run build
```

## 📜 Game Rules

1. **Basic Gameplay**: Control the snake to move, eat food to make it longer and score points
2. **Controls**: Use WASD keys to control the snake's direction
   - W - Up
   - S - Down
   - A - Left
   - D - Right
3. **Game Over Conditions**:
   - Snake head hits the wall
   - Snake head hits its own body
   - Snake head hits an obstacle
   - No food eaten within 10 seconds
4. **Scoring Rules**:
   - Each second survived: +1 point
   - Each food eaten: +10 points
5. **Obstacles**: Random obstacles are generated during gameplay to increase difficulty

## 🎨 Game Interface

The game interface includes the following elements:
- Game title
- Game status display
- Main game interface (15x15 grid)
- Game information panel (score, time, bean count)
- Control buttons (start game, save score, view leaderboard)
- Operation instructions
- Leaderboard dialog

## 🏆 Leaderboard Feature

The game supports saving score records and viewing the leaderboard. After the game ends, you can choose to save your current score and then view the best scores from all players in the leaderboard.

The leaderboard displays the following information:
- Rank
- Score
- Game duration
- Bean count
- Game date

## ⚙️ Configuration Instructions

### 🖥️ Backend Configuration

Main configuration is in the `conf/app.conf` file:

- `httpport`: HTTP service port
- `wall.max`: Maximum number of obstacles (default 6)
- `game.speed`: Game speed (default 200ms)
- MySQL and Redis connection configurations

### 🎯 Game Parameters

- Game map size: 15x15
- Initial snake length: 3
- Maximum obstacle count: 6 (configurable)
- Game update interval: 200ms (configurable)

## 👨‍💻 Development Guide

### 📁 Project Structure

```
snake/
├── backend/              # Backend code
│   ├── conf/             # Configuration files
│   ├── controllers/      # Controllers
│   ├── models/           # Data models
│   ├── routers/          # Router configuration
│   ├── utils/            # Utility functions
│   ├── views/            # Views (optional)
│   ├── main.go           # Entry file
│   └── go.mod            # Go module file
└── frontend/             # Frontend code
    ├── public/           # Static resources
    ├── src/              # Source code
    │   ├── assets/       # Resource files
    │   ├── components/   # Vue components
    │   ├── utils/        # Utility functions
    │   ├── App.vue       # Main component
    │   └── main.js       # Entry file
    ├── index.html        # HTML template
    ├── package.json      # NPM configuration
    └── vite.config.js    # Vite configuration
```

## ❓ Common Issues

### 🚨 Backend Service Startup Failure

1. Check if database connection configuration is correct
2. Ensure MySQL and Redis services are running
3. Check console error messages

### 🔌 Frontend Cannot Connect to Backend

1. Check if backend service is running properly
2. Confirm if frontend API call address is correct
3. Check for potential CORS issues