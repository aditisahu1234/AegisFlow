from collections import Counter

import pandas as pd
from imblearn.over_sampling import (
    RandomOverSampler,
    SMOTE,
)


def resample_training_data(
    X_train: pd.DataFrame,
    y_train: pd.Series,
    random_state: int = 42,
):
    """
    Apply the paper's two-stage class-balancing strategy:

        RandomOverSampler -> SMOTE

    Resampling must only be performed on training data.
    """

    if len(X_train) != len(y_train):
        raise ValueError(
            "X_train and y_train must contain "
            "the same number of observations."
        )

    print(
        "\nClass distribution before resampling:",
        Counter(y_train),
    )

    # Stage 1: Random Over-Sampling
    ros = RandomOverSampler(
        sampling_strategy="minority",
        random_state=random_state,
    )

    X_ros, y_ros = ros.fit_resample(
        X_train,
        y_train,
    )

    print(
        "Class distribution after ROS:",
        Counter(y_ros),
    )

    # Stage 2: SMOTE
    smote = SMOTE(
        sampling_strategy="minority",
        random_state=random_state,
        k_neighbors=5,
    )

    X_resampled, y_resampled = (
        smote.fit_resample(
            X_ros,
            y_ros,
        )
    )

    print(
        "Class distribution after SMOTE:",
        Counter(y_resampled),
    )

    return (
        X_resampled,
        y_resampled,
    )