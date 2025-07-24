import unittest
import requests

class FilmTests(unittest.TestCase):
    BASE_URL = "http://localhost:8080"
    HEADERS = {
        "Content-Type": "application/json",
        "X-API-Key": "secure-dev-key-123"
    }

    def test_get_all_films(self):
        """Test GET /v1/films returns non-empty list"""
        print("\n🎬 Testing: GET /v1/films")
        url = f"{self.BASE_URL}/v1/films"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 200)
        films = response.json()
        self.assertIsInstance(films, list)
        self.assertGreater(len(films), 0, "Expected non-empty films list")
        print("✅ Films list retrieved successfully")

    def test_get_film_by_id(self):
        """Test GET /v1/films/1 returns film details"""
        print("\n🎞️ Testing: GET /v1/films/1")
        url = f"{self.BASE_URL}/v1/films/1"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 200)
        film = response.json()
        self.assertIsInstance(film, dict)
        self.assertGreater(len(film), 0, "Expected non-empty film details")
        print("✅ Film details retrieved successfully")

    def test_search_film_by_title(self):
        """Test GET /v1/films/search?title=ACADEMY DINOSAUR returns results"""
        print("\n🔍 Testing: GET /v1/films/search?title=ACADEMY DINOSAUR")
        url = f"{self.BASE_URL}/v1/films/search?title=ACADEMY%20DINOSAUR"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 200)
        results = response.json()
        self.assertIsInstance(results, list)
        self.assertGreater(len(results), 0, "Expected non-empty search results")
        print("✅ Film search results retrieved successfully")

    def test_get_film_with_actors_and_categories(self):
        """Test GET /v1/films/1/with-actors-categories returns enriched film data"""
        print("\n🎭 Testing: GET /v1/films/1/with-actors-categories")
        url = f"{self.BASE_URL}/v1/films/1/with-actors-categories"
        response = requests.get(url, headers=self.HEADERS, timeout=60)
        self.assertEqual(response.status_code, 200)
        enriched_film = response.json()
        self.assertIsInstance(enriched_film, dict)
        self.assertGreater(len(enriched_film), 0, "Expected non-empty enriched film data")
        print("✅ Enriched film data retrieved successfully")

if __name__ == "__main__":
    print("\n===== STARTING Film Tests =====")
    unittest.main()
    print("\n===== FINISHED Film Tests =====\n")
