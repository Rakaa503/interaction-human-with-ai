
from pathlib import Path
import sys

import joblib
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field


BASE_DIR = Path(__file__).resolve().parent.parent
PROJECT_ROOT = BASE_DIR.parent
MODEL_PATH = BASE_DIR / "models" / "intent_model.joblib"

if str(PROJECT_ROOT) not in sys.path:
    sys.path.insert(0, str(PROJECT_ROOT))


app = FastAPI(
    title="AVIGO NLP and Semantic Service",
    version="1.1.0",
)

model = None
retriever = None


class PredictionRequest(BaseModel):
    input: str


class PredictionResponse(BaseModel):
    intent: str
    confidence: float


class SemanticSearchRequest(BaseModel):
    query: str = Field(min_length=1, max_length=2000)
    top_k: int = Field(default=5, ge=1, le=10)


@app.on_event("startup")
def load_services():
    global model, retriever

    if not MODEL_PATH.exists():
        raise RuntimeError(
            f"Model tidak ditemukan: {MODEL_PATH}"
        )

    model = joblib.load(MODEL_PATH)
    print(f"NLP model loaded: {MODEL_PATH}")

    try:
        from ml.search.semantic.knowledge_retriever import (
            KnowledgeRetriever,
        )

        retriever = KnowledgeRetriever()
        print("AVIGO semantic retriever initialized")
    except Exception as exc:
        # NLP prediction remains available if semantic search
        # cannot initialize. Semantic requests return HTTP 503.
        retriever = None
        print(f"Semantic retriever unavailable: {exc}")


@app.get("/health")
def health():
    return {
        "success": True,
        "service": "avigo-nlp-semantic",
        "status": "ok",
        "semantic_ready": retriever is not None,
    }


@app.post(
    "/predict",
    response_model=PredictionResponse,
)
def predict(request: PredictionRequest):
    if model is None:
        raise HTTPException(
            status_code=503,
            detail="Model NLP belum tersedia",
        )

    text = request.input.strip()

    if not text:
        raise HTTPException(
            status_code=400,
            detail="input tidak boleh kosong",
        )

    prediction = model.predict([text])[0]
    probabilities = model.predict_proba([text])[0]

    return PredictionResponse(
        intent=str(prediction),
        confidence=float(probabilities.max()),
    )


@app.post("/semantic/search")
def semantic_search(request: SemanticSearchRequest):
    if retriever is None:
        raise HTTPException(
            status_code=503,
            detail=(
                "Semantic retriever belum tersedia. "
                "Periksa dependency Python dan koneksi "
                "ke knowledge API Go."
            ),
        )

    query = request.query.strip()

    if not query:
        raise HTTPException(
            status_code=400,
            detail="query tidak boleh kosong",
        )

    try:
        results = retriever.search(
            query=query,
            top_k=request.top_k,
        )
    except Exception as exc:
        raise HTTPException(
            status_code=502,
            detail=f"Semantic retrieval gagal: {exc}",
        ) from exc

    return {
        "success": True,
        "query": query,
        "count": len(results),
        "data": results,
    }
