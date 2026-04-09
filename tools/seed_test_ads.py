#!/usr/bin/env python3
import argparse
import json
import os
import sqlite3
import sys
from datetime import datetime
from pathlib import Path
import requests
try:
    import yaml
except ImportError:
    import ruamel.yaml as yaml


class ApiClient:
    def __init__(self, base_url: str, version:str="v1"):
        self.base_url = base_url.rstrip('/')
        self.api_url = self.base_url + "/api/{:s}".format(version)
        self.session = requests.Session()

    def _make_base_request(self, method: str, endpoint: str, **kwargs):
        url = f"{self.base_url}{endpoint}"
        try:
            resp = self.session.request(method, url, **kwargs)
            if resp.status_code >= 200 and resp.status_code < 300:
                return "OK", resp.json() if resp.content else None
            else:
                warning = f"Request failed: {method} {endpoint} - {resp.status_code} {resp.text}"
                print(f"WARNING: {warning}")
                return "WARNING", None
        except requests.RequestException as e:
            warning = f"Request exception: {method} {endpoint} - {e}"
            print(f"WARNING: {warning}")
            return "WARNING", None

    def _make_request(self, method: str, endpoint: str, **kwargs):
        url = f"{self.api_url}{endpoint}"
        try:
            resp = self.session.request(method, url, **kwargs)
            if resp.status_code >= 200 and resp.status_code < 300:
                return "OK", resp.json() if resp.content else None
            else:
                warning = f"Request failed: {method} {endpoint} - {resp.status_code} {resp.text}"
                print(f"WARNING: {warning}")
                return "WARNING", None
        except requests.RequestException as e:
            warning = f"Request exception: {method} {endpoint} - {e}"
            print(f"WARNING: {warning}")
            return "WARNING", None

    def get(self, endpoint: str, **kwargs):
        return self._make_request("GET", endpoint, **kwargs)

    def post(self, endpoint: str, **kwargs):
        return self._make_request("POST", endpoint, **kwargs)

    def put(self, endpoint: str, **kwargs):
        return self._make_request("PUT", endpoint, **kwargs)

    def delete(self, endpoint: str, **kwargs):
        return self._make_request("DELETE", endpoint, **kwargs)

    def health_check(self):
        # Health endpoint is at root URL, not /api/v1
        root_url = self.base_url.replace('/api/v1', '')
        url = f"{root_url}/health"
        try:
            resp = requests.get(url, timeout=10)
            if resp and resp.status_code >= 200 and resp.status_code < 300:
                print("Health check OK.")
            else:
                status = resp.status_code if resp else "unknown"
                text = resp.text if resp else "no response"
                print(f"ERROR: Health check failed: {status} {text}")
                sys.exit(1)
        except requests.RequestException as e:
            print(f"ERROR: Health check request failed: {e}")
            sys.exit(1)

    def register_user(self, data: dict):
        status, resp = self.post("/auth/register", json=data, timeout=10)
        if status != "OK":
            print("ERROR: User registration failed. Interrupting.")
            sys.exit(1)
        return resp

    def login_user(self, data: dict):
        status, resp = self.post("/auth/login", json=data, timeout=10)
        if status != "OK":
            print("ERROR: User login failed. Interrupting.")
            sys.exit(1)
        return resp

    def create_ad(self, token: str, data: dict):
        headers = {"Authorization": f"Bearer {token}"}
        status, resp = self.post("/ads", json=data, headers=headers)
        if status != "OK":
            print(f"WARNING: Failed to create ad '{data.get('title', 'Unknown')}'")
        return status == "OK"


def parse_args():
    parser = argparse.ArgumentParser(
        description="Seed backend with test users and ads using API."
    )
    parser.add_argument(
        "--config",
        default="backend/config/config.yaml",
        help="Path to backend config file (YAML or JSON).",
    )
    parser.add_argument(
        "--user-count",
        type=int,
        default=10,
        help="Number of test users to create.",
    )
    parser.add_argument(
        "--ads-yaml",
        default="tools/data/test_ads.yaml",
        help="Path to YAML file with list of ads data.",
    )
    parser.add_argument(
        "--category-slug",
        default="test-category",
        help="Category slug to create or reuse for test ads.",
    )
    return parser.parse_args()


def load_config(path: Path):
    if not path.exists():
        raise FileNotFoundError(f"Config file not found: {path}")

    if path.suffix.lower() in (".yaml", ".yml"):
        return load_yaml_config(path)
    if path.suffix.lower() == ".json":
        return load_json_config(path)

    # Try JSON first, then YAML-like fallback.
    try:
        return load_json_config(path)
    except ValueError:
        return load_yaml_config(path)


