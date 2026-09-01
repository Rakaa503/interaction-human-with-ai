from pathlib import Path

import joblib


BASE_DIR = Path(__file__).resolve().parent.parent
MODEL_PATH = BASE_DIR / "models" / "intent_model.joblib"


def load_model():
    if not MODEL_PATH.exists():
        raise FileNotFoundError(
            f"Model tidak ditemukan: {MODEL_PATH}"
        )

    return joblib.load(MODEL_PATH)


def predict(model, text: str):
    prediction = model.predict([text])[0]

    probabilities = model.predict_proba([text])[0]
    confidence = float(probabilities.max())

    return prediction, confidence


def main():
    print("=== AVIGO ML PREDICTION ===")
    print()

    model = load_model()

    print("Model loaded successfully.")
    print()

    test_inputs = [
        "Halo AVIGO",
        "Bagaimana cara membuat REST API?",
        "Aplikasi saya mengalami bug saat login",
        "Tolong bantu analisis kode saya",
        "Saya sedang belajar teknologi",
    ]

    for text in test_inputs:
        intent, confidence = predict(
            model,
            text,
        )

        print(f"Input      : {text}")
        print(f"Intent     : {intent}")
        print(f"Confidence : {confidence:.4f}")
        print("-" * 50)


if __name__ == "__main__":
    main()