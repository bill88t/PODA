"""
Playwright end-to-end tests for the PODA frontend.

These tests run AFTER `make runtests_full` has already populated the database,
so the following accounts are pre-existing:

  Bill (client):
    email:    bill88t@feline.gr
    password: newpassword456        (changed by the makefile tests)
    name:     William Sideris       (changed by the makefile tests)
    phone:    +30 210 1234567
    appointments: 2 remaining (shaving + haircut)

  Tony (barber):
    email:    tony@barber.gr
    password: barberpassword123
"""

import pytest
import time
from playwright.sync_api import Page, expect

BASE_URL = "http://localhost:5173"

# ─── helpers ────────────────────────────────────────────────────────────────


def navigate(page: Page, path: str) -> None:
    """Navigate to a page and wait for the app to settle."""
    page.goto(BASE_URL + path)
    page.wait_for_load_state("networkidle")


def nav_click(page: Page, text: str) -> None:
    """Click a nav-bar link by visible text."""
    page.locator(".nav").get_by_text(text).click()
    page.wait_for_load_state("networkidle")


def login(page: Page, email: str, password: str) -> None:
    """Log in via the Login page."""
    navigate(page, "/login")
    page.get_by_text("Email").locator("+ *").fill(email)
    page.get_by_text("Password").locator("+ *").fill(password)
    page.get_by_role("button", name="Login").click()
    page.wait_for_load_state("networkidle")


def login_via_form(page: Page, email: str, password: str) -> None:
    """Fill login inputs directly (works when already on /login)."""
    page.locator('input[type="text"]').fill(email)
    page.locator('input[type="password"]').fill(password)
    page.get_by_role("button", name="Login").click()
    page.wait_for_load_state("networkidle")


# ─── 1. Registration ─────────────────────────────────────────────────────────


class TestRegistration:
    def test_register_new_user(self, page: Page):
        """A brand-new user can complete the sign-up form."""
        navigate(page, "/sign_up")
        page.locator('input[type="text"]').nth(0).fill("Play")  # First Name
        page.locator('input[type="text"]').nth(1).fill("Wright")  # Last Name
        page.locator('input[type="text"]').nth(2).fill(
            "playwright@test.example"
        )  # Email
        page.locator('input[type="password"]').nth(0).fill("testpass123")
        page.locator('input[type="password"]').nth(1).fill("testpass123")
        page.locator('input[type="text"]').nth(3).fill("6901234567")  # Phone
        page.locator('input[type="text"]').nth(4).fill("123 Test St")  # Address
        page.locator('input[type="date"]').fill("1995-06-15")  # Birthday
        page.get_by_role("button", name="Sign Up").click()
        page.wait_for_load_state("networkidle")
        # On success the SPA navigates to "/"
        expect(page).to_have_url(BASE_URL + "/")

    def test_duplicate_registration_shows_error(self, page: Page):
        """Attempting to re-register the same e-mail shows an error."""
        navigate(page, "/sign_up")
        page.locator('input[type="text"]').nth(0).fill("Play")
        page.locator('input[type="text"]').nth(1).fill("Wright")
        page.locator('input[type="text"]').nth(2).fill("playwright@test.example")
        page.locator('input[type="password"]').nth(0).fill("testpass123")
        page.locator('input[type="password"]').nth(1).fill("testpass123")
        page.locator('input[type="text"]').nth(3).fill("6901234567")
        page.locator('input[type="text"]').nth(4).fill("123 Test St")
        page.locator('input[type="date"]').fill("1995-06-15")
        page.get_by_role("button", name="Sign Up").click()
        page.wait_for_load_state("networkidle")
        expect(page.locator(".error")).to_be_visible()

    def test_register_password_mismatch_shows_error(self, page: Page):
        """Mismatched passwords surface a client-side validation error."""
        navigate(page, "/sign_up")
        page.locator('input[type="text"]').nth(0).fill("Test")
        page.locator('input[type="text"]').nth(1).fill("User")
        page.locator('input[type="text"]').nth(2).fill("mismatch@test.example")
        page.locator('input[type="password"]').nth(0).fill("testpass123")
        page.locator('input[type="password"]').nth(1).fill("wrongpass123")
        page.locator('input[type="text"]').nth(3).fill("6901234567")
        page.locator('input[type="text"]').nth(4).fill("123 Test St")
        page.locator('input[type="date"]').fill("1995-06-15")
        page.get_by_role("button", name="Sign Up").click()
        expect(page.locator(".error")).to_have_text("Passwords don't match")


