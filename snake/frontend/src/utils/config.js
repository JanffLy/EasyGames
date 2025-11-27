// 前端配置文件
export default {
  // API基础URL
  apiBaseUrl: '/api',
  
  // 游戏配置
  game: {
    // 游戏区域大小
    width: 15,
    height: 15,
    
    // 状态轮询间隔（毫秒）
    pollInterval: 100,
    
    // 方向键映射
    keyMapping: {
      'KeyW': 'up',
      'KeyS': 'down',
      'KeyA': 'left',
      'KeyD': 'right'
    }
  }
};
