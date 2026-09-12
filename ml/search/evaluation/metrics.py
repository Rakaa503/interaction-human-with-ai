def precision_at_k(retrieved_documents, relevant_documents, k):
    retrieved = retrieved_documents[:k]

    if not retrieved:
        return 0.0

    relevant_count = sum(
        document in relevant_documents
        for document in retrieved
    )

    return relevant_count / len(retrieved)


def recall_at_k(retrieved_documents, relevant_documents, k):
    if not relevant_documents:
        return None

    retrieved = retrieved_documents[:k]

    relevant_count = sum(
        document in relevant_documents
        for document in retrieved
    )

    return relevant_count / len(relevant_documents)


def reciprocal_rank(retrieved_documents, relevant_documents):
    if not relevant_documents:
        return None

    for rank, document in enumerate(
        retrieved_documents,
        start=1,
    ):
        if document in relevant_documents:
            return 1.0 / rank

    return 0.0


def mean(values):
    valid_values = [
        value
        for value in values
        if value is not None
    ]

    if not valid_values:
        return 0.0

    return sum(valid_values) / len(valid_values)
