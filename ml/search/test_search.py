from pathlib import Path

from .engine import SearchEngine
from .similarity import explain_cosine_similarity


BASE_DIR = Path(__file__).resolve().parent

DATA_PATH = BASE_DIR / "data" / "documents.json"


def print_section(title: str):
    print()
    print("=" * 70)
    print(title)
    print("=" * 70)


def main():

    engine = SearchEngine(str(DATA_PATH))

    query = "what is machine learning"

    # ============================================================
    # QUERY
    # ============================================================

    print_section("QUERY")

    print(query)

    # ============================================================
    # TOKENIZATION
    # ============================================================

    explanation = engine.tfidf.explain_query(query)

    print_section("TOKENS")

    print(explanation["tokens"])

    # ============================================================
    # QUERY TF
    # ============================================================

    print_section("QUERY TF")

    for term, value in explanation["tf"].items():
        print(f"{term:<20} {value:.6f}")

    # ============================================================
    # IDF
    # ============================================================

    print_section("IDF")

    for term, value in engine.tfidf.idf.items():

        if term in explanation["tokens"]:
            print(f"{term:<20} {value:.6f}")

    # ============================================================
    # QUERY TF-IDF
    # ============================================================

    print_section("QUERY TF-IDF")

    for term, value in explanation["tfidf"].items():

        if value > 0:
            print(f"{term:<20} {value:.6f}")

    # ============================================================
    # DOCUMENT EXPLANATION
    # ============================================================

    document_index = 0

    document = engine.documents[document_index]

    document_explanation = (
        engine.tfidf.explain_document(document_index)
    )

    print_section(
        f"DOCUMENT TF-IDF: {document['title']}"
    )

    for term, value in document_explanation["tfidf"].items():

        if value > 0:
            print(f"{term:<20} {value:.6f}")

    # ============================================================
    # COSINE
    # ============================================================

    query_vector = engine.tfidf.transform_query(query)

    document_vector = engine.document_vectors[
        document_index
    ]

    cosine = explain_cosine_similarity(
        query_vector,
        document_vector,
    )

    print_section("COSINE SIMILARITY")

    print(
        f"Dot Product : "
        f"{cosine['dot_product']:.6f}"
    )

    print(
        f"Query Magnitude : "
        f"{cosine['magnitude_a']:.6f}"
    )

    print(
        f"Document Magnitude : "
        f"{cosine['magnitude_b']:.6f}"
    )

    print(
        f"Cosine Similarity : "
        f"{cosine['cosine_similarity']:.6f}"
    )

    # ============================================================
    # FINAL RANKING
    # ============================================================

    print_section("FINAL RANKING")

    results = engine.search(
        query,
        top_k=5,
    )

    for result in results:

        print(
            f"{result['score']:.6f} | "
            f"{result['title']} | "
            f"{result['url']}"
        )


if __name__ == "__main__":
    main()