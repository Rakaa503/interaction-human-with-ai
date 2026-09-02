from pathlib import Path

from .engine import SearchEngine


BASE_DIR = Path(__file__).resolve().parent

DATA_PATH = BASE_DIR / "data" / "documents.json"


def main():
    engine = SearchEngine(str(DATA_PATH))

    query = "what is machine learning"

    results = engine.search(
        query,
        top_k=5,
    )

    print("\nQUERY:")
    print(query)

    print("\nRESULTS:")

    for result in results:
        print(
            f"{result['score']:.6f} | "
            f"{result['title']} | "
            f"{result['url']}"
        )


if __name__ == "__main__":
    main()