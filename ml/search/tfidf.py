import math
from collections import Counter

from .tokenizer import tokenize


class TFIDF:
    def __init__(self, documents: list[str]):
        self.documents = documents

        self.tokenized_documents = [
            tokenize(document)
            for document in documents
        ]

        self.vocabulary = self._build_vocabulary()
        self.idf = self._calculate_idf()

    def _build_vocabulary(self) -> list[str]:
        vocabulary = set()

        for tokens in self.tokenized_documents:
            vocabulary.update(tokens)

        return sorted(vocabulary)

    def _calculate_idf(self) -> dict[str, float]:
        total_documents = len(self.documents)

        idf = {}

        for term in self.vocabulary:
            document_frequency = sum(
                1
                for tokens in self.tokenized_documents
                if term in tokens
            )

            idf[term] = math.log(
                total_documents / document_frequency
            )

        return idf

    def calculate_tf(self, tokens: list[str]) -> dict[str, float]:
        counts = Counter(tokens)

        total_terms = len(tokens)

        if total_terms == 0:
            return {}

        return {
            term: count / total_terms
            for term, count in counts.items()
        }

    def transform_document(
        self,
        document_index: int,
    ) -> dict[str, float]:

        tokens = self.tokenized_documents[document_index]

        tf = self.calculate_tf(tokens)

        return {
            term: tf.get(term, 0.0) * self.idf[term]
            for term in self.vocabulary
        }

    def transform_query(
        self,
        query: str,
    ) -> dict[str, float]:

        tokens = tokenize(query)

        tf = self.calculate_tf(tokens)

        return {
            term: tf.get(term, 0.0) * self.idf.get(term, 0.0)
            for term in self.vocabulary
        }

    def explain_query(self, query: str) -> dict:
        tokens = tokenize(query)

        tf = self.calculate_tf(tokens)

        tfidf = self.transform_query(query)

        return {
            "query": query,
            "tokens": tokens,
            "tf": tf,
            "tfidf": tfidf,
        }

    def explain_document(
        self,
        document_index: int,
    ) -> dict:

        tokens = self.tokenized_documents[document_index]

        tf = self.calculate_tf(tokens)

        tfidf = self.transform_document(document_index)

        df = {
            term: sum(
                1
                for document_tokens in self.tokenized_documents
                if term in document_tokens
            )
            for term in self.vocabulary
        }

        return {
            "tokens": tokens,
            "tf": tf,
            "df": df,
            "idf": self.idf,
            "tfidf": tfidf,
        }