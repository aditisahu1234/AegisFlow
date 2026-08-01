from __future__ import annotations

import pandas as pd

from src.models.hist_gradient_boosting import (
    train_hist_gradient_boosting,
)

from src.evaluation.evaluator import (
    evaluate_classifier,
)


def run_ablation_experiment(
    *,
    experiment_name: str,
    X_train: pd.DataFrame,
    X_test: pd.DataFrame,
    y_train,
    y_test,
    sample_weights,
):
    """
    Train one ablation experiment.

    Returns a dictionary of metrics.
    """

    model = train_hist_gradient_boosting(
        X_train=X_train,
        y_train=y_train,
        sample_weights=sample_weights,
    )

    evaluation = evaluate_classifier(
        model=model,
        X_test=X_test,
        y_test=y_test,
    )

    return {
        "Model": experiment_name,
        "Accuracy": evaluation["accuracy"],
        "Macro F1": evaluation["f1_macro"],
        "Weighted F1": evaluation["f1_weighted"],
        "ROC-AUC": evaluation["roc_auc_macro_ovr"],
    }