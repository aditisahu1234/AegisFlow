import joblib
import pandas as pd

from sklearn.preprocessing import OneHotEncoder

from src.config import (
    CATEGORICAL_COLUMNS,
    PREPROCESSING_ARTIFACTS_DIR,
)


ENCODER_PATH = (
    PREPROCESSING_ARTIFACTS_DIR
    / "encoder.joblib"
)


def fit_encoder(
    train_df: pd.DataFrame,
) -> tuple[pd.DataFrame, OneHotEncoder]:

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
    df: pd.DataFrame,
    encoder: OneHotEncoder,
) -> pd.DataFrame:

    encoded = encoder.transform(
        df[CATEGORICAL_COLUMNS]
    )

    encoded_df = pd.DataFrame(
        encoded,
        columns=encoder.get_feature_names_out(
            CATEGORICAL_COLUMNS
        ),
        index=df.index,
    )

    df = df.drop(
        columns=CATEGORICAL_COLUMNS
    )

    df = pd.concat(
        [
            df,
            encoded_df,
        ],
        axis=1,
    )

    return df


def save_encoder(
    encoder: OneHotEncoder,
) -> None:

    PREPROCESSING_ARTIFACTS_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        encoder,
        ENCODER_PATH,
    )


def load_encoder() -> OneHotEncoder:

    return joblib.load(
        ENCODER_PATH
    )