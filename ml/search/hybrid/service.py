from ml.search.hybrid.corpus import KnowledgeCorpus
from ml.search.hybrid.ranker import HybridRanker
from ml.search.semantic.embedder import SemanticEmbedder
from ml.search.semantic.similarity import (
    cosine_similarity as semantic_cosine_similarity,
)
from ml.search.similarity import (
    cosine_similarity as tfidf_cosine_similarity,
)
from ml.search.tfidf import TFIDF


class HybridSearchService:
    def __init__(
        self,
        alpha: float = 0.5,
        api_url: str = "http://127.0.0.1:8080/api/v1/knowledge",
    ):
        self.corpus_provider = KnowledgeCorpus(
            api_url
        )

        self.embedder = SemanticEmbedder()

        self.ranker = HybridRanker(
            alpha=alpha
        )

    def search(
        self,
        query: str,
        top_k: int = 5,
    ) -> list[dict]:

        query = query.strip()

        if not query:
            return []

        documents = self.corpus_provider.fetch()

        if not documents:
            return []

        contents = [
            document["text"]
            for document in documents
        ]

        # ==============================
        # TF-IDF
        # ==============================

        tfidf = TFIDF(contents)

        query_tfidf = tfidf.transform_query(
            query
        )

        document_tfidf_vectors = [
            tfidf.transform_document(index)
            for index in range(len(contents))
        ]

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

        # ==============================
        # SEMANTIC
        # ==============================

        query_embedding = self.embedder.encode(
            query
        )

        semantic_vectors = {}

        for document in documents:
            document_id = document["id"]

            semantic_vectors[document_id] = (
                self.embedder.encode(
                    document["text"]
                )
            )

        # ==============================
        # HYBRID
        # ==============================

        results = self.ranker.rank(
            query_vector=query_embedding,
            documents=documents,
            document_vectors=semantic_vectors,
            tfidf_scores=tfidf_scores,
        )

        return results[:top_k]