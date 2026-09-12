from ml.search.semantic.knowledge_retriever import (
    KnowledgeRetriever,
)


def main():
    retriever = KnowledgeRetriever()

    query = "apa itu machine learning"

    results = retriever.search(
        query,
        top_k=5,
    )

    print("\n==============================")
    print("SEMANTIC KNOWLEDGE RETRIEVAL")
    print("==============================")

    print(f"\nQuery: {query}")

    print(
        f"\nKnowledge Documents Found: "
        f"{len(results)}"
    )

    print("\n==============================")
    print("RANKING")
    print("==============================")

    for rank, result in enumerate(
        results,
        start=1,
    ):
        print(f"\n#{rank}")

        print(
            f"Score    : "
            f"{result['score']:.6f}"
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