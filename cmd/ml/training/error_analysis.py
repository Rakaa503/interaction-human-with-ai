from pathlib import Path

import joblib
import pandas as pd


BASE_DIR = Path(__file__).resolve().parent.parent

DATASET_PATH = BASE_DIR / "dataset" / "interactions.csv"
MODEL_PATH = BASE_DIR / "models" / "intent_model.joblib"
OUTPUT_PATH = BASE_DIR / "evaluation" / "error_analysis.csv"


def main():
    print("=== AVIGO ML ERROR ANALYSIS ===")
    print()

    if not DATASET_PATH.exists():
        raise FileNotFoundError(
            f"Dataset tidak ditemukan: {DATASET_PATH}"
        )

    if not MODEL_PATH.exists():
        raise FileNotFoundError(
            f"Model tidak ditemukan: {MODEL_PATH}"
        )

    df = pd.read_csv(DATASET_PATH)

    df = df.dropna(
        subset=["text", "intent"]
    ).copy()

    model = joblib.load(
        MODEL_PATH
    )

    texts = df["text"].astype(str)
    expected = df["intent"].astype(str)

    predictions = model.predict(texts)
    probabilities = model.predict_proba(texts)

    confidence = probabilities.max(
        axis=1
    )

    result = pd.DataFrame(
        {
            "input": texts,
            "expected": expected,
            "predicted": predictions,
            "confidence": confidence,
        }
    )

    errors = result[
        result["expected"]
        != result["predicted"]
    ].copy()

    errors = errors.sort_values(
        by="confidence",
        ascending=True,
    )

    OUTPUT_PATH.parent.mkdir(
        parents=True,
        exist_ok=True,
    )

    errors.to_csv(
        OUTPUT_PATH,
        index=False,
    )

    print(
        f"Total samples : {len(result)}"
    )

    print(
        f"Errors        : {len(errors)}"
    )

    print(
        f"Accuracy      : "
        f"{1 - (len(errors) / len(result)):.4f}"
    )

    print()

    if errors.empty:
        print(
            "Tidak ditemukan error."
        )
    else:
        print(
            "=== MISCLASSIFIED SAMPLES ==="
        )
        print()

        for _, row in errors.iterrows():
            print(
                f"Input      : {row['input']}"
            )

            print(
                f"Expected   : {row['expected']}"
            )

            print(
                f"Predicted  : {row['predicted']}"
            )

            print(
                f"Confidence : "
                f"{row['confidence']:.4f}"
            )

            print(
                "-" * 60
            )

    print()

    print(
        f"Output: {OUTPUT_PATH}"
    )

    print()
    print(
        "Error analysis selesai."
    )


if __name__ == "__main__":
    main()