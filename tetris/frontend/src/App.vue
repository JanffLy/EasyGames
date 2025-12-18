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
        
        <!-- 保存记录表单 -->
        <div class="save-record-form" v-if="gameState.gameOver && !recordSaved">
          <h3>Game Over!</h3>
          <p>Final Score: {{ gameState.score }}</p>
          <div class="form-group">
            <label for="playerName">Enter your name:</label>
            <input 
              type="text" 
              id="playerName" 
              v-model="playerName" 
              placeholder="Your name" 
              maxlength="20"
              @keyup.enter="saveRecord"
            >
          </div>
          <button @click="saveRecord" class="save-button">Save Record</button>
        </div>
        
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
    
    <!-- 排行榜区域 -->
    <div class="leaderboard-section">
      <h2>Leaderboard</h2>
      <div class="leaderboard">
        <div class="leaderboard-header">
          <div class="rank">Rank</div>
          <div class="name">Player</div>
          <div class="score">Score</div>
          <div class="lines">Lines</div>
        </div>
        <div class="leaderboard-body">
          <div class="leaderboard-item" v-for="(item, index) in leaderboard" :key="index">
            <div class="rank">{{ index + 1 }}</div>
            <div class="name">{{ item.player_name }}</div>
            <div class="score">{{ item.score }}</div>
            <div class="lines">{{ item.lines }}</div>
          </div>
          <div class="no-records" v-if="leaderboard.length === 0">
            No records yet. Be the first! 
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import GameBoard from './components/GameBoard.vue'
import { startNewGame, getGameState, sendGameAction, endGame, saveGameRecord, getLeaderboard } from './utils/gameService.js'

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
      board: { grid: [] },
      currentPiece: null,
      nextPiece: null
    })
    
    const playerName = ref('')
    const recordSaved = ref(false)
    const leaderboard = ref([])

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
          recordSaved.value = false // 重置记录保存状态
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

    // 保存游戏记录
    const saveRecord = async () => {
      if (!playerName.value.trim() || !gameState.value.gameId) return
      
      try {
        await saveGameRecord(gameState.value.gameId, playerName.value.trim(), gameState.value.score)
        recordSaved.value = true
        playerName.value = ''
        fetchLeaderboard() // 更新排行榜
      } catch (error) {
        console.error('Failed to save record:', error)
      }
    }
    
    // 获取排行榜
    const fetchLeaderboard = async () => {
      try {
        const data = await getLeaderboard(10)
        // 确保leaderboard.value始终是数组，避免null导致的TypeError
        leaderboard.value = data || []
      } catch (error) {
        console.error('Failed to fetch leaderboard:', error)
        // 发生错误时也确保leaderboard.value是数组
        leaderboard.value = []
      }
    }
    
    // 监听游戏结束状态，自动获取最新排行榜
    watch(() => gameState.value.gameOver, (newVal) => {
      if (newVal) {
        fetchLeaderboard()
      }
    })
    
    onMounted(() => {
      window.addEventListener('keydown', handleKeyDown)
      fetchLeaderboard() // 初始化时获取排行榜
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
      playerName,
      recordSaved,
      leaderboard,
      startGame,
      pauseGame,
      saveRecord
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
  
  /* 保存记录表单样式 */
  .save-record-form {
    margin-top: 20px;
    padding: 15px;
    background-color: #f5f5f5;
    border-radius: 5px;
    text-align: center;
  }
  
  .save-record-form h3 {
    margin-top: 0;
    color: #333;
  }
  
  .save-record-form p {
    font-size: 18px;
    font-weight: bold;
    color: #4CAF50;
    margin: 10px 0;
  }
  
  .form-group {
    margin: 15px 0;
  }
  
  .form-group label {
    display: block;
    margin-bottom: 5px;
    color: #666;
    font-size: 14px;
  }
  
  .form-group input {
    padding: 8px;
    border: 1px solid #ddd;
    border-radius: 3px;
    font-size: 16px;
    width: 100%;
    box-sizing: border-box;
  }
  
  .save-button {
    background-color: #2196F3;
    color: white;
    border: none;
    padding: 10px 20px;
    font-size: 16px;
    cursor: pointer;
    border-radius: 5px;
    transition: background-color 0.3s;
  }
  
  .save-button:hover {
    background-color: #0b7dda;
  }
  
  /* 排行榜样式 */
  .leaderboard-section {
    margin-top: 40px;
    padding: 20px;
    background-color: #f9f9f9;
    border-radius: 10px;
  }
  
  .leaderboard-section h2 {
    margin-top: 0;
    color: #333;
    text-align: center;
  }
  
  .leaderboard {
    max-width: 600px;
    margin: 0 auto;
    border-radius: 5px;
    overflow: hidden;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  }
  
  .leaderboard-header {
    display: flex;
    background-color: #333;
    color: white;
    font-weight: bold;
    padding: 10px;
  }
  
  .leaderboard-header .rank, .leaderboard-header .score, .leaderboard-header .lines {
    width: 80px;
    text-align: center;
  }
  
  .leaderboard-header .name {
    flex: 1;
    text-align: center;
  }
  
  .leaderboard-body {
    background-color: white;
  }
  
  .leaderboard-item {
    display: flex;
    padding: 10px;
    border-bottom: 1px solid #eee;
  }
  
  .leaderboard-item:last-child {
    border-bottom: none;
  }
  
  .leaderboard-item .rank, .leaderboard-item .score, .leaderboard-item .lines {
    width: 80px;
    text-align: center;
  }
  
  .leaderboard-item .name {
    flex: 1;
    text-align: center;
    font-weight: bold;
  }
  
  .no-records {
    padding: 20px;
    text-align: center;
    color: #666;
    font-style: italic;
  }
</style>
