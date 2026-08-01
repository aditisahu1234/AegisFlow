from pathlib import Path

import matplotlib.pyplot as plt
import numpy as np

from sklearn.metrics import (
    ConfusionMatrixDisplay,
    confusion_matrix,
)


PLOTS_DIR = Path(
    "artifacts/plots"
)


def save_confusion_matrix(
    model,
    X_test,
    y_test,
    filename: str,
    title: str,
):
    """
    Generate and save a publication-quality
    confusion matrix.

    Parameters
    ----------
    model
        Trained classifier.

    X_test
        Test feature matrix.

    y_test
        Ground-truth labels.

    filename
        Output png filename.

    title
        Plot title.
    """

    PLOTS_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    predictions = model.predict(
        X_test
    )

    labels = np.unique(y_test)

    matrix = confusion_matrix(
        y_test,
        predictions,
        labels=labels,
    )

    fig, ax = plt.subplots(
        figsize=(7, 6)
    )

    display = ConfusionMatrixDisplay(
        confusion_matrix=matrix,
        display_labels=labels,
    )

    display.plot(
        cmap="Blues",
        colorbar=False,
        values_format="d",
        ax=ax,
    )

    ax.set_title(title)

    plt.tight_layout()

    plt.savefig(
        PLOTS_DIR / filename,
        dpi=300,
        bbox_inches="tight",
    )

    plt.close(fig)