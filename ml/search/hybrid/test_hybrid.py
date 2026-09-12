from ml.search.hybrid.service import HybridSearchService


def main():
    service = HybridSearchService(
        alpha=0.5,
    )

    query = "apa itu machine learning"

    results = service.search(
        query,
        top_k=5,
    )

    print("\n==============================")
    print("HYBRID SEARCH")
    print("==============================")

    print(f"\nQuery    : {query}")
    print(f"Alpha    : {service.ranker.alpha}")
    print(
        f"Documents: {len(results)}"
    )

    print("\n==============================")
    print("HYBRID RANKING")
    print("==============================")

    for rank, result in enumerate(
        results,
        start=1,
    ):
        print(f"\n#{rank}")

        print(
            f"Hybrid   : "
            f"{result['hybrid_score']:.6f}"
        )

        print(
            f"TF-IDF   : "
            f"{result['tfidf_score']:.6f}"
        )

        print(
            f"Semantic : "
            f"{result['semantic_score']:.6f}"
        )

        print(
            f"ID       : "
            f"{result['id']}"
        )

        print(
            f"Title    : "
            f"{result['title']}"
        )

        print(
            f"Category : "
            f"{result['category']}"
        )

        print(
            f"Source   : "
            f"{result['source']}"
        )

        print(
            f"Content  : "
            f"{result['content']}"
        )

        print(
            f"URL      : "
            f"{result['url']}"
        )


if __name__ == "__main__":
    main()