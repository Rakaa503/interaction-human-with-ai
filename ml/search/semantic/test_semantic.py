from ml.search.semantic.embedder import (
    SemanticEmbedder,
)

from ml.search.semantic.similarity import (
    cosine_similarity,
    dot_product,
    magnitude,
)


def main():
    embedder = SemanticEmbedder()

    query = "apa itu machine learning"

    documents = [
        "Machine learning adalah teknologi yang memungkinkan komputer belajar dari data.",
        "Artificial intelligence adalah bidang ilmu yang membuat mesin dapat melakukan tugas yang membutuhkan kecerdasan.",
        "Go adalah bahasa pemrograman yang dikembangkan oleh Google.",
        "Computer vision digunakan untuk memahami informasi dari gambar dan video.",
    ]

    query_vector = embedder.encode(
        query
    )

    print("\n==============================")
    print("SEMANTIC SEARCH")
    print("==============================")

    print(
        f"\nQuery: {query}"
    )

    print(
        f"\nEmbedding Dimension: "
        f"{len(query_vector)}"
    )

    print(
        "\nQuery Vector (first 10 values):"
    )

    for index, value in enumerate(
        query_vector[:10]
    ):
        print(
            f"[{index}] {value:.6f}"
        )

    results = []

    for document in documents:
        document_vector = embedder.encode(
            document
        )

        score = cosine_similarity(
            query_vector,
            document_vector,
        )

        results.append(
            {
                "document": document,
                "vector": document_vector,
                "score": score,
            }
        )

    results.sort(
        key=lambda item: item["score"],
        reverse=True,
    )

    print(
        "\n=============================="
    )
    print("SEMANTIC RANKING")
    print(
        "=============================="
    )

    for rank, result in enumerate(
        results,
        start=1,
    ):
        print(
            f"\n#{rank}"
        )

        print(
            f"Score      : "
            f"{result['score']:.6f}"
        )

        print(
            f"Document   : "
            f"{result['document']}"
        )

    best = results[0]

    best_vector = best["vector"]

    print(
        "\n=============================="
    )
    print("COSINE EXPLANATION")
    print(
        "=============================="
    )

    dot = dot_product(
        query_vector,
        best_vector,
    )

    magnitude_query = magnitude(
        query_vector
    )

    magnitude_document = magnitude(
        best_vector
    )

    print(
        f"\nDot Product          : "
        f"{dot:.6f}"
    )

    print(
        f"Query Magnitude      : "
        f"{magnitude_query:.6f}"
    )

    print(
        f"Document Magnitude   : "
        f"{magnitude_document:.6f}"
    )

    print(
        f"Cosine Similarity    : "
        f"{best['score']:.6f}"
    )


if __name__ == "__main__":
    main()