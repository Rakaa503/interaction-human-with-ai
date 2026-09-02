import json
from pathlib import Path

from .tfidf import TFIDF
from .ranker import rank_documents


class SearchEngine:
    def __init__(self, data_path: str):
        self.data_path = Path(data_path)

        with open(
            self.data_path,
            "r",
            encoding="utf-8",
        ) as file:
            self.documents = json.load(file)

        contents = [
            document["content"]
            for document in self.documents
        ]

        self.tfidf = TFIDF(contents)

        self.document_vectors = [
            self.tfidf.transform_document(index)
            for index in range(len(contents))
        ]

    def search(
        self,
        query: str,
        top_k: int = 5,
    ) -> list[dict]:

        query_vector = self.tfidf.transform_query(query)

        ranked = rank_documents(
            query_vector,
            self.document_vectors,
        )

        results = []

        for index, score in ranked[:top_k]:

            document = self.documents[index]

            results.append({
                "id": document["id"],
                "title": document["title"],
                "url": document["url"],
                "score": round(score, 6),
            })

        return results