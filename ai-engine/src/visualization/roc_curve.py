from pathlib import Path

import matplotlib.pyplot as plt
import numpy as np

from sklearn.metrics import (
    roc_curve,
    auc,
)
from sklearn.preprocessing import (
    label_binarize,
)

from src.config import ARTIFACTS_DIR


PLOTS_DIR = ARTIFACTS_DIR / "plots"

PLOTS_DIR.mkdir(
    parents=True,
    exist_ok=True,
)


def save_multiclass_roc_curve(
    model,
    X_test,
    y_test,
    filename: str,
    title: str,
):
    """
    Save a One-vs-Rest multiclass ROC curve.
    """

    classes = model.classes_

    y_score = model.predict_proba(
        X_test
    )

    y_true = label_binarize(
        y_test,
        classes=classes,
    )

    plt.figure(
        figsize=(8, 6)
    )

    for i, class_name in enumerate(classes):

        fpr, tpr, _ = roc_curve(
            y_true[:, i],
            y_score[:, i],
        )

        roc_auc = auc(
            fpr,
            tpr,
        )

        plt.plot(
            fpr,
            tpr,
            linewidth=2,
            label=f"{class_name} (AUC={roc_auc:.3f})",
        )

    plt.plot(
        [0, 1],
        [0, 1],
        linestyle="--",
        linewidth=1,
    )

    plt.title(title)

    plt.xlabel(
        "False Positive Rate"
    )

    plt.ylabel(
        "True Positive Rate"
    )

    plt.legend()

    plt.tight_layout()

    plt.savefig(
        PLOTS_DIR / filename,
        dpi=300,
    )

    plt.close()