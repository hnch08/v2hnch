/// <reference types="vite/client" />

declare interface Window {
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