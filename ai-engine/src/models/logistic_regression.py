import joblib
import numpy as np
import pandas as pd

from sklearn.linear_model import LogisticRegression

from src.config import MODEL_ARTIFACTS_DIR


MODEL_PATH = (
    MODEL_ARTIFACTS_DIR
    / "logistic_regression.joblib"
)


def train_logistic_regression(
    X_train: pd.DataFrame,
    y_train: pd.Series,
    sample_weights: np.ndarray,
) -> LogisticRegression:
    """
    Train the Logistic Regression baseline
    described in Section 3.4.1.

    Pipeline:

        Original Features
                ↓
           ROS + SMOTE
                ↓
        Class-prior weights
                ↓
      Logistic Regression

    Parameters follow the paper where
    specified, with reproducibility enabled.
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

    model = LogisticRegression(
    max_iter=1000,
    random_state=42,
    solver="lbfgs",
    )

    model.fit(
        X_train,
        y_train,
        sample_weight=sample_weights,
    )

    return model


def save_logistic_regression(
    model: LogisticRegression,
) -> None:

    MODEL_ARTIFACTS_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        model,
        MODEL_PATH,
    )


def load_logistic_regression(
) -> LogisticRegression:

    return joblib.load(
        MODEL_PATH
    )
