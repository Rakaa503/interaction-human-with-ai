from pathlib import Path

import pandas as pd

from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.model_selection import StratifiedKFold, cross_validate
from sklearn.pipeline import Pipeline


BASE_DIR = Path(__file__).resolve().parent.parent
DATASET_PATH = BASE_DIR / "dataset" / "interactions.csv"

RANDOM_STATE = 42
N_SPLITS = 5


def load_dataset() -> pd.DataFrame:
    if not DATASET_PATH.exists():
        raise FileNotFoundError(
            f"Dataset tidak ditemukan: {DATASET_PATH}"
        )

    df = pd.read_csv(DATASET_PATH)

    required_columns = {"text", "intent"}

    missing_columns = required_columns - set(df.columns)

    if missing_columns:
        raise ValueError(
            "Kolom dataset tidak lengkap: "
            + ", ".join(sorted(missing_columns))
        )

    df = df.dropna(
        subset=["text", "intent"]
    )

    df["text"] = (
        df["text"]
        .astype(str)
        .str.strip()
    )

    df["intent"] = (
        df["intent"]
        .astype(str)
        .str.strip()
    )

    df = df[
        (df["text"] != "")
        & (df["intent"] != "")
    ]

    if df.empty:
        raise ValueError(
            "Dataset kosong."
        )

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


def main():
    print("=== AVIGO CROSS-VALIDATION ===")
    print()

    df = load_dataset()

    X = df["text"]
    y = df["intent"]

    print(f"Dataset samples : {len(df)}")
    print(f"Number of folds : {N_SPLITS}")
    print()

    model = build_model()

    cv = StratifiedKFold(
        n_splits=N_SPLITS,
        shuffle=True,
        random_state=RANDOM_STATE,
    )

    scoring = [
        "accuracy",
        "precision_macro",
        "recall_macro",
        "f1_macro",
    ]

    results = cross_validate(
        model,
        X,
        y,
        cv=cv,
        scoring=scoring,
    )

    accuracy = results[
        "test_accuracy"
    ]

    precision = results[
        "test_precision_macro"
    ]

    recall = results[
        "test_recall_macro"
    ]

    f1 = results[
        "test_f1_macro"
    ]

    print("=== FOLD RESULTS ===")
    print()

    for index in range(N_SPLITS):
        print(
            f"Fold {index + 1}: "
            f"accuracy={accuracy[index]:.4f} "
            f"precision={precision[index]:.4f} "
            f"recall={recall[index]:.4f} "
            f"f1={f1[index]:.4f}"
        )

    print()

    print("=== CROSS-VALIDATION SUMMARY ===")
    print()

    print(
        f"Accuracy : "
        f"{accuracy.mean():.4f} "
        f"+/- {accuracy.std():.4f}"
    )

    print(
        f"Precision: "
        f"{precision.mean():.4f} "
        f"+/- {precision.std():.4f}"
    )

    print(
        f"Recall   : "
        f"{recall.mean():.4f} "
        f"+/- {recall.std():.4f}"
    )

    print(
        f"F1 Score : "
        f"{f1.mean():.4f} "
        f"+/- {f1.std():.4f}"
    )

    print()
    print(
        "Cross-validation selesai."
    )


if __name__ == "__main__":
    main()