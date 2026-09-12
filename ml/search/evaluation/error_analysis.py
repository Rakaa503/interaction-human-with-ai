import json
from pathlib import Path

from ml.search.hybrid.corpus import KnowledgeCorpus
from ml.search.semantic.embedder import SemanticEmbedder
from ml.search.semantic.similarity import (
    cosine_similarity as semantic_cosine_similarity,
)
from ml.search.similarity import (
    cosine_similarity as tfidf_cosine_similarity,
)
from ml.search.tfidf import TFIDF


ALPHA = 0.25
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


def main():

    queries = load_queries()

    print()
    print("==============================")
    print("AVIGO HYBRID ERROR ANALYSIS")
    print("==============================")
    print()
    print(f"Alpha  : {ALPHA}")
    print(f"Top-K  : {TOP_K}")
    print(f"Queries: {len(queries)}")

    # ==============================
    # CORPUS
    # ==============================

    corpus_provider = KnowledgeCorpus(
        API_URL
    )

    documents = corpus_provider.fetch()

    if not documents:
        raise RuntimeError(
            "Knowledge corpus kosong."
        )

    print(
        f"Documents: {len(documents)}"
    )

    contents = [
        document["text"]
        for document in documents
    ]

    # ==============================
    # TF-IDF
    # ==============================

    print()
    print(
        "Preparing TF-IDF..."
    )

    tfidf = TFIDF(contents)

    document_tfidf_vectors = [
        tfidf.transform_document(index)
        for index in range(len(contents))
    ]

    # ==============================
    # SEMANTIC
    # ==============================

    print(
        "Loading semantic model..."
    )

    embedder = SemanticEmbedder()

    print(
        "Encoding documents..."
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
    # QUERY ANALYSIS
    # ==============================

    print()
    print(
        "Running error analysis..."
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

        # ------------------------------
        # Query Semantic
        # ------------------------------

        query_embedding = embedder.encode(
            query
        )

        results = []

        for index, document in enumerate(
            documents
        ):

            document_id = document["id"]

            tfidf_score = (
                tfidf_cosine_similarity(
                    query_tfidf,
                    document_tfidf_vectors[
                        index
                    ],
                )
            )

            semantic_score = (
                semantic_cosine_similarity(
                    query_embedding,
                    semantic_vectors[
                        document_id
                    ],
                )
            )

            hybrid_score = (
                ALPHA * tfidf_score
                + (1.0 - ALPHA)
                * semantic_score
            )

            results.append(
                {
                    "title": document["title"],
                    "tfidf": tfidf_score,
                    "semantic": semantic_score,
                    "hybrid": hybrid_score,
                    "relevant": (
                        document["title"]
                        in relevant
                    ),
                }
            )

        results.sort(
            key=lambda item: item["hybrid"],
            reverse=True,
        )

        retrieved = results[:TOP_K]

        retrieved_titles = [
            result["title"]
            for result in retrieved
        ]

        # ==============================
        # CHECK RECALL
        # ==============================

        missing = [
            title
            for title in relevant
            if title not in retrieved_titles
        ]

        if not missing:
            continue

        # ==============================
        # PRINT FAILURE
        # ==============================

        print()
        print("=" * 70)
        print("❌ FAILED RECALL@3")
        print("=" * 70)

        print()
        print(
            f"Query     : {query}"
        )

        print(
            f"Relevant  : {relevant}"
        )

        print(
            f"Missing   : {missing}"
        )

        print()
        print("TOP 3 RETRIEVED")
        print("-" * 70)

        for rank, result in enumerate(
            retrieved,
            start=1,
        ):

            marker = (
                "✓"
                if result["relevant"]
                else "✗"
            )

            print(
                f"{rank}. {marker} "
                f"{result['title']}"
            )

            print(
                f"   TF-IDF   : "
                f"{result['tfidf']:.6f}"
            )

            print(
                f"   Semantic : "
                f"{result['semantic']:.6f}"
            )

            print(
                f"   Hybrid   : "
                f"{result['hybrid']:.6f}"
            )

        print()
        print("MISSING RELEVANT DOCUMENTS")
        print("-" * 70)

        for title in missing:

            target = next(
                (
                    result
                    for result in results
                    if result["title"] == title
                ),
                None,
            )

            if target is None:
                continue

            missing_rank = (
                next(
                    (
                        rank
                        for rank, result in enumerate(
                            results,
                            start=1,
                        )
                        if result["title"]
                        == title
                    ),
                    None,
                )
            )

            print(
                f"{title}"
            )

            print(
                f"   Rank     : #{missing_rank}"
            )

            print(
                f"   TF-IDF   : "
                f"{target['tfidf']:.6f}"
            )

            print(
                f"   Semantic : "
                f"{target['semantic']:.6f}"
            )

            print(
                f"   Hybrid   : "
                f"{target['hybrid']:.6f}"
            )

        print()


if __name__ == "__main__":
    main()
