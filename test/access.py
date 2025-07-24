import unittest
import json
import os
import requests


class APITestCase(unittest.TestCase):
    BASE_URL = "http://localhost:8080"
    HEADERS = {
        "Content-Type": "application/json",
    }

    def test_login(self):
        """Test Login"""
        url = f"{self.BASE_URL}/v1/login"
        body = self.read_json("payloads/login.json")
        response = requests.post(url, json=body, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 200)
        token = response.json().get("token")
        print(response)
        self.assertEqual(token, "secure-dev-key-123")
        print(f"\n✅ Login Success {url}: got token: {token}")

    def read_json(self, file_name):
        """Helper to read and parse JSON file"""
        try:
            base_path = os.path.dirname(__file__)
            full_path = os.path.join(base_path, file_name)
            with open(full_path, "r", encoding="utf-8") as file:
                file_content = json.load(file)
        except FileNotFoundError as e:
            self.fail(f"File not found '{file_name}': {e}")

        return file_content

if __name__ == "__main__":
    unittest.main()
