import math


def cosine_similarity(
    vector_a: dict[str, float],
    vector_b: dict[str, float],
) -> float:

    dot_product = sum(
        vector_a.get(term, 0.0) * vector_b.get(term, 0.0)
        for term in vector_a
    )

    magnitude_a = math.sqrt(
        sum(value ** 2 for value in vector_a.values())
    )

    magnitude_b = math.sqrt(
        sum(value ** 2 for value in vector_b.values())
    )

    if magnitude_a == 0 or magnitude_b == 0:
        return 0.0

    return dot_product / (magnitude_a * magnitude_b)