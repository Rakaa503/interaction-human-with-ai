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


def explain_cosine_similarity(
    vector_a: dict[str, float],
    vector_b: dict[str, float],
) -> dict:

    products = {
        term: vector_a.get(term, 0.0)
        * vector_b.get(term, 0.0)
        for term in vector_a
    }

    dot_product = sum(products.values())

    magnitude_a = math.sqrt(
        sum(value ** 2 for value in vector_a.values())
    )

    magnitude_b = math.sqrt(
        sum(value ** 2 for value in vector_b.values())
    )

    if magnitude_a == 0 or magnitude_b == 0:
        similarity = 0.0
    else:
        similarity = (
            dot_product
            / (magnitude_a * magnitude_b)
        )

    return {
        "products": products,
        "dot_product": dot_product,
        "magnitude_a": magnitude_a,
        "magnitude_b": magnitude_b,
        "cosine_similarity": similarity,
    }