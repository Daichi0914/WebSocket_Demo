import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useChat } from './useChat';

// WebSocketのモック
class MockWebSocket {
  onopen: (() => void) | null = null;
  onclose: (() => void) | null = null;
  onmessage: ((ev: any) => void) | null = null;
  close = vi.fn();
  send = vi.fn();

  constructor(public url: string) {
    // コンストラクタでインスタンスを保持し、テストから操作できるようにする
    MockWebSocket.lastInstance = this;
  }
  static lastInstance: MockWebSocket;
}

describe('useChat Hook', () => {
  beforeEach(() => {
    vi.stubGlobal('WebSocket', MockWebSocket);
    vi.stubGlobal('fetch', vi.fn());
    vi.spyOn(console, 'error').mockImplementation(() => {});
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  it('正常系: 初期化時に履歴を取得し、WebSocketを接続すること', async () => {
    const mockMessages = [{ sender: 'System', content: 'Welcome' }];
    (fetch as any).mockResolvedValue({
      json: () => Promise.resolve(mockMessages),
    });

    const { result } = renderHook(() => useChat(true, 'Alice'));

    // 履歴の取得を確認
    await waitFor(() => {
      expect(result.current.messages).toEqual(mockMessages);
    });

    // WebSocketの接続を確認
    const wsInstance = MockWebSocket.lastInstance;
    expect(wsInstance).toBeDefined();
    
    // openイベントをシミュレート
    wsInstance.onopen?.();
    await waitFor(() => expect(result.current.isConnected).toBe(true));
  });

  it('異常系: fetchが失敗してもクラッシュせず、空のメッセージリストを保持すること', async () => {
    (fetch as any).mockRejectedValue(new Error('Network Error'));

    const { result } = renderHook(() => useChat(true, 'Alice'));

    // エラーがコンソールに出力され、メッセージは空のままであること
    await waitFor(() => {
      expect(console.error).toHaveBeenCalled();
      expect(result.current.messages).toEqual([]);
    });
  });

  it('異常系: WebSocketが切断されたら isConnected が false になること', async () => {
    (fetch as any).mockResolvedValue({ json: () => Promise.resolve([]) });

    const { result } = renderHook(() => useChat(true, 'Alice'));
    const wsInstance = MockWebSocket.lastInstance;

    // 接続
    wsInstance.onopen?.();
    await waitFor(() => expect(result.current.isConnected).toBe(true));

    // 切断シミュレート
    wsInstance.onclose?.();
    await waitFor(() => expect(result.current.isConnected).toBe(false));
  });

  it('準正常系: WebSocketから不正なJSONが届いてもエラーにならず、メッセージが追加されないこと', async () => {
    (fetch as any).mockResolvedValue({ json: () => Promise.resolve([]) });

    const { result } = renderHook(() => useChat(true, 'Alice'));
    const wsInstance = MockWebSocket.lastInstance;

    const initialCount = result.current.messages.length;

    // 不正なJSONを送りつける
    wsInstance.onmessage?.({ data: 'invalid json' });

    // メッセージが増えていないことを確認
    expect(result.current.messages.length).toBe(initialCount);
  });

  it('正常系: sendMessageを呼ぶとWebSocket経由でメッセージが送信されること', async () => {
    (fetch as any).mockResolvedValue({ json: () => Promise.resolve([]) });

    const { result } = renderHook(() => useChat(true, 'Alice'));
    const wsInstance = MockWebSocket.lastInstance;

    result.current.sendMessage('Hello');

    expect(wsInstance.send).toHaveBeenCalledWith(
      expect.stringContaining('"content":"Hello"')
    );
    expect(wsInstance.send).toHaveBeenCalledWith(
      expect.stringContaining('"sender":"Alice"')
    );
  });
});
