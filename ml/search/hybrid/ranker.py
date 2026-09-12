from ml.search.semantic.similarity import cosine_similarity


class HybridRanker:
    def __init__(
        self,
        alpha: float = 0.5,
    ):
        if not 0.0 <= alpha <= 1.0:
            raise ValueError(
                "alpha harus berada di antara 0 dan 1"
            )

        self.alpha = alpha

    def combine(
        self,
        tfidf_score: float,
        semantic_score: float,
    ) -> float:
        return (
            self.alpha * tfidf_score
            + (1.0 - self.alpha) * semantic_score
        )

    def rank(
        self,
        query_vector,
        documents: list[dict],
        document_vectors: dict[int, object],
        tfidf_scores: dict[int, float],
    ) -> list[dict]:

        results = []

        for document in documents:
            document_id = document["id"]

            semantic_score = cosine_similarity(
                query_vector,
                document_vectors[document_id],
            )

            tfidf_score = tfidf_scores.get(
                document_id,
                0.0,
            )

            hybrid_score = self.combine(
                tfidf_score,
                semantic_score,
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

        return results