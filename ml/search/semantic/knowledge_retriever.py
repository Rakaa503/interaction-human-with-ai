import requests

from ml.search.semantic.embedder import SemanticEmbedder
from ml.search.semantic.similarity import cosine_similarity


class KnowledgeRetriever:
    def __init__(
        self,
        api_url: str = "http://127.0.0.1:8080/api/v1/knowledge",
    ):
        self.api_url = api_url
        self.embedder = SemanticEmbedder()

    def fetch_documents(self) -> list[dict]:
        response = requests.get(
            self.api_url,
            timeout=10,
        )

        response.raise_for_status()

        payload = response.json()

        return payload.get("data", [])

    def search(
        self,
        query: str,
        top_k: int = 5,
    ) -> list[dict]:

        query = query.strip()

        if not query:
            return []

        documents = self.fetch_documents()

        if not documents:
            return []

        query_vector = self.embedder.encode(query)

        results = []

        for document in documents:
            text = (
                f"{document.get('Title', '')}. "
                f"{document.get('Content', '')}"
            ).strip()

            if not text:
                continue

            document_vector = self.embedder.encode(text)

            score = cosine_similarity(
                query_vector,
                document_vector,
            )

            results.append(
                {
                    "id": document.get("ID"),
                    "title": document.get("Title"),
                    "content": document.get("Content"),
                    "url": document.get("URL"),
                    "source": document.get("Source"),
                    "category": document.get("Category"),
                    "score": score,
                }
            )

        results.sort(
            key=lambda item: item["score"],
            reverse=True,
        )

        return results[:top_k]