def load_json_config(path: Path):
    with path.open("r", encoding="utf-8") as f:
        data = json.load(f)
    if "database" not in data or "path" not in data["database"]:
        raise ValueError("Missing database.path in JSON config")
    return data


def load_yaml_config(path: Path):
    try:
        import yaml
    except ImportError:
        return parse_yaml_without_pyyaml(path)

    with path.open("r", encoding="utf-8") as f:
        data = yaml.safe_load(f)
    if not isinstance(data, dict) or "database" not in data or "path" not in data["database"]:
        raise ValueError("Missing database.path in YAML config")
    return data


def parse_yaml_without_pyyaml(path: Path):
    data = {}
    current_section = None
    with path.open("r", encoding="utf-8") as f:
        for raw_line in f:
            line = raw_line.strip()
            if not line or line.startswith("#"):
                continue
            if line.endswith(":"):
                current_section = line[:-1].strip()
                data[current_section] = {}
                continue
            if ":" not in line:
                continue
            key, value = [part.strip() for part in line.split(":", 1)]
            if current_section and key and current_section not in data:
                data[current_section] = {}
            if current_section and current_section in data:
                data[current_section][key] = strip_quotes(value)
            else:
                data[key] = strip_quotes(value)
    if "database" not in data or "path" not in data["database"]:
        raise ValueError("Missing database.path in YAML config")
    return data


def strip_quotes(value: str):
    if (value.startswith('"') and value.endswith('"')) or (value.startswith("'") and value.endswith("'")):
        return value[1:-1]
    return value


def load_ads_yaml(path: Path):
    if not path.exists():
        raise FileNotFoundError(f"Ads YAML file not found: {path}")
    with path.open("r", encoding="utf-8") as f:
        data = yaml.safe_load(f)
    return data.get("ads", [])


def make_db_path(config_data, config_path: Path) -> Path:
    db_path = config_data["database"]["path"]
    db_path = str(db_path)
    candidate = Path(db_path)
    if not candidate.is_absolute():
        candidate = (config_path.parent.parent / candidate).resolve()
    return candidate


def ensure_category(cursor, slug: str):
    cursor.execute("SELECT id FROM categories WHERE slug = ? LIMIT 1", (slug,))
    row = cursor.fetchone()
    if row:
        return row[0]

    now = datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S")
    cursor.execute(
        "INSERT INTO categories (slug, name_fr, name_ar, name_en, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, 1, ?, ?)",
        (slug, "Catégorie de test", "فئة اختبار", "Test category", now, now),
    )
    return cursor.lastrowid








def main():
    args = parse_args()
    repo_root = Path(__file__).resolve().parent.parent
    config_path = (repo_root / args.config).resolve() if not Path(args.config).is_absolute() else Path(args.config)
    ads_yaml_path = (repo_root / args.ads_yaml).resolve() if not Path(args.ads_yaml).is_absolute() else Path(args.ads_yaml)

    config_data = load_config(config_path)

    ads_data = load_ads_yaml(ads_yaml_path)
    if not ads_data:
        print("No ads data found in YAML")
        sys.exit(1)

    print(f"Using config: {config_path}")
    print(f"Using ads YAML: {ads_yaml_path}")

    # Now use API to register users and create ads
    base_url = "http://localhost:3000"
    client = ApiClient(base_url)
    client.health_check()

    password = "TestPass123!"
    total_created = 0
    category_id = 1  # Assume category ID 1 exists

    for i in range(1, args.user_count + 1):
        user_id = f"tester{i}"
        phone = f"+2129000000{i:02d}"  # +212600000001 to 010
        display_name = f"Tester {i}"

        # Try to login first
        login_data = {"identifier": phone, "password": password}
        status, login_resp = client.post("/auth/login", json=login_data, timeout=10)
        if status == "OK":
            token = login_resp["access_token"]
            print(f"Login successful for existing user {user_id}")
        else:
            # Login failed, register new user
            register_data = {
                "phone": phone,
                "password": password,
                "display_name": display_name,
                "country": "MA"
            }
            client.register_user(register_data)
            print(f"Registered new user {user_id}")

            # Now login
            login_resp = client.login_user(login_data)
            token = login_resp["access_token"]
            print(f"Login successful for {user_id}")

        # Create ads for this user
        created = 0
        for ad in ads_data:
            ad_data = {
                "category_id": category_id,
                "title": ad["title"],
                "body": ad["body"],
                "price": ad.get("price"),
                "currency": ad.get("currency", "MAD"),
                "city": ad.get("city", "Casablanca"),
                "status": ad.get("status", "active")
            }
            if client.create_ad(token, ad_data):
                created += 1
        total_created += created
        print(f"Created {created} ads for {user_id}")

    print(f"Total ads created: {total_created} for {args.user_count} users.")


if __name__ == "__main__":
    main()
