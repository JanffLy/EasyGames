<template>
  <div class="game-board" @keydown="handleKeyDown" tabindex="0" ref="gameBoard">
    <div 
      v-for="(row, rowIndex) in grid" 
      :key="rowIndex" 
      class="row"
    >
      <div 
        v-for="(cell, colIndex) in row" 
        :key="colIndex" 
        class="cell"
        :class="getCellClass(colIndex, rowIndex)"
      ></div>
    </div>
    <div v-if="gameState.status === 'ended'" class="game-over">
      <div class="game-over-content">
        <h2>游戏结束</h2>
        <p>最终得分: {{ gameState.score }}</p>
        <p>游戏时间: {{ gameState.time }}秒</p>
        <p>吃到豆子: {{ gameState.foodCount }}个</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted, watch, nextTick } from 'vue'

// 定义props
const props = defineProps({
  gameState: {
    type: Object,
    required: true
  }
})

// 定义emits
const emit = defineEmits(['direction-change'])

// 响应式数据
const cellSize = ref(20) // 每个格子的大小（像素）
const gameBoard = ref(null)

// 计算游戏网格 - 确保15x15的游戏区域
const grid = computed(() => {
  const width = 15 // 固定宽度为15
  const height = 15 // 固定高度为15
  const result = []
  for (let y = 0; y < height; y++) {
    const row = []
    for (let x = 0; x < width; x++) {
      row.push({ x, y })
    }
    result.push(row)
  }
  return result
})

// 获取单元格的CSS类
const getCellClass = (x, y) => {
  const classes = []
  
  // 检查是否是蛇的身体
  const snakeBody = props.gameState?.snake?.body || []
  for (let i = 0; i < snakeBody.length; i++) {
    if (snakeBody[i].x === x && snakeBody[i].y === y) {
      if (i === 0) {
        classes.push('snake-head')
        // 根据方向添加额外的类
        const direction = props.gameState.snake.direction
        switch(direction) {
          case 0: classes.push('direction-up'); break // 上
          case 1: classes.push('direction-down'); break // 下
          case 2: classes.push('direction-left'); break // 左
          case 3: classes.push('direction-right'); break // 右
        }
      } else {
        classes.push('snake-body')
      }
      return classes
    }
  }
  
  // 检查是否是食物（黄色小球）
  const food = props.gameState?.food?.position
  if (food && food.x === x && food.y === y) {
    classes.push('food')
    return classes
  }
  
  // 检查是否是墙体
  const walls = props.gameState?.walls || []
  for (const wall of walls) {
    if (wall.position.x === x && wall.position.y === y) {
      classes.push('wall')
      return classes
    }
  }
  
  // 普通单元格
  classes.push('empty')
  return classes
}

// 处理键盘事件 - WASD控制
const handleKeyDown = (event) => {
  // 阻止默认行为，防止页面滚动
  event.preventDefault()
  
  const key = event.key.toLowerCase()
  let direction = null
  
  // 根据WASD键位映射到方向
  switch(key) {
    case 'w': direction = 'up'; break
    case 's': direction = 'down'; break
    case 'a': direction = 'left'; break
    case 'd': direction = 'right'; break
    default: return // 不是有效的控制键
  }
  
  // 发射方向变化事件
  emit('direction-change', direction)
}

// 聚焦游戏区域
const focusGameBoard = () => {
  if (gameBoard.value) {
    gameBoard.value.focus()
  }
}

// 组件挂载时，聚焦到游戏区域以捕获键盘事件
onMounted(() => {
  focusGameBoard()
})

// 监听游戏状态变化
watch(() => props.gameState, () => {
  nextTick(() => {
    focusGameBoard()
  })
}, { deep: true })
</script>

<style scoped>
.game-board {
  display: inline-block;
  border: 2px solid #333;
  background-color: #f5f5f5;
  outline: none;
  position: relative;
  user-select: none;
}

.row {
  display: flex;
}

.cell {
  width: 20px;
  height: 20px;
  border: 1px solid #ddd;
  box-sizing: border-box;
}

/* 蛇头样式 - 使用方块表示 */
.snake-head {
  background-color: #4CAF50;
  border: 1px solid #2E7D32;
}

/* 蛇身体样式 - 使用方块表示 */
.snake-body {
  background-color: #2E7D32;
  border: 1px solid #1B5E20;
}

/* 食物样式 - 黄色小球 */
.food {
  background-color: #FFEB3B;
  border-radius: 50%;
  border: 1px solid #FDD835;
  box-shadow: 0 0 3px rgba(255, 235, 59, 0.8);
}

/* 墙体样式 */
.wall {
  background-color: #757575;
  border: 1px solid #616161;
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.3);
}

/* 空单元格样式 */
.empty {
  background-color: #f5f5f5;
}

/* 游戏结束遮罩 */
.game-over {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  color: white;
}

.game-over-content {
  text-align: center;
  background-color: #333;
  padding: 20px;
  border-radius: 10px;
}

.game-over-content h2 {
  margin-bottom: 10px;
  color: #f44336;
}

.game-over-content p {
  margin: 8px 0;
  font-size: 16px;
}

/* 方向指示器 */
.direction-up::before {
  content: '▲';
  color: white;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 14px;
}

.direction-down::before {
  content: '▼';
  color: white;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 14px;
}

.direction-left::before {
  content: '◀';
  color: white;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 14px;
}

.direction-right::before {
  content: '▶';
  color: white;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 14px;
}
</style>