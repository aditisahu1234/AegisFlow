from src.preprocessing.loader import load_dataset

from src.preprocessing.cleaning import (
    remove_leakage_columns,
    fill_missing_values,
)

from src.preprocessing.scaler import (
    log_transform,
    fit_scaler,
    transform_scaler,
    save_scaler,
)

from src.preprocessing.encoder import (
    fit_encoder,
    transform_encoder,
    save_encoder,
)

from src.preprocessing.splitter import (
    split_dataset,
)

import pandas as pd

df = load_dataset()

df = remove_leakage_columns(df)

df = fill_missing_values(df)

(
    X_train,
    X_test,
    y_train,
    y_test,
) = split_dataset(df)

X_train = log_transform(X_train)

X_test = log_transform(X_test)

X_train, encoder = fit_encoder(X_train)

save_encoder(
    encoder,
)

X_test = transform_encoder(
    X_test,
    encoder,
)

X_train, scaler = fit_scaler(
    X_train,
)

save_scaler(
    scaler,
)

X_test = transform_scaler(
    X_test,
    scaler,
)

train_df = pd.concat(
    [
        X_train,
        y_train,
    ],
    axis=1,
)

test_df = pd.concat(
    [
        X_test,
        y_test,
    ],
    axis=1,
)

print(train_df.head())

print()

print(train_df.shape)

print(test_df.shape)