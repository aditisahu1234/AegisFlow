import joblib
import numpy as np
import pandas as pd

from sklearn.preprocessing import StandardScaler

from src.config import (
    NUMERIC_COLUMNS,
    PREPROCESSING_ARTIFACTS_DIR,
)


SCALER_PATH = (
    PREPROCESSING_ARTIFACTS_DIR
    / "scaler.joblib"
)


def log_transform(
    df: pd.DataFrame,
) -> pd.DataFrame:
    """
    Apply the paper's distribution-stabilising transformation:

        x' = log(1 + max(x, 0))

    Only numerical features are transformed.
    """

    df = df.copy()

    for column in NUMERIC_COLUMNS:
        df[column] = np.log1p(
            df[column].clip(lower=0)
        )

    return df


def fit_scaler(
    train_df: pd.DataFrame,
) -> tuple[pd.DataFrame, StandardScaler]:
    """
    Fit StandardScaler using training data only
    and transform the training data.
    """

    scaler = StandardScaler()

    train_df = train_df.copy()

    train_df[NUMERIC_COLUMNS] = scaler.fit_transform(
        train_df[NUMERIC_COLUMNS]
    )

    return train_df, scaler


def transform_scaler(
    df: pd.DataFrame,
    scaler: StandardScaler,
) -> pd.DataFrame:
    """
    Transform data using an already-fitted scaler.
    """

    df = df.copy()

    df[NUMERIC_COLUMNS] = scaler.transform(
        df[NUMERIC_COLUMNS]
    )

    return df


def save_scaler(
    scaler: StandardScaler,
) -> None:
    """
    Persist the fitted training scaler for inference.
    """

    PREPROCESSING_ARTIFACTS_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        scaler,
        SCALER_PATH,
    )


def load_scaler() -> StandardScaler:
    """
    Load the exact scaler fitted during training.
    """

    return joblib.load(
        SCALER_PATH
    )