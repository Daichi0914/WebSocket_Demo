import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ChatHeader } from './ChatHeader';

describe('ChatHeader Component', () => {
  it('正常系: 接続済みの場合は🟢が表示されること', () => {
    // Given: 接続状態がtrue
    // When: コンポーネントをレンダリング
    render(<ChatHeader isConnected={true} />);

    // Then: 🟢が表示される
    const heading = screen.getByRole('heading', { level: 1 });
    expect(heading.textContent).toContain('🟢');
  });

  it('異常系: 切断されている場合は🔴が表示されること', () => {
    // Given: 接続状態がfalse
    // When: コンポーネントをレンダリング
    render(<ChatHeader isConnected={false} />);

    // Then: 🔴が表示される
    const heading = screen.getByRole('heading', { level: 1 });
    expect(heading.textContent).toContain('🔴');
  });
});
