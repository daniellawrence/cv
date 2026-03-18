import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import puppeteer, { Browser, Page } from 'puppeteer';

// Base URL for tests - configurable via TEST_BASE_URL env var or CLI --define
// Examples:
//   npm run test:dev          -> http://localhost:5173 (Vite dev server)
//   npm run test:prod         -> http://localhost (production build)
//   TEST_BASE_URL=http://example.com npx vitest  -> custom URL
const baseUrl = process.env.TEST_BASE_URL || 'http://localhost:5173';

describe('Basic Page Load Test', () => {
  let browser: Browser;

  beforeAll(async () => {
    console.log(`[TEST] Running tests against: ${baseUrl}`);
    if (!baseUrl.startsWith('http://') && !baseUrl.startsWith('https://')) {
      throw new Error(`Invalid BASE_URL: ${baseUrl}. Must start with http:// or https://`);
    }
    browser = await puppeteer.launch({ headless: true });
  });

  afterAll(async () => {
    if (browser) {
      await browser.close();
    }
  });

  it('should load the page without JavaScript errors', async () => {
    const page = await browser.newPage();
    
    // Set up error collection on the page
    await page.evaluate(() => {
      (window as any).errorLog = [];
      
      const originalOnError = window.onerror || (() => false);
      window.onerror = (msg, source, line, column, error) => {
        const errorMsg = msg || String(error?.message);
        if (!errorMsg.includes('Failed to load resource')) {
          (window as any).errorLog.push(errorMsg);
        }
        return originalOnError(msg, source, line, column, error);
      };
    });

    await page.goto(baseUrl, { waitUntil: 'networkidle0' });

    // Check for errors collected on the page
    const errorLogs = await page.evaluate(() => {
      return (window as any).errorLog || [];
    });
    
    expect(errorLogs).toHaveLength(0, `JavaScript errors detected: ${errorLogs.join(', ')}`);

    await page.close();
  });

  it('should have valid HTML structure', async () => {
    const page = await browser.newPage();
    const html = await page.content();
    
    expect(html).toContain('<html');
    expect(html).toContain('</html>');
    
    await page.close();
  });

  it('should load React app container in DOM', async () => {
    const page = await browser.newPage();
    await page.goto(baseUrl, { waitUntil: 'networkidle0' });
    
    const rootExists = await page.$('#root');
    expect(rootExists).toBeTruthy();
    
    await page.close();
  });

  it('should have a valid title', async () => {
    const page = await browser.newPage();
    await page.goto(baseUrl, { waitUntil: 'networkidle0' });
    
    const title = await page.title();
    expect(title).toContain('Daniel Lawrence');
    
    await page.close();
  });
});
