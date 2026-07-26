import joblib
import numpy as np
import pandas as pd

from sklearn.ensemble import HistGradientBoostingClassifier

from src.config import MODEL_ARTIFACTS_DIR


MODEL_PATH = (
    MODEL_ARTIFACTS_DIR
    / "hist_gradient_boosting.joblib"
)


def train_hist_gradient_boosting(
    X_train: pd.DataFrame,
    y_train: pd.Series,
    sample_weights: np.ndarray,
    learning_rate: float = 0.1,
    max_iter: int = 100,
    random_state: int = 42,
) -> HistGradientBoostingClassifier:
    """
    Train the imbalance-aware HGB classifier described
    in Section 3.3.5.

    Training representation:
        D_tilde = {(x_tilde_i, y_i, omega_i)}

    omega_i is supplied through sample_weight.
    """

    if len(X_train) != len(y_train):
        raise ValueError(
            "X_train and y_train have different lengths."
        )

    if len(X_train) != len(sample_weights):
        raise ValueError(
            "X_train and sample_weights have different lengths."
        )

    if not np.all(np.isfinite(sample_weights)):
        raise ValueError(
            "Sample weights contain NaN or infinity."
        )

    if np.any(sample_weights <= 0):
        raise ValueError(
            "All sample weights must be positive."
        )

    model = HistGradientBoostingClassifier(
        learning_rate=learning_rate,
        max_iter=max_iter,
        random_state=random_state,
    )

    model.fit(
        X_train,
        y_train,
        sample_weight=sample_weights,
    )

    return model


def save_hist_gradient_boosting(
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


def load_hist_gradient_boosting(
) -> HistGradientBoostingClassifier:

    return joblib.load(
        MODEL_PATH
    )