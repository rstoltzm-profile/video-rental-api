import unittest
import requests
import json
import os

class APITestCase(unittest.TestCase):
    BASE_URL = "http://localhost:8080"
    HEADERS = {
        "Content-Type": "application/json",
        "X-API-Key": "secure-dev-key-123"
    }

    def test_get_health(self):
        """Test GET health request"""
        url = f"{self.BASE_URL}/health"
        response = requests.get(url, headers=self.HEADERS)
        self.assertEqual(response.status_code, 200)
        status = response.json().get("status")
        self.assertEqual(status, "ok")
        print(f"\n✅ API /health: {status}")

    def test_rent_and_return_movie(self):
        """Test renting and returning a movie"""
        # Step 1: Rent a movie
        rent_url = f"{self.BASE_URL}/v1/rentals"
        body = self.read_json("payloads/rental.json")
        rent_response = requests.post(rent_url, json=body, headers=self.HEADERS)

        self.assertEqual(rent_response.status_code, 201, f"Rent failed: {rent_response.text}")
        rental_id = rent_response.json().get("id")
        self.assertIsNotNone(rental_id, "Rental ID not found in response")
        self.assertIsInstance(rental_id, int, "Rental ID is not an integer")

        print(f"\n✅ Rental created with ID: {rental_id}")

        # Step 2: Return the movie
        return_url = f"{self.BASE_URL}/v1/rentals/{rental_id}/return"
        return_response = requests.post(return_url, headers=self.HEADERS)

        self.assertEqual(return_response.status_code, 204, f"Return failed: {return_response.status_code}")
        print(f"✅ Movie returned successfully (status {return_response.status_code})")

    def test_rent_checked_out_movie(self):
        """Test renting a checked out movie"""
        rent_url = f"{self.BASE_URL}/v1/rentals"
        body = self.read_json("payloads/rental-2.json")

        # First rental should succeed
        rent_response_1 = requests.post(rent_url, json=body, headers=self.HEADERS)
        self.assertEqual(rent_response_1.status_code, 201, f"First rent failed: {rent_response_1.text}")
        rental_id = rent_response_1.json().get("id")
        self.assertIsNotNone(rental_id, "Rental ID not found in response")

        print(f"\n✅ First rental created with ID: {rental_id}")
        # Second rental should fail
        rent_response_2 = requests.post(rent_url, json=body, headers=self.HEADERS)
        self.assertNotEqual(rent_response_2.status_code, 201, "Second rental should not succeed")
        self.assertIn(rent_response_2.status_code, [400, 409, 422, 500], f"Unexpected error code: {rent_response_2.status_code}")
        print(f"❌ Second rental failed as expected (status {rent_response_2.status_code})")

        # Return the movie
        return_url = f"{self.BASE_URL}/v1/rentals/{rental_id}/return"
        return_response = requests.post(return_url, headers=self.HEADERS)
        self.assertEqual(return_response.status_code, 204, f"Return failed: {return_response.status_code}")
        print(f"✅ Movie returned successfully (status {return_response.status_code})")


    def test_get_all_rentals(self):
        """Test GET /v1/rentals returns non-empty list"""
        print("\n🔍 Testing: GET /v1/rentals")
        url = f"{self.BASE_URL}/v1/rentals"
        response = requests.get(url, headers=self.HEADERS)
        self.assertEqual(response.status_code, 200)
        rentals = response.json()
        self.assertIsInstance(rentals, list)
        self.assertGreater(len(rentals), 0, "Expected non-empty rentals list")
        print("✅ Rentals list retrieved successfully")

    def test_get_late_rentals(self):
        """Test GET /v1/rentals?late=true returns non-empty list"""
        print("\n🔍 Testing: GET /v1/rentals?late=true")
        url = f"{self.BASE_URL}/v1/rentals?late=true"
        response = requests.get(url, headers=self.HEADERS)
        self.assertEqual(response.status_code, 200)
        rentals = response.json()
        self.assertIsInstance(rentals, list)
        self.assertGreater(len(rentals), 0, "Expected non-empty late rentals list")
        print("✅ Late rentals list retrieved successfully")

    def test_get_rentals_by_customer(self):
        """Test GET /v1/rentals?customer_id=373 returns non-empty list"""
        print("\n🔍 Testing: GET /v1/rentals?customer_id=373")
        url = f"{self.BASE_URL}/v1/rentals?customer_id=373"
        response = requests.get(url, headers=self.HEADERS)
        self.assertEqual(response.status_code, 200)
        rentals = response.json()
        self.assertIsInstance(rentals, list)
        self.assertGreater(len(rentals), 0, "Expected non-empty rentals for customer")
        print("✅ Customer rentals retrieved successfully")

    def test_get_late_rentals_by_customer(self):
        """Test GET /v1/rentals?customer_id=373&late=true returns non-empty list"""
        print("\n🔍 Testing: GET /v1/rentals?customer_id=373&late=true")
        url = f"{self.BASE_URL}/v1/rentals?customer_id=373&late=true"
        response = requests.get(url, headers=self.HEADERS)
        self.assertEqual(response.status_code, 200)
        rentals = response.json()
        self.assertIsInstance(rentals, list)
        self.assertGreater(len(rentals), 0, "Expected non-empty late rentals for customer")
        print("✅ Late customer rentals retrieved successfully")


    def read_json(self, file_name):
        """Helper to read and parse JSON file"""
        try:
            base_path = os.path.dirname(__file__)
            full_path = os.path.join(base_path, file_name)
            with open(full_path, "r") as file:
                return json.load(file)
        except Exception as e:
            self.fail(f"Failed to read JSON file '{file_name}': {e}")

if __name__ == "__main__":
    unittest.main()