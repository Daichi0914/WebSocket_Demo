import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MessageInput } from './MessageInput';

describe('MessageInput Component', () => {
  it('異常系: 入力が空の場合はonSendが呼ばれないこと', () => {
    // Given: 入力が空の状態でレンダリング
    const mockOnSend = vi.fn();
    render(<MessageInput disabled={false} onSend={mockOnSend} />);

    // When: Sendボタンをクリック
    const button = screen.getByRole('button', { name: /send/i });
    fireEvent.click(button);

    // Then: onSendは発火しない
    expect(mockOnSend).not.toHaveBeenCalled();
  });

  it('正常系: 有効なメッセージが入力された場合、onSendが呼ばれて入力欄がクリアされること', () => {
    // Given: テキストが入力されている状態
    const mockOnSend = vi.fn();
    render(<MessageInput disabled={false} onSend={mockOnSend} />);

    const input = screen.getByPlaceholderText(/type a message/i) as HTMLInputElement;
    fireEvent.change(input, { target: { value: 'Hello World' } });
    
    // When: Sendボタンをクリック
    const button = screen.getByRole('button', { name: /send/i });
    fireEvent.click(button);

    // Then: onSendが正しい値で呼ばれ、入力欄がクリアされること
    expect(mockOnSend).toHaveBeenCalledWith('Hello World');
    expect(input.value).toBe('');
  });

  it('異常系: disabledがtrueの場合、入力欄とボタンが操作不可になること', () => {
    // Given: disabled=true でレンダリング
    const mockOnSend = vi.fn();
    
    // When
    render(<MessageInput disabled={true} onSend={mockOnSend} />);

    // Then: 入力欄とボタンの両方が操作不可になっていること
    const input = screen.getByPlaceholderText(/type a message/i) as HTMLInputElement;
    const button = screen.getByRole('button', { name: /send/i }) as HTMLButtonElement;
    
    expect(input.disabled).toBe(true);
    expect(button.disabled).toBe(true);
  });
});
