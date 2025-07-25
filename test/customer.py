import unittest
import json
import os
import requests


class APITestCase(unittest.TestCase):
    BASE_URL = "http://localhost:8080"
    HEADERS = {
        "Content-Type": "application/json",
        "X-API-Key": "secure-dev-key-123"
    }

    def test_get_health(self):
        """Test GET health request"""
        url = f"{self.BASE_URL}/health"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 200)
        status = response.json().get("status")
        self.assertEqual(status, "ok")
        print(f"\n✅ API Health {url}: {status}")

    def test_get_customers(self):
        """Test get all customers"""
        # Step 1: Rent a movie
        url = f"{self.BASE_URL}/v1/customers"
        response = requests.get(url, headers=self.HEADERS, timeout=60)

        customers = response.json()
        self.assertEqual(response.status_code, 200, f"get customers failed: {response.text}")
        self.assertIsInstance(customers, list, "Expected a list of customers")
        self.assertGreater(len(customers), 0, "Customer list is empty")
        print(f"\n✅ Customer list is not empty from {url}")

        first_customer = customers[0]
        self.assertIn("first_name", first_customer, "Missing 'first_name' in customer data")
        print("\n✅ First customer Exists in list of customers")


    def test_get_customer(self):
        """Test get a customer"""
        url = f"{self.BASE_URL}/v1/customers/1"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        customer = response.json()
        self.assertIn("first_name", customer, "Missing 'first_name' in customer data")
        print(f"\n✅ First customer returned from {url}: {customer['first_name']}")

    def test_create_and_delete_customer(self):
        """Test create and delete a customer"""
        customer_id = self.create_customer()
        self.assertIsNotNone(customer_id)
        print(f"\n✅ Created customer: {customer_id}")
        delete_response = self.delete_customer(customer_id)

        self.assertEqual(delete_response, 204)
        print(f"\n✅ Deleted customer {customer_id}: response {delete_response}")

    def create_customer(self):
        url = f"{self.BASE_URL}/v1/customers"
        body = self.read_json("payloads/customer.json")
        response = requests.post(url, json=body, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 201, f"Create customer failed: {response.text}")
        customer_id = response.json().get("id")
        return customer_id

    def delete_customer(self, customer_id):
        url = f"{self.BASE_URL}/v1/customers/{customer_id}"
        response = requests.delete(url, headers=self.HEADERS, timeout=60)
        return response.status_code

    def test_create_customer_validation_failure(self):
        """Test create customer with invalid data returns validation error"""
        url = f"{self.BASE_URL}/v1/customers"

        # Test with missing required fields
        invalid_body = {"invalid": "json"}

        response = requests.post(url, json=invalid_body, headers=self.HEADERS, timeout=60)

        # Should return 400 Bad Request for validation failure
        error = f"Expected 400 for validation failure, got {response.status_code}"
        self.assertEqual(response.status_code, 400, error)

        # Check that response contains validation error message
        response_text = response.text
        self.assertIn("validation", response_text.lower(),
                       "Response should contain validation error")
        self.assertIn("required", response_text.lower(),
                       "Response should mention required fields")

        print(f"\n✅ Validation failure test passed: {response.status_code}")
        print(f"   Error message: {response_text[:100]}...")

    def test_create_customer_partial_validation_failure(self):
        """Test create customer with partially valid data"""
        url = f"{self.BASE_URL}/v1/customers"

        # Test with some fields but missing others
        partial_body = {
            "first_name": "John",
            "last_name": "Doe"
            # Missing email, store_id, address
        }

        response = requests.post(url, json=partial_body, headers=self.HEADERS, timeout=60)

        self.assertEqual(response.status_code, 400,
                         f"Expected 400 for validation failure, got {response.status_code}")

        response_text = response.text
        self.assertIn("email", response_text.lower(), "Should complain about missing email")
        self.assertIn("storeid", response_text.lower(), "Should complain about missing storeid")
        self.assertIn("address", response_text.lower(), "Should complain about missing address")

        print(f"\n✅ Partial validation failure test passed: {response.status_code}")
        print(f"   Missing fields detected in: {response_text[:100]}...")

    def test_get_customer_rentals(self):
        """Test get customer's rentals"""
        url = f"{self.BASE_URL}/v1/customers/373/rentals"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        customer_rentals = response.json()[0]
        self.assertIn("first_name", customer_rentals, "Missing 'first_name' in customer data")
        print(f"\n✅ customer_rentals returned from {url}: {customer_rentals['first_name']}")

    def test_get_customer_late_rentals(self):
        """Test get customer's late rentals"""
        url = f"{self.BASE_URL}/v1/customers/373/rentals?late=true"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        customer_rentals = response.json()[0]
        self.assertIn("first_name", customer_rentals, "Missing 'first_name' in customer data")
        print(f"\n✅ Customer late rentals returned from {url}: {customer_rentals['first_name']}")

    def test_make_payment(self):
        url = f"{self.BASE_URL}/v1/payments"
        body = self.read_json("payloads/payment.json")
        response = requests.post(url, json=body, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 201, f"Make payment failed: {response.text}")
        payment_id = response.json()
        print(f"\n✅ Make payment succeeded: {payment_id}")

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
