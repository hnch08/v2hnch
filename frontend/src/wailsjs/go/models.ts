// 定义 Go 绑定的类型
declare global {
  interface Window {
    go: {
      main: {
        App: {
          Login: (username: string, password: string) => Promise<boolean>;
          HideWindow: () => Promise<void>;
          ShowWindow: () => Promise<void>;
          Quit: () => Promise<void>;
        }
      }
    }
  }
}

export {}; 