# ─── 2. Login / Logout ───────────────────────────────────────────────────────


class TestAuth:
    def test_login_playwright_user(self, page: Page):
        """The newly registered playwright user can log in."""
        navigate(page, "/login")
        login_via_form(page, "playwright@test.example", "testpass123")
        # After login the navbar should show "Log out" instead of "Login"
        expect(page.locator(".nav")).to_contain_text("Log out")

    def test_invalid_login_shows_error(self, page: Page):
        """Wrong credentials show an error message."""
        navigate(page, "/login")
        login_via_form(page, "playwright@test.example", "wrongpassword")
        expect(page.locator(".error")).to_be_visible()

    def test_logout(self, page: Page):
        """Logging out returns the user to an unauthenticated state."""
        navigate(page, "/login")
        login_via_form(page, "playwright@test.example", "testpass123")
        expect(page.locator(".nav")).to_contain_text("Log out")
        page.locator(".nav").get_by_text("Log out").click()
        page.wait_for_load_state("networkidle")
        expect(page.locator(".nav")).to_contain_text("Login")

    def test_login_bill_existing_user(self, page: Page):
        """Bill's account (created by makefile tests) still works with the changed password."""
        navigate(page, "/login")
        login_via_form(page, "bill88t@feline.gr", "newpassword456")
        expect(page.locator(".nav")).to_contain_text("Log out")

    def test_login_tony_barber(self, page: Page):
        """Tony's barber account (created by makefile tests) can log in."""
        navigate(page, "/login")
        login_via_form(page, "tony@barber.gr", "barberpassword123")
        expect(page.locator(".nav")).to_contain_text("Log out")


# ─── 3. Appointments (playwright test user) ──────────────────────────────────


class TestAppointments:
    @pytest.fixture(autouse=True)
    def logged_in(self, page: Page):
        """Log in as the playwright test user before each test in this class."""
        navigate(page, "/login")
        login_via_form(page, "playwright@test.example", "testpass123")
        page.wait_for_load_state("networkidle")
        expect(page.locator(".nav")).to_contain_text("Log out")

    def test_events_page_accessible(self, page: Page):
        """Authenticated user can navigate to the Events page."""
        nav_click(page, "Events")
        expect(page).to_have_url(BASE_URL + "/events")
        expect(page.locator("h1")).to_contain_text("Events")

    def test_create_appointment(self, page: Page):
        """User can create a new appointment via the Events UI."""
        nav_click(page, "Events")
        page.get_by_role("button", name="Create").click()
        page.wait_for_selector("#appointment-kind")
        page.locator("#appointment-kind").select_option("haircut")
        page.locator('input[type="datetime-local"]').fill("2027-06-01T10:00")
        page.get_by_role("button", name="New Appointment").click()
        page.wait_for_load_state("networkidle")

    def test_appointment_list_not_empty_after_create(self, page: Page):
        """After creating an appointment, the events list has at least one entry."""
        nav_click(page, "Events")
        page.wait_for_load_state("networkidle")
        delete_links = page.locator("a", has_text="Delete")
        count_before = delete_links.count()
        assert (
            count_before == 1
        ), f"Expected exactly one appointment, got {count_before}"

    # def test_delete_appointment(self, page: Page):
    #     """User can delete their own appointment."""
    #     nav_click(page, "Events")
    #     page.wait_for_load_state("networkidle")
    #     delete_links = page.locator("a", has_text="Delete")
    #     count_before = delete_links.count()
    #     assert count_before > 0, "Expected at least one appointment to delete"
    #     delete_links.first.click()
    #     page.wait_for_load_state("networkidle")
    #     time.sleep(0.3)
    #     count_after = page.locator("a", has_text="Delete").count()
    #     assert count_after == count_before - 1


