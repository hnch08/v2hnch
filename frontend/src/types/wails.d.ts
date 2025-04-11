declare module 'wailsjs/go/models' {
  export interface App {
    Login: (username: string, password: string) => Promise<boolean>;
    HideWindow: () => Promise<void>;
    ShowWindow: () => Promise<void>;
    Quit: () => Promise<void>;
  }
}

declare global {
  interface Window {
    go: {
      main: {
        App: App;
      }
    }
  }
} 