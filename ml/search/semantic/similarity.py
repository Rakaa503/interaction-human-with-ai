import math


def dot_product(
    vector_a,
    vector_b,
) -> float:
    return sum(
        a * b
        for a, b in zip(vector_a, vector_b)
    )


def magnitude(vector) -> float:
    return math.sqrt(
        sum(
            value * value
            for value in vector
        )
    )


def cosine_similarity(
    vector_a,
    vector_b,
) -> float:
    dot = dot_product(
        vector_a,
        vector_b,
    )

    magnitude_a = magnitude(
        vector_a
    )

    magnitude_b = magnitude(
        vector_b
    )

    if magnitude_a == 0 or magnitude_b == 0:
        return 0.0

    return dot / (
        magnitude_a * magnitude_b
    )