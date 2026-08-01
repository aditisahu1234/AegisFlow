import joblib
import numpy as np
import pandas as pd

from sklearn.ensemble import (
    HistGradientBoostingClassifier,
)

from src.config import (
    MODEL_ARTIFACTS_DIR,
)


MODEL_PATH = (
    MODEL_ARTIFACTS_DIR
    / "hist_gradient_boosting_baseline.joblib"
)


def train_hgb_baseline(
    X_train: pd.DataFrame,
    y_train: pd.Series,
    sample_weights: np.ndarray,
) -> HistGradientBoostingClassifier:
    """
    Train the paper's Histogram Gradient Boosting
    baseline.

    This baseline DOES NOT use:

        • Isolation Forest
        • Deviation score
        • ECDF
        • Normality Fusion

    It is trained directly on the original
    behavioural features after preprocessing,
    ROS and SMOTE.

    Hyperparameters follow Table 3 of the paper.
    """

    if len(X_train) != len(y_train):
        raise ValueError(
            "X_train and y_train have different lengths."
        )

    if len(X_train) != len(sample_weights):
        raise ValueError(
            "Sample weights have different length."
        )

    if not np.all(
        np.isfinite(sample_weights)
    ):
        raise ValueError(
            "Sample weights contain NaN."
        )

    if np.any(sample_weights <= 0):
        raise ValueError(
            "Sample weights must be positive."
        )

    model = HistGradientBoostingClassifier(
        max_depth=7,
        learning_rate=0.05,
        max_iter=700,
        min_samples_leaf=30,
        early_stopping=False,
        random_state=42,
    )

    model.fit(
        X_train,
        y_train,
        sample_weight=sample_weights,
    )

    return model


def save_hgb_baseline(
    model: HistGradientBoostingClassifier,
) -> None:

    MODEL_ARTIFACTS_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        model,
        MODEL_PATH,
    )


def load_hgb_baseline(
) -> HistGradientBoostingClassifier:

    return joblib.load(
        MODEL_PATH
    )