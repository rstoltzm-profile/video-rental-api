import unittest
import requests

class StoreTests(unittest.TestCase):
    BASE_URL = "http://localhost:8080"
    HEADERS = {
        "Content-Type": "application/json",
        "X-API-Key": "secure-dev-key-123"
    }

    def test_get_store_inventory_summary(self):
        """Test GET /v1/stores/1/inventory/summary returns non-empty list"""
        print("\n🏪 Testing: GET /v1/stores/1/inventory/summary")
        url = f"{self.BASE_URL}/v1/stores/1/inventory/summary"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 200)
        summary = response.json()
        self.assertIsInstance(summary, list)
        self.assertGreater(len(summary), 0, "Expected non-empty inventory summary list")
        print("✅ Store inventory summary retrieved successfully")

if __name__ == "__main__":
    unittest.main()
