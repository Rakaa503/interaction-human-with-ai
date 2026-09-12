from sentence_transformers import SentenceTransformer


class SemanticEmbedder:
    def __init__(
        self,
        model_name: str = "paraphrase-multilingual-MiniLM-L12-v2",
    ):
        self.model = SentenceTransformer(model_name)

    def encode(self, text: str):
        return self.model.encode(
            text,
            normalize_embeddings=False,
        )