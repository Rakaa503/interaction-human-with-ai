import json
from pathlib import Path

from ml.search.hybrid.service import HybridSearchService
from ml.search.evaluation.metrics import (
    precision_at_k,
    recall_at_k,
    reciprocal_rank,
    mean,
)


def main():
    base_dir = Path(__file__).resolve().parent
    query_path = base_dir / "queries.json"

    with open(
        query_path,
        "r",
        encoding="utf-8-sig",
    ) as file:
        queries = json.load(file)

    service = HybridSearchService(alpha=0.5)

    k = 3

    precision_scores = []
    recall_scores = []
    reciprocal_ranks = []

    print()
    print("==============================")
    print("AVIGO HYBRID SEARCH EVALUATION")
    print("==============================")
    print()
    print(f"Alpha : {service.ranker.alpha}")
    print(f"Top-K : {k}")
    print(f"Queries : {len(queries)}")

    for index, item in enumerate(
        queries,
        start=1,
    ):
        query = item["query"]
        relevant = item["relevant"]

        results = service.search(
            query=query,
            top_k=k,
        )

        retrieved_titles = [
            result["title"]
            for result in results
        ]

        precision = precision_at_k(
            retrieved_titles,
            relevant,
            k,
        )

        recall = recall_at_k(
            retrieved_titles,
            relevant,
            k,
        )

        rr = reciprocal_rank(
            retrieved_titles,
            relevant,
        )

        if relevant:
            precision_scores.append(precision)
            recall_scores.append(recall)
            reciprocal_ranks.append(rr)

        print()
        print(f"[{index}] {query}")
        print(f"Relevant : {relevant}")
        print(f"Retrieved: {retrieved_titles}")
        print(f"P@{k}    : {precision:.4f}")

        if recall is not None:
            print(f"R@{k}    : {recall:.4f}")

        if rr is not None:
            print(f"RR      : {rr:.4f}")

    print()
    print("==============================")
    print("FINAL METRICS")
    print("==============================")

    print(
        f"Precision@{k} : "
        f"{mean(precision_scores):.4f}"
    )

    print(
        f"Recall@{k}    : "
        f"{mean(recall_scores):.4f}"
    )

    print(
        f"MRR           : "
        f"{mean(reciprocal_ranks):.4f}"
    )

    print()
    print("==============================")
    print("UNKNOWN QUERIES")
    print("==============================")

    print(
        "Known queries   : "
        f"{len(precision_scores)}"
    )

    print(
        "Unknown queries : "
        f"{len(queries) - len(precision_scores)}"
    )


if __name__ == "__main__":
    main()
