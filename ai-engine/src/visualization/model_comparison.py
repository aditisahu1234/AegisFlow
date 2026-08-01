from pathlib import Path

import matplotlib.pyplot as plt
import pandas as pd


PLOTS_DIR = Path("artifacts/plots")
PLOTS_DIR.mkdir(parents=True, exist_ok=True)


def save_metric_barplot(
    comparison_df: pd.DataFrame,
    metric: str,
    filename: str,
) -> None:
    """
    Save a publication-style comparison plot.
    """

    plt.figure(figsize=(7, 5))

    plt.bar(
        comparison_df["Model"],
        comparison_df[metric],
    )

    plt.ylabel(metric)
    plt.title(f"Model Comparison ({metric})")

    plt.ylim(
        max(0, comparison_df[metric].min() - 0.02),
        1.01,
    )

    plt.grid(
        axis="y",
        alpha=0.3,
    )

    plt.tight_layout()

    plt.savefig(
        PLOTS_DIR / filename,
        dpi=300,
    )

    plt.close()


def save_all_metric_plots(
    comparison_df: pd.DataFrame,
):
    save_metric_barplot(
        comparison_df,
        "Accuracy",
        "accuracy.png",
    )

    save_metric_barplot(
        comparison_df,
        "Macro F1",
        "macro_f1.png",
    )

    save_metric_barplot(
        comparison_df,
        "Weighted F1",
        "weighted_f1.png",
    )

    save_metric_barplot(
        comparison_df,
        "ROC-AUC",
        "roc_auc.png",
    )