import json
from pathlib import Path

from ml.search.evaluation.metrics import (
    mean,
    precision_at_k,
    recall_at_k,
    reciprocal_rank,
)
from ml.search.hybrid.corpus import KnowledgeCorpus
from ml.search.semantic.embedder import SemanticEmbedder
from ml.search.similarity import (
    cosine_similarity as tfidf_cosine_similarity,
)
from ml.search.tfidf import TFIDF
from ml.search.semantic.similarity import (
    cosine_similarity as semantic_cosine_similarity,
)


ALPHAS = [
    0.00,
    0.05,
    0.10,
    0.15,
    0.20,
    0.25,
    0.30,
    0.35,
    0.40,
    0.45,
    0.50,
    0.55,
    0.60,
    0.65,
    0.70,
    0.75,
    0.80,
    0.85,
    0.90,
    0.95,
    1.00,
]

TOP_K = 3

API_URL = (
    "http://127.0.0.1:8080/api/v1/knowledge"
)


def load_queries() -> list[dict]:
    base_dir = Path(__file__).resolve().parent
    query_path = base_dir / "queries_challenging.json"

    with open(
        query_path,
        "r",
        encoding="utf-8-sig",
    ) as file:
        return json.load(file)


def prepare_scores(
    queries: list[dict],
) -> tuple[list[dict], dict[str, list[dict]]]:

    corpus_provider = KnowledgeCorpus(
        API_URL
    )

    documents = corpus_provider.fetch()

    if not documents:
        raise RuntimeError(
            "Knowledge corpus kosong."
        )

    contents = [
        document["text"]
        for document in documents
    ]

    print()
    print(
        f"Knowledge documents : {len(documents)}"
    )

    # ==============================
    # TF-IDF
    # ==============================

    print(
        "Preparing TF-IDF scores..."
    )

    tfidf = TFIDF(contents)

    document_tfidf_vectors = [
        tfidf.transform_document(index)
        for index in range(len(contents))
    ]

    # ==============================
    # SEMANTIC MODEL
    # ==============================

    print(
        "Loading semantic model once..."
    )

    embedder = SemanticEmbedder()

    print(
        "Encoding document embeddings once..."
    )

    semantic_vectors = {}

    for document in documents:
        document_id = document["id"]

        semantic_vectors[document_id] = (
            embedder.encode(
                document["text"]
            )
        )

    # ==============================
    # QUERY SCORES
    # ==============================

    prepared_queries = {}

    print(
        "Preparing query scores..."
    )

    for item in queries:
        query = item["query"]
        relevant = item["relevant"]

        if not relevant:
            continue

        # ------------------------------
        # Query TF-IDF
        # ------------------------------

        query_tfidf = tfidf.transform_query(
            query
        )

        tfidf_scores = {}

        for index, vector in enumerate(
            document_tfidf_vectors
        ):
            document_id = documents[index]["id"]

            tfidf_scores[document_id] = (
                tfidf_cosine_similarity(
                    query_tfidf,
                    vector,
                )
            )

        # ------------------------------
        # Query Semantic
        # ------------------------------

        query_embedding = embedder.encode(
            query
        )

        semantic_scores = {}

        for document in documents:
            document_id = document["id"]

            semantic_scores[document_id] = (
                semantic_cosine_similarity(
                    query_embedding,
                    semantic_vectors[
                        document_id
                    ],
                )
            )

        prepared_queries[query] = {
            "relevant": relevant,
            "tfidf_scores": tfidf_scores,
            "semantic_scores": semantic_scores,
        }

    return documents, prepared_queries


def rank_for_alpha(
    alpha: float,
    documents: list[dict],
    prepared_query: dict,
) -> list[dict]:

    tfidf_scores = prepared_query[
        "tfidf_scores"
    ]

    semantic_scores = prepared_query[
        "semantic_scores"
    ]

    results = []

    for document in documents:
        document_id = document["id"]

        tfidf_score = tfidf_scores.get(
            document_id,
            0.0,
        )

        semantic_score = semantic_scores.get(
            document_id,
            0.0,
        )

        hybrid_score = (
            alpha * tfidf_score
            + (1.0 - alpha)
            * semantic_score
        )

        results.append(
            {
                **document,
                "tfidf_score": tfidf_score,
                "semantic_score": semantic_score,
                "hybrid_score": hybrid_score,
            }
        )

    results.sort(
        key=lambda item: item["hybrid_score"],
        reverse=True,
    )

    return results[:TOP_K]


def evaluate_alpha(
    alpha: float,
    documents: list[dict],
    prepared_queries: dict[str, dict],
) -> dict:

    precision_scores = []
    recall_scores = []
    reciprocal_ranks = []

    for query, prepared in (
        prepared_queries.items()
    ):
        relevant = prepared["relevant"]

        results = rank_for_alpha(
            alpha=alpha,
            documents=documents,
            prepared_query=prepared,
        )

        retrieved_titles = [
            result["title"]
            for result in results
        ]

        precision_scores.append(
            precision_at_k(
                retrieved_titles,
                relevant,
                TOP_K,
            )
        )

        recall_scores.append(
            recall_at_k(
                retrieved_titles,
                relevant,
                TOP_K,
            )
        )

        reciprocal_ranks.append(
            reciprocal_rank(
                retrieved_titles,
                relevant,
            )
        )

    return {
        "alpha": alpha,
        "precision_at_3": mean(
            precision_scores
        ),
        "recall_at_3": mean(
            recall_scores
        ),
        "mrr": mean(
            reciprocal_ranks
        ),
    }


def main():

    queries = load_queries()

    print()
    print("==============================")
    print("AVIGO ALPHA FINE-TUNING")
    print("==============================")
    print()
    print(f"Top-K  : {TOP_K}")
    print(f"Queries: {len(queries)}")
    print()
    print("Testing alpha values:")
    print(ALPHAS)

    # ==============================
    # PREPARE ONCE
    # ==============================

    documents, prepared_queries = (
        prepare_scores(
            queries
        )
    )

    print()
    print(
        "Score preparation completed."
    )

    # ==============================
    # ALPHA EXPERIMENT
    # ==============================

    results = []

    for alpha in ALPHAS:

        print()
        print(
            f"Testing alpha = {alpha:.2f} ..."
        )

        result = evaluate_alpha(
            alpha=alpha,
            documents=documents,
            prepared_queries=prepared_queries,
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

    # ==============================
    # FINAL RESULTS
    # ==============================

    print()
    print("==============================")
    print("FINE-TUNING RESULTS")
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

    # ==============================
    # BEST ALPHA
    # ==============================

    best_mrr = max(
        result["mrr"]
        for result in results
    )

    best_results = [
        result
        for result in results
        if result["mrr"] == best_mrr
    ]

    print()
    print("==============================")
    print("BEST ALPHA CANDIDATES")
    print("==============================")
    print()

    print(
        f"Best MRR: {best_mrr:.4f}"
    )

    for result in best_results:

        print(
            f"Alpha {result['alpha']:.2f} "
            f"-> "
            f"P@3={result['precision_at_3']:.4f}, "
            f"R@3={result['recall_at_3']:.4f}, "
            f"MRR={result['mrr']:.4f}"
        )


if __name__ == "__main__":
    main()
