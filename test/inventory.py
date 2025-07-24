import unittest
import requests

class InventoryTests(unittest.TestCase):
    BASE_URL = "http://localhost:8080"
    HEADERS = {
        "Content-Type": "application/json",
        "X-API-Key": "secure-dev-key-123"
    }

    def test_get_all_inventory(self):
        """Test GET /v1/inventory returns non-empty list"""
        print("\n📦 Testing: GET /v1/inventory")
        url = f"{self.BASE_URL}/v1/inventory"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 200)
        inventory = response.json()
        self.assertIsInstance(inventory, list)
        self.assertGreater(len(inventory), 0, "Expected non-empty inventory list")
        print("✅ Inventory list retrieved successfully")

    def test_get_inventory_by_store(self):
        """Test GET /v1/inventory?store_id=1 returns non-empty list"""
        print("\n🏬 Testing: GET /v1/inventory?store_id=1")
        url = f"{self.BASE_URL}/v1/inventory?store_id=1"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 200)
        inventory = response.json()
        self.assertIsInstance(inventory, list)
        self.assertGreater(len(inventory), 0, "Expected non-empty inventory list for store")
        print("✅ Store inventory retrieved successfully")

    def test_get_available_inventory(self):
        """Test GET /v1/inventory/available?film_id=1&store_id=2 returns non-empty list"""
        print("\n🎞️ Testing: GET /v1/inventory/available?film_id=1&store_id=2")
        url = f"{self.BASE_URL}/v1/inventory/available?film_id=1&store_id=2"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 200)
        inventory = response.json()
        self.assertIsInstance(inventory, dict)
        self.assertGreater(len(inventory), 0, "Expected non-empty available inventory list")
        print("✅ Available inventory retrieved successfully")

if __name__ == "__main__":
    unittest.main()