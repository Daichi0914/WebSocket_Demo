import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { NamePrompt } from './NamePrompt';

describe('NamePrompt Component', () => {
  it('異常系: 名前が空の場合はonJoinが呼ばれないこと', () => {
    // Given: 初期状態でコンポーネントがレンダリングされ、名前が空
    const mockSetSenderName = vi.fn();
    const mockOnJoin = vi.fn();
    
    render(
      <NamePrompt 
        senderName="" 
        setSenderName={mockSetSenderName} 
        onJoin={mockOnJoin} 
      />
    );

    // When: そのままJoinボタンを押す
    const button = screen.getByRole('button', { name: /join/i });
    fireEvent.click(button);

    // Then: onJoinイベントが発火しないこと
    expect(mockOnJoin).not.toHaveBeenCalled();
  });

  it('正常系: 有効な名前が入力された場合はonJoinが呼ばれること', () => {
    // Given: 名前が入力されている状態
    const mockSetSenderName = vi.fn();
    const mockOnJoin = vi.fn();
    
    render(
      <NamePrompt 
        senderName="Alice" 
        setSenderName={mockSetSenderName} 
        onJoin={mockOnJoin} 
      />
    );

    // When: Joinボタンを押す
    const button = screen.getByRole('button', { name: /join/i });
    fireEvent.click(button);

    // Then: onJoinイベントが発火すること
    expect(mockOnJoin).toHaveBeenCalledTimes(1);
  });

  it('正常系: 入力フィールドに文字を入力した際にsetSenderNameが呼ばれること', () => {
    // Given
    const mockSetSenderName = vi.fn();
    const mockOnJoin = vi.fn();
    
    render(
      <NamePrompt 
        senderName="" 
        setSenderName={mockSetSenderName} 
        onJoin={mockOnJoin} 
      />
    );

    // When: 入力フィールドに文字を入力する
    const input = screen.getByPlaceholderText(/your name/i);
    fireEvent.change(input, { target: { value: 'Bob' } });

    // Then: setSenderNameが入力した値で呼ばれること
    expect(mockSetSenderName).toHaveBeenCalledWith('Bob');
  });
});
