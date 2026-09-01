from pathlib import Path

import joblib
import pandas as pd

from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import (
    accuracy_score,
    classification_report,
    f1_score,
    precision_score,
    recall_score,
)
from sklearn.model_selection import train_test_split
from sklearn.pipeline import Pipeline


BASE_DIR = Path(__file__).resolve().parent.parent

DATASET_PATH = BASE_DIR / "dataset" / "interactions.csv"
MODEL_PATH = BASE_DIR / "models" / "intent_model.joblib"

TEST_SIZE = 0.20
RANDOM_STATE = 42


def rule_based_predict(text: str) -> str:
    text = text.lower()

    if "error" in text or "bug" in text:
        return "problem_solving"

    if "halo" in text or "hai" in text:
        return "greeting"

    question_words = [
        "bagaimana",
        "apa",
        "kenapa",
        "mengapa",
        "apakah",
    ]

    if any(word in text for word in question_words):
        return "question"

    request_words = [
        "tolong",
        "bantu",
        "buatkan",
        "buat",
    ]

    if any(word in text for word in request_words):
        return "request"

    return "general"


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


def calculate_metrics(y_true, predictions):
    return {
        "accuracy": accuracy_score(
            y_true,
            predictions,
        ),
        "precision": precision_score(
            y_true,
            predictions,
            average="macro",
            zero_division=0,
        ),
        "recall": recall_score(
            y_true,
            predictions,
            average="macro",
            zero_division=0,
        ),
        "f1": f1_score(
            y_true,
            predictions,
            average="macro",
            zero_division=0,
        ),
    }


def print_metrics(name, metrics):
    print(f"=== {name} ===")
    print()
    print(
        f"Accuracy : {metrics['accuracy']:.4f}"
    )
    print(
        f"Precision: {metrics['precision']:.4f}"
    )
    print(
        f"Recall   : {metrics['recall']:.4f}"
    )
    print(
        f"F1 Score : {metrics['f1']:.4f}"
    )
    print()


def main():
    print("=== AVIGO FAIR HELD-OUT BENCHMARK ===")
    print()

    if not DATASET_PATH.exists():
        raise FileNotFoundError(
            f"Dataset tidak ditemukan: {DATASET_PATH}"
        )

    df = pd.read_csv(DATASET_PATH)

    df = df.dropna(
        subset=["text", "intent"]
    )

    X = df["text"].astype(str)
    y = df["intent"].astype(str)

    X_train, X_test, y_train, y_test = train_test_split(
        X,
        y,
        test_size=TEST_SIZE,
        random_state=RANDOM_STATE,
        stratify=y,
    )

    print(f"Total dataset : {len(df)}")
    print(f"Training data : {len(X_train)}")
    print(f"Test data     : {len(X_test)}")
    print()

    # =========================
    # TRAIN ML ONLY ON TRAIN SET
    # =========================

    ml_model = build_model()

    print("Training ML model...")
    ml_model.fit(
        X_train,
        y_train,
    )

    print("Training completed.")
    print()

    # =========================
    # RULE-BASED ON HELD-OUT SET
    # =========================

    rule_predictions = [
        rule_based_predict(text)
        for text in X_test
    ]

    rule_metrics = calculate_metrics(
        y_test,
        rule_predictions,
    )

    # =========================
    # ML ON HELD-OUT SET
    # =========================

    ml_predictions = ml_model.predict(
        X_test
    )

    ml_metrics = calculate_metrics(
        y_test,
        ml_predictions,
    )

    # =========================
    # RESULTS
    # =========================

    print_metrics(
        "RULE-BASED",
        rule_metrics,
    )

    print_metrics(
        "MACHINE LEARNING",
        ml_metrics,
    )

    # =========================
    # COMPARISON
    # =========================

    print("=== FAIR COMPARISON ===")
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

    # =========================
    # PER-SAMPLE COMPARISON
    # =========================

    print("=== HELD-OUT TEST RESULTS ===")
    print()

    for text, actual, rule_pred, ml_pred in zip(
        X_test,
        y_test,
        rule_predictions,
        ml_predictions,
    ):
        rule_status = (
            "PASS"
            if rule_pred == actual
            else "FAIL"
        )

        ml_status = (
            "PASS"
            if ml_pred == actual
            else "FAIL"
        )

        print(f"Input    : {text}")
        print(f"Expected : {actual}")
        print(
            f"Rule     : {rule_pred} "
            f"[{rule_status}]"
        )
        print(
            f"ML       : {ml_pred} "
            f"[{ml_status}]"
        )
        print("-" * 60)

    print()
    print("=== ML CLASSIFICATION REPORT ===")
    print()

    print(
        classification_report(
            y_test,
            ml_predictions,
            zero_division=0,
        )
    )

    # =========================
    # SAVE TRAINED MODEL
    # =========================

    MODEL_PATH.parent.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        ml_model,
        MODEL_PATH,
    )

    print(
        f"Model saved: {MODEL_PATH}"
    )

    print()
    print(
        "Fair held-out benchmark selesai."
    )


if __name__ == "__main__":
    main()