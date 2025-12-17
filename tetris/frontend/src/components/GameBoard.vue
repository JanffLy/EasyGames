<template>
  <div class="game-board-container">
    <div class="main-board">
      <div class="board-grid">
        <div 
          v-for="(row, rowIndex) in gameBoard" 
          :key="'row-' + rowIndex" 
          class="board-row"
        >
          <div 
            v-for="(cell, colIndex) in row" 
            :key="'cell-' + rowIndex + '-' + colIndex" 
            class="board-cell"
            :style="getCellClass(cell, rowIndex, colIndex)"
          ></div>
        </div>
      </div>
      
      <div class="game-over-overlay" v-if="gameOver">
        <div class="game-over-message">
          <h2>Game Over!</h2>
        </div>
      </div>
    </div>
    
    <div class="next-piece-container">
      <h3>Next Piece</h3>
      <div class="next-piece-grid">
        <div 
          v-for="(row, rowIndex) in 4" 
          :key="'next-row-' + rowIndex" 
          class="next-piece-row"
        >
          <div 
            v-for="(cell, colIndex) in 4" 
            :key="'next-cell-' + rowIndex + '-' + colIndex" 
            class="next-piece-cell"
            :style="getNextPieceCellClass(rowIndex, colIndex)"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'GameBoard',
  props: {
    gameBoard: {
      type: Array,
      required: true
    },
    currentPiece: {
      type: Object,
      default: null
    },
    nextPiece: {
      type: Object,
      default: null
    },
    gameOver: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    // 获取当前单元格的CSS类
    getCellClass(cell, rowIndex, colIndex) {
      let cellClass = ''
      
      // 检查当前单元格是否是当前方块的一部分
      if (this.currentPiece && !this.gameOver) {
        const { shape, position, color } = this.currentPiece
        
        // 遍历方块的形状，检查是否覆盖当前单元格
        for (let i = 0; i < shape.length; i++) {
          for (let j = 0; j < shape[i].length; j++) {
            if (shape[i][j] === 1) {
              const pieceRow = position.row + i
              const pieceCol = position.col + j
              
              if (pieceRow === rowIndex && pieceCol === colIndex) {
                return { backgroundColor: color }
              }
            }
          }
        }
      }
      
      // 如果是已固定的方块
      if (cell !== 0) {
        // 根据cell的值（方块类型）获取对应的颜色
        const colors = {
          0: '',
          1: '#00FF00',  // S
          2: '#FF0000',  // Z
          3: '#FFA500',  // L
          4: '#0000FF',  // J
          5: '#00FFFF',  // I
          6: '#FFFF00',  // O
          7: '#800080'   // T
        }
        return { backgroundColor: colors[cell] || '' }
      }
      
      return {}
    },
    
    // 获取下一个方块单元格的CSS类
    getNextPieceCellClass(rowIndex, colIndex) {
      if (!this.nextPiece) return ''
      
      const { shape, color } = this.nextPiece
      
      // 遍历方块的形状，检查是否覆盖当前单元格
      for (let i = 0; i < shape.length; i++) {
        for (let j = 0; j < shape[i].length; j++) {
          if (shape[i][j] === 1) {
            // 居中显示下一个方块
            const offsetRow = Math.floor((4 - shape.length) / 2)
            const offsetCol = Math.floor((4 - shape[0].length) / 2)
            
            const displayRow = offsetRow + i
            const displayCol = offsetCol + j
            
            if (displayRow === rowIndex && displayCol === colIndex) {
              return { backgroundColor: color }
            }
          }
        }
      }
      
      return {}
    }
  }
}
</script>

<style scoped>
.game-board-container {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.main-board {
  position: relative;
  border: 2px solid #333;
  border-radius: 5px;
  background-color: #111;
}

.board-grid {
  display: flex;
  flex-direction: column;
}

.board-row {
  display: flex;
}

.board-cell {
  width: 30px;
  height: 30px;
  border: 1px solid #333;
  box-sizing: border-box;
}

/* 方块颜色 */
.cell-I { background-color: #00FFFF; }
.cell-J { background-color: #0000FF; }
.cell-L { background-color: #FFA500; }
.cell-O { background-color: #FFFF00; }
.cell-S { background-color: #00FF00; }
.cell-T { background-color: #800080; }
.cell-Z { background-color: #FF0000; }

.next-piece-container {
  background-color: #222;
  padding: 15px;
  border-radius: 5px;
  color: white;
  min-width: 150px;
}

.next-piece-container h3 {
  margin-top: 0;
  margin-bottom: 10px;
  text-align: center;
  color: #fff;
}

.next-piece-grid {
  display: flex;
  flex-direction: column;
  background-color: #111;
  border: 1px solid #333;
  border-radius: 3px;
}

.next-piece-row {
  display: flex;
}

.next-piece-cell {
  width: 25px;
  height: 25px;
  border: 1px solid #333;
  box-sizing: border-box;
}

.game-over-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 3px;
}

.game-over-message {
  text-align: center;
  color: white;
  animation: pulse 1s infinite alternate;
}

.game-over-message h2 {
  font-size: 36px;
  margin: 0;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.8);
}

@keyframes pulse {
  from { opacity: 0.7; transform: scale(0.9); }
  to { opacity: 1; transform: scale(1); }
}
</style>
