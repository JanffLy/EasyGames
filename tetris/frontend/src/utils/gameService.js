// 游戏服务，处理与后端的通信
export class GameService {
  constructor() {
    this.apiBaseUrl = "/api";
    // 请求配置
    this.maxRetries = 2;
    this.timeout = 3000;
  }

  // 带重试和超时的fetch请求
  async fetchWithRetry(url, options = {}, retries = this.maxRetries) {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(url, {
        ...options,
        signal: controller.signal,
      });
      clearTimeout(timeoutId);
      return response;
    } catch (error) {
      clearTimeout(timeoutId);

      // 如果是AbortError（超时）或网络错误，并且还有重试次数
      if (
        (error.name === "AbortError" || error.message === "Failed to fetch") &&
        retries > 0
      ) {
        console.log(`请求失败，${retries}次重试中...`);
        // 指数退避策略，等待一段时间后重试
        await new Promise((resolve) =>
          setTimeout(resolve, (this.maxRetries - retries + 1) * 500)
        );
        return this.fetchWithRetry(url, options, retries - 1);
      }
      throw error;
    }
  }

  // 创建新游戏
  async createGame() {
    const url = `${this.apiBaseUrl}/game`;
    console.log(`正在创建新游戏: ${url}`);
    try {
      const response = await this.fetchWithRetry(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });
      console.log(
        `创建游戏API响应状态: ${response.status}, ${response.statusText}`
      );

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        console.error(`创建游戏API错误数据:`, errorData);
        throw new Error(`创建游戏失败: ${response.statusText}`);
      }

      const data = await response.json();
      console.log(`成功创建游戏:`, data);
      return data;
    } catch (error) {
      console.error(`创建游戏时发生异常:`, error);
      throw error;
    }
  }

  // 获取游戏状态
  async getGameState(gameId) {
    const url = `${this.apiBaseUrl}/game/${gameId}`;
    console.log(`正在请求游戏状态: ${url}`);
    try {
      const response = await this.fetchWithRetry(url);
      console.log(`API响应状态: ${response.status}, ${response.statusText}`);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        console.error(`API错误数据:`, errorData);
        throw new Error(`获取游戏状态失败: ${response.statusText}`);
      }

      const data = await response.json();
      console.log(`成功获取游戏状态数据:`, data);
      return data;
    } catch (error) {
      console.error(`请求游戏状态时发生异常:`, error);
      throw error;
    }
  }

  // 发送游戏操作
  async sendAction(gameId, action) {
    const url = `${this.apiBaseUrl}/game/${gameId}/action`;
    console.log(`正在发送游戏操作: ${url}, 操作: ${action}`);
    try {
      // 确保action是小写字符串
      const normalizedAction = action.toLowerCase();
      
      const response = await this.fetchWithRetry(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ action: normalizedAction }),
      });
      console.log(
        `更新游戏API响应状态: ${response.status}, ${response.statusText}`
      );

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        console.error(`API错误数据:`, errorData);
        throw new Error(`更新游戏失败: ${response.statusText}`);
      }

      const data = await response.json();
      console.log(`成功更新游戏:`, data);
      return data;
    } catch (error) {
      console.error(`更新游戏时发生异常:`, error);
      throw error;
    }
  }

  // 结束游戏
  async endGame(gameId) {
    const url = `${this.apiBaseUrl}/game/${gameId}`;
    console.log(`正在结束游戏: ${url}`);
    try {
      const response = await this.fetchWithRetry(url, {
        method: "DELETE",
      });
      console.log(`结束游戏API响应状态: ${response.status}, ${response.statusText}`);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        console.error(`API错误数据:`, errorData);
        throw new Error(`结束游戏失败: ${response.statusText}`);
      }

      console.log(`成功结束游戏`);
      return true;
    } catch (error) {
      console.error(`结束游戏时发生异常:`, error);
      throw error;
    }
  }
}

// 导出单例实例
const gameService = new GameService();

// 导出便捷方法
export const startNewGame = () => gameService.createGame();
export const getGameState = (gameId) => gameService.getGameState(gameId);
export const sendGameAction = (gameId, action) => gameService.sendAction(gameId, action);
export const endGame = (gameId) => gameService.endGame(gameId);
