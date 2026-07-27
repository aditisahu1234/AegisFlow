from typing import Any

import numpy as np

from sklearn.metrics import (
    accuracy_score,
    classification_report,
    confusion_matrix,
    precision_score,
    recall_score,
    f1_score,
    roc_auc_score,
)


def evaluate_classifier(
    model,
    X_test,
    y_test,
) -> dict[str, Any]:
    """
    Evaluate a multiclass classifier.

    Reports:
        - accuracy
        - macro precision/recall/F1
        - weighted precision/recall/F1
        - multiclass ROC-AUC (OvR)
        - confusion matrix
        - per-class classification report

    Macro metrics are particularly important for AegisFlow because
    the dataset is extremely imbalanced. They give every class equal
    importance regardless of its frequency.
    """

    y_pred = model.predict(X_test)

    y_proba = model.predict_proba(X_test)

    classes = model.classes_

    # ==========================================
    # GLOBAL METRICS
    # ==========================================

    accuracy = accuracy_score(
        y_test,
        y_pred,
    )

    # ==========================================
    # MACRO METRICS
    # ==========================================

    precision_macro = precision_score(
        y_test,
        y_pred,
        average="macro",
        zero_division=0,
    )

    recall_macro = recall_score(
        y_test,
        y_pred,
        average="macro",
        zero_division=0,
    )

    f1_macro = f1_score(
        y_test,
        y_pred,
        average="macro",
        zero_division=0,
    )

    # ==========================================
    # WEIGHTED METRICS
    # ==========================================

    precision_weighted = precision_score(
        y_test,
        y_pred,
        average="weighted",
        zero_division=0,
    )

    recall_weighted = recall_score(
        y_test,
        y_pred,
        average="weighted",
        zero_division=0,
    )

    f1_weighted = f1_score(
        y_test,
        y_pred,
        average="weighted",
        zero_division=0,
    )

    # ==========================================
    # MULTICLASS ROC-AUC
    # ==========================================

    try:
        roc_auc_macro_ovr = roc_auc_score(
            y_test,
            y_proba,
            labels=classes,
            multi_class="ovr",
            average="macro",
        )

        roc_auc_weighted_ovr = roc_auc_score(
            y_test,
            y_proba,
            labels=classes,
            multi_class="ovr",
            average="weighted",
        )

    except ValueError:
        roc_auc_macro_ovr = np.nan
        roc_auc_weighted_ovr = np.nan

    # ==========================================
    # CONFUSION MATRIX
    # ==========================================

    cm = confusion_matrix(
        y_test,
        y_pred,
        labels=classes,
    )

    # ==========================================
    # CLASSIFICATION REPORT
    # ==========================================

    report = classification_report(
        y_test,
        y_pred,
        labels=classes,
        zero_division=0,
    )

    return {
        "accuracy": accuracy,

        "precision_macro": precision_macro,
        "recall_macro": recall_macro,
        "f1_macro": f1_macro,

        "precision_weighted": precision_weighted,
        "recall_weighted": recall_weighted,
        "f1_weighted": f1_weighted,

        "roc_auc_macro_ovr": roc_auc_macro_ovr,
        "roc_auc_weighted_ovr": roc_auc_weighted_ovr,

        "classes": classes,
        "confusion_matrix": cm,
        "classification_report": report,
    }