import requests


class KnowledgeCorpus:
    def __init__(
        self,
        api_url: str = "http://127.0.0.1:8080/api/v1/knowledge",
    ):
        self.api_url = api_url

    def fetch(self) -> list[dict]:
        response = requests.get(
            self.api_url,
            timeout=10,
        )

        response.raise_for_status()

        payload = response.json()

        documents = payload.get("data", [])

        corpus = []

        for document in documents:
            document_id = document.get("ID")
            title = document.get("Title", "")
            content = document.get("Content", "")

            if document_id is None:
                continue

            text = f"{title}. {content}".strip()

            if not text:
                continue

            corpus.append(
                {
                    "id": document_id,
                    "title": title,
                    "content": content,
                    "url": document.get("URL"),
                    "source": document.get("Source"),
                    "category": document.get("Category"),
                    "text": text,
                }
            )

        return corpus