"""Playwright test fixtures for the PODA frontend."""

import pytest
from playwright.sync_api import sync_playwright, Browser, BrowserContext, Page

BASE_URL = "http://localhost:5173"


@pytest.fixture(scope="session")
def browser_instance():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        yield browser
        browser.close()


@pytest.fixture
def context(browser_instance: Browser):
    ctx = browser_instance.new_context(base_url=BASE_URL)
    yield ctx
    ctx.close()


@pytest.fixture
def page(context: BrowserContext):
    p = context.new_page()
    yield p
    p.close()
