#!/usr/bin/env python3
"""
seed_test_ads.py — API-only seeder and integration smoke-test.

Registers test users (or logs in if they already exist), then posts ads
through the real backend REST endpoints.  Every HTTP call is recorded as a
PASS / FAIL so the script doubles as a lightweight smoke-test suite.

Usage:
    python tools/seed_test_ads.py [options]

Key options:
    --base-url      Backend base URL          (default: http://localhost:3000)
    --user-count    Number of test users       (default: 5)
    --ads-per-user  Max ads per user           (default: all)
    --ads-yaml      Ads data file              (default: tools/data/test_ads.yaml)
    --admin-token   Admin Bearer token — used to auto-create categories when
                    the database has none yet
    --password      Password for all test users (default: TestPass123!)
    --dry-run       Print the plan without making any requests
"""

import argparse
import sys
import textwrap
from pathlib import Path

import requests
try:
    import yaml
except ImportError:
    print("ERROR: PyYAML is required.  Run: pip install pyyaml")
    sys.exit(1)


# ── Terminal colours ──────────────────────────────────────────────────────────

def _supports_color() -> bool:
    import os
    return hasattr(sys.stdout, "isatty") and sys.stdout.isatty() and os.name != "nt"

_COLOR = _supports_color()

def _c(code: str, text: str) -> str:
    return f"\033[{code}m{text}\033[0m" if _COLOR else text

def green(t):  return _c("32", t)
def red(t):    return _c("31", t)
def yellow(t): return _c("33", t)
def bold(t):   return _c("1",  t)
def dim(t):    return _c("2",  t)


# ── Test result tracker ───────────────────────────────────────────────────────

class Results:
    def __init__(self):
        self._items: list[dict] = []

    def record(self, label: str, passed: bool, detail: str = ""):
        self._items.append({"label": label, "passed": passed, "detail": detail})
        icon = green("PASS") if passed else red("FAIL")
        line = f"  [{icon}] {label}"
        if detail:
            line += f"  {dim(detail)}"
        print(line)

    def summary(self):
        total  = len(self._items)
        passed = sum(1 for r in self._items if r["passed"])
        failed = total - passed
        print()
        print(bold("── Summary " + "─" * 50))
        print(f"  Total : {total}")
        print(f"  {green('Passed')}: {passed}")
        if failed:
            print(f"  {red('Failed')}: {failed}")
            print()
            print(bold("Failed tests:"))
            for r in self._items:
                if not r["passed"]:
                    print(f"  • {r['label']}  {dim(r['detail'])}")
        print()
        return failed == 0


# ── HTTP client ───────────────────────────────────────────────────────────────

class ApiClient:
    """Thin requests wrapper.  All calls return (ok: bool, body: dict|None)."""

    def __init__(self, base_url: str, version: str = "v1"):
        self.base_url = base_url.rstrip("/")
        self.api_url  = f"{self.base_url}/api/{version}"
        self.session  = requests.Session()

    # ── raw helpers ──────────────────────────────────────────────────────────

    def _call(self, method: str, url: str, **kwargs) -> tuple[bool, dict | None]:
        try:
            resp = self.session.request(method, url, timeout=15, **kwargs)
        except requests.RequestException as exc:
            return False, {"_error": str(exc)}

        try:
            body = resp.json() if resp.content else {}
        except ValueError:
            body = {"_raw": resp.text}

        ok = 200 <= resp.status_code < 300
        if not ok:
            body["_status"] = resp.status_code
        return ok, body

    def _api(self, method: str, endpoint: str, **kwargs):
        return self._call(method, f"{self.api_url}{endpoint}", **kwargs)

    # ── public methods ───────────────────────────────────────────────────────

    def get(self, endpoint: str, **kw):
        return self._api("GET", endpoint, **kw)

    def post(self, endpoint: str, **kw):
        return self._api("POST", endpoint, **kw)

    def put(self, endpoint: str, **kw):
        return self._api("PUT", endpoint, **kw)

    def delete(self, endpoint: str, **kw):
        return self._api("DELETE", endpoint, **kw)

    def auth_headers(self, token: str) -> dict:
        return {"Authorization": f"Bearer {token}"}

    # ── health ───────────────────────────────────────────────────────────────

    def health_check(self, results: Results) -> bool:
        ok, body = self._call("GET", f"{self.base_url}/health")
        results.record("GET /health", ok, body.get("_error", ""))
        return ok


# ── Auth helpers ──────────────────────────────────────────────────────────────

def _unwrap_auth(body: dict) -> dict | None:
    """Return the inner {user, tokens} dict from a { data: {…} } response."""
    return body.get("data")


