from pathlib import Path

import joblib
import matplotlib.pyplot as plt
import pandas as pd

from sklearn.metrics import ConfusionMatrixDisplay
from sklearn.model_selection import train_test_split


BASE_DIR = Path(__file__).resolve().parent.parent

DATASET_PATH = BASE_DIR / "dataset" / "interactions.csv"
MODEL_PATH = BASE_DIR / "models" / "intent_model.joblib"
OUTPUT_PATH = BASE_DIR / "evaluation" / "confusion_matrix.png"

RANDOM_STATE = 42
TEST_SIZE = 0.20


def main():
    print("=== AVIGO ML EVALUATION ===")
    print()

    df = pd.read_csv(DATASET_PATH)

    df = df.dropna(
        subset=["text", "intent"]
    )

    X = df["text"].astype(str)
    y = df["intent"].astype(str)

    _, X_test, _, y_test = train_test_split(
        X,
        y,
        test_size=TEST_SIZE,
        random_state=RANDOM_STATE,
        stratify=y,
    )

    model = joblib.load(MODEL_PATH)

    predictions = model.predict(X_test)

    print(f"Test samples: {len(X_test)}")
    print()

    print("=== PREDICTIONS ===")

    for text, actual, predicted in zip(
        X_test,
        y_test,
        predictions,
    ):
        status = "PASS" if actual == predicted else "FAIL"

        print(f"[{status}]")
        print(f"Input     : {text}")
        print(f"Expected  : {actual}")
        print(f"Predicted : {predicted}")
        print()

    display = ConfusionMatrixDisplay.from_predictions(
        y_test,
        predictions,
        xticks_rotation="vertical",
    )

    display.figure_.tight_layout()

    OUTPUT_PATH.parent.mkdir(
        parents=True,
        exist_ok=True,
    )

    display.figure_.savefig(
        OUTPUT_PATH,
        dpi=150,
    )

    plt.close(display.figure_)

    print("=== OUTPUT ===")
    print(f"Confusion matrix: {OUTPUT_PATH}")


if __name__ == "__main__":
    main()