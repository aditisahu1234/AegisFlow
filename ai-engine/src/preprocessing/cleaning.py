import pandas as pd

from src.config import (
    LEAKAGE_COLUMNS,
    NUMERIC_COLUMNS,
    CATEGORICAL_COLUMNS,
)


def remove_leakage_columns(
    df: pd.DataFrame,
) -> pd.DataFrame:
    """
    Remove identifier columns.
    """

    columns = [
        col
        for col in LEAKAGE_COLUMNS
        if col in df.columns
    ]

    return df.drop(columns=columns)


def fill_missing_values(
    df: pd.DataFrame,
) -> pd.DataFrame:

    df = df.copy()

    for column in NUMERIC_COLUMNS:

        median = df[column].median()

        df[column] = df[column].fillna(median)

    for column in CATEGORICAL_COLUMNS:

        mode = df[column].mode()[0]

        df[column] = df[column].fillna(mode)

    return df