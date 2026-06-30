import joblib
import pandas as pd
import numpy as np

from sklearn.preprocessing import StandardScaler

from src.config import (
    NUMERIC_COLUMNS,
    ARTIFACTS_DIR,

)
def log_transform(
    df: pd.DataFrame,
) -> pd.DataFrame:

    df = df.copy()

    for column in NUMERIC_COLUMNS:

        df[column] = np.log1p(
            df[column].clip(lower=0)
        )

    return df

def fit_scaler(
    train_df: pd.DataFrame,
):

    scaler = StandardScaler()

    train_df = train_df.copy()

    train_df[NUMERIC_COLUMNS] = scaler.fit_transform(
        train_df[NUMERIC_COLUMNS]
    )

    return (
        train_df,
        scaler,
    )


def transform_scaler(
    test_df: pd.DataFrame,
    scaler,
):

    test_df = test_df.copy()

    test_df[NUMERIC_COLUMNS] = scaler.transform(
        test_df[NUMERIC_COLUMNS]
    )

    return test_df

def save_scaler(
    scaler,
):
    ARTIFACTS_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        scaler,
        ARTIFACTS_DIR /
        "standard_scaler.joblib",
    )

def load_scaler():

    return joblib.load(
        ARTIFACTS_DIR /
        "standard_scaler.joblib"
    )