def register_or_login(client: ApiClient, results: Results,
                      phone: str, display_name: str, password: str,
                      country: str = "MA") -> str | None:
    """
    Try to log in.  If the account doesn't exist yet, register it first.
    Returns the access token on success, None on failure.
    """
    login_data = {"identifier": phone, "password": password}
    ok, body = client.post("/auth/login", json=login_data)
    if ok:
        auth = _unwrap_auth(body)
        if auth and auth.get("tokens", {}).get("access_token"):
            results.record(f"POST /auth/login ({display_name})", True, "existing user")
            return auth["tokens"]["access_token"]

    # Not found — register
    reg_data = {
        "phone":        phone,
        "password":     password,
        "display_name": display_name,
        "country":      country,
    }
    ok, body = client.post("/auth/register", json=reg_data)
    if not ok:
        err = body.get("error") or body.get("_raw") or str(body)
        results.record(f"POST /auth/register ({display_name})", False, err)
        return None
    auth = _unwrap_auth(body)
    if not auth or not auth.get("tokens", {}).get("access_token"):
        results.record(
            f"POST /auth/register ({display_name})", False,
            f"unexpected response shape: {body}"
        )
        return None
    results.record(f"POST /auth/register ({display_name})", True, "new user")

    # Verify we can log in right after registering
    ok2, body2 = client.post("/auth/login", json=login_data)
    auth2 = _unwrap_auth(body2) if ok2 else None
    token = auth2["tokens"]["access_token"] if auth2 else auth["tokens"]["access_token"]
    results.record(f"POST /auth/login ({display_name})", ok2, "post-register login")
    return token


def fetch_me(client: ApiClient, results: Results, token: str, label: str) -> dict | None:
    ok, body = client.get("/auth/me", headers=client.auth_headers(token))
    user = body.get("data") if ok else None
    results.record(f"GET /auth/me ({label})", ok and user is not None,
                   f"user_id={user.get('ID')} display={user.get('display_name')}" if user else str(body))
    return user


# ── Category helpers ──────────────────────────────────────────────────────────

def fetch_categories(client: ApiClient, results: Results) -> dict[str, int]:
    """Return {slug: id} for all active categories."""
    ok, body = client.get("/categories")
    cats = body.get("data", []) if ok else []
    results.record("GET /categories", ok, f"{len(cats)} categories returned")

    slug_to_id: dict[str, int] = {}
    def _recurse(items):
        for c in items:
            slug_to_id[c["slug"]] = c["ID"]
            if c.get("children"):
                _recurse(c["children"])
    _recurse(cats)
    return slug_to_id


def ensure_categories_via_admin(
    client: ApiClient,
    results: Results,
    admin_token: str,
    categories: list[dict],
) -> dict[str, int]:
    """Create missing categories using the admin endpoint, return slug->id map."""
    slug_to_id: dict[str, int] = {}
    headers = client.auth_headers(admin_token)
    for cat in categories:
        ok, body = client.post("/admin/categories", json=cat, headers=headers)
        cat_id = body.get("ID") if ok else None
        label  = f"POST /admin/categories ({cat['slug']})"
        results.record(label, ok and cat_id is not None,
                       f"id={cat_id}" if cat_id else str(body))
        if cat_id:
            slug_to_id[cat["slug"]] = cat_id
    return slug_to_id


# ── Ad helpers ────────────────────────────────────────────────────────────────

def create_ad(client: ApiClient, results: Results,
              token: str, ad: dict, category_id: int) -> bool:
    """POST /ads with multipart/form-data (as the backend requires)."""
    data = {
        "category_id": str(category_id),
        "title":       ad["title"],
        "body":        ad.get("body", ""),
        "city":        ad.get("city", "Casablanca"),
        "currency":    ad.get("currency", "MAD"),
    }
    if ad.get("price") is not None:
        data["price"] = str(ad["price"])

    ok, body = client.post("/ads", data=data, headers=client.auth_headers(token))
    ad_slug = body.get("data", {}).get("slug", "") if ok else ""
    results.record(
        f"POST /ads \"{ad['title'][:40]}\"",
        ok,
        f"slug={ad_slug}" if ok else (body.get("error") or str(body)),
    )
    return ok


def verify_ads_listed(client: ApiClient, results: Results,
                      category_slug: str, expected_min: int):
    """GET /ads?cat=<slug> and assert at least expected_min results."""
    ok, body = client.get(f"/ads?cat={category_slug}&limit=100")
    total = body.get("total", 0) if ok else 0
    passed = ok and total >= expected_min
    results.record(
        f"GET /ads?cat={category_slug}",
        passed,
        f"total={total} (expected >= {expected_min})",
    )


# ── Data loading ──────────────────────────────────────────────────────────────

def load_yaml(path: Path) -> dict:
    if not path.exists():
        print(red(f"ERROR: file not found: {path}"))
        sys.exit(1)
    with path.open("r", encoding="utf-8") as f:
        return yaml.safe_load(f) or {}


# ── CLI ───────────────────────────────────────────────────────────────────────

