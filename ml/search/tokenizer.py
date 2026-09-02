import re


def tokenize(text: str) -> list[str]:
    text = text.lower()

    tokens = re.findall(r"\b[a-z0-9]+\b", text)

    return tokens