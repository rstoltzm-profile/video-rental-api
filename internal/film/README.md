# film domain

## 
```bash
curl -s $BASE_URL/v1/films/1/with-actors-categories
```
```json
{
    "actors": [
        "SANDRA PECK",
        "JOHNNY CAGE",
        "OPRAH KILMER",
        "MARY KEITEL",
        "PENELOPE GUINESS",
        "LUCILLE TRACY",
        "MENA TEMPLE",
        "WARREN NOLTE",
        "ROCK DUKAKIS",
        "CHRISTIAN GABLE"
    ],
    "categories": [
        "Games",
        "New",
        "Travel"
    ],
    "description": "A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rockies",
    "language": "English",
    "rating": "PG",
    "release_year": 2012,
    "title": "ACADEMY DINOSAUR"
}
```