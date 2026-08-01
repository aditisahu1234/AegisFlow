from pathlib import Path

import pandas as pd

from src.config import ARTIFACTS_DIR


RESULTS_DIR = ARTIFACTS_DIR / "results"

RESULTS_DIR.mkdir(
    parents=True,
    exist_ok=True,
)


def save_error_analysis(
    model,
    X_test,
    y_test,
):
    """
    Save every prediction error as a table.

    Output:

    artifacts/results/error_analysis.csv
    """

    predictions = model.predict(
        X_test
    )

    errors = pd.DataFrame(
        {
            "True Class": y_test,
            "Predicted Class": predictions,
        }
    )

    errors = (
        errors.groupby(
            [
                "True Class",
                "Predicted Class",
            ]
        )
        .size()
        .reset_index(name="Count")
        .sort_values(
            "Count",
            ascending=False,
        )
    )

    errors.to_csv(
        RESULTS_DIR / "error_analysis.csv",
        index=False,
    )

    return errors