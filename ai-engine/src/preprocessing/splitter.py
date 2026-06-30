from sklearn.model_selection import train_test_split

from src.config import TARGET_COLUMN


def split_dataset(
    df,
    test_size=0.2,
    random_state=42,
):

    X = df.drop(columns=[TARGET_COLUMN])

    y = df[TARGET_COLUMN]

    X_train, X_test, y_train, y_test = train_test_split(
        X,
        y,
        test_size=test_size,
        stratify=y,
        random_state=random_state,
    )

    return (
        X_train,
        X_test,
        y_train,
        y_test,
    )