def parse_args():
    p = argparse.ArgumentParser(
        description=textwrap.dedent("""\
            Seed the backend with realistic test data via REST API.
            Also validates every endpoint it touches (PASS / FAIL).
        """),
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    p.add_argument("--base-url",    default="http://localhost:3000",
                   help="Backend base URL (default: http://localhost:3000)")
    p.add_argument("--user-count",  type=int, default=5,
                   help="Number of test users to create (default: 5)")
    p.add_argument("--ads-per-user", type=int, default=0,
                   help="Max ads per user; 0 = all ads in YAML (default: 0)")
    p.add_argument("--ads-yaml",    default="tools/data/test_ads.yaml",
                   help="Path to ads YAML (default: tools/data/test_ads.yaml)")
    p.add_argument("--admin-token", default="",
                   help="Admin Bearer token (needed to auto-create categories)")
    p.add_argument("--password",    default="TestPass123!",
                   help="Password for all test users (default: TestPass123!)")
    p.add_argument("--dry-run",     action="store_true",
                   help="Print plan without making requests")
    return p.parse_args()


# ── Main ──────────────────────────────────────────────────────────────────────

def main():
    args    = parse_args()
    results = Results()
    repo_root = Path(__file__).resolve().parent.parent

    ads_yaml_path = (
        Path(args.ads_yaml) if Path(args.ads_yaml).is_absolute()
        else (repo_root / args.ads_yaml).resolve()
    )

    print(bold("\n── Mezian API seeder / smoke-test ───────────────────────────────"))
    print(f"  Base URL  : {args.base_url}")
    print(f"  Users     : {args.user_count}")
    print(f"  Ads YAML  : {ads_yaml_path}")
    print(f"  Dry run   : {args.dry_run}")
    print()

    raw = load_yaml(ads_yaml_path)
    ads_data: list[dict]           = raw.get("ads", [])
    seed_categories: list[dict]    = raw.get("categories", [])

    if not ads_data:
        print(red("ERROR: no ads found in YAML file"))
        sys.exit(1)

    ads_slice = ads_data if args.ads_per_user == 0 else ads_data[:args.ads_per_user]

    if args.dry_run:
        print(yellow("DRY RUN — no requests will be made"))
        print(f"  Would register/login {args.user_count} users")
        print(f"  Would post {len(ads_slice)} ads per user")
        print(f"  Total ads: {args.user_count * len(ads_slice)}")
        sys.exit(0)

    client = ApiClient(args.base_url)

    # ── 1. Health check ───────────────────────────────────────────────────────
    print(bold("── Health ───────────────────────────────────────────────────────"))
    if not client.health_check(results):
        print(red("\nBackend is unreachable — aborting."))
        sys.exit(1)
    print()

    # ── 2. Resolve categories ─────────────────────────────────────────────────
    print(bold("── Categories ───────────────────────────────────────────────────"))
    slug_to_id = fetch_categories(client, results)

    if not slug_to_id and seed_categories and args.admin_token:
        print(yellow("  No categories found — creating via admin token…"))
        slug_to_id = ensure_categories_via_admin(
            client, results, args.admin_token, seed_categories
        )
    elif not slug_to_id:
        print(yellow(
            "  WARNING: no categories in the database.\n"
            "  Pass --admin-token to auto-create them, or seed categories manually.\n"
            "  Ads will be skipped if their category_slug cannot be resolved."
        ))
    print()

    # ── 3. Register/login users and post ads ──────────────────────────────────
    total_ads_created = 0

    for i in range(1, args.user_count + 1):
        phone        = f"+2126{i:08d}"
        display_name = f"Testeur {i}"

        print(bold(f"── User {i}/{args.user_count}: {display_name} ({phone}) ───"))

        token = register_or_login(
            client, results, phone, display_name, args.password
        )
        if not token:
            print(red(f"  Could not obtain token for {display_name} — skipping ads"))
            print()
            continue

        # Verify /auth/me works with this token
        fetch_me(client, results, token, display_name)

        # Post ads
        user_ads_created = 0
        for ad in ads_slice:
            cat_slug = ad.get("category_slug", "")
            cat_id   = slug_to_id.get(cat_slug)
            if cat_id is None:
                results.record(
                    f"POST /ads \"{ad['title'][:40]}\"", False,
                    f"category_slug '{cat_slug}' not found in DB — skipped"
                )
                continue
            if create_ad(client, results, token, ad, cat_id):
                user_ads_created += 1

        total_ads_created += user_ads_created
        print(dim(f"  → {user_ads_created} ad(s) created for {display_name}"))
        print()

    # ── 4. Verify listings per category ──────────────────────────────────────
    if slug_to_id and total_ads_created > 0:
        print(bold("── Listing verification ─────────────────────────────────────────"))
        # Check each slug that was actually used in the seed data
        used_slugs = {ad.get("category_slug") for ad in ads_slice if ad.get("category_slug") in slug_to_id}
        for slug in sorted(used_slugs):
            verify_ads_listed(client, results, slug, expected_min=1)
        print()

    # ── 5. Summary ────────────────────────────────────────────────────────────
    print(f"  Total ads posted: {bold(str(total_ads_created))}")
    all_passed = results.summary()
    sys.exit(0 if all_passed else 1)


if __name__ == "__main__":
    main()
