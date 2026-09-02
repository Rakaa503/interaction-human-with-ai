from .similarity import cosine_similarity


def rank_documents(
    query_vector: dict[str, float],
    document_vectors: list[dict[str, float]],
) -> list[tuple[int, float]]:

    scores = []

    for index, document_vector in enumerate(document_vectors):
        score = cosine_similarity(
            query_vector,
            document_vector,
        )

        scores.append((index, score))

    scores.sort(
        key=lambda item: item[1],
        reverse=True,
    )

    return scores