# ─── 4. Profile changes ──────────────────────────────────────────────────────


# class TestProfile:
#     @pytest.fixture(autouse=True)
#     def logged_in(self, page: Page):
#         navigate(page, "/login")
#         login_via_form(page, "playwright@test.example", "testpass123")
#         page.wait_for_load_state("networkidle")

#     def _go_profile(self, page: Page) -> None:
#         nav_click(page, "Profile")
#         expect(page).to_have_url(BASE_URL + "/profile")

#     def test_profile_page_shows_name(self, page: Page):
#         """Profile page displays the user's name."""
#         self._go_profile(page)
#         expect(page.locator(".form")).to_contain_text("Play")
#         expect(page.locator(".form")).to_contain_text("Wright")

#     def test_change_personal_info(self, page: Page):
#         """User can update their first and last name."""
#         self._go_profile(page)
#         page.get_by_role("button", name="Edit Personal Information").click()
#         fname_input = page.locator('input[type="text"]').nth(0)
#         fname_input.fill("PlayUpdated")
#         page.get_by_role("button", name="Confirm Information Changes").click()
#         page.wait_for_load_state("networkidle")
#         # Confirm the form collapsed (edit mode exited) without an error
#         expect(page.locator(".error")).to_have_count(0)
#         expect(page.locator(".form")).to_contain_text("PlayUpdated")

#     def test_change_contact_info(self, page: Page):
#         """User can update their email."""
#         self._go_profile(page)
#         page.get_by_role("button", name="Edit Contact Information").click()
#         email_input = page.locator('input[type="text"]').first
#         email_input.fill("playwright@test.example")  # keep same email
#         page.get_by_role("button", name="Confirm Contact changes").click()
#         page.wait_for_load_state("networkidle")
#         expect(page.locator(".error")).to_have_count(0)

#     def test_change_password(self, page: Page):
#         """User can change their password and re-login with the new one."""
#         self._go_profile(page)
#         page.get_by_role("button", name="Change Password").click()
#         page.locator('input[type="password"]').nth(0).fill("newpass456")
#         page.locator('input[type="password"]').nth(1).fill("newpass456")
#         page.get_by_role("button", name="Confirm Password Change").click()
#         page.wait_for_load_state("networkidle")
#         expect(page.locator(".error")).to_have_count(0)

#         # Log out and back in with the new password
#         page.locator(".nav").get_by_text("Log out").click()
#         page.wait_for_load_state("networkidle")
#         navigate(page, "/login")
#         login_via_form(page, "playwright@test.example", "newpass456")
#         expect(page.locator(".nav")).to_contain_text("Log out")


# ─── 5. Existing data validation (Bill's account) ────────────────────────────


# class TestExistingData:
#     """Validate that the data seeded by `make runtests_full` is visible in the UI."""

#     def test_bill_profile_shows_updated_name(self, page: Page):
#         """Bill's name was changed to 'William Sideris' by the makefile tests."""
#         navigate(page, "/login")
#         login_via_form(page, "bill88t@feline.gr", "newpassword456")
#         nav_click(page, "Profile")
#         expect(page.locator(".form")).to_contain_text("William")
#         expect(page.locator(".form")).to_contain_text("Sideris")

#     def test_bill_appointments_visible(self, page: Page):
#         """Bill has at least one remaining appointment after the makefile tests."""
#         navigate(page, "/login")
#         login_via_form(page, "bill88t@feline.gr", "newpassword456")
#         nav_click(page, "Events")
#         # At least one appointment row (with a Delete link) should be visible
#         expect(page.locator("a", has_text="Delete").first).to_be_visible()

#     def test_barber_sees_all_appointments(self, page: Page):
#         """Tony (barber) sees appointments from all users on the Events page."""
#         navigate(page, "/login")
#         login_via_form(page, "tony@barber.gr", "barberpassword123")
#         nav_click(page, "Events")
#         # Barbers see ALL appointments; at least two should exist (Bill's)
#         delete_links = page.locator("a", has_text="Delete")
#         assert (
#             delete_links.count() >= 2
#         ), f"Expected barber to see ≥2 appointments, found {delete_links.count()}"
