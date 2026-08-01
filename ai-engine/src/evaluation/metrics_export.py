from pathlib import Path
import json

from src.config import ARTIFACTS_DIR


RESULTS_DIR = (
    ARTIFACTS_DIR
    / "results"
)


RESULTS_DIR.mkdir(
    parents=True,
    exist_ok=True,
)


def save_metrics(
    metrics: dict,
    filename: str,
) -> None:
    """
    Save evaluation metrics as JSON.

    Example:
        fusion_metrics.json
        hgb_metrics.json
        logistic_metrics.json
    """

    metrics_to_save = {
        key: value
        for key, value in metrics.items()
        if isinstance(
            value,
            (
                int,
                float,
                str,
                bool,
            ),
        )
    }

    with open(
        RESULTS_DIR / filename,
        "w",
    ) as f:

        json.dump(
            metrics_to_save,
            f,
            indent=4,
        )