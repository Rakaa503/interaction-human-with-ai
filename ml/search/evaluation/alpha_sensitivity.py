import json
from pathlib import Path

from ml.search.evaluation.metrics import (
    mean,
    precision_at_k,
    recall_at_k,
    reciprocal_rank,
)
from ml.search.hybrid.service import HybridSearchService


ALPHAS = [
    0.0,
    0.25,
    0.5,
    0.75,
    1.0,
]

TOP_K = 3


def load_queries() -> list[dict]:
    base_dir = Path(__file__).resolve().parent
    query_path = base_dir / "queries_challenging.json"

    with open(
        query_path,
        "r",
        encoding="utf-8-sig",
    ) as file:
        return json.load(file)


def evaluate_alpha(
    alpha: float,
    queries: list[dict],
) -> dict:

    service = HybridSearchService(alpha=alpha)

    precision_scores = []
    recall_scores = []
    reciprocal_ranks = []

    for item in queries:
        query = item["query"]
        relevant = item["relevant"]

        results = service.search(
            query=query,
            top_k=TOP_K,
        )

        retrieved_titles = [
            result["title"]
            for result in results
        ]

        if not relevant:
            continue

        precision = precision_at_k(
            retrieved_titles,
            relevant,
            TOP_K,
        )

        recall = recall_at_k(
            retrieved_titles,
            relevant,
            TOP_K,
        )

        rr = reciprocal_rank(
            retrieved_titles,
            relevant,
        )

        precision_scores.append(precision)
        recall_scores.append(recall)
        reciprocal_ranks.append(rr)

    return {
        "alpha": alpha,
        "precision_at_3": mean(precision_scores),
        "recall_at_3": mean(recall_scores),
        "mrr": mean(reciprocal_ranks),
    }


def main():
    queries = load_queries()

    print()
    print("==============================")
    print("AVIGO ALPHA SENSITIVITY")
    print("==============================")
    print()
    print(f"Top-K  : {TOP_K}")
    print(f"Queries: {len(queries)}")

    results = []

    for alpha in ALPHAS:
        print()
        print(f"Testing alpha = {alpha:.2f} ...")

        result = evaluate_alpha(
            alpha=alpha,
            queries=queries,
        )

        results.append(result)

        print(
            f"Precision@3 : "
            f"{result['precision_at_3']:.4f}"
        )

        print(
            f"Recall@3    : "
            f"{result['recall_at_3']:.4f}"
        )

        print(
            f"MRR         : "
            f"{result['mrr']:.4f}"
        )

    print()
    print("==============================")
    print("ALPHA COMPARISON")
    print("==============================")
    print()

    print(
        f"{'Alpha':<10}"
        f"{'Precision@3':<18}"
        f"{'Recall@3':<15}"
        f"{'MRR':<10}"
    )

    print("-" * 53)

    for result in results:
        print(
            f"{result['alpha']:<10.2f}"
            f"{result['precision_at_3']:<18.4f}"
            f"{result['recall_at_3']:<15.4f}"
            f"{result['mrr']:<10.4f}"
        )

    best_mrr = max(
        results,
        key=lambda result: result["mrr"],
    )

    print()
    print("==============================")
    print("BEST ALPHA")
    print("==============================")
    print()
    print(
        f"Alpha : {best_mrr['alpha']:.2f}"
    )
    print(
        f"MRR   : {best_mrr['mrr']:.4f}"
    )


if __name__ == "__main__":
    main()
