import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MessageList } from './MessageList';

// jsdomではscrollIntoViewが実装されていないためモック化
window.HTMLElement.prototype.scrollIntoView = vi.fn();

describe('MessageList Component', () => {
  it('正常系: メッセージ一覧がレンダリングされ、自分と他人のメッセージのスタイルが区別されること', () => {
    // Given: 異なる送信者のメッセージが存在する
    const mockMessages = [
      { id: 1, sender: 'Alice', content: 'Hello Bob!' },
      { id: 2, sender: 'Bob', content: 'Hi Alice!' }
    ];

    // When: 現在のユーザーを"Alice"としてコンポーネントをレンダリングする
    render(<MessageList messages={mockMessages} currentUser="Alice" />);

    // Then: 自分のメッセージには'mine'、他人のメッセージには'other'のクラスが付与される
    const aliceMsg = screen.getByText('Hello Bob!').parentElement;
    const bobMsg = screen.getByText('Hi Alice!').parentElement;

    expect(aliceMsg?.className).toContain('mine');
    expect(bobMsg?.className).toContain('other');
  });
});
