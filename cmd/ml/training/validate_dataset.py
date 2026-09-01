import csv
from collections import Counter
from pathlib import Path


DATASET_PATH = (
    Path(__file__).resolve().parent.parent
    / "dataset"
    / "interactions.csv"
)


REQUIRED_COLUMNS = {
    "text",
    "intent",
    "emotion",
    "topic",
}


def load_dataset():
    with open(
        DATASET_PATH,
        "r",
        encoding="utf-8",
        newline="",
    ) as file:
        return list(csv.DictReader(file))


def validate_rows(rows):
    errors = []

    for index, row in enumerate(rows, start=2):
        for column in REQUIRED_COLUMNS:
            if not row.get(column):
                errors.append(
                    f"Row {index}: missing {column}"
                )

        if row.get("text"):
            if not row["text"].strip():
                errors.append(
                    f"Row {index}: empty text"
                )

    return errors


def main():
    print("=== AVIGO DATASET VALIDATION ===")
    print()

    if not DATASET_PATH.exists():
        print("Dataset not found:")
        print(DATASET_PATH)
        return

    rows = load_dataset()

    print(f"Dataset : {DATASET_PATH}")
    print(f"Rows    : {len(rows)}")
    print()

    if not rows:
        print("ERROR: dataset is empty")
        return

    columns = set(rows[0].keys())

    missing_columns = REQUIRED_COLUMNS - columns

    if missing_columns:
        print(
            "ERROR: missing columns:",
            ", ".join(sorted(missing_columns)),
        )
        return

    errors = validate_rows(rows)

    if errors:
        print("VALIDATION FAILED")

        for error in errors:
            print("-", error)

        return

    print("Validation : PASS")
    print()

    print("=== INTENT DISTRIBUTION ===")

    intents = Counter(
        row["intent"]
        for row in rows
    )

    for label, count in intents.items():
        print(f"{label}: {count}")

    print()

    print("=== EMOTION DISTRIBUTION ===")

    emotions = Counter(
        row["emotion"]
        for row in rows
    )

    for label, count in emotions.items():
        print(f"{label}: {count}")

    print()

    print("=== TOPIC DISTRIBUTION ===")

    topics = Counter(
        row["topic"]
        for row in rows
    )

    for label, count in topics.items():
        print(f"{label}: {count}")

    print()
    print("Dataset validation completed.")


if __name__ == "__main__":
    main()