import json
from pathlib import Path

import joblib
import pandas as pd

from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import (
    accuracy_score,
    classification_report,
)
from sklearn.model_selection import train_test_split
from sklearn.pipeline import Pipeline


BASE_DIR = Path(__file__).resolve().parent.parent

DATASET_PATH = BASE_DIR / "dataset" / "interactions.csv"
MODEL_DIR = BASE_DIR / "models"
MODEL_PATH = MODEL_DIR / "intent_model.joblib"
METRICS_PATH = BASE_DIR / "evaluation" / "intent_metrics.json"


TEST_SIZE = 0.20
RANDOM_STATE = 42


def load_dataset() -> pd.DataFrame:
    if not DATASET_PATH.exists():
        raise FileNotFoundError(
            f"Dataset tidak ditemukan: {DATASET_PATH}"
        )

    df = pd.read_csv(DATASET_PATH)

    required_columns = {
        "text",
        "intent",
    }

    missing_columns = required_columns - set(df.columns)

    if missing_columns:
        raise ValueError(
            "Kolom dataset tidak lengkap: "
            + ", ".join(sorted(missing_columns))
        )

    df = df.dropna(subset=["text", "intent"])

    df["text"] = df["text"].astype(str).str.strip()
    df["intent"] = df["intent"].astype(str).str.strip()

    df = df[
        (df["text"] != "")
        & (df["intent"] != "")
    ]

    if df.empty:
        raise ValueError("Dataset kosong setelah preprocessing.")

    return df


def build_model() -> Pipeline:
    return Pipeline(
        [
            (
                "tfidf",
                TfidfVectorizer(
                    lowercase=True,
                    ngram_range=(1, 2),
                ),
            ),
            (
                "classifier",
                LogisticRegression(
                    max_iter=1000,
                    random_state=RANDOM_STATE,
                ),
            ),
        ]
    )


def main() -> None:
    print("=== AVIGO ML TRAINING ===")
    print()

    df = load_dataset()

    print(f"Dataset : {DATASET_PATH}")
    print(f"Samples : {len(df)}")
    print()

    print("=== INTENT DISTRIBUTION ===")

    distribution = df["intent"].value_counts()

    for intent, count in distribution.items():
        print(f"{intent}: {count}")

    print()

    X = df["text"]
    y = df["intent"]

    X_train, X_test, y_train, y_test = train_test_split(
        X,
        y,
        test_size=TEST_SIZE,
        random_state=RANDOM_STATE,
        stratify=y,
    )

    print(f"Training samples : {len(X_train)}")
    print(f"Testing samples  : {len(X_test)}")
    print()

    model = build_model()

    print("Training model...")

    model.fit(
        X_train,
        y_train,
    )

    print("Training completed.")
    print()

    predictions = model.predict(X_test)

    accuracy = accuracy_score(
        y_test,
        predictions,
    )

    print("=== EVALUATION ===")
    print()
    print(f"Accuracy: {accuracy:.4f}")
    print()

    report = classification_report(
        y_test,
        predictions,
        output_dict=True,
        zero_division=0,
    )

    print(
        classification_report(
            y_test,
            predictions,
            zero_division=0,
        )
    )

    MODEL_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    METRICS_PATH.parent.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        model,
        MODEL_PATH,
    )

    metrics = {
        "model": "TF-IDF + Logistic Regression",
        "target": "intent",
        "samples": len(df),
        "training_samples": len(X_train),
        "testing_samples": len(X_test),
        "accuracy": accuracy,
        "classification_report": report,
    }

    with open(
        METRICS_PATH,
        "w",
        encoding="utf-8",
    ) as file:
        json.dump(
            metrics,
            file,
            indent=2,
        )

    print("=== OUTPUT ===")
    print(f"Model   : {MODEL_PATH}")
    print(f"Metrics : {METRICS_PATH}")
    print()
    print("AVIGO ML training selesai.")


if __name__ == "__main__":
    main()