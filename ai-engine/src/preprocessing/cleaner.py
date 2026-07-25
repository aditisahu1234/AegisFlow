import joblib
import pandas as pd

from src.config import (
    LEAKAGE_COLUMNS,
    NUMERIC_COLUMNS,
    CATEGORICAL_COLUMNS,
    PREPROCESSING_ARTIFACTS_DIR,
)


IMPUTER_PATH = (
    PREPROCESSING_ARTIFACTS_DIR
    / "imputer.joblib"
)


def remove_leakage_columns(
    df: pd.DataFrame,
) -> pd.DataFrame:
    """
    Remove identifier/leakage columns that must not
    participate in model training.
    """

    columns = [
        column
        for column in LEAKAGE_COLUMNS
        if column in df.columns
    ]

    return df.drop(columns=columns)


def fit_imputer(
    train_df: pd.DataFrame,
) -> tuple[pd.DataFrame, dict]:
    """
    Learn missing-value replacement values from
    training data only.

    Numeric columns -> median
    Categorical columns -> mode
    """

    train_df = train_df.copy()

    imputer_values = {
        "numeric": {},
        "categorical": {},
    }

    for column in NUMERIC_COLUMNS:

        median = train_df[column].median()

        imputer_values["numeric"][column] = median

        train_df[column] = train_df[column].fillna(
            median
        )

    for column in CATEGORICAL_COLUMNS:

        mode = train_df[column].mode()[0]

        imputer_values["categorical"][column] = mode

        train_df[column] = train_df[column].fillna(
            mode
        )

    return train_df, imputer_values


def transform_imputer(
    df: pd.DataFrame,
    imputer_values: dict,
) -> pd.DataFrame:
    """
    Apply missing-value replacements learned
    from the training dataset.
    """

    df = df.copy()

    for column in NUMERIC_COLUMNS:

        value = imputer_values["numeric"][column]

        df[column] = df[column].fillna(
            value
        )

    for column in CATEGORICAL_COLUMNS:

        value = imputer_values["categorical"][column]

        df[column] = df[column].fillna(
            value
        )

    return df


def save_imputer(
    imputer_values: dict,
) -> None:
    """
    Persist training-time imputation values.
    """

    PREPROCESSING_ARTIFACTS_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        imputer_values,
        IMPUTER_PATH,
    )


def load_imputer() -> dict:
    """
    Load the imputation values learned during training.
    """

    return joblib.load(
        IMPUTER_PATH
    )