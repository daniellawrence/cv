"""
Integration tests to validate that each backend responds with valid JSON.

Uses pytest for simple, parametrized testing.
"""

import json
import pytest
import urllib.request
import urllib.error


# Base URLs for the backends - can be overridden via environment variables
BACKEND_URLS = {
    "identity": "http://identity.localhost/identity",
    "experience": "http://experience.localhost/experience",
    "education": "http://education.localhost/education",
    "interest": "http://interest.localhost/interest",
    "qrcode": "http://qrcode.localhost/qrcode",
}


def fetch_json(url: str) -> dict:
    """Fetch URL and return parsed JSON response."""
    try:
        with urllib.request.urlopen(url, timeout=10) as response:
            assert response.status == 200, f"Expected status 200 from {url}, got {response.status}"
            data = json.loads(response.read().decode('utf-8'))
            return data
    except urllib.error.HTTPError as e:
        pytest.fail(f"HTTP {e.code} error fetching {url}: {e.reason}")
    except urllib.error.URLError as e:
        pytest.fail(f"Failed to fetch {url}: {e}")


def test_identity_endpoint():
    """Test that identity endpoint returns valid JSON."""
    data = fetch_json(BACKEND_URLS["identity"])
    assert isinstance(data, dict), "Response should be a JSON object"
    assert "identity" in data, "Response should contain 'identity' field"


def test_experience_endpoint():
    """Test that experience endpoint returns valid JSON."""
    data = fetch_json(BACKEND_URLS["experience"])
    assert isinstance(data, dict), "Response should be a JSON object"
    assert "experience" in data, "Response should contain 'experience' field"


def test_experience_pagination():
    """Test that experience pagination endpoint returns valid JSON."""
    url = f"{BACKEND_URLS['experience']}/0/4"
    data = fetch_json(url)
    assert isinstance(data, dict), "Response should be a JSON object"
    assert "experience" in data, "Response should contain 'experience' field"


def test_education_endpoint():
    """Test that education endpoint returns valid JSON."""
    data = fetch_json(BACKEND_URLS["education"])
    assert isinstance(data, dict), "Response should be a JSON object"
    assert "education" in data, "Response should contain 'education' field"


def test_interest_endpoint():
    """Test that interest endpoint returns valid JSON."""
    data = fetch_json(BACKEND_URLS["interest"])
    assert isinstance(data, dict), "Response should be a JSON object"
    assert "interest" in data, "Response should contain 'interest' field"


def test_qrcode_endpoint():
    """Test that qrcode endpoint returns valid JSON."""
    url = f"{BACKEND_URLS['qrcode']}?url=https://example.com"
    try:
        with urllib.request.urlopen(url, timeout=10) as response:
            assert response.status == 200, f"Expected status 200 from {url}, got {response.status}"
            data = json.loads(response.read().decode('utf-8'))
            assert isinstance(data, dict), "Response should be a JSON object"
            assert "url" in data, "Response should contain 'url' field"
            assert "imageBase64" in data, "Response should contain 'imageBase64' field"
    except urllib.error.HTTPError as e:
        pytest.fail(f"HTTP {e.code} error fetching QR code endpoint: {e.reason}")


def test_qrcode_missing_url():
    """Test that qrcode without url parameter returns an error."""
    url = BACKEND_URLS['qrcode']
    try:
        with urllib.request.urlopen(url, timeout=10) as response:
            data = json.loads(response.read().decode('utf-8'))
            assert isinstance(data, dict), "Error response should be JSON"
    except urllib.error.HTTPError as e:
        # HTTP 400 is expected for missing URL parameter
        if e.code == 400:
            pass  # Expected behavior
        else:
            pytest.fail(f"Unexpected HTTP error {e.code}")


# Parametrized test for all backends (excluding qrcode which requires parameters)
@pytest.mark.parametrize("name, url", [
    ("identity", "http://identity.localhost/identity"),
    ("experience", "http://experience.localhost/experience"),
    ("education", "http://education.localhost/education"),
    ("interest", "http://interest.localhost/interest"),

])
def test_all_backends_json(name, url):
    """Parametrized test to validate each backend returns valid JSON."""
    data = fetch_json(url)
    assert isinstance(data, dict), f"Backend {name} should return a JSON object"


# Health check endpoints for each service (for reference)
HEALTH_CHECKS = {
    "identity": "http://identity.localhost/healthz",
    "experience": "http://experience.localhost/healthz",
    "education": "http://education.localhost/healthz",
    "interest": "http://interest.localhost/healthz",
    "qrcode": "http://qrcode.localhost/healthz",
}


# Parametrized test for all healthz endpoints
@pytest.mark.parametrize("name, url", [
    ("identity", "http://identity.localhost/healthz"),
    ("experience", "http://experience.localhost/healthz"),
    ("education", "http://education.localhost/healthz"),
    ("interest", "http://interest.localhost/healthz"),
    ("qrcode", "http://qrcode.localhost/healthz"),
])
def test_all_healthz_endpoints(name, url):
    """Parametrized test to validate each service healthz endpoint returns ok."""
    try:
        with urllib.request.urlopen(url, timeout=10) as response:
            assert response.status == 200, f"Expected status 200 from {url}, got {response.status}"
            data = response.read().decode('utf-8').strip()
            # Healthz can return plain text "ok" or JSON {"status": "ok"}
            if data.startswith('{'):
                parsed = json.loads(data)
                assert isinstance(parsed, dict), f"{name} healthz should return a JSON object"
                assert parsed.get("status") == "ok", f"{name} healthz expected status 'ok', got '{parsed.get('status')}'"
            else:
                assert data == "ok", f"{name} healthz expected plain text 'ok', got '{data}'"
    except urllib.error.HTTPError as e:
        pytest.fail(f"HTTP {e.code} error fetching {name} /healthz endpoint: {e.reason}")
