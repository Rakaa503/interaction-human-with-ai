import json
from pathlib import Path

from ml.search.engine import SearchEngine


BASE_DIR = Path(__file__).resolve().parent
QUERY_FILE = BASE_DIR / "queries.json"

DOCUMENT_FILE = (
    BASE_DIR.parent
    / "data"
    / "documents.json"
)


def load_queries():
    with open(QUERY_FILE, "r", encoding="utf-8") as file:
        return json.load(file)


def precision_at_k(results, relevant_documents, k):
    top_k = results[:k]

    if not top_k:
        return 0.0

    relevant_count = sum(
        1
        for result in top_k
        if result["title"] in relevant_documents
    )

    return relevant_count / k


def recall_at_k(results, relevant_documents, k):
    if not relevant_documents:
        return 0.0

    top_k = results[:k]

    relevant_count = sum(
        1
        for result in top_k
        if result["title"] in relevant_documents
    )

    return relevant_count / len(relevant_documents)


def reciprocal_rank(results, relevant_documents):
    for index, result in enumerate(results, start=1):
        if result["title"] in relevant_documents:
            return 1 / index

    return 0.0


def main():
    queries = load_queries()

    engine = SearchEngine(str(DOCUMENT_FILE))

    precision_scores = []
    recall_scores = []
    reciprocal_rank_scores = []

    print("=" * 70)
    print("AVIGO SEARCH ENGINE EVALUATION")
    print("=" * 70)

    for item in queries:
        query = item["query"]
        relevant_documents = item["relevant_documents"]

        results = engine.search(query, top_k=5)

        precision = precision_at_k(
            results,
            relevant_documents,
            3,
        )

        recall = recall_at_k(
            results,
            relevant_documents,
            3,
        )

        rr = reciprocal_rank(
            results,
            relevant_documents,
        )

        precision_scores.append(precision)
        recall_scores.append(recall)
        reciprocal_rank_scores.append(rr)

        print()
        print("-" * 70)
        print(f"QUERY      : {query}")
        print(f"RELEVANT   : {relevant_documents}")
        print("-" * 70)

        for rank, result in enumerate(results, start=1):
            print(
                f"{rank}. "
                f"{result['score']:.6f} | "
                f"{result['title']}"
            )

        print()
        print(f"Precision@3 : {precision:.4f}")
        print(f"Recall@3    : {recall:.4f}")
        print(f"RR          : {rr:.4f}")

    mean_precision = (
        sum(precision_scores)
        / len(precision_scores)
    )

    mean_recall = (
        sum(recall_scores)
        / len(recall_scores)
    )

    mrr = (
        sum(reciprocal_rank_scores)
        / len(reciprocal_rank_scores)
    )

    print()
    print("=" * 70)
    print("FINAL EVALUATION")
    print("=" * 70)

    print(f"Precision@3 : {mean_precision:.4f}")
    print(f"Recall@3    : {mean_recall:.4f}")
    print(f"MRR         : {mrr:.4f}")


if __name__ == "__main__":
    main()