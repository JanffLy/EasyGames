<template>
  <div class="game-container">
    <h1>贪吃蛇游戏</h1>
    <!-- 游戏状态显示 -->
    <div v-if="gameState" class="game-status">
      <div class="status-text">{{ getStatusText() }}</div>
    </div>
    <div v-else class="loading">
      <p>请点击开始游戏</p>
    </div>
    <!-- 游戏主界面 -->
    <GameBoard 
      v-if="gameState" 
      :game-state="gameState" 
      @direction-change="handleDirectionChange"
    />
    <div class="game-info">
      <div class="score">得分: {{ gameState?.score || 0 }}</div>
      <div class="time">时间: {{ gameState?.time || 0 }}秒</div>
      <div class="beans">豆子: {{ gameState?.foodCount || 0 }}</div>
    </div>
    <div class="controls">
      <button @click="startNewGame" class="btn">开始新游戏</button>
      <button @click="saveScore" v-if="gameState && gameState.status === 'ended'" class="btn">保存得分</button>
      <button @click="showLeaderboard" class="btn btn-secondary">查看排行榜</button>
    </div>
    <div class="instructions">
      <p>使用WASD键控制蛇的移动方向</p>
      <p>积分规则：时间(秒)*1 + 豆子数量*10</p>
    </div>
    
    <!-- 排行榜对话框 -->
    <div v-if="showLeaderboardDialog" class="modal-overlay" @click="closeLeaderboard">
      <div class="modal-content" @click.stop>
        <h2>排行榜</h2>
        <div v-if="loadingLeaderboard" class="loading">
          <p>加载排行榜中...</p>
        </div>
        <div v-else-if="leaderboardData && leaderboardData.length > 0" class="leaderboard">
          <table>
            <thead>
              <tr>
                <th>排名</th>
                <th>得分</th>
                <th>游戏时长(秒)</th>
                <th>豆子数量</th>
                <th>日期</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in leaderboardData" :key="index">
                <td>{{ index + 1 }}</td>
                <td>{{ item.score }}</td>
                <td>{{ item.time_played }}</td>
                <td>{{ item.food_count }}</td>
                <td>{{ formatDate(item.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="empty-leaderboard">
          <p>暂无排行榜数据</p>
        </div>
        <button @click="closeLeaderboard" class="btn btn-close">关闭</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import GameBoard from './components/GameBoard.vue'
import { gameService } from './utils/gameService.js'

// 响应式数据
const gameState = ref(null)
const gameLoop = ref(null)
const showLeaderboardDialog = ref(false)
const leaderboardData = ref([])
const loadingLeaderboard = ref(false)

// 开始新游戏
const startNewGame = async () => {
  try {
    const newGame = await gameService.createGame()
    gameState.value = newGame
    startGameLoop()
  } catch (error) {
    console.error('创建游戏失败:', error)
  }
}

// 开始游戏循环，定时更新游戏状态
const startGameLoop = () => {
  // 清除之前的循环
  if (gameLoop.value) {
    clearInterval(gameLoop.value)
  }

  // 每200ms更新一次游戏状态
  gameLoop.value = setInterval(async () => {
    // 确保gameState和id存在
    if (!gameState.value || !gameState.value.id) {
      console.warn('游戏状态不存在，清除游戏循环');
      clearInterval(gameLoop.value);
      gameLoop.value = null;
      return;
    }
    
    // 只在游戏运行时更新
    if (gameState.value.status === 'running') {
      try {
        console.log('正在获取游戏状态，ID:', gameState.value.id)
        const updatedState = await gameService.getGameState(gameState.value.id)
        console.log('成功获取游戏状态:', updatedState)
        gameState.value = updatedState

        // 检查游戏是否结束
        if (updatedState.status === 'ended') {
          console.log('游戏已结束，停止游戏循环');
          clearInterval(gameLoop.value);
          gameLoop.value = null;
        }
      } catch (error) {
        // 如果是404错误或游戏不存在错误，停止游戏循环
        if (error.message?.includes('Not Found') || error.message?.includes('游戏不存在')) {
          console.error('游戏不存在，停止游戏循环:', error);
          clearInterval(gameLoop.value);
          gameLoop.value = null;
          // 标记游戏为结束状态
          if (gameState.value) {
            gameState.value.status = 'ended';
          }
        } else {
          console.error('更新游戏状态失败:', error)
          console.error('错误详情:', error.message, error.stack)
        }
      }
    } else if (gameState.value.status === 'ended') {
      // 如果游戏已经结束，停止游戏循环
      console.log('游戏状态为ended，停止游戏循环');
      clearInterval(gameLoop.value);
      gameLoop.value = null;
    }
  }, 200)
}

// 处理方向变化
const handleDirectionChange = async (direction) => {
  if (gameState.value && gameState.value.status === 'running') {
    console.log(`处理方向变化: ${direction}, 游戏ID: ${gameState.value.id}`);
    try {
      const result = await gameService.updateDirection(gameState.value.id, direction);
      console.log(`方向更新成功:`, result);
    } catch (error) {
      console.error('更新方向失败:', error);
      console.error('错误详情:', error.message, error.stack);
    }
  } else {
    console.warn(`无法更新方向: 游戏状态不是running，当前状态: ${gameState.value?.status || 'undefined'}`);
  }
}

// 保存得分
const saveScore = async () => {
  if (gameState.value) {
    try {
      const playerName = 'Player'
      await gameService.saveScore(gameState.value.id, playerName, gameState.value.score)
      alert('得分保存成功！')
    } catch (error) {
      console.error('保存得分失败:', error)
    }
  }
}

// 获取游戏状态文本
const getStatusText = () => {
  if (!gameState.value) return ''
  if (gameState.value.status === 'running') {
    return '游戏进行中'
  } else if (gameState.value.status === 'ended') {
    return '游戏结束'
  }
  return '未知状态'
}

// 显示排行榜
const showLeaderboard = async () => {
  showLeaderboardDialog.value = true
  loadingLeaderboard.value = true
  try {
    const data = await gameService.getLeaderboard()
    leaderboardData.value = data
  } catch (error) {
    console.error('获取排行榜失败:', error)
    alert('获取排行榜失败，请稍后重试')
  } finally {
    loadingLeaderboard.value = false
  }
}

// 关闭排行榜
const closeLeaderboard = () => {
  showLeaderboardDialog.value = false
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 组件生命周期
onMounted(() => {
  // 组件挂载时的初始化
})

onUnmounted(() => {
  // 组件卸载时清理
  if (gameLoop.value) {
    clearInterval(gameLoop.value)
  }
})
</script>

<style>
.game-status {
  margin-bottom: 10px;
  font-size: 18px;
  color: #4CAF50;
  font-weight: bold;
}

.status-text {
  padding: 5px 10px;
  background-color: #e8f5e9;
  border-radius: 5px;
  display: inline-block;
}

.game-container {
  text-align: center;
  background-color: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
  max-width: 700px;
  width: 100%;
  margin: 0 auto;
}

h1 {
  color: #333;
  margin-bottom: 20px;
}

.game-info {
  display: flex;
  justify-content: space-around;
  margin: 20px 0;
  font-size: 18px;
  font-weight: bold;
}

.controls {
  margin-top: 20px;
}

.btn {
  background-color: #4CAF50;
  border: none;
  color: white;
  padding: 10px 20px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  margin: 4px 2px;
  cursor: pointer;
  border-radius: 5px;
  transition: background-color 0.3s;
}

.btn:hover {
  background-color: #45a049;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.loading {
  padding: 50px;
  font-size: 18px;
  color: #666;
}

.instructions {
  margin-top: 20px;
  font-size: 14px;
  color: #666;
  line-height: 1.5;
}

.instructions p {
  margin: 5px 0;
}

/* 排行榜样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 10px;
  max-width: 800px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.2);
}

.modal-content h2 {
  margin-top: 0;
  color: #333;
  text-align: center;
  margin-bottom: 20px;
}

.leaderboard {
  margin-bottom: 20px;
}

.leaderboard table {
  width: 100%;
  border-collapse: collapse;
}

.leaderboard th, .leaderboard td {
  padding: 12px 8px;
  text-align: center;
  border-bottom: 1px solid #ddd;
}

.leaderboard th {
  background-color: #4CAF50;
  color: white;
  font-weight: bold;
}

.leaderboard tr:hover {
  background-color: #f5f5f5;
}

.empty-leaderboard {
  text-align: center;
  padding: 40px;
  color: #666;
}

.btn-secondary {
  background-color: #2196F3;
}

.btn-secondary:hover {
  background-color: #1976D2;
}

.btn-close {
  background-color: #f44336;
  width: 100%;
}

.btn-close:hover {
  background-color: #d32f2f;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .game-container {
    padding: 10px;
  }
  
  h1 {
    font-size: 24px;
  }
  
  .game-info {
    font-size: 16px;
  }
  
  .modal-content {
    width: 95%;
    padding: 15px;
  }
  
  .leaderboard th, .leaderboard td {
    padding: 8px 4px;
    font-size: 14px;
  }
}

@media (max-width: 480px) {
  .leaderboard th, .leaderboard td {
    font-size: 12px;
    padding: 6px 2px;
  }
}
</style>