<template>
  <div class="game-container">
    <h1>Tetris Game</h1>
    <div class="game-info">
      <div class="info-item">
        <span class="label">Score:</span>
        <span class="value">{{ gameState.score }}</span>
      </div>
      <div class="info-item">
        <span class="label">Lines:</span>
        <span class="value">{{ gameState.lines }}</span>
      </div>
      <div class="info-item">
        <span class="label">Level:</span>
        <span class="value">{{ gameState.level }}</span>
      </div>
    </div>
    
    <div class="game-main">
      <GameBoard 
        :game-board="gameState.board.grid" 
        :current-piece="gameState.currentPiece" 
        :next-piece="gameState.nextPiece"
        :game-over="gameState.gameOver"
      />
      
      <div class="game-controls">
        <button @click="startGame" v-if="!gameState.gameStarted || gameState.gameOver" class="control-button">
          {{ gameState.gameOver ? 'Restart Game' : 'Start Game' }}
        </button>
        <button @click="pauseGame" v-if="gameState.gameStarted && !gameState.gameOver" class="control-button">
          {{ gameState.paused ? 'Resume Game' : 'Pause Game' }}
        </button>
        
        <div class="controls-info">
          <h3>Controls:</h3>
          <ul>
            <li><strong>←/→</strong> Move Left/Right</li>
            <li><strong>↓</strong> Soft Drop</li>
            <li><strong>C</strong> Hard Drop</li>
            <li><strong>Z</strong> Rotate</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import GameBoard from './components/GameBoard.vue'
import { startNewGame, getGameState, sendGameAction, endGame } from './utils/gameService.js'

export default {
  name: 'App',
  components: {
    GameBoard
  },
  setup() {
    const gameState = ref({
      gameId: null,
      gameStarted: false,
      gameOver: false,
      paused: false,
      score: 0,
      lines: 0,
      level: 1,
      board: [],
      currentPiece: null,
      nextPiece: null
    })

    let gameLoop = null

    const startGame = async () => {
      try {
        // 创建新游戏
        const newGame = await startNewGame()
        gameState.value.gameId = newGame.gameId
        gameState.value.gameStarted = true
        gameState.value.gameOver = false
        gameState.value.paused = false
        gameState.value.score = 0
        gameState.value.lines = 0
        gameState.value.level = 1
        gameState.value.board = newGame.board
        gameState.value.currentPiece = newGame.currentPiece
        gameState.value.nextPiece = newGame.nextPiece

        // 开始游戏循环
        startGameLoop()
      } catch (error) {
        console.error('Failed to start game:', error)
      }
    }

    const pauseGame = () => {
      gameState.value.paused = !gameState.value.paused
    }

    const startGameLoop = () => {
      // 清除之前的循环
      if (gameLoop) {
        clearInterval(gameLoop)
      }

      // 每100毫秒更新一次游戏状态
      gameLoop = setInterval(async () => {
        if (!gameState.value.gameStarted || gameState.value.gameOver || gameState.value.paused) {
          return
        }

        try {
          // 获取最新游戏状态
          const updatedState = await getGameState(gameState.value.gameId)
          gameState.value.score = updatedState.score
          gameState.value.lines = updatedState.lines
          gameState.value.level = updatedState.level
          gameState.value.board = updatedState.board
          gameState.value.currentPiece = updatedState.currentPiece
          gameState.value.nextPiece = updatedState.nextPiece
          gameState.value.gameOver = updatedState.gameOver

          if (updatedState.gameOver) {
            endGameLoop()
          }
        } catch (error) {
          console.error('Failed to update game state:', error)
        }
      }, 100)
    }

    const endGameLoop = () => {
      if (gameLoop) {
        clearInterval(gameLoop)
        gameLoop = null
      }
    }

    // 键盘事件处理
    const handleKeyDown = async (event) => {
      if (!gameState.value.gameStarted || gameState.value.gameOver || gameState.value.paused) {
        return
      }

      let action = ''

      switch (event.key) {
        case 'ArrowLeft':
          action = 'left'
          break
        case 'ArrowRight':
          action = 'right'
          break
        case 'ArrowDown':
          action = 'down'
          break
        case 'z':
        case 'Z':
          action = 'rotate'
          break
        case 'c':
        case 'C':
          action = 'hardDrop'
          break
        default:
          return
      }

      // 阻止默认行为
      event.preventDefault()

      try {
        await sendGameAction(gameState.value.gameId, action)
      } catch (error) {
        console.error('Failed to send game action:', error)
      }
    }

    onMounted(() => {
      window.addEventListener('keydown', handleKeyDown)
    })

    onUnmounted(() => {
      window.removeEventListener('keydown', handleKeyDown)
      endGameLoop()
      if (gameState.value.gameId) {
        endGame(gameState.value.gameId)
      }
    })

    return {
      gameState,
      startGame,
      pauseGame
    }
  }
}
</script>

<style>
.game-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  text-align: center;
  font-family: Arial, sans-serif;
}

.game-info {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
  gap: 30px;
}

.info-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.label {
  font-size: 14px;
  color: #666;
  margin-bottom: 5px;
}

.value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.game-main {
  display: flex;
  justify-content: center;
  gap: 30px;
  flex-wrap: wrap;
}

.control-button {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
  border-radius: 5px;
  margin-bottom: 20px;
  transition: background-color 0.3s;
}

.control-button:hover {
  background-color: #45a049;
}

.controls-info {
  text-align: left;
  background-color: #f5f5f5;
  padding: 15px;
  border-radius: 5px;
}

.controls-info h3 {
  margin-top: 0;
  margin-bottom: 10px;
  color: #333;
}

.controls-info ul {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.controls-info li {
  margin-bottom: 8px;
  color: #666;
}

.controls-info strong {
  color: #333;
}
</style>
