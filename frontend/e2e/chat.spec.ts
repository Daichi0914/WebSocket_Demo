import { test, expect } from '@playwright/test';

/**
 * E2Eテスト: リアルタイムチャットの送受信
 *
 * 前提条件:
 *   - バックエンド（Go + MySQL + Redis）が起動していること（`make up`）
 *   - フロントエンドはPlaywrightが自動起動します（playwright.config.ts の webServer 設定）
 */
test.describe('リアルタイムチャット E2E', () => {
  test('正常系: User Aが送信したメッセージがUser Bの画面にリアルタイムで表示されること', async ({ browser }) => {
    // === Given: 2つのブラウザコンテキスト（= 2人のユーザー）を作成する ===
    const contextA = await browser.newContext();
    const contextB = await browser.newContext();
    const pageA = await contextA.newPage();
    const pageB = await contextB.newPage();

    // === When: 両者がチャットルームに参加する ===

    // User A: 名前を入力してJoin
    await pageA.goto('/');
    await pageA.getByPlaceholder('Your Name...').fill('Alice');
    await pageA.getByRole('button', { name: 'Join' }).click();

    // User A: WebSocket接続を待機（🟢が表示されるまで）
    await expect(pageA.getByRole('heading', { level: 1 })).toContainText('🟢', { timeout: 10000 });

    // User B: 名前を入力してJoin
    await pageB.goto('/');
    await pageB.getByPlaceholder('Your Name...').fill('Bob');
    await pageB.getByRole('button', { name: 'Join' }).click();

    // User B: WebSocket接続を待機（🟢が表示されるまで）
    await expect(pageB.getByRole('heading', { level: 1 })).toContainText('🟢', { timeout: 10000 });

    // === When: User Aがメッセージを送信する ===
    await pageA.getByPlaceholder('Type a message...').fill('こんにちは！');
    await pageA.getByRole('button', { name: 'Send' }).click();

    // === Then: User A自身の画面にメッセージが表示されること ===
    await expect(pageA.locator('.message-bubble').last()).toHaveText('こんにちは！', { timeout: 5000 });

    // === Then: User Bの画面にもリアルタイムでメッセージが表示されること ===
    await expect(pageB.locator('.message-bubble').last()).toHaveText('こんにちは！', { timeout: 5000 });

    // === Then: User Bの画面で送信者が「Alice」であることが確認できること ===
    await expect(pageB.locator('.message-sender').last()).toHaveText('Alice');

    // クリーンアップ
    await contextA.close();
    await contextB.close();
  });
});
