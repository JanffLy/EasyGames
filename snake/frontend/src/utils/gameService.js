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

  // 更新游戏方向
  async updateDirection(gameId, direction) {
    const url = `${this.apiBaseUrl}/game/${gameId}/direction`;
    console.log(`正在更新游戏方向: ${url}, 方向: ${direction}`);
    try {
      // 确保direction是小写字符串
      const normalizedDirection = direction.toLowerCase();
      console.log("标准化后的方向:", normalizedDirection);

      // 创建符合后端GameRequest结构体的请求体（使用小写direction）
      const requestBody = { direction: normalizedDirection };
      console.log("请求体对象(使用小写direction):", requestBody);

      // 直接使用JSON.stringify，不额外赋值
      console.log("JSON请求体:", JSON.stringify(requestBody));

      const requestOptions = {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(requestBody),
      };
      console.log("发送请求选项:", requestOptions);

      const response = await this.fetchWithRetry(url, requestOptions);
      console.log(
        `更新方向API响应状态: ${response.status}, ${response.statusText}`
      );

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        console.error(`更新方向API错误数据:`, errorData);
        throw new Error(`更新方向失败: ${response.statusText}`);
      }

      const data = await response.json();
      console.log(`成功更新游戏方向:`, data);
      return data;
    } catch (error) {
      console.error(`更新方向时发生异常:`, error);
      throw error;
    }
  }

  // 保存游戏记录
  async saveScore(gameId, playerName, score) {
    const url = `${this.apiBaseUrl}/game/${gameId}/record`;
    console.log(
      `正在保存游戏记录: ${url}, 玩家: ${playerName}, 得分: ${score}`
    );
    try {
      const response = await this.fetchWithRetry(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ playerName, score }),
      });
      console.log(
        `保存记录API响应状态: ${response.status}, ${response.statusText}`
      );

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        console.error(`保存记录API错误数据:`, errorData);
        throw new Error(`保存得分失败: ${response.statusText}`);
      }

      const data = await response.json();
      console.log(`成功保存游戏记录:`, data);
      return data;
    } catch (error) {
      console.error(`保存游戏记录时发生异常:`, error);
      throw error;
    }
  }

  // 获取排行榜
  async getLeaderboard() {
    const url = `${this.apiBaseUrl}/leaderboard`;
    console.log(`正在获取排行榜: ${url}`);
    try {
      const response = await this.fetchWithRetry(url);
      console.log(
        `获取排行榜API响应状态: ${response.status}, ${response.statusText}`
      );

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        console.error(`获取排行榜API错误数据:`, errorData);
        throw new Error(`获取排行榜失败: ${response.statusText}`);
      }

      const data = await response.json();
      console.log(`成功获取排行榜数据:`, data);
      return data;
    } catch (error) {
      console.error(`获取排行榜时发生异常:`, error);
      throw error;
    }
  }
}

// 创建并导出游戏服务实例
export const gameService = new GameService();
