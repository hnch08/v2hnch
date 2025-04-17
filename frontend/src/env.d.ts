/// <reference types="vite/client" />

declare interface Window {
  go: {
    main: {
      App: {
        Login: (username: string, password: string) => Promise<boolean>;
        HideWindow: () => Promise<void>;
        ShowWindow: () => Promise<void>;
        Quit: () => Promise<void>;
        StartProxy: () => Promise<boolean>;
        StopProxy: () => Promise<boolean>;
        SetAddress: (address: string) => Promise<boolean>;
        GetConfig: () => Promise<config.Config>;
        GetStatus: () => Promise<number>;
        GetLoginStatus: () => Promise<boolean>;
        CheckURL: () => Promise<boolean>;
      }
    }
  },
  runtime: {
    EventsOn: (eventName: string, callback: (...data: any) => void) => void;
  }
} 