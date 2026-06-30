import joblib
import pandas as pd

from sklearn.preprocessing import OneHotEncoder

from src.config import (
    CATEGORICAL_COLUMNS,
    ARTIFACTS_DIR,
)


def fit_encoder(
    train_df: pd.DataFrame,
):

    encoder = OneHotEncoder(
        sparse_output=False,
        handle_unknown="ignore",
    )

    encoded = encoder.fit_transform(
        train_df[CATEGORICAL_COLUMNS]
    )

    encoded_df = pd.DataFrame(
        encoded,
        columns=encoder.get_feature_names_out(
            CATEGORICAL_COLUMNS
        ),
        index=train_df.index,
    )

    train_df = train_df.drop(
        columns=CATEGORICAL_COLUMNS
    )

    train_df = pd.concat(
        [
            train_df,
            encoded_df,
        ],
        axis=1,
    )

    return train_df, encoder


def transform_encoder(
    test_df: pd.DataFrame,
    encoder,
):

    encoded = encoder.transform(
        test_df[CATEGORICAL_COLUMNS]
    )

    encoded_df = pd.DataFrame(
        encoded,
        columns=encoder.get_feature_names_out(
            CATEGORICAL_COLUMNS
        ),
        index=test_df.index,
    )

    test_df = test_df.drop(
        columns=CATEGORICAL_COLUMNS
    )

    test_df = pd.concat(
        [
            test_df,
            encoded_df,
        ],
        axis=1,
    )

    return test_df


def save_encoder(
    encoder,
):

    ARTIFACTS_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        encoder,
        ARTIFACTS_DIR /
        "one_hot_encoder.joblib",
    )


def load_encoder():

    return joblib.load(
        ARTIFACTS_DIR /
        "one_hot_encoder.joblib",
    )