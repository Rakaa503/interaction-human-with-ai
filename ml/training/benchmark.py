from pathlib import Path

import joblib
import pandas as pd

from sklearn.metrics import (
    accuracy_score,
    classification_report,
    f1_score,
    precision_score,
    recall_score,
)


BASE_DIR = Path(__file__).resolve().parent.parent

DATASET_PATH = BASE_DIR / "dataset" / "interactions.csv"
MODEL_PATH = BASE_DIR / "models" / "intent_model.joblib"


def rule_based_predict(text: str) -> str:
    text = text.lower()

    # Technical/problem solving
    if "error" in text or "bug" in text:
        return "problem_solving"

    # Greeting
    if "halo" in text or "hai" in text:
        return "greeting"

    # Question
    question_words = [
        "bagaimana",
        "apa",
        "kenapa",
        "mengapa",
        "apakah",
    ]

    if any(
        word in text
        for word in question_words
    ):
        return "question"

    # Request
    request_words = [
        "tolong",
        "bantu",
        "buatkan",
        "buat",
    ]

    if any(
        word in text
        for word in request_words
    ):
        return "request"

    return "general"


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

    missing = (
        required_columns
        - set(df.columns)
    )

    if missing:
        raise ValueError(
            "Kolom dataset tidak lengkap: "
            + ", ".join(sorted(missing))
        )

    df = df.dropna(
        subset=["text", "intent"]
    )

    return df


def evaluate(
    name: str,
    y_true,
    y_pred,
):
    accuracy = accuracy_score(
        y_true,
        y_pred,
    )

    precision = precision_score(
        y_true,
        y_pred,
        average="macro",
        zero_division=0,
    )

    recall = recall_score(
        y_true,
        y_pred,
        average="macro",
        zero_division=0,
    )

    f1 = f1_score(
        y_true,
        y_pred,
        average="macro",
        zero_division=0,
    )

    print(f"=== {name} ===")
    print()
    print(f"Accuracy : {accuracy:.4f}")
    print(f"Precision: {precision:.4f}")
    print(f"Recall   : {recall:.4f}")
    print(f"F1 Score : {f1:.4f}")
    print()

    return {
        "accuracy": accuracy,
        "precision": precision,
        "recall": recall,
        "f1": f1,
    }


def main():
    print("=== AVIGO MODEL BENCHMARK ===")
    print()

    df = load_dataset()

    print(f"Dataset samples: {len(df)}")
    print()

    X = df["text"]
    y_true = df["intent"]

    # =========================
    # RULE-BASED
    # =========================

    rule_predictions = [
        rule_based_predict(text)
        for text in X
    ]

    rule_metrics = evaluate(
        "RULE-BASED ANALYZER",
        y_true,
        rule_predictions,
    )

    # =========================
    # MACHINE LEARNING
    # =========================

    if not MODEL_PATH.exists():
        raise FileNotFoundError(
            f"Model tidak ditemukan: {MODEL_PATH}"
        )

    model = joblib.load(
        MODEL_PATH
    )

    ml_predictions = model.predict(X)

    ml_metrics = evaluate(
        "ML ANALYZER",
        y_true,
        ml_predictions,
    )

    # =========================
    # COMPARISON
    # =========================

    print("=== COMPARISON ===")
    print()

    print(
        f"Accuracy : "
        f"Rule-Based={rule_metrics['accuracy']:.4f} | "
        f"ML={ml_metrics['accuracy']:.4f}"
    )

    print(
        f"Precision: "
        f"Rule-Based={rule_metrics['precision']:.4f} | "
        f"ML={ml_metrics['precision']:.4f}"
    )

    print(
        f"Recall   : "
        f"Rule-Based={rule_metrics['recall']:.4f} | "
        f"ML={ml_metrics['recall']:.4f}"
    )

    print(
        f"F1 Score : "
        f"Rule-Based={rule_metrics['f1']:.4f} | "
        f"ML={ml_metrics['f1']:.4f}"
    )

    print()

    print("=== ML CLASSIFICATION REPORT ===")
    print()

    print(
        classification_report(
            y_true,
            ml_predictions,
            zero_division=0,
        )
    )

    print(
        "Benchmark selesai."
    )


if __name__ == "__main__":
    main()