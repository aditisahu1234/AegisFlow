import pandas as pd

from imblearn.over_sampling import (
    RandomOverSampler,
    SMOTE,
)


ATTACK_CLASS = "attack"
BOT_CLASS = "bot"
NORMAL_CLASS = "normal"
OUTLIER_CLASS = "outlier"


def apply_ros_smote(
    X_train: pd.DataFrame,
    y_train: pd.Series,
    random_state: int = 42,
):
    """
    Apply the two-stage imbalance correction
    described in Section 3.2.4:

        1. Random Over-Sampling (ROS)
           attack -> majority-class count

        2. SMOTE
           bot and normal -> majority-class count

    Resampling must be applied to TRAINING DATA ONLY.

    The held-out test set must never be passed here.
    """

    if len(X_train) != len(y_train):
        raise ValueError(
            "X_train and y_train have different lengths."
        )

    class_counts = y_train.value_counts()

    required_classes = {
        ATTACK_CLASS,
        BOT_CLASS,
        NORMAL_CLASS,
        OUTLIER_CLASS,
    }

    missing_classes = (
        required_classes
        - set(class_counts.index)
    )

    if missing_classes:
        raise ValueError(
            f"Missing training classes: "
            f"{sorted(missing_classes)}"
        )

    majority_count = int(
        class_counts.max()
    )

    # ==========================================
    # STAGE 1: RANDOM OVER-SAMPLING
    # ==========================================

    ros = RandomOverSampler(
        sampling_strategy={
            ATTACK_CLASS: majority_count,
        },
        random_state=random_state,
    )

    X_ros, y_ros = ros.fit_resample(
        X_train,
        y_train,
    )

    ros_counts = y_ros.value_counts()

    # ==========================================
    # STAGE 2: SMOTE
    # ==========================================

    smote = SMOTE(
        sampling_strategy={
            BOT_CLASS: majority_count,
            NORMAL_CLASS: majority_count,
        },
        random_state=random_state,
        k_neighbors=5,
    )

    X_resampled, y_resampled = (
        smote.fit_resample(
            X_ros,
            y_ros,
        )
    )

    final_counts = (
        y_resampled.value_counts()
    )

    return (
        X_resampled,
        y_resampled,
        {
            "before": class_counts,
            "after_ros": ros_counts,
            "after_smote": final_counts,
        },
    )