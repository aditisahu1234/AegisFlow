import pandas as pd

from src.config import DATASET_PATH


def load_dataset() -> pd.DataFrame:
    """
    Load the supervised dataset.
    """

    df = pd.read_csv(DATASET_PATH)

    return df