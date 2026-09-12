from ml.search.semantic.knowledge_retriever import KnowledgeRetriever


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
    print(f"\nKnowledge Documents Found: {len(results)}")

    print("\n==============================")
    print("RANKING")
    print("==============================")

    for rank, result in enumerate(results, start=1):
        print(f"\n#{rank}")
        print(f"Score    : {result['score']:.6f}")
        print(f"ID       : {result['id']}")
        print(f"Title    : {result['title']}")
        print(f"Category : {result['category']}")
        print(f"Source   : {result['source']}")
        print(f"Content  : {result['content']}")
        print(f"URL      : {result['url']}")


if __name__ == "__main__":
    main()
