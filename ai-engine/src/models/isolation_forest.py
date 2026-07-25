from pathlib import Path

import joblib
import pandas as pd

from sklearn.ensemble import IsolationForest

from src.config import MODEL_ARTIFACTS_DIR
import numpy as np
import pandas as pd


def compute_deviation_scores(
    model: IsolationForest,
    X: pd.DataFrame,
) -> np.ndarray:
    """
    Compute a continuous deviation-from-normality score.

    sklearn IsolationForest.score_samples():
        higher value -> more normal
        lower value  -> more anomalous

    Normality Fusion requires the opposite orientation:
        higher deviation -> more abnormal

    Therefore:
        deviation(x) = -score_samples(x)
    """

    normality_scores = model.score_samples(X)

    deviation_scores = -normality_scores

    return deviation_scores

ISOLATION_FOREST_PATH = (
    MODEL_ARTIFACTS_DIR
    / "isolation_forest.joblib"
)


def extract_normal_samples(
    X_train: pd.DataFrame,
    y_train: pd.Series,
) -> pd.DataFrame:
    """
    Extract only normal training samples.

    Paper Eq. (12):

        X_N = {x_i | y_i = normal}
    """

    normal_mask = y_train == "normal"

    X_normal = X_train.loc[normal_mask]

    if X_normal.empty:
        raise ValueError(
            "No normal samples found."
        )

    return X_normal


def fit_isolation_forest(
    X_normal: pd.DataFrame,
) -> IsolationForest:
    """
    Learn the normal behaviour manifold.

    Hyperparameters follow the paper.
    """

    model = IsolationForest(
        n_estimators=800,
        contamination=0.01,
        random_state=42,
        n_jobs=-1,
    )

    model.fit(X_normal)

    return model


def save_isolation_forest(
    model: IsolationForest,
) -> None:

    MODEL_ARTIFACTS_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        model,
        ISOLATION_FOREST_PATH,
    )


def load_isolation_forest() -> IsolationForest:

    return joblib.load(
        ISOLATION_FOREST_PATH
    )