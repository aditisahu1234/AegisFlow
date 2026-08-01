from pathlib import Path

import matplotlib.pyplot as plt
import pandas as pd

from sklearn.inspection import permutation_importance

from src.config import ARTIFACTS_DIR


PLOTS_DIR = ARTIFACTS_DIR / "plots"

PLOTS_DIR.mkdir(
    parents=True,
    exist_ok=True,
)


def save_feature_importance(
    model,
    X_test,
    y_test,
    filename: str = "feature_importance.png",
):
    """
    Compute permutation importance for the
    Normality Fusion HGB model.

    Saves:
        artifacts/plots/feature_importance.png
    """

    result = permutation_importance(
        estimator=model,
        X=X_test,
        y=y_test,
        n_repeats=10,
        random_state=42,
        n_jobs=-1,
        scoring="f1_weighted",
    )

    importance = pd.DataFrame(
        {
            "Feature": X_test.columns,
            "Importance": result.importances_mean,
        }
    )

    importance = (
        importance
        .sort_values(
            "Importance",
            ascending=False,
        )
        .head(20)
    )

    plt.figure(
        figsize=(10, 8)
    )

    plt.barh(
        importance["Feature"],
        importance["Importance"],
    )

    plt.gca().invert_yaxis()

    plt.title(
        "Top 20 Permutation Feature Importances"
    )

    plt.xlabel(
        "Decrease in Weighted F1"
    )

    plt.tight_layout()

    plt.savefig(
        PLOTS_DIR / filename,
        dpi=300,
    )

    plt.close()

    return importance