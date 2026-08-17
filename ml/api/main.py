from pathlib import Path

import joblib
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel


BASE_DIR = Path(__file__).resolve().parent.parent
MODEL_PATH = BASE_DIR / "models" / "intent_model.joblib"


app = FastAPI(
    title="AVIGO ML Service",
    version="1.0.0",
)


model = None


class PredictionRequest(BaseModel):
    input: str


class PredictionResponse(BaseModel):
    intent: str
    confidence: float


@app.on_event("startup")
def load_model():
    global model

    if not MODEL_PATH.exists():
        raise RuntimeError(
            f"Model tidak ditemukan: {MODEL_PATH}"
        )

    model = joblib.load(MODEL_PATH)

    print(
        f"AVIGO ML model loaded: {MODEL_PATH}"
    )


@app.get("/health")
def health():
    return {
        "success": True,
        "service": "avigo-ml",
        "status": "ok",
    }


@app.post(
    "/predict",
    response_model=PredictionResponse,
)
def predict(
    request: PredictionRequest,
):
    if model is None:
        raise HTTPException(
            status_code=503,
            detail="ML model belum tersedia",
        )

    text = request.input.strip()

    if not text:
        raise HTTPException(
            status_code=400,
            detail="input tidak boleh kosong",
        )

    prediction = model.predict(
        [text]
    )[0]

    probabilities = model.predict_proba(
        [text]
    )[0]

    confidence = float(
        probabilities.max()
    )

    return PredictionResponse(
        intent=str(prediction),
        confidence=confidence,
    )